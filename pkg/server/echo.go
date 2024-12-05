package server

import (
	"net/http"

	"github.com/faruqputraaa/go-ticket/config"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/pkg/response"
	"github.com/faruqputraaa/go-ticket/pkg/route"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(cfg *config.Config,publicRoutes, privateRoutes []route.Route) *Server {
	e := echo.New()
	v1 := e.Group("api/v1")
	if len(publicRoutes) > 0 {
		for _, route := range publicRoutes {
			v1.Add(route.Method, route.Path, route.Handler)
		}
	}
	if len(privateRoutes) > 0 {
		for _, route := range privateRoutes {
			v1.Add(route.Method, route.Path, route.Handler, JWTMiddleware(cfg.JWTConfig.SecretKey), RBACMiddleware(route.Roles))
		}
	}
	return &Server{e}
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func (c echo.Context) jwt.Claims {
			return new(entity.JWTCustomClaims)
		},
		SigningKey: []byte(secretKey),
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "anda harus login untuk mengakses resource ini"))
		},
	})
}


func RBACMiddleware(roles []string) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*entity.JWTCustomClaims)

			allowed := false
			for _, role := range roles {
				if role == claims.Role {
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "anda tidak memiliki akses untuk resource ini"))
			}
			return next(c)
  }
 }
}

