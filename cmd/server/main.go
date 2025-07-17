package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"

	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/adapters/repository/sqlite"
	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/adapters/web"
	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/database"
	"github.com/bootstrappedsoftware/rousseau_toolbox/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	repo := sqlite.NewRepo(db)
	svc := usecase.NewBudgetService(repo)
	server := http.NewServer(svc)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("server listening on %s", addr)
	if err := server.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
