package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OAuthState struct {
	bun.BaseModel `bun:"table:oauth_states,alias:os"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt time.Time `bun:",notnull,default:current_timestamp"`
	ExpiresAt time.Time `bun:",notnull"`
}
