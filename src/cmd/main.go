package main

/*
	Изучите мой проект https://mailhub.su (https://github.com/go-park-mail-ru/2024_1_Refugio) и оцените мои навыки.
	Я мечтаю попасть в Авито!
*/

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	_ "github.com/jackc/pgx/stdlib"

	"avito_2024/src/internal/delivery/middleware"
	"avito_2024/src/internal/repository/postgresql"

	hand "avito_2024/src/internal/delivery/http"

	migrate "github.com/rubenv/sql-migrate"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "avito_2024/docs"
)

// @title API Avito
// @version 1.0
// @description API server for Avito

// @host localhost:8080
// @BasePath /
func main() {
	setLocalTime()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := initializeDatabase()
	defer db.Close()

	migrateDatabase(db)

	router := setupRouter(db)
	startServer(router)
}

func setLocalTime() {
	local, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("Error with setting the local time on the server.", err)
		return
	}

	time.Local = local
}

// initializeDatabase database initialization.
func initializeDatabase() *sql.DB {
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")
	sslmode := "disable"

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		username, database, password, host, port, sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Can't parse config: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is not available: %v", err)
	}

	db.SetMaxOpenConns(10)

	return db
}

func migrateDatabase(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "src/db/migration",
	}

	_, errMigration := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if errMigration != nil {
		log.Fatalf("Failed to apply migrations: %v", errMigration)
	}
}

func initializeTender(db *sql.DB) *hand.TenderHandler {
	tenderRepository := postgresql.NewTenderRepository(sqlx.NewDb(db, "pqx"))

	return hand.NewTenderHandler(tenderRepository)
}

func initializeProposal(db *sql.DB) *hand.ProposalHandler {
	proposalRepository := postgresql.NewProposalRepository(sqlx.NewDb(db, "pqx"))

	return hand.NewProposalHandler(proposalRepository)
}

func setupRouter(db *sql.DB) http.Handler {
	router := mux.NewRouter()

	tender := setupTenderRouter(db)
	router.PathPrefix("/api").Handler(tender)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return middleware.RequestLogger(router)
}

func setupTenderRouter(db *sql.DB) http.Handler {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	tenderHandler := initializeTender(db)
	proposalHandler := initializeProposal(db)

	router.HandleFunc("/ping", tenderHandler.Ping).Methods("GET", "OPTIONS")
	router.HandleFunc("/tenders/new", tenderHandler.CreateTender).Methods("POST", "OPTIONS")
	router.HandleFunc("/tenders/{tenderId}/edit", tenderHandler.EditTender).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/tenders/my", tenderHandler.GetMyTenders).Methods("GET", "OPTIONS")
	router.HandleFunc("/tenders", tenderHandler.GetTenders).Methods("GET", "OPTIONS")
	router.HandleFunc("/tenders/{tenderId}/rollback/{version}", tenderHandler.RollbackTender).Methods("PUT", "OPTIONS")
	router.HandleFunc("/tenders/{tenderId}/publish", tenderHandler.PublishTender).Methods("PUT", "OPTIONS")
	router.HandleFunc("/tenders/{tenderId}/close", tenderHandler.CloseTender).Methods("PUT", "OPTIONS")
	router.HandleFunc("/tenders/status", tenderHandler.GetTenderStatus).Methods("GET", "OPTIONS")

	router.HandleFunc("/bids/new", proposalHandler.CreateProposal).Methods("POST", "OPTIONS")
	router.HandleFunc("/bids/my", proposalHandler.GetMyProposals).Methods("GET", "OPTIONS")
	router.HandleFunc("/bids/{tenderId}/list", proposalHandler.GetProposalsByTender).Methods("GET", "OPTIONS")
	router.HandleFunc("/bids/{bidId}/edit", proposalHandler.EditProposal).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/bids/{bidId}/rollback/{version}", proposalHandler.RollbackProposal).Methods("PUT", "OPTIONS")
	router.HandleFunc("/bids/{bidId}/publish", proposalHandler.PublishProposal).Methods("PUT", "OPTIONS")
	router.HandleFunc("/bids/{bidId}/cancel", proposalHandler.CancelProposal).Methods("PUT", "OPTIONS")
	router.HandleFunc("/bids/status", proposalHandler.GetProposalStatus).Methods("GET", "OPTIONS")

	return router
}

func startServer(router http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5000", "http://localhost:5000", "http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Csrf-Token", "Content-Type", "AuthToken"},
		ExposedHeaders:   []string{"X-Csrf-Token", "AuthToken"},
	})

	corsHandler := c.Handler(router)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	fmt.Printf("The server is running on http://localhost:%d\n", port)
	fmt.Printf("Swagger is running on http://localhost:%d/swagger/index.html\n", port)

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), corsHandler)
	if err != nil {
		log.Fatalf("Error when starting the server: %v", err)
		return
	}
}
