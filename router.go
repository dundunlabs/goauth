package goauth

import (
	"github.com/dundunlabs/goauth/common"
	v1 "github.com/dundunlabs/goauth/handler/v1"
	"github.com/dundunlabs/goauth/httperror"
	jwtservice "github.com/dundunlabs/goauth/service/jwt_service"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func (app *App) initRouter() *Router {
	br := bunrouter.New(
		bunrouter.Use(
			reqlog.NewMiddleware(reqlog.FromEnv()),
			httperror.ErrorMiddleware,
		),
	)
	baseHandler := &common.BaseHanlder{
		Config: app.config,
	}

	baseService := &common.BaseService{
		Config: app.config,
	}
	jwtService := &jwtservice.JWTService{
		BaseService: baseService,
	}

	r := &Router{
		Router:      br,
		signupHdrV1: &v1.SignupHandler{BaseHanlder: baseHandler},
		signinHdrV1: &v1.SigninHandler{
			BaseHanlder: baseHandler,
			JWTService:  jwtService,
		},
		signoutHdrV1: &v1.SignoutHandler{BaseHanlder: baseHandler},
		oauthHdrV1:   &v1.OAuthHandler{BaseHanlder: baseHandler},
		authHdrV1: &v1.AuthHandler{
			BaseHanlder: baseHandler,
			JWTService:  jwtService,
		},
	}

	r.initRoutes()

	return r
}

type Router struct {
	*bunrouter.Router
	signupHdrV1  *v1.SignupHandler
	signinHdrV1  *v1.SigninHandler
	signoutHdrV1 *v1.SignoutHandler
	oauthHdrV1   *v1.OAuthHandler
	authHdrV1    *v1.AuthHandler
}
