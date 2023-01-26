package input

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type command interface {
	InOrStdin() io.Reader
}

func ExpandGlobs(args []string, ignore []string) ([]string, error) {
	var files []string

	if len(ignore) > 0 {
		var err error
		ignore, err = ExpandGlobs(ignore, []string{})
		if err != nil {
			return files, err
		}
	}

	for _, source := range args {
		if strings.Contains(source, "*") {
			matches, err := glob(source)
			if err != nil {
				return files, err
			}

			files = append(files, filter(matches, ignore)...)
		} else if !contains(ignore, source) {
			files = append(files, source)
		}
	}

	return unique(files), nil
}

func filter(in []string, ignore []string) []string {
	var out []string

	for _, s := range in {
		if !contains(ignore, s) {
			out = append(out, s)
		}
	}

	return out
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func glob(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		// passthru to core package if no double-star
		return filepath.Glob(pattern)
	}
	return expand(pattern)
}

// expand finds matches for the provided Globs
func expand(pattern string) ([]string, error) {
	globs := strings.Split(pattern, "**")

	var matches = []string{""}

	for _, glob := range globs {
		var hits []string

		for _, match := range matches {
			paths, err := filepath.Glob(match + glob)
			if err != nil {
				return nil, err
			}

			for _, path := range paths {
				err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					hits = append(hits, strings.Replace(path, `//`, `/`, -1))

					return nil
				})

				if err != nil {
					return nil, err
				}
			}
		}

		matches = hits
	}

	return matches, nil
}
