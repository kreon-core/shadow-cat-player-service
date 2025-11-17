package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kreon-core/shadow-cat-common/logc"
	"github.com/kreon-core/shadow-cat-common/postgres"

	"sc-player-service/infrastructure/config"
	"sc-player-service/repository/playersqlc"
)

type Manager struct {
	HTTPServer *HTTPServer
	Container  *Container

	PlayerDBPool    *pgxpool.Pool
	PlayerDBQueries *playersqlc.Queries

	// RedisClient *redis.Client

	// KafkaConsumer *kafka.Consumer
	// KafkaProducer *kafka.Producer

	// Workers ...
}

func New(ctx context.Context, cfg *config.Config) (*Manager, error) {
	playerDBPool, err := postgres.NewConnection(ctx, &cfg.DB.Player)
	if err != nil {
		return nil, fmt.Errorf("init_player_db -> %w", err)
	}
	logc.Info().Str("db", "player").Msg("DB connection pool established")

	playerDBQueries := playersqlc.New(playerDBPool)

	// Redis
	// KafkaProducer
	// KafkaConsumer

	// Workers ...

	container := NewContainer(cfg, playerDBQueries)
	httpServer := NewHTTPServer(&cfg.HTTP, container)

	return &Manager{
		HTTPServer:      httpServer,
		Container:       container,
		PlayerDBPool:    playerDBPool,
		PlayerDBQueries: playerDBQueries,
	}, nil
}

func (a *Manager) Start() {
	var wg sync.WaitGroup
	wg.Go(func() { a.HTTPServer.Run() })
	// Workers
	wg.Wait()
}

func (a *Manager) Stop(ctx context.Context) error {
	err := a.HTTPServer.Stop(ctx)
	if err != nil {
		return fmt.Errorf("stop_http_server -> %w", err)
	}

	// Workers
	// Kafka Consumer
	// Kafka Producer
	// Redis

	if a.PlayerDBPool != nil {
		a.PlayerDBPool.Close()
	}

	return nil
}
