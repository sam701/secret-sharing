# secret-sharing
*Secrets sharing with Shamir's secret scheme*

This is a wrapper around the [Vault](https://github.com/hashicorp/vault)'s implementation of the [Shamir's scheme](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing).

## Build
```sh
glide install
go install
```

## Usage

### Split file
To split a file into 3 parts with threshold (number of parts required to combine the original file) 2:
```sh
secret-sharing split --parts 3 --threshold 2 <file to split>
```

### Combine file
To combine a file from 2 parts:
```sh
secret-sharing combine --output <file name> <part1> <part2>
```
