package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/dprio/redis-limiter/internal/domain"
)

type (
	TokenRepository interface {
		Save(ctx context.Context, totalRequests int64) (domain.Token, error)
		Update(ctx context.Context, token string, totalRequests int64) error
		Get(ctx context.Context, token string) (domain.Token, error)
	}
)

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (repo *tokenRepository) Save(ctx context.Context, totalRequests int64) (domain.Token, error) {
	token, err := generateToken()
	if err != nil {
		return domain.Token{}, fmt.Errorf("erro ao gerar token: %w", err)
	}

	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO tokens (hash, total_requests) VALUES (?, ?)")
	if err != nil {
		return domain.Token{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, token, totalRequests)
	if err != nil {
		return domain.Token{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		ID:            id,
		Hash:          token,
		TotalRequests: totalRequests,
	}, nil
}

func (repo *tokenRepository) Update(ctx context.Context, token string, totalRequests int64) error {
	stmt, err := repo.db.PrepareContext(ctx, "UPDATE tokens SET total_requests = ? WHERE hash = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, totalRequests, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (repo *tokenRepository) Get(ctx context.Context, token string) (domain.Token, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT id, hash, total_requests FROM tokens WHERE hash = ?", token)

	var t domain.Token
	if err := row.Scan(&t.ID, &t.Hash, &t.TotalRequests); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Token{}, domain.ErrTokenNotFound
		}
		return domain.Token{}, err
	}

	return t, nil
}

func generateToken() (string, error) {
	buf := make([]byte, 32) // 256 bits de entropia
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf), nil
}
