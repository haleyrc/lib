// Package conflate provides general-purpose string parsing, transformation, and
// template merging. It extracts named capture groups from a regexp, applies
// caller-supplied transforms, and merges the results into a text/template.
//
// The three functions — Parse, Transform, Merge — are independently useful and
// designed to be composed by the caller.
package conflate

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

// TransformFunc converts a single captured field value.
type TransformFunc func(value string) (string, error)

// Parse compiles pattern into a regular expression and then calls [ParseRegexp]
// with it and input.
func Parse(pattern, input string) (map[string]string, error) {
	exp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	results, err := ParseRegexp(exp, input)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return results, nil
}

// ParseRegexp extracts named capture groups from input using exp.
// Returns raw (untransformed) string values keyed by group name.
// An error is returned if input does not match exp.
func ParseRegexp(exp *regexp.Regexp, input string) (map[string]string, error) {
	matches := exp.FindStringSubmatch(input)
	if matches == nil {
		return nil, fmt.Errorf("parse regexp: input %q did not match expression", input)
	}

	result := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	return result, nil
}

// Transform applies transformers to matching keys in groups.
// Keys without a matching transformer pass through unchanged.
// Returns a new map; does not mutate the input.
func Transform(groups map[string]string, transformers map[string]TransformFunc) (map[string]string, error) {
	result := make(map[string]string, len(groups))
	for k, v := range groups {
		fn, ok := transformers[k]
		if !ok {
			result[k] = v
			continue
		}
		transformed, err := fn(v)
		if err != nil {
			return nil, fmt.Errorf("transform %q: %w", k, err)
		}
		result[k] = transformed
	}
	return result, nil
}

// Merge compiles format into a template and then calls [MergeTemplate] with it
// and data.
func Merge(format string, data map[string]string) (string, error) {
	tmpl, err := template.New("").Parse(format)
	if err != nil {
		return "", fmt.Errorf("merge: %w", err)
	}
	result, err := MergeTemplate(tmpl, data)
	if err != nil {
		return "", fmt.Errorf("merge: %w", err)
	}
	return result, nil
}

// MergeTemplate executes tmpl against data and returns the result string.
func MergeTemplate(tmpl *template.Template, data map[string]string) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("merge template: %w", err)
	}
	return buf.String(), nil
}
