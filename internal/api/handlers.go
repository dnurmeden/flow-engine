package api

import (
	"github.com/dnurmeden/flow-engine/internal/models"
	"github.com/dnurmeden/flow-engine/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	process *service.ProcessService
}

func NewHandler(process *service.ProcessService) *Handler {
	return &Handler{process: process}
}

func (h *Handler) StartProcess(c *gin.Context) {
	var req models.StartProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.process.StartProcess(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetInstance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp, err := h.process.GetInstance(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	// если нужно отдать «красиво» с форматированием:
	// c.IndentedJSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ClaimTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.ClaimTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user required"})
		return
	}
	if err := h.process.ClaimTask(c, taskID, req.User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) CompleteTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.CompleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user required"})
		return
	}
	if err := h.process.CompleteTask(c, taskID, req.User, req.Output); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
