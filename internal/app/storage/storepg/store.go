package storepg

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func Ping(ps string) error {

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
