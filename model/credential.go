package model

import "github.com/uptrace/bun"

type CredentialMethod string

const (
	CredentialMethodPassword CredentialMethod = "PASSWORD"
	CredentialMethodOAuth    CredentialMethod = "OAUTH"
)

type Credential struct {
	bun.BaseModel `bun:"table:credentials,alias:c"`

	ID         int64            `bun:",pk,autoincrement"`
	Method     CredentialMethod `bun:",notnull,default:PASSWORD"`
	Provider   string           `bun:",notnull"`
	Secret     string           `bun:",notnull"`
	IdentityID int64            `bun:",notnull"`
	Identity   *Identity        `bun:"rel:belongs-to,join:identity_id=id"`
	Sessions   []*Session       `bun:"rel:has-many,join:id=credential_id"`
}
