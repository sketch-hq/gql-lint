# GQL-lint

`gql-lint` is a tool for finding deprecated fields in graphql queries and unused deprecated fields in schemas. While there are many tools that does similar things we wanted a tool that was extremely easy and fast to install, run (locally and in CI) and debug. Most tools in this category are nodejs based and thus fiddly for non-JS developers to install and use.

## Install

The latest version:
```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/main/install.sh | /bin/bash -s -- latest
```

Get a specific version:
```
curl https://raw.githubusercontent.com/sketch-hq/gql-lint/main/install.sh | /bin/bash -s -- v1
```

## Usage

```
gql-lint is a tool to lint GraphQL queries and mutations

Usage:
  gql-lint [command]

Available Commands:
  deprecation Find deprecated fields in queries and mutations given a list of files
  diff        Find deprecated fields present in the first file but not in the second
  help        Help about any command
  unused      Find unused deprecated fields

Flags:
  -h, --help            help for gql-lint
      --output string   Output format. Choose between stdout, json, xcode. (default "stdout")
  -v, --verbose         Verbose mode. Will print debug messages
      --version         version for gql-lint

Use "gql-lint [command] --help" for more information about a command.
```

### Providing a schema

For commands that take a schema one or more schemas can be provided by either repeating the `--schema` parameter:

```
gql-lint deprecation --schema schema1.gql --schema schema2.gql
```

or specifying a comma seperated list of files/urls:

```
gql-lint deprecation --schema schema1.gql,schema2.gql
```

The schema can either be a file or a URL. If a URL is provided `gql-lint` will do an introspection query to the provided URL to extract the schema.

### Providing graphql files

#### Globbing

`gql-lint` supports `**/*` type globbing when specifying graphql files:
```
gql-lint unused  --schema <schema> testdata/queries/unused/**/*.gql
```

Depending on your shell you might have to enclose the argument in single quotes:
```
gql-lint unused  --schema <schema> 'testdata/queries/unused/**/*.gql'
```

#### Include/ignore files

If you need to ignore some files or only match on some files you can do that with the `--ignore` and `--include` parameters. These parameteres also supports the same globbing as above.

```
gql-lint unused --schema <schema> --ignore testdata/queries/unused/without_title/*.gql testdata/queries/unused/**/*.gql
```

```
$ gql-lint unused --include ./**/album.gql --schema testdata/schemas/with_deprecations.gql testdata/queries/unused/**/*.gql
```

#### Piping arguments

We suggest you use [xargs](https://man7.org/linux/man-pages/man1/xargs.1.html) if you want to pipe arguments into this CLI.  This is very useful for reading from files for example :point_down:

```
$ cat query_file_list.txt | xargs gql-lint unused --schema path/to/schema.gql
```

### Find deprecated query fields

```
gql-lint deprecation --schema https://rickandmortyapi.com/graphql <query files>
```

You can also check queries against multiple schemas, assuming there's no overlap between the schemas:

```
gql-lint deprecation --schema https://swapi-graphql.netlify.app/.netlify/functions/index --schema https://rickandmortyapi.com/graphql <query files>
```

### Compare two git branches

A typical use case for `gql-lint` is to find and stop developers adding more deprecated fields to their queries. You can diff two git branches using the `diff` command like this:

```bash
git checkout main
git deprecation --output json --schema <schema file/url> <query files> > base.json
git checkout <pr branch>
git deprecation --output json --schema <schema file/url> <query files> > pr.json
git diff base.json pr.json
Schema:  schema.graphql
Article.isPublic (Please migrate to publicAccessLevel)
  query.gql:7
```

### Find unused deprecated fields

Use the `unused` command to find deprecated fields in the schema that are not used by a set of query files:

```
gql-lint unused --schema cmd/testdata/schemas/with_deprecations.gql 'cmd/testdata/queries/unused/without_title/*.gql'
Schema: cmd/testdata/schemas/with_deprecations.gql
  Book.title (line 1) is unused and can be removed
```

## Testing

Run all tests using `make test`.

### Updating cmdtest tests

The "external interface" of the CLI is tested using `cmdtest-go`. If you makes changes to the output such as the help text it can be bothersome to update all the test files manually. Instead you can set an env var that will make `cmdtest-go` recreate the files with the current output:

```sh
UPDATE_CMD_TESTS=true make test
```

Remember to manually inspect any difference to ensure the new output is correct.