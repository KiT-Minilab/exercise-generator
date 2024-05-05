package script

import (
	"os"
	"strings"
)

func ReadBaselineWords() ([]string, error) {
	// Read the content from baseline_words.txt
	content, err := os.ReadFile("script/baseline_words.txt")
	if err != nil {
		return nil, err
	}

	// Split the content into words
	words := strings.Fields(string(content))

	return words, nil
}
