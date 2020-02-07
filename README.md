# yeah

[![crates.io](https://img.shields.io/crates/v/yeah.svg)](https://crates.io/crates/yeah)

yeah is a command-line tool to return the vendor name for a given MAC address. 

Queries are ran against the [IEEE OUI vendor list](http://standards-oui.ieee.org/oui.txt).

Functionality:

- Supports complete and partial MAC address queries

Todo:

- Add ability to store IEEE OUI list locally for offline queries
- General optimizations

## Install yeah

- With crates.io:

```bash
cargo install yeah
```

- From source:

```bash
cargo install --path /path/to/yeah/repo
```

- Or, download and run from `.zip` or `.tar.gz` in releases

## Use yeah

- Print vendor of given MAC address

```bash
$ yeah <mac>
```

- Options:

```
  -h, --help          Prints help information
  -v, --version       Prints version information
```