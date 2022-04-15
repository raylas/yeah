# yeah

[![crates.io](https://img.shields.io/crates/v/yeah.svg)](https://crates.io/crates/yeah)

yeah is a command-line tool to return the vendor name for a given MAC address. 

Queries are ran against the [IEEE OUI vendor list](http://standards-oui.ieee.org/oui.txt).

Functionality:

- Complete and partial MAC address queries
- Option to pretty-print an ASCII table for results

Todo:

- [ ] Caching of IEEE OUI list (~5MB)
- [ ] Ability to pass multiple MAC addresses and return multiple match groups

## Install yeah

With crates.io:
```bash
cargo install yeah
```

From source:
```bash
cargo install --path /path/to/yeah/repo
```

Or, download and run the binary from [the latest release](https://github.com/raylas/yeah/releases).

## Use yeah

```bash
yeah - return the vendor name for a given MAC address
usage: yeah [options...] <mac>
 -t, --table        Print output as table
 -h, --help         Get help for commands
 -v, --version      Show version and quit
```
