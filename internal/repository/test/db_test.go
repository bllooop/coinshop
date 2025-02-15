package repository

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
	/*
	   	req := testcontainers.ContainerRequest{
	   		Image:        "postgres:15",
	   		ExposedPorts: []string{"5432/tcp"},
	   		Env: map[string]string{
	   			"POSTGRES_USER":     "testuser",
	   			"POSTGRES_PASSWORD": "testpassword",
	   			"POSTGRES_DB":       "testdb",
	   		},
	   		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port int) string {
	   			return fmt.Sprintf("postgres://testuser:testpassword@%s:%d/testdb?sslmode=disable", host, port)
	   		}),
	   	}

	   	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	   		ContainerRequest: req,
	   		Started:          true,
	   	})

	   	if err != nil {
	   		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	   	}

	   host, err := container.Host(ctx)

	   	if err != nil {
	   		log.Fatalf("Failed to get container host: %v", err)
	   	}

	   port, err := container.MappedPort(ctx, "5432")

	   	if err != nil {
	   		log.Fatalf("Failed to get container port: %v", err)
	   	}

	   dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())

	   testDB, err = sqlx.Open("postgres", dsn)

	   	if err != nil {
	   		log.Fatalf("Failed to connect to test database: %v", err)
	   	}

	   time.Sleep(2 * time.Second)

	   code := m.Run()

	   _ = testDB.Close()
	   _ = container.Terminate(ctx)

	   log.Printf("Integration tests completed with exit code %d", code)

	*/

}
