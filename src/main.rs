/// Data structure derived from command-line arguments.
#[derive(Debug)]
struct Arguments {
    query: String,
    table: bool
}

fn print_usage() {
    println!("yeah - return the vendor name for a given MAC address");
    println!("usage: yeah [options...] <mac>
 -t, --table        Print output as table
 -h, --help         Get help for commands
 -v, --version      Show version and quit");
    std::process::exit(1);
}

fn print_version() {
    println!("yeah version {}", env!("CARGO_PKG_VERSION"));
    std::process::exit(1);
}

/// Parse command-line arguments, and only accept a single free-floating argument.
/// 
/// Return `Arguments` data structure.
fn parse_args() -> Arguments {
    let mut args = pico_args::Arguments::from_env();

    if args.contains(["-h", "--help"]) {
        print_usage();
    }

    if args.contains(["-v", "--version"]) {
        print_version();
    }

    let table = args.contains(["-t", "--table"]);

    let query: String = args.free_from_str().unwrap();

    let orphans = args.finish();
    if !orphans.is_empty() {
        eprintln!("warning: unused arguments: {:?}", orphans);
    }

    Arguments {
        query,
        table
    }
}

use std::error::Error;

/// Retrieve text file of IEEE OUI vendors.
/// 
/// Return a String containing said contents.
fn get_vendors() -> Result<String, Box<dyn Error>> {
    let url = "https://standards-oui.ieee.org/oui.txt";

    let vendors = reqwest::blocking::get(url)?.text()?;
        
    Ok(vendors)
}

#[test]
fn test_get_vendors() {
    assert_eq!(get_vendors().unwrap().is_empty(), false);
}

use regex::RegexBuilder;

/// Given a MAC address and string of OUI vendor mappings.
/// Function will trim a given MAC address to a max length of 6 characters.
/// 
/// Return a vector of vectors of str containing the mapped `OUI` and `vendor`.
fn find_matches<'a>(mac: &str, vendors: &'a str) -> Result<Vec<Vec<&'a str>>, Box<dyn Error>> {
    let oui = &mac
        .replace(':', "")
        .replace('-', "");

    let len: usize = if oui.len() < 6 { oui.len() } else { 6 };   

    let expr = format!(r"^{}.*\(base 16\)", &oui[..len]);

    let re = RegexBuilder::new(&expr)
        .case_insensitive(true)
        .build()
        .unwrap();
    
    let mut results = Vec::<Vec<&str>>::new();

    vendors
        .lines()
        .filter(|line| re.is_match(line))
        .for_each(|line| {
            let result: Vec<&str> = line.split("(base 16)").collect();
            let oui = result[0].trim();
            let vendor = result[1].trim();
            results.push(vec![oui, vendor]);
        });

    Ok(results)
}

#[test]
fn test_find_matches() {
    let result = vec![vec!["F492BF", "Ubiquiti Networks Inc."]];
    let vendors = "\
54BF64     (base 16)	Dell Inc.
8417EF     (base 16)	Technicolor CH USA Inc.
F492BF     (base 16)    Ubiquiti Networks Inc.";

    assert_eq!(find_matches("F4-92-BF-AC-36-9F", vendors).unwrap(), result);
    assert_eq!(find_matches("F4:92:BF:AC:36:9F", vendors).unwrap(), result);
    assert_eq!(find_matches("f4-92-bf-ac-36", vendors).unwrap(), result);
    assert_eq!(find_matches("f492bfac369f", vendors).unwrap(), result);
    assert_eq!(find_matches("f4:92:bf:ac", vendors).unwrap(), result);
    assert_eq!(find_matches("F4-92-BF", vendors).unwrap(), result);
    assert_eq!(find_matches("F4:92:BF", vendors).unwrap(), result);
    assert_eq!(find_matches("F492BF", vendors).unwrap(), result);
    assert_eq!(find_matches("F492", vendors).unwrap(), result);
    assert_eq!(find_matches("F49", vendors).unwrap(), result);
    assert_eq!(find_matches("f4", vendors).unwrap(), result);
    assert_ne!(find_matches("54BF64", vendors).unwrap(), result);
    assert_ne!(find_matches("E492BF", vendors).unwrap(), result);
}

use prettytable::{Table, row, cell};

fn print_table(results: &[Vec<&str>]) {
    let mut table = Table::new();

    for result in results {
        table.add_row(row![result[0], result[1]]);
    }

    table.printstd();
} 

fn main() {
    let args = parse_args();

    let vendors = match get_vendors() {
        Ok(v) => v,
        Err(e) => {
            eprintln!("error: failed to retrieve vendors list: {:?}", e);
            std::process::exit(1);
        }
    };

    let results = match find_matches(&args.query, &vendors) {
        Ok(v) => v,
        Err(e) => {
            eprintln!("error: failed to match vendors for {}: {:?}", &args.query, e);
            std::process::exit(1);
        }
    };

    if !results.is_empty() {
        match &args.table {
            true => print_table(&results),
            false => {
                for result in &results {
                    println!("{}    {}", &result[0], &result[1]);
                } 
            }
        }
    } else {
        println!("no OUI matches found")
    }
}
