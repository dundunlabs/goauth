package v1

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/dundunlabs/goauth/brutil"
	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/httperror"
	"github.com/dundunlabs/goauth/model"
	jwtservice "github.com/dundunlabs/goauth/service/jwt_service"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"
)

type SigninHandler struct {
	*common.BaseHanlder
	JWTService *jwtservice.JWTService
}

type SigninBody struct {
	Method   model.CredentialMethod `validate:"oneof=PASSWORD OAUTH"`
	Key      *string                `validate:"required_if=Method PASSWORD,omitempty"`
	Password *string                `validate:"required_if=Method PASSWORD,omitempty"`
	Provider *string                `validate:"required_if=Method OAUTH,omitempty"`
	Code     *string                `validate:"required_if=Method OAUTH,omitempty"`
	State    *uuid.UUID             `validate:"required_if=Method OAUTH,omitempty"`
}

func (h *SigninHandler) Signin(w http.ResponseWriter, req bunrouter.Request) error {
	var body SigninBody
	parser, err := brutil.NewBodyParser(req.Body)
	if err != nil {
		return err
	}
	if err := parser.ParseJSON(&body); err != nil {
		return err
	}

	cred := new(model.Credential)
	switch body.Method {
	case model.CredentialMethodPassword:
		if err := h.singinWithPassword(req.Context(), body, cred); err != nil {
			return err
		}
	case model.CredentialMethodOAuth:
		if err := h.singinWithOAuth(req.Context(), body, cred); err != nil {
			return err
		}
	default:
		return httperror.New(http.StatusBadRequest, "bad_request", "invalid methods")
	}

	sess := new(model.Session)
	sess.CredentialID = cred.ID
	sess.Credential = cred
	if _, err := h.Config.DB.NewInsert().Model(sess).Exec(req.Context()); err != nil {
		return err
	}
	token, err := h.JWTService.Sign(sess)
	if err != nil {
		return err
	}
	return brutil.SendJSON(w, http.StatusOK, bunrouter.H{
		"session": sess,
		"token":   token,
	})
}

func (h *SigninHandler) singinWithPassword(ctx context.Context, body SigninBody, cred *model.Credential) error {
	if err := h.Config.DB.NewSelect().Model(cred).
		Relation("Identity").
		Where("method = ?", model.CredentialMethodPassword).
		Where("i.id is not null").
		WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("i.traits ->> 'email' = ?", body.Key).
				WhereOr("i.traits ->> 'phone' = ?", body.Key).
				WhereOr("i.traits ->> 'username' = ?", body.Key)
		}).
		Scan(ctx); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.Secret), []byte(*body.Password)); err != nil {
		return httperror.New(http.StatusUnprocessableEntity, "invalid_data", "invalid password")
	}

	return nil
}

func (h *SigninHandler) singinWithOAuth(ctx context.Context, body SigninBody, cred *model.Credential) error {
	conf := h.Config.OAuth[*body.Provider]
	if conf == nil {
		return httperror.New(http.StatusBadRequest, "bad_request", "invalid provider")
	}
	oauthState := &model.OAuthState{
		ID: *body.State,
	}
	if err := h.Config.DB.NewSelect().Model(oauthState).WherePK().Scan(ctx); err != nil {
		return httperror.New(http.StatusUnprocessableEntity, "invalid_state", "invalid state")
	}
	if oauthState.ExpiresAt.Before(time.Now()) {
		return httperror.New(http.StatusUnprocessableEntity, "expired_state", "expired state")
	}

	if _, err := h.Config.DB.NewDelete().Model(oauthState).WherePK().Exec(ctx); err != nil {
		return err
	}

	auth, err := conf.ExchangeAuthInfo(ctx, *body.Code)
	if err != nil {
		return err
	}

	if err := h.Config.DB.NewSelect().Model(cred).
		Where("method = ?", model.CredentialMethodOAuth).
		Where("provider = ?", *body.Provider).
		Where("secret = ?", auth.ID).
		Scan(ctx); err == nil {
		return nil
	}

	if err := h.Config.DB.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		i := &model.Identity{
			Traits: model.IdentityTraits{
				model.IdentityTraitsEmail:  auth.Email,
				model.IdentityTraitsName:   auth.Name,
				model.IdentityTraitsAvatar: auth.Picture,
			},
		}
		if _, err := tx.NewInsert().Model(i).Exec(ctx); err != nil {
			return err
		}

		cred.Method = model.CredentialMethodOAuth
		cred.Provider = *body.Provider
		cred.Secret = auth.ID
		cred.IdentityID = i.ID
		_, err := tx.NewInsert().Model(cred).Exec(ctx)
		return err
	}); err != nil {
		return httperror.New(http.StatusUnprocessableEntity, "unprocessable_enitty", err.Error())
	}

	return nil
}
