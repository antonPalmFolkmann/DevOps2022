package templates

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var user M

// Hack for an array of maps in golang:
// https://stackoverflow.com/questions/47130003/how-can-i-declare-list-of-maps-in-golang
type M map[string]interface{}

// Return the gravatar image for the given email address.
// Converting string to bytes: https://stackoverflow.com/questions/42541297/equivalent-of-pythons-encodeutf8-in-golang
// Converting bytes to hexadecimal s%}tring: https://pkg.go.dev/encoding/hex#EncodeToString
func GravatarUrl(email interface{}, size int) string {
	strEmail := email.(string)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d",
		hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(strEmail)))), size)
}

func initTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"gravatar":       func(size int, email interface{}) string { return GravatarUrl(email, size) },
		"datetimeformat": FormatDatetime,
	})
}

func AddMessageTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("addmessage.html").ParseFiles("templates/layout.html", "templates/addmessage.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render login template with err: %v", err)
	}
}

func LoginTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("login.html").ParseFiles("templates/layout.html", "templates/login.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render login template with err: %v", err)
	}
}

func RegisterTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("register.html").ParseFiles("templates/layout.html", "templates/register.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func TimelineTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("timeline.html").ParseFiles("templates/layout.html", "templates/timeline.html")
	if err != nil {
		log.Printf("Failed to parse the templates with err: %v", err)
	}

	err = tmpl.ExecuteTemplate(w, "timeline.html", data)
	if err != nil {
		log.Printf("Failed to render the template with err: %v", err)
	}
}

func FormatDatetime(timestamp int64) string {
	timeUnix := time.Unix(timestamp, 0)
	return timeUnix.Format("2006-01-02 15:04")
}