package handlers

import (
        "github.com/gin-gonic/gin"
        "net/http"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct{}

// NewEventHandler creates a new event handler
func NewEventHandler() *EventHandler {
        return &EventHandler{}
}

// GetLatestEvents retrieves latest events
func (eh *EventHandler) GetLatestEvents(c *gin.Context) {
        c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}

// MarkEventTriggered marks an event as triggered
func (eh *EventHandler) MarkEventTriggered(c *gin.Context) {
        c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
}