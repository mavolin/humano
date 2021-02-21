package replier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoSplit(t *testing.T) {
	expect := []string{"abc def jkl"}
	assert.Equal(t, expect, NoSplit(expect[0]))
}

func TestFieldsFuncSplitter(t *testing.T) {
	message := "abc\n\ndef\njkl"
	expect := []string{"abc", "def", "jkl"}

	splitter := FieldsFuncSplitter(func(r rune) bool {
		return r == '\n'
	})

	actual := splitter(message)
	assert.Equal(t, expect, actual)
}

func TestStaticDelay(t *testing.T) {
	expect := 123 * time.Millisecond
	assert.Equal(t, expect, StaticDelay(expect)("abc"))
}
