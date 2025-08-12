package db

import (
    "context"
    "fmt"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func Connect() error {
    godotenv.Load()

    dsn := fmt.Sprintf(
        "postgresql://%s:%s@%s:%s/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        return err
    }

    DB = pool
    return nil
}

func Close() {
    DB.Close()
}
