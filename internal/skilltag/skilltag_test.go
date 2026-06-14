package skilltag

import (
	"reflect"
	"testing"
)

func TestNormalizeStripsHTMLAndLowercases(t *testing.T) {
	got := normalize("<div><p>Senior <b>Go</b> Engineer</p></div>")
	// Two spaces between words: each surrounding tag is replaced by one space.
	want := "senior  go  engineer"
	if got != want {
		t.Fatalf("normalize = %q, want %q", got, want)
	}
}

func TestWordTokens(t *testing.T) {
	got := wordTokens("go, node.js & c++17")
	want := []string{"go", "node", "js", "c", "17"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("wordTokens = %#v, want %#v", got, want)
	}
}
