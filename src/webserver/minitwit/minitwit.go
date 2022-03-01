package minitwit

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/antonPalmFolkmann/DevOps2022/templates"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	PER_PAGE = 30
)

type messageData struct {
	Request *http.Request
	Message string
	User    interface{}
	Error   string
}

// Registers a new message for the user.
func AddMessage(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	userError := ""
	/* if _, found := session["user_id"]; !found {
		log.Fatalln("Abort 401")
	} */

	r.ParseForm()
	if _, found := r.Form["message"]; found {
		// Avoid SQL injections
		err := storage.AddMessageQuery(r)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("SHOULD FLASH: Your message was recorded")
		http.Redirect(w, r, "/public", http.StatusFound)
	}

	message := ""
	if len(r.Form["message"]) != 0 {
		message = r.Form["message"][0]
	}

	data := messageData{
		Request: r,
		Message: message,
		User:    storage.UserM,
		Error:   userError,
	}
	templates.AddMessageTemplate(w, data)
}

type loginData struct {
	Request  *http.Request
	Username string
	User     interface{}
	Error    string
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	if _, found := storage.Session["user_id"]; found {
		log.Printf("Session is: %v", storage.Session)
		http.Redirect(w, r, "/", http.StatusMultipleChoices)
		return
	}

	userError := ""
	if r.Method == "POST" {
		r.ParseForm()

		if _, found := r.Form["username"]; found {
			//We concatenate like this because variable assignment with % doesn't seem to work here
			queryResult := storage.LoginQuery(r)
			log.Println(queryResult)
			log.Printf("Query result: %v", queryResult)
			storage.UserM = queryResult[0]

			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])
			formPwHash := fmt.Sprintf("%x", hash.Sum(nil))

			if queryResult == nil {
				userError = "Invalid username"
			} else if queryResult[0]["pw_hash"].(string) != formPwHash {
				userError = "Invalid password"
			} else {
				log.Printf("SHOULD FLASH: You were logged in")
				queryUserID := queryResult[0]["user_id"].(int64)
				storage.Session["user_id"] = strconv.Itoa(int(queryUserID))
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
	}

	username := ""
	if len(r.Form["username"]) != 0 {
		username = r.Form["username"][0]
	}

	data := loginData{
		Request:  r,
		Username: username,
		User:     storage.UserM,
		Error:    userError,
	}

	templates.LoginTemplate(w, data)
}

type registerData struct {
	Request  *http.Request
	Username string
	Email    string
	User     interface{}
	Error    string
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	if _, found := storage.Session["user_id"]; found {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	registerError := ""
	if r.Method == "POST" {
		r.ParseForm()
		if _, found := r.Form["username"]; !found {
			registerError = "Please enter a username"
		} else if _, found := r.Form["email"]; !found {
			registerError = "Please enter a valid e-mail address"
		} else if !strings.Contains(r.Form["email"][0], "@") {
			registerError = "Please enter a valid e-mail address"
		} else if _, found := r.Form["password"]; !found {
			registerError = "Please enter a password"
		} else if _, err := storage.UserNameExistsInDB(r.Form["username"][0]); err != nil {
			registerError = "Username already taken"
		} else {
			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])

			err := storage.CreateUserQuery(r, hash)

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("SHOULD FLASH: You were successfully registered and can login now")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	username := ""
	if len(r.Form["username"]) != 0 {
		username = r.Form["username"][0]
	}

	email := ""
	if len(r.Form["email"]) != 0 {
		email = r.Form["email"][0]
	}

	data := &registerData{
		Request:  r,
		User:     storage.UserM,
		Username: username,
		Email:    email,
		Error:    registerError,
	}

	templates.RegisterTemplate(w, &data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	log.Printf("SHOULD FLASH: You were logged out")
	delete(storage.Session, "user_id")
	http.Redirect(w, r, "/public", http.StatusOK)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	vars := mux.Vars(r)
	username := vars["username"]

	if _, found := storage.Session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		err := storage.CreateNewFollowingQuery(r)
		if err != nil {
			log.Fatalln(err.Error())
		}

		redirectTo := fmt.Sprintf("/user/%s", username)
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	vars := mux.Vars(r)
	username := vars["username"]

	if _, found := storage.Session["user_id"]; !found {
		log.Fatalln("Abort 401")
	}

	r.ParseForm()
	if _, found := r.Form["text"]; found {
		// Avoid SQL injections
		err := storage.DeleteFollowerQuery(r)
		if err != nil {
			log.Fatalln(err.Error())
		}

		redirectTo := fmt.Sprintf("/user/%s", username)
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}

func GetMessagesFromURL(url string) []storage.M {
	var resultMap []storage.M
	split := strings.Split(url, "/")

	if split[1] == "public" {
		resultMap = storage.GetAllNonFlaggedMessages()
	} else if split[1] == "" {
		resultMap = storage.GetAllMessages()
	} else if split[1] == "user_timeline" {
		resultMap = storage.GetAllNonFlaggedMessagesFromUser(split[2])
	}
	return resultMap
}

type timelineData struct {
	Title       string
	Request     *http.Request
	Messages    []storage.M
	UserId      string
	User        storage.M
	Followed    bool
	ProfileUser storage.M
	PerPage     int
}

// Shows a users timeline or if no user is logged in it will
// redirect to the public timeline.  This timeline shows the user's
// messages as well as all the messages of followed users.
func Timeline(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	log.Printf("We got a vistor from %s", r.RemoteAddr)
	log.Printf("User is: %v", storage.UserM)

	if storage.UserM == nil {
		http.Redirect(w, r, "/public", http.StatusMultipleChoices)
		return
	}

	log.Printf("User is: %v", storage.UserM)

	_ = r.URL.Query().Get("offset")

	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id and ( user.user_id = ? or user.user_id in (select whom_id from follower where who_id = ?)) order by message.pub_date desc limit ?"

	data := timelineData{
		Title:    "Public  Timeline",
		Request:  r,
		Messages: storage.QueryDb(messageQuery, false, storage.Session["user_id"], storage.Session["user_id"], PER_PAGE),
		UserId:   storage.Session["user_id"],
		User:     storage.UserM,
		PerPage:  PER_PAGE,
	}

	templates.TimelineTemplate(w, data)
}

// Displays the latest messages of all users.
func PublicTimeline(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	messageQuery := "select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30"

	data := timelineData{
		Title:    "Public Timeline",
		Request:  r,
		Messages: storage.QueryDb(messageQuery, false),
		User:     storage.UserM,
		PerPage:  PER_PAGE,
	}

	templates.TimelineTemplate(w, data)
}

// Displays a user's tweets
func UserTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	log.Println("User is:", username)

	ProfileUser := storage.QueryDb("select * from user where username = ?", true, username)[0]
	log.Println(ProfileUser)
	if ProfileUser == nil {
		w.Write([]byte("404 Not Found"))
	}

	followed := false
	if storage.UserM != nil {
		followed = len(storage.QueryDb("select 1 from follower where follower.who_id = ? and follower.whom_id = ?", true, storage.Session["user_id"], ProfileUser["user_id"])) == 0
	}

	messages := storage.QueryDb("select * from message limit 50", false)
	log.Println("messages: ", messages)

	data := timelineData{
		Title:       "User Timeline",
		Request:     r,
		Messages:    storage.QueryDb("select message.*, user.* from message, user where user.user_id = message.author_id and user.user_id = ? order by message.pub_date desc limit ?", false, ProfileUser["user_id"], PER_PAGE),
		ProfileUser: ProfileUser,
		Followed:    followed,
		PerPage:     PER_PAGE,
		User:        storage.UserM,
	}

	templates.TimelineTemplate(w, data)
}

func ServeCSS(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	http.ServeFile(w, r, "static/style.css")
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	defer storage.AfterRequest()
	w.Write([]byte("Gorilla!\n"))
}

func SetupRoutes(r *mux.Router) {
	r.Use(storage.BeforeRequest)

	r.HandleFunc("/static/style.css", ServeCSS)

	r.HandleFunc("/", Timeline)
	r.HandleFunc("/public", PublicTimeline)
	r.HandleFunc("/user/{username}", UserTimeline)

	r.HandleFunc("/user/{username}/follow", FollowUser)
	r.HandleFunc("/user/{username}/unfollow", UnfollowUser)
	r.HandleFunc("/addmessage", AddMessage)

	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/register", Register)
}