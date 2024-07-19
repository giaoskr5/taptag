package controllers

import (
	"database/sql"
	_ "fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"taptag/auth"
	"taptag/user"
)

// e.GET("/user/signin", controllers.SignInForm()).Name = "userSignInForm"
// e.POST("/user/signin", controllers.SignIn())

func SignUpForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
			"Title":    "Sign Up",
			"Endpoint": "/user/signup",
		})
	}
}

func SignUp(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(user.User)
		// Parse the submitted data and fill the User struct
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// i really miss match from Rust, this is horrendous
		if u.Name == "" {
			if u.Password == "" {
				return c.String(http.StatusBadRequest, "Name and password can't be empty")
			}
			return c.String(http.StatusBadRequest, "Name can't be empty")
		}
		if u.Password == "" {
			return c.String(http.StatusBadRequest, "Password can't be empty")
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
		if err != nil {
			return err
		}
		// this is safe, using '?'
		_, error := db.Exec(
			// don't insert id, will be automatically inserted for us
			"INSERT INTO users (name, password) VALUES (?, ?)", u.Name, hash,
		)
		if error != nil {
			return c.String(http.StatusConflict, "User already exists")
		}
		// val, err := res.RowsAffected()
		// if err != nil {
		// fmt.Println("ERROR", err)
		// 	return err
		// } else if val == 0 {
		// 	return c.String(http.StatusConflict, "User already exists")
		// }

		err = auth.GenerateTokensAndSetCookies(u, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
		}

		c.Response().Header().Add("HX-Redirect", "/dashboard")
		return c.Redirect(http.StatusOK, "/dashboard")
		// return c.Render(http.StatusOK, "success.html", dict)
	}
}
