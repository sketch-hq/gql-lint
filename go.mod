module github.com/sketch-hq/gql-lint

go 1.19

// while we wait for https://github.com/wundergraph/graphql-go-tools/pull/488 to be merged and released
replace github.com/wundergraph/graphql-go-tools => github.com/sketch-hq/graphql-go-tools v0.0.0-20230119134629-3dd6fc431bee

require (
	github.com/google/go-cmdtest v0.4.0
	github.com/matryer/is v1.4.0
	github.com/spf13/cobra v1.6.1
	github.com/vektah/gqlparser/v2 v2.5.1
	github.com/wundergraph/graphql-go-tools v1.60.6
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/renameio v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/gjson v1.11.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/sjson v1.0.4 // indirect
)
