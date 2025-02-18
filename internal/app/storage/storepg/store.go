package storepg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *Store) AddURL(ctx context.Context, url string) (string, error) {
	logger.Log.Info("start function AddURL")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	s.mu.Lock()
	defer s.mu.Unlock()

	xmlID := storage.GenerateXMLID(storage.LengthXMLID)
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return "", auth.ErrorUserNotFound
	}
	//INSERT INTO urls (url, xmlid) SELECT 'dfgddfgd', '123' WHERE NOT EXISTS(SELECT 1 FROM urls WHERE xmlid = '123');
	_, err := s.db.ExecContext(ctx, "INSERT INTO urls (url, user_id, xmlid) VALUES ($1, $2, $3);", url, userID, xmlID)
	if err != nil {
		var errPG *pgconn.PgError
		if errors.As(err, &errPG) && pgerrcode.IsIntegrityConstraintViolation(errPG.Code) {

			logger.Log.Error(" AddURL error : дублирование URL")
			row := s.db.QueryRowContext(ctx, "SELECT xmlid FROM urls WHERE url = $1;", url)
			if err = row.Scan(&xmlID); err != nil {
				return "", err
			}
			return xmlID, storage.ErrDouble
		}
		logger.Log.Error(" AddURL error :", err)
		return "", err
	}

	return xmlID, nil
}

func (s *Store) GetURL(ctx context.Context, xmlID string) (string, error) {
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
		"url VARCHAR(255) NOT NULL UNIQUE,"+
		"user_id VARCHAR(255) NOT NULL,"+
		"xmlid VARCHAR(50) NOT NULL UNIQUE)")
	if err != nil {
		logger.Log.Error("CreateTableForDB error", err)
		return err
	}
	logger.Log.Info(result)

	return nil
}
func (s *Store) AddURLBatch(ctx context.Context, data []models.RequestBatch) ([]models.ResponseBatch, error) {

	var result []models.ResponseBatch

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}

	stmt, err := tx.Prepare("INSERT INTO urls (url, user_id, xmlid) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM urls WHERE xmlid = $4)")
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		s.mu.Lock()
		_, err := stmt.Exec(v.LongURL, userID, v.ID, v.ID)
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

func (s *Store) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {

	result := []models.UserURL{}
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}
	fmt.Println("token :", userID)
	rows, err := s.db.QueryContext(ctx, "SELECT url, xmlid FROM urls WHERE user_id = $1 ORDER BY id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tmp := models.UserURL{}
		err = rows.Scan(&tmp.OriginalURL, &tmp.ShortURL)
		if err != nil {
			return nil, err
		}
		tmp.ShortURL = s.url + "/" + tmp.ShortURL

		fmt.Println(tmp)
		result = append(result, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}
