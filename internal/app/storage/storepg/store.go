package storepg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"sync"
	"time"
)

type Store struct {
	mu  sync.Mutex
	db  *sql.DB
	url string
}

func New(cfg *config.Config) (*Store, error) {

	db, err := sql.Open("pgx", cfg.DBPath)
	if err != nil {
		return nil, err
	}
	return &Store{
		mu:  sync.Mutex{},
		db:  db,
		url: cfg.URL,
	}, nil
}

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

func (s *Store) AddURL(url string) (string, error) {
	logger.Log.Info("start function AddURL")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.mu.Lock()

	xmlID := storage.GenerateXMLID(storage.LengthXMLID)

	//INSERT INTO urls (url, xmlid) SELECT 'dfgddfgd', '123' WHERE NOT EXISTS(SELECT 1 FROM urls WHERE xmlid = '123');
	r, err := s.db.ExecContext(ctx, "INSERT INTO urls (url, xmlid) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM urls WHERE xmlid = $3)", url, xmlID, xmlID)
	if err != nil {
		logger.Log.Error(" AddURL error :", err)
		return "", err
	}
	fmt.Println(r)
	s.mu.Unlock()
	return xmlID, nil
}

func (s *Store) GetURL(xmlID string) (string, error) {
	logger.Log.Info("start function GetURL")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var URL string

	row := s.db.QueryRowContext(ctx, "SELECT url FROM urls WHERE xmlid = $1", xmlID)
	err := row.Scan(&URL)
	if err != nil {
		logger.Log.Error("GetURL scan error", err)
		return "", err
	}
	err = row.Err()
	if err != nil {
		logger.Log.Error("GetURL error", err)
		return "", err
	}

	return URL, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) CreateTableForDB(ctx context.Context) error {
	logger.Log.Info("start CreateTableForDB")
	result, err := s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls ("+
		"id SERIAL PRIMARY KEY,"+
		"url VARCHAR(255) NOT NULL,"+
		"xmlid VARCHAR(50) NOT NULL UNIQUE)")
	if err != nil {
		logger.Log.Error("CreateTableForDB error", err)
		return err
	}
	logger.Log.Info(result)

	return nil
}
func (s *Store) AddURLBatch(data []models.RequestBatch) ([]models.ResponseBatch, error) {

	var result []models.ResponseBatch

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO urls (url, xmlid) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM urls WHERE xmlid = $3)")
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		s.mu.Lock()
		_, err := stmt.Exec(v.LongURL, v.ID, v.ID)
		s.mu.Unlock()

		if err != nil {
			return nil, err
		}

		var tmp = models.ResponseBatch{
			ID:       v.ID,
			ShortURL: s.url + "/" + v.ID,
		}
		result = append(result, tmp)

	}
	tx.Commit()

	return result, nil
}
