package goauth

import "github.com/uptrace/bunrouter"

func (r *Router) initRoutes() {
	r.WithGroup("/v1", func(g *bunrouter.Group) {
		g.POST("/signup", r.signupHdrV1.Signup)
		g.POST("/signin", r.signinHdrV1.Signin)
		g.POST("/oauth/auth-code-url", r.oauthHdrV1.AuthCodeURL)

		ag := g.WithMiddleware(r.authHdrV1.Middleware)

		ag.GET("/me", r.authHdrV1.Me)
		ag.DELETE("/signout", r.signoutHdrV1.Signount)
	})
}
