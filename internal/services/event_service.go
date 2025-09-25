package services

import (
        "fmt"
        "twitch-rpg/internal/database"
        "twitch-rpg/internal/models"
        "twitch-rpg/internal/storage"
)

// EventService handles event-related operations
type EventService struct{}

// NewEventService creates a new event service
func NewEventService() *EventService {
        return &EventService{}
}

// GetLatestEvents retrieves the most recent events
func (es *EventService) GetLatestEvents(limit int) ([]models.Event, error) {
        if database.DB == nil {
                return storage.Memory.GetLatestEvents(limit)
        }

        query := `
                SELECT id, type, title, description, data, is_triggered, created_at, expires_at
                FROM events 
                ORDER BY created_at DESC 
                LIMIT ?`

        rows, err := database.DB.Query(query, limit)
        if err != nil {
                return nil, fmt.Errorf("failed to get events: %v", err)
        }
        defer rows.Close()

        var events []models.Event
        for rows.Next() {
                var event models.Event
                err := rows.Scan(
                        &event.ID, &event.Type, &event.Title, &event.Description,
                        &event.Data, &event.IsTriggered, &event.CreatedAt, &event.ExpiresAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan event: %v", err)
                }
                events = append(events, event)
        }

        return events, nil
}

// MarkEventTriggered marks an event as triggered
func (es *EventService) MarkEventTriggered(eventID int) error {
        if database.DB == nil {
                return storage.Memory.MarkEventTriggered(eventID)
        }

        query := `UPDATE events SET is_triggered = true WHERE id = ?`
        
        result, err := database.DB.Exec(query, eventID)
        if err != nil {
                return fmt.Errorf("failed to mark event as triggered: %v", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get rows affected: %v", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("event not found")
        }

        return nil
}

// CreateEvent creates a new event
func (es *EventService) CreateEvent(eventType, title, description string, data map[string]interface{}) (*models.Event, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }

        // Convert data to JSON string (simplified)
        dataJSON := "{}"
        if len(data) > 0 {
                // In a real implementation, you'd use json.Marshal here
                dataJSON = fmt.Sprintf("%v", data)
        }

        query := `
                INSERT INTO events (type, title, description, data, is_triggered)
                VALUES (?, ?, ?, ?, false)`

        result, err := database.DB.Exec(query, eventType, title, description, dataJSON)
        if err != nil {
                return nil, fmt.Errorf("failed to create event: %v", err)
        }

        id, err := result.LastInsertId()
        if err != nil {
                return nil, fmt.Errorf("failed to get event ID: %v", err)
        }

        // Return the created event
        return es.GetEventByID(int(id))
}

// GetEventByID retrieves an event by ID
func (es *EventService) GetEventByID(id int) (*models.Event, error) {
        if database.DB == nil {
                return nil, fmt.Errorf("database connection not available")
        }

        query := `
                SELECT id, type, title, description, data, is_triggered, created_at, expires_at
                FROM events 
                WHERE id = ?`

        event := &models.Event{}
        err := database.DB.QueryRow(query, id).Scan(
                &event.ID, &event.Type, &event.Title, &event.Description,
                &event.Data, &event.IsTriggered, &event.CreatedAt, &event.ExpiresAt,
        )

        if err != nil {
                return nil, fmt.Errorf("failed to get event: %v", err)
        }

        return event, nil
}