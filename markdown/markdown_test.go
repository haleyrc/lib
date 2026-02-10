package markdown_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/haleyrc/lib/assert"
	"github.com/haleyrc/lib/markdown"
)

func TestDecoder_Decode(t *testing.T) {
	dec := markdown.NewDecoder()

	t.Run("with frontmatter", func(t *testing.T) {
		content := readTestFile(t, "post.md")
		var fm postFrontmatter
		var buff bytes.Buffer

		err := dec.Decode(
			bytes.NewReader(content),
			&fm,
			&buff,
		)
		assert.OK(t, err).Fatal()

		assert.Equal(t, "title", "This Is The Title", fm.Title)
		assert.True(t, "draft", fm.Draft)

		got := buff.String()
		assert.StringContains(t, "body", got, "This is the content")
		assert.StringContains(t, "body", got, "Heading 2")
		assert.StringContains(t, "body", got, "Some more content")
	})

	t.Run("without frontmatter", func(t *testing.T) {
		content := readTestFile(t, "post-no-frontmatter.md")
		var buff bytes.Buffer

		err := dec.Decode(
			bytes.NewReader(content),
			nil,
			&buff,
		)
		assert.OK(t, err).Fatal()

		got := buff.String()
		assert.StringContains(t, "body", got, "This is the content")
		assert.StringContains(t, "body", got, "Heading 2")
		assert.StringContains(t, "body", got, "Some more content")
	})
}

func TestEncoder_Encode(t *testing.T) {
	enc := markdown.NewEncoder()
	fm := postFrontmatter{
		Title: "This Is The Title",
		Draft: true,
	}
	body := "\nThis is the content\n\n## Heading 2\n\nSome more content\n"
	want := readTestFile(t, "post.md")

	var buff bytes.Buffer
	err := enc.Encode(&buff, fm, strings.NewReader(body))
	assert.OK(t, err).Fatal()

	assert.DeepEqual(t, "post", want, buff.Bytes())
}

type postFrontmatter struct {
	Title string `yaml:"title"`
	Draft bool   `yaml:"draft"`
}

func readTestFile(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	bytes, err := os.ReadFile(path)
	assert.OK(t, err).Fatal()
	return bytes
}
