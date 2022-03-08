package templates

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
)

func AddMessageTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("addmessage.html").ParseFiles("templates/layout.html", "templates/addmessage.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.ExecuteTemplate(w, "addmessage.html", data); err != nil {
		log.Printf("Failed to render login template with err: %v", err)
	}
}

func LoginTemplate(w http.ResponseWriter, data interface{}) {
	tmpl, err := initTemplate("login.html").ParseFiles("templates/layout.html", "templates/login.html")
	if err != nil {
		log.Printf("Failed to parse login template with err: %v", err)
	}

	if err := tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
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

// Return the gravatar image for the given email address.
// Converting string to bytes: https://stackoverflow.com/questions/42541297/equivalent-of-pythons-encodeutf8-in-golang
// Converting bytes to hexadecimal s%}tring: https://pkg.go.dev/encoding/hex#EncodeToString
func initTemplate(name string) *template.Template {
	return template.New(name).Funcs(template.FuncMap{
		"gravatar":       func(size int, email interface{}) string { return GravatarUrl(email, size) },
		"datetimeformat": FormatDatetime,
		"urlContains":    urlContains,
		"printStuff":     printStuff,
		"followUser":     followUser,
	})
}

func printStuff() string {
	return "print stuff"
}

func followUser(followee_id string) bool {
	err := storage.CreateNewFollowingQuery(nil, followee_id)
	if err != nil {
		log.Fatalln("Error")
	}
	return false
}

func urlContains(url string, check string) bool {
	return strings.Contains(url, check)
}

/* func follows(userOne int32, userTwo int32) bool {
	defer storage.AfterRequest()
	storage.IsUserFollowed()
} */

func FormatDatetime(timestamp int64) string {
	timeUnix := time.Unix(timestamp, 0)
	return timeUnix.Format("2006-01-02 15:04")
}

func GravatarUrl(email interface{}, size int) string {
	strEmail := email.(string)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d",
		hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(strEmail)))), size)
}
