package conflate_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/haleyrc/lib/assert"
	"github.com/haleyrc/lib/conflate"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    map[string]string
		wantErr string
	}{
		{
			name:    "extracts named capture groups",
			pattern: `^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})$`,
			input:   "2026-02-26",
			want:    map[string]string{"year": "2026", "month": "02", "day": "26"},
		},
		{
			name:    "returns error for no match",
			pattern: `^(?P<year>\d{4})$`,
			input:   "not-a-year",
			wantErr: "did not match expression",
		},
		{
			name:    "returns error for invalid pattern",
			pattern: `(unclosed`,
			input:   "input",
			wantErr: "parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conflate.Parse(tt.pattern, tt.input)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.OK(t, err).Fatal()
			assert.DeepEqual(t, "groups", tt.want, got)
		})
	}
}

func TestParseRegexp(t *testing.T) {
	tests := []struct {
		name    string
		exp     *regexp.Regexp
		input   string
		want    map[string]string
		wantErr string
	}{
		{
			name:  "extracts named capture groups",
			exp:   regexp.MustCompile(`^(?P<first>\w+)\s+(?P<last>\w+)$`),
			input: "Jane Doe",
			want:  map[string]string{"first": "Jane", "last": "Doe"},
		},
		{
			name:    "returns error for no match",
			exp:     regexp.MustCompile(`^(?P<digit>\d+)$`),
			input:   "abc",
			wantErr: "did not match expression",
		},
		{
			name:  "ignores unnamed capture groups",
			exp:   regexp.MustCompile(`^(\w+)\s+(?P<last>\w+)$`),
			input: "Jane Doe",
			want:  map[string]string{"last": "Doe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conflate.ParseRegexp(tt.exp, tt.input)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.OK(t, err).Fatal()
			assert.DeepEqual(t, "groups", tt.want, got)
		})
	}
}

func TestTransform(t *testing.T) {
	tests := []struct {
		name         string
		groups       map[string]string
		transformers map[string]conflate.TransformFunc
		want         map[string]string
		wantErr      string
	}{
		{
			name:   "applies matching transformers",
			groups: map[string]string{"name": "alice", "age": "30"},
			transformers: map[string]conflate.TransformFunc{
				"name": func(v string) (string, error) {
					return strings.ToUpper(v), nil
				},
			},
			want: map[string]string{"name": "ALICE", "age": "30"},
		},
		{
			name:   "passes through keys without transformers",
			groups: map[string]string{"a": "1", "b": "2"},
			want:   map[string]string{"a": "1", "b": "2"},
		},
		{
			name:   "returns error from transformer",
			groups: map[string]string{"key": "val"},
			transformers: map[string]conflate.TransformFunc{
				"key": func(v string) (string, error) {
					return "", fmt.Errorf("bad value")
				},
			},
			wantErr: `transform "key"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conflate.Transform(tt.groups, tt.transformers)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.OK(t, err).Fatal()
			assert.DeepEqual(t, "groups", tt.want, got)
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		data    map[string]string
		want    string
		wantErr string
	}{
		{
			name:   "merges data into template",
			format: "{{.greeting}}, {{.name}}!",
			data:   map[string]string{"greeting": "Hello", "name": "World"},
			want:   "Hello, World!",
		},
		{
			name:    "returns error for invalid template",
			format:  "{{.unclosed",
			wantErr: "merge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conflate.Merge(tt.format, tt.data)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.OK(t, err).Fatal()
			assert.Equal(t, "result", tt.want, got)
		})
	}
}

func TestMergeTemplate(t *testing.T) {
	tests := []struct {
		name string
		tmpl *template.Template
		data map[string]string
		want string
	}{
		{
			name: "executes template with data",
			tmpl: template.Must(template.New("").Parse("{{.item}} costs {{.price}}")),
			data: map[string]string{"item": "Widget", "price": "$9.99"},
			want: "Widget costs $9.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conflate.MergeTemplate(tt.tmpl, tt.data)
			assert.OK(t, err).Fatal()
			assert.Equal(t, "result", tt.want, got)
		})
	}
}
