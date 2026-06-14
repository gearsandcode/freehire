// Package skilltag derives a job's technology tags deterministically from its
// free-text (HTML) description.
//
// Like internal/location, it is a curated dictionary, not an extractor: it
// resolves a known vocabulary of languages, frameworks, datastores, and infra by
// alias, and emits nothing for anything it cannot resolve (it never guesses).
// Tokens are lowercase slugs (go, postgresql, react, kubernetes), the same shape
// the enrichment contract's skills field uses, so the parser and the LLM payload
// speak one vocabulary and union cleanly at read time.
package skilltag

import (
	"regexp"
	"strings"
)

// htmlTagRE matches an HTML tag; descriptions are raw ATS HTML, so tags are
// replaced with a space before matching to keep markup tokens (div, href) out of
// the result and to avoid gluing words across a tag boundary.
var htmlTagRE = regexp.MustCompile(`<[^>]*>`)

// wordTokenRE splits normalized text into bare alphanumeric tokens for the word
// pass. Punctuated terms (c++, node.js) are handled separately by the phrase pass.
var wordTokenRE = regexp.MustCompile(`[a-z0-9]+`)

// normalize strips HTML tags, lowercases the text, and trims leading/trailing
// whitespace. Tags are replaced with a space (not empty) to preserve word
// boundaries so "<b>Go</b>Engineer" cannot fuse.
func normalize(text string) string {
	return strings.TrimSpace(strings.ToLower(htmlTagRE.ReplaceAllString(text, " ")))
}

// wordTokens returns the alphanumeric tokens of already-normalized text, in order.
func wordTokens(norm string) []string {
	return wordTokenRE.FindAllString(norm, -1)
}
