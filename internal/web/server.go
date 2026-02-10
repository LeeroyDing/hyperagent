package web

import (
"context"
"embed"
"io/fs"
"net/http"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/LeeroyDing/hyperagent/internal/memory"
"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticAssets embed.FS

type Server struct {
Agent   *agent.Agent
History history.History
Memory  memory.Memory
router  *gin.Engine
}

func NewServer(a *agent.Agent, h history.History, m memory.Memory) *Server {
gin.SetMode(gin.ReleaseMode)
r := gin.Default()

s := &Server{
Agent:   a,
History: h,
Memory:  m,
router:  r,
}

s.setupRoutes()
return s
}

func (s *Server) setupRoutes() {
api := s.router.Group("/api")
{
api.GET("/sessions", s.getSessions)
api.POST("/sessions", s.createSession)
api.GET("/sessions/:id/messages", s.getMessages)
api.POST("/sessions/:id/messages", s.sendMessage)
api.GET("/memory", s.searchMemory)
api.DELETE("/memory/:id", s.deleteMemory)
}

// Serve embedded static files
sub, _ := fs.Sub(staticAssets, "static")
s.router.StaticFS("/ui", http.FS(sub))

// Redirect root to /ui/
s.router.GET("/", func(c *gin.Context) {
c.Redirect(http.StatusMovedPermanently, "/ui/")
})
}

func (s *Server) Run(addr string) error {
return s.router.Run(addr)
}

func (s *Server) getSessions(c *gin.Context) {
sessions, err := s.History.ListSessions()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, sessions)
}

func (s *Server) createSession(c *gin.Context) {
var req struct {
Name string `json:"name"`
}
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

id, err := s.History.CreateSession(req.Name)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s *Server) getMessages(c *gin.Context) {
id := c.Param("id")
messages, err := s.History.LoadHistory(id)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, messages)
}

func (s *Server) sendMessage(c *gin.Context) {
id := c.Param("id")
var req struct {
Content string `json:"content"`
}
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

response, err := s.Agent.Run(context.Background(), id, req.Content)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"response": response})
}

func (s *Server) searchMemory(c *gin.Context) {
query := c.Query("q")
results, err := s.Memory.Search(context.Background(), query, 10)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, results)
}

func (s *Server) deleteMemory(c *gin.Context) {
id := c.Param("id")
if err := s.Memory.Forget(context.Background(), id); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
