package v1

import (
	"net/http"
	"time"

	"github.com/dundunlabs/goauth/brutil"
	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/httperror"
	"github.com/dundunlabs/goauth/model"
	"github.com/uptrace/bunrouter"
)

type OAuthHandler struct {
	*common.BaseHanlder
}

type AuthCodeURLBody struct {
	Provider string `validate:"required"`
}

func (h *OAuthHandler) AuthCodeURL(w http.ResponseWriter, req bunrouter.Request) error {
	parser, err := brutil.NewBodyParser(req.Body)
	if err != nil {
		return err
	}
	var body AuthCodeURLBody
	if err := parser.ParseJSON(&body); err != nil {
		return err
	}
	conf := h.Config.OAuth[body.Provider]
	if conf == nil {
		return httperror.New(http.StatusBadRequest, "bad_request", "invalid provider")
	}
	now := time.Now()
	oauthState := &model.OAuthState{
		CreatedAt: now,
		ExpiresAt: now.Add(5 * time.Minute),
	}
	if _, err := h.Config.DB.NewInsert().Model(oauthState).
		Exec(req.Context()); err != nil {
		return err
	}

	return brutil.SendJSON(w, http.StatusOK, bunrouter.H{
		"url": conf.AuthCodeURL(oauthState.ID.String()),
	})
}
