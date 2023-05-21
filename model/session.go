package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions,alias:s"`

	ID              int64     `bun:",pk,autoincrement"`
	AuthenticatedAt time.Time `bun:",notnull,default:current_timestamp"`
	ExpiresAt       *time.Time
	CredentialID    int64       `bun:",notnull"`
	Credential      *Credential `bun:"rel:belongs-to,join:credential_id=id"`
}
