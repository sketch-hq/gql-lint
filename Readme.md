# GQL-lint

## Install

```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/lab/releases/install.sh | /bin/bash -s -- <version>
```

```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/lab/releases/install.sh | /bin/bash -s -- v0
```

You can also use `latest` as the version to get whatever is the latest available version.

## Usage

### Piping arguments

We suggest you use [xargs](https://man7.org/linux/man-pages/man1/xargs.1.html) if you want to pipe arguments into this CLI.
This is very useful for reading from files for example :point_down:

```
$ cat query_file_list.txt | xargs gql-lint unused --schema path/to/schema.gql
```

## Testing

Run all tests using `make test`.

### Updating cmdtest tests

The "external interface" of the CLI is tested using `cmdtest-go`. If you makes changes to the output such as the help text it can be bothersome to update all the test files manually. Instead you can set an env var that will make `cmdtest-go` recreate the files with the current output:

```sh
UPDATE_CMD_TESTS=true make test
```

Remember to manually inspect any difference to ensure the new output is correct.
