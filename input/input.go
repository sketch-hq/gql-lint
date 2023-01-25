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

func QueryFiles(args []string) ([]string, error) {
	var files []string
	for _, source := range args {
		if strings.Contains(source, "*") {
			matches, err := glob(source)
			if err != nil {
				return files, err
			}

			files = append(files, matches...)
		} else {
			files = append(files, source)
		}
	}

	return unique(files), nil
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

					hits = append(hits, path)

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
