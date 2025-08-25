package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type PostgreSQL struct {
}

func (pg *PostgreSQL) Connect() (*pgxpool.Pool, error) {

	pool, err := pgxpool.Connect(context.Background(), pg.ConnectionURI())
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (pg *PostgreSQL) Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

func (pg *PostgreSQL) ConnectionURI() string {

	_ = godotenv.Load()

	// os.Setenv("PG_HOST", "128.199.97.209")
	// os.Setenv("PG_PORT", "5432")
	// os.Setenv("PG_USER", "postgres")
	// os.Setenv("PG_PASS", "postgres")
	// os.Setenv("PG_DB", "postgres")

	// host := os.Getenv("PG_HOST")
	// port := os.Getenv("PG_PORT")
	// user := os.Getenv("PG_USER")
	// pass := os.Getenv("PG_PASS")
	// db := os.Getenv("PG_DB")

	host := "128.199.97.209"
	port := "5432"
	user := "postgres"
	pass := "postgres"
	db := "postgres"
	// if port == "" {
	// 	port = "5432"
	// }

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}
