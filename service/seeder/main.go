package main

import (
	"context"
	"log"

	"github.com/MamangRust/microservice-point-of-sale-pkg/database"
	"github.com/MamangRust/microservice-point-of-sale-pkg/database/seeder"
	"github.com/MamangRust/microservice-point-of-sale-pkg/dotenv"
	"github.com/MamangRust/microservice-point-of-sale-pkg/hash"
	"github.com/MamangRust/microservice-point-of-sale-pkg/logger"
)

func main() {
	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("seeder")
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	gormDB, err := database.NewGormClient(l)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	ctx := context.Background()

	s := seeder.NewSeeder(seeder.Deps{
		DB:     gormDB,
		Ctx:    ctx,
		Logger: l,
		Hash:   hash.NewHashingPassword(),
	})

	if err := s.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	l.Info("Seeding completed successfully.")
}
