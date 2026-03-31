package utils

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LogAudit asynchronously creates an audit log entry in the database
func LogAudit(db *pgxpool.Pool, userId string, action string, entityType string, entityId string, description string, ipAddress string) {
	go func() {
		ctx := context.Background()
		query := `
			INSERT INTO audit_logs (user_id, action, entity_type, entity_id, description, ip_address)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		
		var nullableIP *string
		if ipAddress != "" {
			nullableIP = &ipAddress
		}

		_, err := db.Exec(ctx, query, userId, action, entityType, entityId, description, nullableIP)
		if err != nil {
			log.Printf("⚠️ Failed to write audit log: %v", err)
		}
	}()
}
