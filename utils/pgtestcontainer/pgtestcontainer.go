package pgtestcontainer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tupyy/gophoto/internal/repo/postgres"
)

const (
	postgresUser       = "postgres"
	postgresPassword   = "postgres"
	postgresDatabase   = "postgres"
	postgresImage      = "postgres"
	postgresDefaultTag = "13"
	postgresPort       = "5432/tcp"
)

// PostgreSQLContainerRequest completes ContainerRequest.
// with PostgreSQL specific parameters.
type PostgreSQLContainerRequest struct {
	BindMounts map[string]string
	Env        map[string]string
	ShowLog    bool
	User       string
	Password   string
	Database   string
	Image      string
}

// PostgreSQLContainer should always be created via NewPostgreSQLContainer.
type PostgreSQLContainer struct {
	Container testcontainers.Container
	req       PostgreSQLContainerRequest
}

type PostgresLogConsumer int

func (l *PostgresLogConsumer) Accept(log testcontainers.Log) {
	fmt.Print(string(log.Content))
}

//nolint:funlen
// NewPostgreSQLContainer creates and (optionally) starts a Postgres database.
// If autostarted, the function returns only after a successful execution of a query
// (confirming that the database is ready).
func NewPostgreSQLContainer(ctx context.Context, req PostgreSQLContainerRequest) (*PostgreSQLContainer, error) {

	if req.Env == nil {
		req.Env = map[string]string{}
	}

	// Set the default values if none were provided in the request
	if req.Image == "" {
		req.Image = fmt.Sprintf("%s:%s", postgresImage, postgresDefaultTag)
	}

	if req.User == "" {
		req.User = postgresUser
	}

	if req.Password == "" {
		req.Password = postgresPassword
	}

	if req.Database == "" {
		req.Database = postgresDatabase
	}

	connectorVars := map[string]string{
		"host":     "localhost",
		"user":     req.User,
		"password": req.Password,
		"database": req.Database,
	}

	containerReq := testcontainers.ContainerRequest{
		Image:        req.Image,
		ExposedPorts: []string{postgresPort},
		BindMounts:   req.BindMounts,
		Env: map[string]string{
			"POSTGRES_USER":     req.User,
			"POSTGRES_PASSWORD": req.Password,
			"POSTGRES_DB":       req.Database,
		},
		WaitingFor: wait.ForSQL(nat.Port(postgresPort), "postgres", postgresURL(connectorVars)),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create container")
	}

	postgresC := &PostgreSQLContainer{
		Container: c,
		req:       req,
	}

	if req.ShowLog {
		// set log consumer
		var pl PostgresLogConsumer = 0

		err := c.StartLogProducer(ctx)
		if err != nil {
			logrus.WithError(err).Warn("cannot start log producer")
		} else {
			c.FollowOutput(&pl)
		}
	}

	return postgresC, nil
}

// GetInitialClient returs a postgres.Client of the initial DB.
func (c *PostgreSQLContainer) GetInitialClient(ctx context.Context) (postgres.Client, error) {
	host, err := c.Container.Host(ctx)
	if err != nil {
		return postgres.Client{}, err
	}

	mappedPort, err := c.Container.MappedPort(ctx, postgresPort)
	if err != nil {
		return postgres.Client{}, err
	}

	client, _ := postgres.NewClient(postgres.ClientParams{
		Host:     host,
		Port:     uint(mappedPort.Int()),
		User:     c.req.User,
		Password: c.req.Password,
		DBName:   c.req.Database,
	})

	return client, nil
}

// GetClient returns a postgres.Client which could be used to connect to another database than the one created by default.
func (c *PostgreSQLContainer) GetClient(ctx context.Context, username, password, db string) (postgres.Client, error) {
	host, err := c.Container.Host(ctx)
	if err != nil {
		return postgres.Client{}, err
	}

	mappedPort, err := c.Container.MappedPort(ctx, postgresPort)
	if err != nil {
		return postgres.Client{}, err
	}

	client, _ := postgres.NewClient(postgres.ClientParams{
		Host:     host,
		Port:     uint(mappedPort.Int()),
		User:     username,
		Password: password,
		DBName:   db,
	})

	return client, nil
}

// GetDriver returns a sql.DB connecting to the previously started Postgres DB.
// All the parameters are taken from the previous PostgreSQLContainerRequest
// The driver can be used to make raw SQL requests.
func (c *PostgreSQLContainer) GetDriver(ctx context.Context) (*sql.DB, error) {
	host, err := c.Container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := c.Container.MappedPort(ctx, postgresPort)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		mappedPort.Int(),
		c.req.User,
		c.req.Password,
		c.req.Database,
	))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func postgresURL(variables map[string]string) func(nat.Port) string {
	return func(p nat.Port) string {
		connString := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			variables["host"],
			p.Int(),
			variables["user"],
			variables["password"],
			variables["database"],
		)

		return connString
	}
}
