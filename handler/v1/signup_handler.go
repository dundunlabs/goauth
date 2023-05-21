package v1

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/dundunlabs/goauth/brutil"
	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/httperror"
	"github.com/dundunlabs/goauth/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	*common.BaseHanlder
}

type SignupBody struct {
	Email    *string `validate:"required_without_all=Username Phone,omitempty,email"`
	Phone    *string `validate:"required_without_all=Email Username,omitempty,e164"`
	Username *string `validate:"required_without_all=Email Phone,omitempty"`
	Password string  `validate:"required"`
}

func (h *SignupHandler) Signup(w http.ResponseWriter, req bunrouter.Request) error {
	var body SignupBody
	parser, err := brutil.NewBodyParser(req.Body)
	if err != nil {
		return err
	}
	if err := parser.ParseJSON(&body); err != nil {
		return err
	}

	secret, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return err
	}

	if err := h.Config.DB.RunInTx(req.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		identity := &model.Identity{
			Traits: model.IdentityTraits{
				model.IdentityTraitsEmail:    body.Email,
				model.IdentityTraitsPhone:    body.Phone,
				model.IdentityTraitsUsername: body.Username,
			},
		}
		if _, err := tx.NewInsert().Model(identity).Exec(ctx); err != nil {
			return err
		}
		credential := &model.Credential{
			Method:     model.CredentialMethodPassword,
			Secret:     string(secret),
			IdentityID: identity.ID,
		}
		_, err := tx.NewInsert().Model(credential).Exec(ctx)
		return err
	}); err != nil {
		return httperror.New(http.StatusUnprocessableEntity, "unprocessable_enitty", err.Error())
	}

	return brutil.Send(w, http.StatusCreated, nil)
}
