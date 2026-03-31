package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

func SeedDefaultAdmin(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// 1. Seed Admin User
	var count int
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE role = 'ADMIN'").Scan(&count)
	if count == 0 {
		hash, err := utils.HashPassword("admin123")
		if err == nil {
			_, err = pool.Exec(ctx,
				`INSERT INTO users (username, password_hash, full_name, role, is_active) 
				 VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
				"admin", hash, "Administrator", "ADMIN", true,
			)
			if err == nil {
				log.Println("✅ Default admin user created (username: admin, password: admin123)")
			}
		}
	} else {
		log.Println("ℹ️  Admin user already exists, skipping user seed")
	}

	// 2. Seed default categories
	categories := []struct {
		Name      string
		Desc      string
		ColorCode string
	}{
		{"Elektronik", "Peralatan elektronik (proyektor, laptop, dll)", "#2196F3"},
		{"Furniture", "Meja, kursi, lemari", "#795548"},
		{"Alat Lab", "Peralatan laboratorium", "#9C27B0"},
		{"Olahraga", "Peralatan olahraga", "#4CAF50"},
		{"ATK", "Alat tulis kantor", "#FF9800"},
		{"Lainnya", "Kategori lain-lain", "#607D8B"},
	}

	catCount := 0
	for _, c := range categories {
		_, err := pool.Exec(ctx,
			`INSERT INTO categories (name, description, color_code) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING`,
			c.Name, c.Desc, c.ColorCode,
		)
		if err == nil {
			catCount++
		}
	}
	if catCount > 0 {
		log.Printf("✅ %d categories seeded", catCount)
	}

	// 3. Seed default students
	students := []struct {
		NIS   string
		Name  string
		Class string
	}{
		{"12345", "Budi Santoso", "12 IPA 1"},
		{"12346", "Siti Aminah", "12 IPA 2"},
		{"12347", "Andi Wijaya", "11 IPS 1"},
		{"12348", "Dewi Lestari", "10 TG 1"},
	}

	studCount := 0
	for _, s := range students {
		_, err := pool.Exec(ctx,
			`INSERT INTO students (nis, full_name, class) VALUES ($1, $2, $3) ON CONFLICT (nis) DO NOTHING`,
			s.NIS, s.Name, s.Class,
		)
		if err == nil {
			studCount++
		}
	}
	// 4. Seed a few initial transactions for testing
	var itemId string
	var adminId string
	_ = pool.QueryRow(ctx, "SELECT id FROM items LIMIT 1").Scan(&itemId)
	_ = pool.QueryRow(ctx, "SELECT id FROM users WHERE username = 'admin'").Scan(&adminId)

	if itemId != "" && adminId != "" {
		trxCount := 0
		// Transaction for Budi (12345)
		_, err := pool.Exec(ctx, `
			INSERT INTO transactions 
			(item_id, borrower_type, student_nis, student_name, student_class, borrowed_by, borrowed_at, due_date, status, purpose)
			VALUES ($1, 'STUDENT', '12345', 'Budi Santoso', '12 IPA 1', $2, NOW() - INTERVAL '1 day', NOW() + INTERVAL '5 hours', 'ACTIVE', 'Praktikum Fisika')
			ON CONFLICT DO NOTHING`,
			itemId, adminId,
		)
		if err == nil { trxCount++ }

		// Past transaction for Siti (12346) - Returned
		_, err = pool.Exec(ctx, `
			INSERT INTO transactions 
			(item_id, borrower_type, student_nis, student_name, student_class, borrowed_by, borrowed_at, due_date, returned_at, status, purpose)
			VALUES ($1, 'STUDENT', '12346', 'Siti Aminah', '12 IPA 2', $2, NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days' + INTERVAL '2 hours', 'RETURNED', 'Tugas Kelompok')
			ON CONFLICT DO NOTHING`,
			itemId, adminId,
		)
		if err == nil { trxCount++ }

		if trxCount > 0 {
			log.Printf("✅ %d initial transactions seeded", trxCount)
		}
	}

	return nil
}
