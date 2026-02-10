// Package markdown provides utilities for parsing and creating Markdown files
// with optional frontmatter.
package markdown

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var defaultDecoder = NewDecoder()
var defaultEncoder = NewEncoder()

// The DefaultDelimiter is used to determine the bounds of frontmatter if no
// other delimiter is set.
const DefaultDelimiter = "---"

// Decode uses the default Decoder to parse the contents of r into its
// respective frontmatter (if any) and body content.
func Decode(r io.Reader, fm any, body io.Writer) error {
	return defaultDecoder.Decode(r, fm, body)
}

// Encode uses the default Encoder to write the provided frontmatter (in YAML)
// and body contents to w.
func Encode(w io.Writer, fm any, body io.Reader) error {
	return defaultEncoder.Encode(w, fm, body)
}

// Decoder parses markdown content into frontmatter and body components.
type Decoder struct {
	// The delimiter used to determine the bounds of frontmatter.
	Delimiter []byte
}

// NewDecoder returns a Decoder configured with the default options.
func NewDecoder() *Decoder {
	p := &Decoder{
		Delimiter: []byte(DefaultDelimiter),
	}
	return p
}

// Decode reads the contents of r and extracts them into their respective
// frontmatter (if any) and body contents. If the delimiter is found at the
// start of the contents of r, it is assumed that the file contains YAML
// frontmatter and it MUST have a corresponding closing delimiter or the method
// will return an error.
//
// If fm is provided but no frontmatter is present in r, nothing will be written
// to fm. Likewise, if fm is not provided but frontmatter exists, the
// frontmatter will be ignored.
//
// Similarly, if body is not provided, it will be skipped.
func (dec *Decoder) Decode(r io.Reader, fm any, body io.Writer) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading contents: %w", err)
	}

	frontmatter, bodyContent, err := dec.splitContents(contents)
	if err != nil {
		return fmt.Errorf("splitting contents: %w", err)
	}

	if body != nil {
		if _, err := io.Copy(body, bytes.NewReader(bodyContent)); err != nil {
			return fmt.Errorf("copying body: %w", err)
		}
	}

	if frontmatter != nil && fm != nil {
		if err := yaml.Unmarshal(frontmatter, fm); err != nil {
			return fmt.Errorf("unmarshaling frontmatter: %w", err)
		}
	}

	return nil
}

func (dec *Decoder) hasFrontmatter(cs []byte) bool {
	return bytes.HasPrefix(cs, dec.Delimiter)
}

func (dec *Decoder) splitContents(cs []byte) ([]byte, []byte, error) {
	if !dec.hasFrontmatter(cs) {
		return nil, cs, nil
	}

	parts := bytes.SplitN(
		bytes.TrimPrefix(cs, dec.Delimiter),
		dec.Delimiter,
		2,
	)

	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid markdown: no closing delimiter found")
	}

	return parts[0], parts[1], nil
}

// Encoder writes markdown files with optional YAML frontmatter.
type Encoder struct {
	// The delimiter used to determine the bounds of frontmatter.
	Delimiter []byte
}

// NewEncoder returns a Encoder configured with the default options.
func NewEncoder() *Encoder {
	enc := &Encoder{
		Delimiter: []byte(DefaultDelimiter),
	}
	return enc
}

// Encode writes the provided frontmatter (in YAML) and body contents to w. If
// no frontmatter is provided, an empy frontmatter block will be added to the
// beginning of the output.
func (enc *Encoder) Encode(w io.Writer, fm any, body io.Reader) error {
	if err := enc.writeFrontmatter(w, fm); err != nil {
		return fmt.Errorf("writing frontmatter: %w", err)
	}

	if _, err := io.Copy(w, body); err != nil {
		return fmt.Errorf("writing body: %w", err)
	}

	return nil
}

func (enc *Encoder) writeFrontmatter(w io.Writer, fm any) error {
	if _, err := fmt.Fprintln(w, string(enc.Delimiter)); err != nil {
		return err
	}

	if fm != nil {
		// Close over the encoding portion to ensure the encoder closes before
		// writing the next delimiter
		if err := func() error {
			yamlEnc := yaml.NewEncoder(w)
			defer yamlEnc.Close()
			return yamlEnc.Encode(fm)
		}(); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w, string(enc.Delimiter)); err != nil {
		return err
	}

	return nil
}
