package v1

import (
	"context"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dundunlabs/goauth/brutil"
	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/httperror"
	"github.com/dundunlabs/goauth/model"
	jwtservice "github.com/dundunlabs/goauth/service/jwt_service"
	"github.com/uptrace/bunrouter"
)

type AuthHandler struct {
	*common.BaseHanlder
	JWTService *jwtservice.JWTService
}

var bearerRegex = regexp.MustCompile(`^Bearer\s+`)

func (h *AuthHandler) Middleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		token := bearerRegex.ReplaceAllString(req.Header.Get("Authorization"), "")
		if token == "" {
			return httperror.New(http.StatusUnauthorized, "unauthenticated", "invalid token")
		}
		claims, err := h.JWTService.Parse(token)
		if err != nil {
			return err
		}
		sid, err := strconv.Atoi(claims.ID)
		if err != nil {
			return err
		}
		s := &model.Session{
			ID: int64(sid),
		}
		ctx := req.Context()
		if err := h.Config.DB.NewSelect().Model(s).
			Relation("Credential.Identity").
			WherePK().
			Scan(ctx); err != nil {
			return err
		}
		return next(w, req.WithContext(context.WithValue(ctx, "session", s)))
	}
}

func (h *AuthHandler) Me(w http.ResponseWriter, req bunrouter.Request) error {
	s := req.Context().Value("session").(*model.Session)
	return brutil.SendJSON(w, http.StatusOK, s.Credential.Identity)
}
