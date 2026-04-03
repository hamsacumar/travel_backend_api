package db

import (
    "database/sql"
    "time"

    "github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type TokenRepo struct {
    DB *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
    return &TokenRepo{DB: db}
}

func (r *TokenRepo) Save(t entity.Token) error {
    query := `
        INSERT INTO tokens (id, user_id, role, token, expires_at, revoked, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    var expires interface{}
    if t.ExpiresAt != nil {
        expires = *t.ExpiresAt
    } else {
        expires = nil
    }
    _, err := r.DB.Exec(query,
        t.ID, t.UserID, t.Role, t.Token, expires, t.Revoked, time.Now(),
    )
    return err
}

func (r *TokenRepo) FindByToken(token string) (*entity.Token, error) {
    query := `
        SELECT id, user_id, role, token, expires_at, revoked, created_at
        FROM tokens WHERE token = $1
    `
    row := r.DB.QueryRow(query, token)
    var t entity.Token
    var nt sql.NullTime
    err := row.Scan(&t.ID, &t.UserID, &t.Role, &t.Token, &nt, &t.Revoked, &t.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    if nt.Valid {
        ts := nt.Time
        t.ExpiresAt = &ts
    } else {
        t.ExpiresAt = nil
    }
    return &t, nil
}

func (r *TokenRepo) Revoke(token string) error {
    query := `UPDATE tokens SET revoked = true WHERE token = $1`
    _, err := r.DB.Exec(query, token)
    return err
}
