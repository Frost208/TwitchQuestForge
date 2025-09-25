package handlers

import (
        "net/http"
        "strconv"
        "twitch-rpg/internal/services"

        "github.com/gin-gonic/gin"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
        eventService *services.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler() *EventHandler {
        return &EventHandler{
                eventService: services.NewEventService(),
        }
}

// GetLatestEvents retrieves latest events
func (eh *EventHandler) GetLatestEvents(c *gin.Context) {
        limit := 10 // default
        if limitStr := c.Query("limit"); limitStr != "" {
                if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
                        limit = l
                }
        }

        events, err := eh.eventService.GetLatestEvents(limit)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// MarkEventTriggered marks an event as triggered
func (eh *EventHandler) MarkEventTriggered(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
                return
        }

        err = eh.eventService.MarkEventTriggered(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Event marked as triggered"})
}