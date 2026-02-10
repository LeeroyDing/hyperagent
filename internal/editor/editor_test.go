package editor

import (
"os"
"testing"

"github.com/stretchr/testify/assert"
)

func TestFileEditor_ReadLines(t *testing.T) {
editor := NewFileEditor()
tmpfile, err := os.CreateTemp("", "test_read_lines")
assert.NoError(t, err)
defer os.Remove(tmpfile.Name())

content := "line1\nline2\nline3\nline4\nline5"
_, err = tmpfile.WriteString(content)
assert.NoError(t, err)
tmpfile.Close()

t.Run("ReadAll", func(t *testing.T) {
lines, err := editor.ReadLines(tmpfile.Name(), 1, 0)
assert.NoError(t, err)
assert.Equal(t, []string{"line1", "line2", "line3", "line4", "line5"}, lines)
})

t.Run("ReadRange", func(t *testing.T) {
lines, err := editor.ReadLines(tmpfile.Name(), 2, 4)
assert.NoError(t, err)
assert.Equal(t, []string{"line2", "line3", "line4"}, lines)
})

t.Run("FileNotFound", func(t *testing.T) {
_, err := editor.ReadLines("nonexistent", 1, 0)
assert.Error(t, err)
})
}

func TestFileEditor_Replace(t *testing.T) {
editor := NewFileEditor()
tmpfile, err := os.CreateTemp("", "test_replace")
assert.NoError(t, err)
defer os.Remove(tmpfile.Name())

content := "Hello World\nThis is a test."
_, err = tmpfile.WriteString(content)
assert.NoError(t, err)
tmpfile.Close()

t.Run("Success", func(t *testing.T) {
err := editor.Replace(tmpfile.Name(), "World", "Go")
assert.NoError(t, err)
newContent, _ := os.ReadFile(tmpfile.Name())
assert.Contains(t, string(newContent), "Hello Go")
})

t.Run("NotFound", func(t *testing.T) {
err := editor.Replace(tmpfile.Name(), "Missing", "New")
assert.Error(t, err)
assert.Contains(t, err.Error(), "not found")
})

t.Run("MultipleFound", func(t *testing.T) {
os.WriteFile(tmpfile.Name(), []byte("test test test"), 0644)
err := editor.Replace(tmpfile.Name(), "test", "check")
assert.Error(t, err)
assert.Contains(t, err.Error(), "multiple times")
})

t.Run("FileNotFound", func(t *testing.T) {
err := editor.Replace("nonexistent", "a", "b")
assert.Error(t, err)
})
}
