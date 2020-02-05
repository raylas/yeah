extern crate regex;
extern crate reqwest;

use std::process;
use std::error::Error;
use regex::Regex;

const VERSION: &'static str = env!("CARGO_PKG_VERSION");

pub struct Config {
    pub query: String,
}

pub fn help() {
    println!("usage: yeah [complete mac|partial mac]
       yeah [-h|--help] [-v|--version]")
}

pub fn version() {
    println!("yeah version {}", VERSION)
}

impl Config {
    pub fn new(args: &[String]) -> Result<Config, &'static str> {
        if args.len() < 2 {
            return Err("not enough arguments");
        } else if args.len() > 2 {
            return Err("too many arguments");
        } else if args[1] == "-h" || args[1] == "--help" {
            help();
            process::exit(1);
        } else if args[1] == "-v" || args[1] == "--version" {
            version();
            process::exit(1);
        }
        
        let query = args[1].clone();

        Ok(Config { query })
    }
}

pub fn run(config: Config) -> Result<(), Box<dyn Error>> {
    let url = "http://standards-oui.ieee.org/oui.txt";
    let contents = reqwest::get(url)?.text()?;

    let re = Regex::new(r"\s+").unwrap();
    let results = search(&config.query, &contents);

    if results.len() >= 1 {
        for line in results {
            let result = re.replace_all(line, " ").to_string().replace("(base 16)", "::");
            println!("{}", result);
        }
    } else {
        println!{"No OUI matches found"}
    }
        
    Ok(())
}

pub fn search<'a>(query: &str, contents: &'a str) -> Vec<&'a str> {
    let re = Regex::new("[^a-zA-Z0-9]+").unwrap();
    let mut query = re.replace_all(&query, "").to_string().to_uppercase();
    if query.len() > 6 {
        query = query[..6].to_string();
    }

    let mut results = Vec::new();

    for line in contents.lines() {
        if line.to_uppercase().contains(&query) {
            results.push(line);
        }
    }

    results
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn one_result() {
        let query = "F4-92-BF";
        let contents = "\
54BF64     (base 16)	Dell Inc.
8417EF     (base 16)	Technicolor CH USA Inc.
F492BF     (base 16)    Ubiquiti Networks Inc.";

        assert_eq!(
            vec!["F492BF     (base 16)    Ubiquiti Networks Inc."],
            search(query, contents)
        );
    }

}