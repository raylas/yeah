//
// OUI == oui == yes == yeah
//
use std::env;
use std::process;

use yeah::Config;
use yeah::help;

fn main() {
    let args: Vec<String> = env::args().collect();

    let config = Config::new(&args).unwrap_or_else(|err| {
        eprintln!("yeah: {}", err);
        help();
        process::exit(1);
    });

    if let Err(e) = yeah::run(config) {
        eprintln!("Application error: {}", e);
        process::exit(1);
    }
}