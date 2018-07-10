package mw

import (
	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	mvcc "github.com/lpisces/bootstrap/cmd/serve/mvc/c"
	"net/http"
	//"github.com/labstack/gommon/log"
)

// CasbinAuth
func CasbinAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		e := casbin.NewEnforcer("./cmd/serve/casbin/model.conf", "./cmd/serve/casbin/policy.csv")
		user, err := mvcc.CurrentUser(c)

		// role
		role := "anonymous"
		if err != nil || user == nil {
			role = "anonymous"
		} else if user.IsAdmin() {
			role = "admin"
		} else {
			role = "member"
		}

		// resource
		res := c.Request().URL.Path

		// act
		act := c.Request().Method

		// authorize
		if !e.Enforce(role, res, act) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}
