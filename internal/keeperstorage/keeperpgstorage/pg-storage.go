package keeperpgstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"yudinsv/gophkeeper/internal/models"

	_ "github.com/lib/pq"
)

// PostgresStorage represents a PostgreSQL database connection.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new instance of PostgresStorage.
func NewPostgresStorage(uri string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}
func (s *PostgresStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	// Create the secrets table if it does not already exist
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS public.secrets (
		id SERIAL PRIMARY KEY,
		owner_id TEXT NOT NULL,
		value BYTEA NOT NULL,
		secret_type TEXT NOT NULL,
		description TEXT NOT NULL,
		is_deleted BOOLEAN NOT NULL
		ver TIMESTAMP NOT NULL,
	)`)
	if err != nil {
		return fmt.Errorf("unable to create secrets table: %v", err)
	}
	return nil
}

func (s *PostgresStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

// PutSecret adds a new secret to the database.
func (s *PostgresStorage) PutSecret(ctx context.Context, secret models.Secret) error {
	//TODO add on coflict (update)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO secrets (owner_id, value, secret_type, description, is_deleted)
		VALUES ($1, $2, $3, $4, $5)
	`, secret.OwnerID, secret.Value, secret.Type, secret.Description, secret.IsDeleted)

	if err != nil {
		return err
	}

	return nil
}

// GetSecret retrieves the first secret found in the database for a given owner ID.
func (s *PostgresStorage) GetSecret(ctx context.Context, userID string) (models.Secret, error) {
	var secret models.Secret

	err := s.db.QueryRowContext(ctx, `
		SELECT id, owner_id, value, secret_type, description, is_deleted, ver
		FROM secrets
		WHERE owner_id = $1 AND is_deleted = false
		LIMIT 1
	`, userID).Scan(
		&secret.ID,
		&secret.OwnerID,
		&secret.Value,
		&secret.Type,
		&secret.Description,
		&secret.IsDeleted,
		&secret.Ver,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Secret{}, errors.New("secret not found")
		}
		return models.Secret{}, err
	}

	return secret, nil
}

// DeleteSecret removes the first secret found in the database for a given owner ID.
func (s *PostgresStorage) DeleteSecret(ctx context.Context, userID string) error {
	err := s.db.QueryRowContext(ctx, `
		UPDATE secrets
		set owner_id = $1 AND is_deleted = true
	`, userID).Err()

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("secret not found")
		}
		return err
	}
	return nil
}

func (s *PostgresStorage) SyncSecret(ctx context.Context, userID string) ([]models.LiteSecret, error) {
	var liteSecrets []models.LiteSecret
	rows, err := s.db.QueryContext(ctx, "SELECT id, md5(value) , md5(description), is_deleted, ver FROM secrets WHERE owner_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	for rows.Next() {
		var liteSecret models.LiteSecret
		if err := rows.Scan(&liteSecret.ID, &liteSecret.ValueHash, &liteSecret.DescriptionHash, &liteSecret.IsDeleted, &liteSecret.Ver); err != nil {
			return nil, err
		}
		liteSecrets = append(liteSecrets, liteSecret)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return liteSecrets, nil
}
