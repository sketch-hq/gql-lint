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

func ExpandGlobs(args []string, include []string, ignore []string) ([]string, error) {
	var files []string

	if len(include) > 0 {
		var err error
		include, err = ExpandGlobs(include, []string{}, []string{})
		if err != nil {
			return files, err
		}
	}

	if len(ignore) > 0 {
		var err error
		ignore, err = ExpandGlobs(ignore, []string{}, []string{})
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

			filtered := filterIncludes(matches, include)
			filtered = filterIgnored(filtered, ignore)
			files = append(files, filtered...)
			continue
		}

		// skip any files not matching our include filter
		if len(include) > 0 && !contains(include, source) {
			continue
		}

		// skip any files on the ignore list
		if contains(ignore, source) {
			continue
		}

		files = append(files, source)
	}

	return unique(files), nil
}

func filterIncludes(in []string, includes []string) []string {
	if len(includes) == 0 {
		return in
	}

	var out []string

	for _, s := range in {
		if contains(includes, s) {
			out = append(out, s)
		}
	}

	return out
}

func filterIgnored(in []string, ignore []string) []string {
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
		if filepath.Clean(v) == filepath.Clean(str) {
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

					hits = append(hits, filepath.Clean(path))

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
