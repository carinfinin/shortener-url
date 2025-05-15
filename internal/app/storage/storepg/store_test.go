package storepg

import (
	"context"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(t *testing.T) (string, func()) {
	ctx := context.Background()

	// Запускаем контейнер
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Получаем строку подключения
	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Функция для остановки контейнера
	cleanup := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	return connStr, cleanup
}

func TestWithContainer(t *testing.T) {
	connStr, cleanup := setupPostgresContainer(t)
	defer cleanup()

	// Используем connStr в вашем хранилище
	cfg := &config.Config{DBPath: connStr}
	store, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(store)

	// Далее ваши тесты...
}
