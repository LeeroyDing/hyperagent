package editor

import (
"bufio"
"fmt"
"os"
"strings"
)

// FileEditor provides methods for safe file manipulation.
type FileEditor struct{}

// NewFileEditor creates a new FileEditor.
func NewFileEditor() *FileEditor {
return &FileEditor{}
}

// ReadLines reads specific lines from a file (1-indexed).
func (e *FileEditor) ReadLines(path string, start, end int) ([]string, error) {
file, err := os.Open(path)
if err != nil {
return nil, err
}
defer file.Close()

var lines []string
scanner := bufio.NewScanner(file)
lineNum := 0
for scanner.Scan() {
lineNum++
if lineNum >= start && (end <= 0 || lineNum <= end) {
lines = append(lines, scanner.Text())
}
if end > 0 && lineNum >= end {
break
}
}

if err := scanner.Err(); err != nil {
return nil, err
}

return lines, nil
}

// Replace replaces oldText with newText in the file.
// It returns an error if oldText is not found or found multiple times (to be safe).
func (e *FileEditor) Replace(path string, oldText, newText string) error {
content, err := os.ReadFile(path)
if err != nil {
return err
}

strContent := string(content)
count := strings.Count(strContent, oldText)
if count == 0 {
return fmt.Errorf("old text not found in file")
}
if count > 1 {
return fmt.Errorf("old text found multiple times (%d), please be more specific", count)
}

newContent := strings.Replace(strContent, oldText, newText, 1)
return os.WriteFile(path, []byte(newContent), 0644)
}
