# GQL-lint

## Install

```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/lab/releases/install.sh | /bin/bash -s -- <version>
```

```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/lab/releases/install.sh | /bin/bash -s -- v0
```

## Usage

### Piping arguments

We suggest you use [xargs](https://man7.org/linux/man-pages/man1/xargs.1.html) if you want to pipe arguments into this CLI.
This is very useful for reading from files for example :point_down:

```
$ cat query_file_list.txt | xargs gql-lint unused --schema path/to/schema.gql
```
