package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vannyezha/sis-inv/internal/utils"
)

func SeedDefaultAdmin(pool *pgxpool.Pool) error {
	var count int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE role = 'ADMIN'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("ℹ️  Admin user already exists, skipping seed")
		return nil
	}

	hash, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(),
		`INSERT INTO users (username, password_hash, full_name, role, is_active) 
		 VALUES ($1, $2, $3, $4, $5)`,
		"admin", hash, "Administrator", "ADMIN", true,
	)
	if err != nil {
		return err
	}

	log.Println("✅ Default admin user created (username: admin, password: admin123)")

	// Seed default categories
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

	for _, c := range categories {
		_, err := pool.Exec(context.Background(),
			`INSERT INTO categories (name, description, color_code) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
			c.Name, c.Desc, c.ColorCode,
		)
		if err != nil {
			log.Printf("⚠️  Failed to seed category %s: %v", c.Name, err)
		}
	}

	log.Println("✅ Default categories seeded")

	// Seed default students
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

	for _, s := range students {
		_, err := pool.Exec(context.Background(),
			`INSERT INTO students (nis, full_name, class) VALUES ($1, $2, $3) ON CONFLICT (nis) DO NOTHING`,
			s.NIS, s.Name, s.Class,
		)
		if err != nil {
			log.Printf("⚠️  Failed to seed student %s: %v", s.Name, err)
		}
	}
	log.Println("✅ Default students seeded")

	return nil
}
