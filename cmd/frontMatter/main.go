package main

import (
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// ExtractFrontMatter isolates the YAML block and returns the YAML, the remaining content, and a boolean status
func ExtractFrontMatter(input string) (yaml string, content string, found bool) {
	// Standardize line endings if needed, then check if it starts with the markdown fence
	if !strings.HasPrefix(input, "---\n") && !strings.HasPrefix(input, "---\r\n") {
		return "", input, false
	}

	// Split by the fence. SplitN with 3 limits the slice size:
	// parts[0] will be empty (before the first ---)
	// parts[1] will be the front matter block
	// parts[2] will be the remainder of the markdown file
	parts := strings.SplitN(input, "---", 3)
	if len(parts) < 3 {
		return "", input, false
	}

	// Trim surrounding whitespace from the block
	yamlBlock := strings.TrimSpace(parts[1])
	remainingContent := strings.TrimSpace(parts[2])

	return yamlBlock, remainingContent, true
}

func GetFileContents() string {
	content, err := os.ReadFile("sample.md")
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Convert the byte slice to a string to print it
	return string(content)
}

func main() {
	contents := GetFileContents()
	y, c, f := ExtractFrontMatter(contents)
	logrus.Info("---------------------------------------")
	logrus.Infof(" YAML %v, %v", f, y)
	logrus.Info("---------------------------------------")
	logrus.Info("---------------------------------------")
	logrus.Infof(" Contents %v", c)
	logrus.Info("---------------------------------------")
}
