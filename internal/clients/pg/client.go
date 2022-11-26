package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4" // Load pgx sql driver.
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	notRunning uint32 = 0
	running    uint32 = 1

	defaultIdleConns int = 1
	defaultOpenConns int = 2

	defaultCheckInterval time.Duration = time.Second
)

type Client struct {
	db *sql.DB

	isAvailable      bool
	reconnectRunning uint32
	checkInterval    time.Duration
}

type ClientParams struct {
	User     string
	Password string
	Host     string
	Port     uint
	DBName   string

	MaxIdleConns *int
	MaxOpenConns *int

	CircuitBreakerCheckInterval *time.Duration
}

func (cp ClientParams) String() string {
	return fmt.Sprintf("Host:%s, Port:%d, User:%s, Has Password:%v, DBName:%s", cp.Host, cp.Port, cp.User, len(cp.Password) > 0, cp.DBName)
}

type CircuitBreaker interface {
	IsAvailable() bool
	Break()
	BreakOnNetworkError(err error) bool
}

func New(params ClientParams) (Client, error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v",
		params.Host, params.Port, params.User, params.Password, params.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return Client{}, errors.Wrapf(err, "could not create postgres sql driver: %v:%v ", params.Host, params.Port)
	}

	idleConns := defaultIdleConns
	if params.MaxIdleConns != nil {
		idleConns = *params.MaxIdleConns
	}

	openConns := defaultOpenConns
	if params.MaxOpenConns != nil {
		openConns = *params.MaxOpenConns
	}

	db.SetMaxIdleConns(idleConns)
	db.SetMaxOpenConns(openConns)

	checkInterval := defaultCheckInterval
	if params.CircuitBreakerCheckInterval != nil {
		checkInterval = *params.CircuitBreakerCheckInterval
	}

	ret := Client{
		db:            db,
		isAvailable:   true,
		checkInterval: checkInterval,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return Client{}, errors.Wrapf(err, "could not ping postgres database at: %v:%v ", params.Host, params.Port)
	}

	return ret, nil
}

func (c Client) Open(config gorm.Config) (*gorm.DB, error) {
	config.DisableAutomaticPing = true
	config.Logger = driverLogger{}

	return gorm.Open(postgresql.New(postgresql.Config{Conn: c.db}), &config)
}

func (c Client) Shutdown(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- c.db.Close()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Circuit Breaker part

func (c Client) GetCircuitBreaker() CircuitBreaker {
	return c
}

func (c Client) IsAvailable() bool {
	return c.isAvailable
}

func (c Client) BreakOnNetworkError(err error) bool {
	isNetworkError := IsNetworkError(err)
	if isNetworkError {
		c.Break()
	}

	return isNetworkError
}

func (c Client) Break() {
	c.isAvailable = false

	// Not very idiomatic.
	if atomic.CompareAndSwapUint32(&c.reconnectRunning, notRunning, running) {
		go c.tryReconnect()
	}
}

func (c Client) tryReconnect() {
	logrus.Debug("trying to reconnect...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := c.db.PingContext(ctx)
	if err != nil {
		logrus.WithError(err).Debug("failed to reconnect")

		time.Sleep(c.checkInterval)

		go c.tryReconnect()

		return
	}

	logrus.Info("postgres connection is up again")
	atomic.StoreUint32(&c.reconnectRunning, notRunning)
}

func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		// Network error https://github.com/jackc/pgerrcode/blob/master/errcode.go
		if len(pgErr.Code) > 1 && (pgErr.Code[0:2] == "08" || pgErr.Code[0:2] == "53" || pgErr.Code[0:2] == "57") {
			return true
		}
	}

	return false
}
