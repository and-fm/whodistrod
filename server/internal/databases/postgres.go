package databases

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/and-fm/whodistrod/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"
)

type postgres struct {
	dig.In `ignore-unexported:"true"`

	pgPool *pgxpool.Pool // not injected, it is created in NewPostgres

	Logger logging.Logger
	Config *config.Config
}

type Postgres interface {
	GetPool() *pgxpool.Pool
	Health() map[string]string
}

func NewPostgres(p postgres) Postgres {
	pgPool, err := p.Connect()
	if err != nil {
		p.Logger.LogFatal("Failed to connect to PostgreSQL: %w", nil, err)
	}

	p.pgPool = pgPool

	return &p
}

func (p *postgres) Connect() (*pgxpool.Pool, error) {
	pgPool, err := pgxpool.New(utils.Ctb(), os.Getenv("PG_CONN"))
	return pgPool, err
}

func (p *postgres) GetPool() *pgxpool.Pool {
	return p.pgPool
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *postgres) Health() map[string]string {
	timeout := 1 * time.Second

	// sometimes my connection to the db is slow, may not need this after we switch to a zero trust setup
	if s.Config.Env == "local" {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.pgPool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.pgPool.Stat()
	stats["acquired_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["idle_connections"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["constructing_connections"] = strconv.Itoa(int(dbStats.ConstructingConns()))
	stats["total_connections"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["wait_count"] = strconv.Itoa(int(dbStats.EmptyAcquireCount()))
	stats["max_conns"] = strconv.Itoa(int(dbStats.MaxConns()))
	stats["max_idle_closed"] = strconv.Itoa(int(dbStats.MaxIdleDestroyCount()))
	stats["max_lifetime_closed"] = strconv.Itoa(int(dbStats.MaxLifetimeDestroyCount()))

	// Evaluate stats to provide a health message
	if dbStats.AcquiredConns() > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.EmptyAcquireCount() > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleDestroyCount() > int64(dbStats.AcquiredConns())/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeDestroyCount() > int64(dbStats.AcquiredConns())/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}
