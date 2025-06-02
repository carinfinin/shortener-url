package storepg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/carinfinin/shortener-url/internal/app/auth"
	"github.com/carinfinin/shortener-url/internal/app/config"
	"github.com/carinfinin/shortener-url/internal/app/logger"
	"github.com/carinfinin/shortener-url/internal/app/models"
	"github.com/carinfinin/shortener-url/internal/app/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

// Store sql хранилище реализует интерфейс Repository.
type Store struct {
	db  *sql.DB
	url string
}

// New конструктор для  Store.
func New(cfg *config.Config) (*Store, error) {

	db, err := sql.Open("pgx", cfg.DBPath)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:  db,
		url: cfg.URL,
	}, nil
}

// Ping проверяет достуаность хранилища
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

// AddURL записывает в хранилище урл.
func (s *Store) AddURL(ctx context.Context, url string) (string, error) {
	logger.Log.Info("start function AddURL")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ID := storage.GenerateXMLID(storage.LengthXMLID)
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return "", auth.ErrorUserNotFound
	}

	_, err := s.db.ExecContext(ctx, "INSERT INTO urls (url, user_id, xmlid) VALUES ($1, $2, $3);", url, userID, ID)
	if err != nil {
		var errPG *pgconn.PgError
		if errors.As(err, &errPG) && pgerrcode.IsIntegrityConstraintViolation(errPG.Code) {

			logger.Log.Error(" AddURL error : дублирование URL")
			row := s.db.QueryRowContext(ctx, "SELECT xmlid FROM urls WHERE url = $1 AND is_deleted = FALSE;", url)
			if err = row.Scan(&ID); err != nil {
				return "", err
			}
			return ID, storage.ErrDouble
		}
		logger.Log.Error(" AddURL error :", err)
		return "", err
	}

	return ID, nil
}

// GetURL получает из хранилища урл.
func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	logger.Log.Info("start function GetURL")
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var URL string
	var deleted bool

	row := s.db.QueryRowContext(ctx, "SELECT url, is_deleted FROM urls WHERE xmlid = $1", id)
	err := row.Scan(&URL, &deleted)
	if err != nil {
		logger.Log.Error("GetURL scan error", err)
		return "", err
	}
	err = row.Err()
	if err != nil {
		logger.Log.Error("GetURL error", err)
		return "", err
	}
	if deleted {
		logger.Log.Error("deleted url")
		return "", storage.ErrDeleteURL
	}

	return URL, nil
}

// Close закрывает хранилище.
func (s *Store) Close() error {
	return s.db.Close()
}

// CreateTableForDB создаёт необходимые таблицы.
func (s *Store) CreateTableForDB(ctx context.Context) error {
	logger.Log.Info("start CreateTableForDB")
	result, err := s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS urls ("+
		"id SERIAL PRIMARY KEY,"+
		"url VARCHAR(255) NOT NULL UNIQUE,"+
		"user_id VARCHAR(255) NOT NULL,"+
		"xmlid VARCHAR(50) NOT NULL UNIQUE,"+
		"is_deleted BOOLEAN NOT NULL DEFAULT FALSE)")
	if err != nil {
		logger.Log.Error("CreateTableForDB error", err)
		return err
	}
	logger.Log.Info(result)

	return nil
}

// AddURLBatch добавляет добавляет пачку урлов.
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

	stmt, err := tx.Prepare("INSERT INTO urls (url, user_id, xmlid) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM urls WHERE xmlid = $4 AND is_deleted = FALSE)")
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		_, err := stmt.Exec(v.LongURL, userID, v.ID, v.ID)

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

// GetUserURLs получает урлы пользователя.
func (s *Store) GetUserURLs(ctx context.Context) ([]models.UserURL, error) {

	result := []models.UserURL{}
	userID, ok := ctx.Value(auth.NameCookie).(string)
	if !ok {
		return nil, auth.ErrorUserNotFound
	}

	rows, err := s.db.QueryContext(ctx, "SELECT url, xmlid FROM urls WHERE user_id = $1 AND is_deleted = FALSE ORDER BY id", userID)
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

		result = append(result, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUserURLs удаляет  урлы пользователя.
func (s *Store) DeleteUserURLs(ctx context.Context, data []models.DeleteURLUser) error {

	tx, err := s.db.Begin()
	if err != nil {
		logger.Log.Debug("tx.Begin", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE urls SET is_deleted = TRUE WHERE xmlid = $1 AND user_id = $2")
	if err != nil {
		logger.Log.Debug("tx.PrepareContext", err)
		return err
	}
	defer stmt.Close()

	for _, v := range data {
		_, err = stmt.ExecContext(ctx, v.Data, v.USerID)
		if err != nil {
			logger.Log.Debug("tx.ExecContext", err)

			return err
		}
	}
	return tx.Commit()
}

// Stat получение общей статистики
func (s *Store) Stat(ctx context.Context) (*models.Stat, error) {
	var stat models.Stat
	row := s.db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) as user_count, COUNT(url) as url_count FROM urls WHERE is_deleted = FALSE")

	err := row.Scan(&stat.Users, &stat.URLs)
	if err != nil {
		return nil, err
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return &stat, nil
}
