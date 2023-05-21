package model

import "github.com/uptrace/bun"

type IdentityTraits map[string]any

const (
	IdentityTraitsEmail    = "email"
	IdentityTraitsPhone    = "phone"
	IdentityTraitsUsername = "username"
	IdentityTraitsName     = "name"
	IdentityTraitsAvatar   = "avatar"
)

type Identity struct {
	bun.BaseModel `bun:"table:identities,alias:i"`

	ID          int64          `bun:",pk,autoincrement"`
	Traits      IdentityTraits `bun:"type:jsonb"`
	Credentials []*Credential  `bun:"rel:has-many,join:id=identity_id"`
}
