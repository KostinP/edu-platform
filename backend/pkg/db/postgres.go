package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var pool *pgxpool.Pool

// Connect инициализирует подключение к БД с пулом соединений
func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Тестируем подключение
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	logrus.Infof("✅ Connected to database (pool size: %d)", pool.Stat().MaxConns())
	return nil
}

// DB возвращает экземпляр пула соединений
func DB() *pgxpool.Pool {
	if pool == nil {
		panic("database not initialized")
	}
	return pool
}

// Close закрывает соединение
func Close() {
	if pool != nil {
		pool.Close()
		logrus.Info("🛑 Database connection closed")
	}
}
