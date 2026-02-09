package web

import (
"net/http"
"net/http/httptest"
"testing"

"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/stretchr/testify/assert"
)

func TestWebRoutes(t *testing.T) {
tmpDir := t.TempDir()
histMgr, _ := history.NewHistoryManager(tmpDir)
s := NewServer(nil, histMgr, nil)

t.Run("ServeUI", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/ui/", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
assert.Contains(t, w.Body.String(), "Hyperagent | OS Companion")
})

t.Run("RedirectRoot", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusMovedPermanently, w.Code)
})

t.Run("APISessions", func(t *testing.T) {
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/api/sessions", nil)
s.router.ServeHTTP(w, req)
assert.Equal(t, http.StatusOK, w.Code)
})
}
