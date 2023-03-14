package keepsqlstorage

import (
	"context"
	"database/sql"
	"log"

	"yudinsv/gophkeeper/internal/constants"
	"yudinsv/gophkeeper/internal/models"
	"yudinsv/gophkeeper/internal/utils"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(dbPath string) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS secrets (
			id UUID PRIMARY KEY,
			value BLOB,
			secret_type TEXT,
			description TEXT,
			owner_id TEXT,
			is_deleted INTEGER DEFAULT 0,
			ver TIMESTAMP ,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return nil, err
	}
	return &SqliteStorage{db: db}, nil
}

func (s *SqliteStorage) Ping() error {
	return s.db.Ping()
}

func (s *SqliteStorage) Close() error {
	return s.db.Close()
}

func (s *SqliteStorage) PutSecret(ctx context.Context, secret models.Secret) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO secrets (id, owner_id, value, secret_type, description, is_deleted, ver)
		VALUES (?, ?,?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			owner_id = ?,
			value = ?,
			description = ?,
			is_deleted = ?,
			ver = ?`,
		secret.ID, secret.OwnerID, secret.Value, secret.Type, secret.Description, secret.IsDeleted, secret.Ver,
		secret.OwnerID, secret.Value, secret.Description, secret.IsDeleted, secret.Ver,
	)
	return err
}

// GetSecret retrieves the first secret found in the store for a given secret ID.
func (s *SqliteStorage) GetSecret(ctx context.Context, secretID uuid.UUID) (models.Secret, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, value, secret_type, description, owner_id, is_deleted, ver FROM secrets WHERE id = ? ORDER BY created_at DESC`, secretID)
	var secret models.Secret
	err := row.Scan(&secret.ID, &secret.Value, &secret.Type, &secret.Description, &secret.OwnerID, &secret.IsDeleted, &secret.Ver)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Secret{}, constants.ErrSecretNotFound
		}
		return models.Secret{}, err
	}
	return secret, nil
}

func (s *SqliteStorage) DeleteSecret(ctx context.Context, secretID uuid.UUID) error {
	res, err := s.db.ExecContext(ctx, `UPDATE secrets SET is_deleted = 1 WHERE id = ? AND is_deleted = 0`, secretID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return constants.ErrSecretNotFound
	}
	return nil
}

func (s *SqliteStorage) SyncSecret(ctx context.Context, userID string) ([]models.LiteSecret, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, value, description, is_deleted, ver FROM secrets WHERE owner_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var liteSecrets []models.LiteSecret
	for rows.Next() {
		var liteSecret models.LiteSecret
		err := rows.Scan(&liteSecret.ID, &liteSecret.ValueHash, &liteSecret.DescriptionHash, &liteSecret.IsDeleted, &liteSecret.Ver)
		if err != nil {
			return nil, err
		}
		liteSecret.ValueHash = utils.GetMD5Hash([]byte(liteSecret.ValueHash))
		liteSecret.DescriptionHash = utils.GetMD5Hash([]byte(liteSecret.DescriptionHash))
		liteSecrets = append(liteSecrets, liteSecret)
	}

	return liteSecrets, rows.Err()
}
