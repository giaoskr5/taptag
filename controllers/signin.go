package controllers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"taptag/auth"
	"taptag/user"
)

func Dashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		userCookie, _ := c.Cookie("user")

		dict := map[string]string{
			"Username": userCookie.Value,
		}
		return c.Render(http.StatusOK, "dashboard.html", dict)
	}
}

// SignInForm responsible for signIn Form rendering.
func SignInForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
			"Title":    "Sign In",
			"Endpoint": "/user/signin",
		})
	}
}

// SignIn will be executed after SignInForm submission.
func SignIn(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Initiate a new User struct.
		u := new(user.User)

		// Parse the submitted data and fill the User struct with the data from the SignIn form.
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		//? is safe to use in queries
		row := db.QueryRow("SELECT password FROM users WHERE name = ?", u.Name)
		var hash []byte
		if err := row.Scan(&hash); err != nil {
			// return echo.NewHTTPError(http.StatusUnauthorized, "No such user")
			return c.String(http.StatusUnauthorized, "No such user")
		}

		// Compare the stored hashed password, with the hashed version of the password that was received.
		if err := bcrypt.CompareHashAndPassword(hash, []byte(u.Password)); err != nil {
			// If the two passwords don't match, return a 401 status.

			// return echo.NewHTTPError(http.StatusUnauthorized, "Password is incorrect")
			return c.String(http.StatusUnauthorized, "Password is incorrect")
		}
		// If password is correct, generate tokens and set cookies.
		err := auth.GenerateTokensAndSetCookies(u, c)
		if err != nil {
			// return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
			return c.String(http.StatusUnauthorized, "Token is incorrect")
		}
		c.Response().Header().Add("HX-Redirect", "/dashboard")
		return c.Redirect(http.StatusOK, "/dashboard")
	}
}
