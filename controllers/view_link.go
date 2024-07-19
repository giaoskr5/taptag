package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Data struct {
	Links []string
}

func ViewLink(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		res, err := db.Query(
			"SELECT ul.link from user_links ul JOIN users u ON ul.user_id = u.id WHERE u.name = ?",
			name,
		)
		if err != nil {
			fmt.Println("ERROR", err)
			return c.String(http.StatusInternalServerError, "Server error")
		}
		var links []string
		for res.Next() {

			var link string
			res.Scan(&link)
			links = append(links, link)

		}
		data := Data{
			Links: links,
		}
		return c.Render(http.StatusOK, "userpage.html", data)
	}
}

func AddLink(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, err := c.Cookie("user")
		if err != nil {
			return c.String(http.StatusInternalServerError, "server error")
		}
		link := c.FormValue("link")
		_, error := db.Exec(
			"INSERT INTO user_links(link, user_id) VALUES (?, (SELECT id FROM users WHERE name = ? LIMIT 1))",
			link,
			username.Value,
		)
		if error != nil {
			fmt.Println("ERROR", error)
			return c.String(http.StatusInternalServerError, "server error")

		}
		return c.String(http.StatusOK, "All done!")
	}
}
