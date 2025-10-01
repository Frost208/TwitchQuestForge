# Twitch RPG System - Complete Source Code Snapshot

**Generated:** October 01, 2025  
**Version:** 1.0.0  
**Language:** Go 1.24  
**Purpose:** Complete Twitch chat-based RPG system with channel points integration

---

## Table of Contents

1. [Project Structure](#project-structure)
2. [Entry Point](#entry-point---cmdservermain.go)
3. [Database Layer](#database-layer)
4. [HTTP Routing](#http-routing)
5. [HTTP Handlers](#http-handlers)
6. [Data Models](#data-models)
7. [Business Logic Services](#business-logic-services)
8. [Storage Layer](#storage-layer)
9. [Database Schema](#database-schema)
10. [Dependencies](#dependencies)
11. [Deployment Guide](#deployment-guide)

---

## Project Structure

```
twitch-rpg/
├── cmd/
│   └── server/
│       └── main.go                      # Server entry point & initialization
├── internal/
│   ├── database/
│   │   └── connection.go                # MySQL database connection management
│   ├── handlers/
│   │   ├── routes.go                    # API route registration
│   │   ├── character_handler.go         # Character HTTP endpoints
│   │   ├── item_handler.go              # Item HTTP endpoints
│   │   ├── combat_handler.go            # Combat HTTP endpoints
│   │   ├── merchant_handler.go          # Merchant HTTP endpoints
│   │   └── event_handler.go             # Event HTTP endpoints
│   ├── models/
│   │   ├── character.go                 # Character data structures
│   │   ├── item.go                      # Item data structures
│   │   ├── combat.go                    # Combat data structures
│   │   └── events.go                    # Event data structures
│   ├── services/
│   │   ├── character_service.go         # Character business logic
│   │   ├── item_service.go              # Item business logic
│   │   ├── merchant_service.go          # Merchant business logic
│   │   ├── combat_service.go            # Combat business logic
│   │   ├── combat_memory.go             # Combat memory fallback
│   │   └── event_service.go             # Event business logic
│   └── storage/
│       └── memory.go                    # In-memory storage fallback system
├── scripts/
│   └── schema.sql                       # Complete database schema
├── go.mod                               # Go module dependencies
├── go.sum                               # Dependency checksums
└── .env                                 # Environment configuration (not in repo)
```

---

