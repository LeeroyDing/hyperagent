package web

import (
"context"
"embed"
"fmt"
"io/fs"
"net/http"
"os"
"strings"

"github.com/LeeroyDing/hyperagent/internal/agent"
"github.com/LeeroyDing/hyperagent/internal/config"
"github.com/LeeroyDing/hyperagent/internal/gemini"
"github.com/LeeroyDing/hyperagent/internal/history"
"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"
"github.com/google/uuid"
)

//go:embed all:static
var staticFiles embed.FS

type Server struct {
Agent      *agent.Agent
History    *history.HistoryManager
Router     *gin.Engine
ConfigPath string
}

func NewServer(a *agent.Agent, h *history.HistoryManager) *Server {
r := gin.Default()
r.Use(cors.Default())

s := &Server{
Agent:      a,
History:    h,
Router:     r,
ConfigPath: config.GetDefaultConfigPath(),
}

s.setupRoutes()
return s
}

func (s *Server) setupRoutes() {
api := s.Router.Group("/api")
{
api.GET("/sessions", s.listSessions)
api.POST("/sessions", s.createSession)
api.GET("/sessions/:id/messages", s.getMessages)
api.POST("/sessions/:id/messages", s.sendMessage)
api.GET("/config", s.getConfig)
api.POST("/config", s.updateConfig)
}

sub, _ := fs.Sub(staticFiles, "static")
s.Router.StaticFS("/ui", http.FS(sub))
s.Router.GET("/", func(c *gin.Context) {
c.Redirect(http.StatusMovedPermanently, "/ui/")
})
}

func (s *Server) listSessions(c *gin.Context) {
sessions, err := s.History.ListSessions()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, sessions)
}

func (s *Server) createSession(c *gin.Context) {
id := uuid.New().String()
c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s *Server) getMessages(c *gin.Context) {
id := c.Param("/id")
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

if err := s.History.AddMessage(id, "user", req.Content); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

response, err := s.Agent.Run(context.Background(), id, req.Content)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

if err := s.History.AddMessage(id, "assistant", response); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

if s.History.GetSessionName(id) == "New Conversation" {
go s.autoNameSession(id, req.Content)
}

messages, _ := s.History.LoadHistory(id)
if len(messages)%10 == 0 {
go s.Agent.Distill(context.Background(), id)
}

c.JSON(http.StatusOK, gin.H{"role": "assistant", "content": response})
}

func (s *Server) autoNameSession(id, firstMessage string) {
prompt := fmt.Sprintf("Generate a short, concise title (max 5 words) for a conversation starting with: '%s'. Return ONLY the title.", firstMessage)
name, _, err := s.Agent.Gemini.GenerateContent(context.Background(), []gemini.Message{
{Role: "user", Content: prompt},
}, nil)
if err == nil {
s.History.SetSessionName(id, strings.TrimSpace(name))
}
}

func (s *Server) getConfig(c *gin.Context) {
data, err := os.ReadFile(s.ConfigPath)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.Data(http.StatusOK, "text/yaml", data)
}

func (s *Server) updateConfig(c *gin.Context) {
data, err := c.GetRawData()
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

if err := os.WriteFile(s.ConfigPath, data, 0644); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (s *Server) Run(addr string) error {
return s.Router.Run(addr)
}
