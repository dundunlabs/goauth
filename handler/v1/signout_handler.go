package v1

import (
	"net/http"

	"github.com/dundunlabs/goauth/common"
	"github.com/dundunlabs/goauth/model"
	"github.com/uptrace/bunrouter"
)

type SignoutHandler struct {
	*common.BaseHanlder
}

func (h *SignoutHandler) Signount(w http.ResponseWriter, req bunrouter.Request) error {
	s, ok := req.Context().Value("session").(*model.Session)
	if !ok {
		return common.ErrRecordNotFound
	}
	if _, err := h.Config.DB.NewDelete().Model(s).
		WherePK().
		Exec(req.Context()); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
