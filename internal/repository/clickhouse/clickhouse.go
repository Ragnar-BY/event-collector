package clickhouse

import (
	"context"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// ClickhouseSettings are settings for clickhouse
type ClickhouseSettings struct {
	Addr     string
	Database string
	Username string
	Password string
}

// ClickhouseClient is client fro clickhouse
type ClickhouseClient struct {
	conn driver.Conn
}

// NewClickhouseClient creates new instance of Clickhouse client
func NewClickhouseClient(settings ClickhouseSettings) (*ClickhouseClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{settings.Addr},
		Auth: clickhouse.Auth{
			Database: settings.Database,
			Username: settings.Username,
			Password: settings.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, err
	}

	client := &ClickhouseClient{
		conn: conn,
	}
	return client, client.conn.Ping(context.Background())
}

// SaveEvents save events in database
func (c *ClickhouseClient) SaveEvents(ctx context.Context, events []domain.Event) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO events")
	if err != nil {
		return err
	}

	for _, e := range events {
		err = batch.AppendStruct(&e)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}

// Close closes connection
func (c *ClickhouseClient) Close() error {
	return c.conn.Close()
}

// Ping pings database
func (c *ClickhouseClient) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}
