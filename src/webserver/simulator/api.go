package simulator

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/antonPalmFolkmann/DevOps2022/minitwit"
	"github.com/gorilla/mux"
)

var (
	LATEST int
)

func NotReqFromSimulator(r *http.Request) []byte {
	fromSimulator := r.Header.Get("Authorization")
	if fromSimulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		error := "You are not authorized to use this resource!"
		jsonify := "{\"status\": 403, \"error_msg\":" + error + "}"
		return []byte(jsonify)
	}
	return nil
}

func UpdateLatest(r *http.Request) {
	vars := mux.Vars(r)

	latest := -1
	if latestQuery, found := vars["latest"]; found {
		asInt, _ := strconv.Atoi(latestQuery)
		latest = asInt
	}

	if latest != -1 {
		LATEST = latest
	}
}

func LatestHandler(w http.ResponseWriter, r *http.Request) {
	respMsg := fmt.Sprintf("{\"latest\": \"%v\"}", LATEST)

	var jsonData = []byte(respMsg)
	w.Write(jsonData)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	UpdateLatest(r)

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var data map[string]interface{}

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	error := ""

	if r.Method != "POST" {
		if data["username"] == "" {
			error = "You have to enter a username"
		} else if data["email"] == "" || !strings.Contains(data["email"].(string), "@") {
			error = "You have to enter a valid email address"
		} else if data["pwd"] == "" {
			error = "You have to enter a password"
		} else if _, err := minitwit.UserNameExistsInDB(r.Form["username"][0]); err != nil {
			error = "Username already taken"
		} else {
			hash := md5.New()
			io.WriteString(hash, r.Form["password"][0])

			insertMessageSQL := "INSERT INTO user (username, email, pw_hash) values (?, ?, ?)"
			stmt, _ := minitwit.Db.Prepare(insertMessageSQL)
			defer stmt.Close()
			_, err = stmt.Exec(r.Form["username"][0], r.Form["email"][0], fmt.Sprintf("%x", hash.Sum(nil)))

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("SHOULD FLASH: You were successfully registered and can login now")
			http.Redirect(w, r, "/", http.StatusNoContent)
			return
		}
	}
	if error != "" {
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		http.Redirect(w, r, "/", http.StatusNoContent)
	}
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	UpdateLatest(r)

	notFromSimResponse := NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Write(notFromSimResponse)
		return
	}

	noMessages := 100
	if arg, found := r.URL.Query()["no"]; found {
		noMessages, _ = strconv.Atoi(arg[0])
	}

	if r.Method == "GET" {
		query := "SELECT message.*, user.* FROM message, user WHERE message.flagged = 0 AND message.author_id = user.user_id ORDER BY message.pub_date DESC LIMIT ?"

		messages := minitwit.QueryDb(query, false, noMessages)

		filteredMsgs := make([]minitwit.M, 0)
		for _, msg := range messages {
			filteredMsg := make(minitwit.M, 0)
			filteredMsg["content"] = msg["text"]
			filteredMsg["pub_date"] = msg["pub_date"]
			filteredMsg["user"] = msg["username"]
			filteredMsgs = append(filteredMsgs, filteredMsg)
		}

		jsonify, _ := json.Marshal(filteredMsgs)
		w.Write(jsonify)
	}
}

func MessagesPerUsernameHandler(w http.ResponseWriter, r *http.Request) {
	UpdateLatest(r)

	vars := mux.Vars(r)
	username := vars["username"]

	notFromSimResponse := NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Write(notFromSimResponse)
		return
	}

	noMessages := 100
	if arg, found := r.URL.Query()["no"]; found {
		noMessages, _ = strconv.Atoi(arg[0])
	}

	if r.Method == "GET" {
		userId := minitwit.GetUserId(username)

		if userId == nil {
			w.WriteHeader(404)
			return
		}

		query := "SELECT message.*, user.* FROM message, user WHERE message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ? ORDER BY message.pub_date DESC LIMIT ?"
		messages := minitwit.QueryDb(query, false, userId, noMessages)

		filteredMsgs := make([]minitwit.M, 0)
		for _, msg := range messages {
			filteredMsg := make(minitwit.M)
			filteredMsg["content"] = msg["text"]
			filteredMsg["pub_date"] = msg["pub_date"]
			filteredMsgs = append(filteredMsgs, filteredMsg)
		}

		jsonify, _ := json.Marshal(filteredMsgs)
		w.Write(jsonify)
	} else if r.Method == "POST" {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		var requestData map[string]interface{}

		json.Unmarshal(body, requestData)

		query := "INSERT INTO message (author_id, text, pub_date, flagged) VALUES (?, ?, ?, 0)"
		minitwit.Db.Exec(query, requestData["content"], time.Now().Unix())

		w.WriteHeader(204)
		w.Write([]byte(""))
	}
}

func FollowsHandler(w http.ResponseWriter, r *http.Request) {
	UpdateLatest(r)

	vars := mux.Vars(r)
	username := vars["username"]

	notFromSimResponse := NotReqFromSimulator(r)
	if notFromSimResponse != nil {
		w.Write(notFromSimResponse)
		return
	}

	userId := minitwit.GetUserId(username)

	if userId == nil {
		w.WriteHeader(404)
		return
	}

	noFollowers := 100
	if arg, found := r.URL.Query()["no"]; found {
		noFollowers, _ = strconv.Atoi(arg[0])
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var data map[string]interface{}
	json.Unmarshal(body, data)

	_, hasFollowKey := data["follow"]
	_, hasUnfollowKey := data["unfollow"]
	if r.Method == "POST" && hasFollowKey {
		followsUsername := data["follow"].(string)
		followsUserId := minitwit.GetUserId(followsUsername)
		if followsUserId == nil {
			w.WriteHeader(404)
			return
		}

		query := "INSERT INTO follower (who_id, whom_id) VALUES (?, ?)"

		minitwit.Db.Exec(query, userId, followsUserId)
		// TODO: Unsure what to do with g.db.commit line

		w.WriteHeader(204)
		w.Write([]byte(""))
	} else if r.Method == "POST" && hasUnfollowKey {
		unfollowsUsername := data["unfollow"].(string)
		unfollowsUserId := minitwit.GetUserId(unfollowsUsername)
		if unfollowsUserId == nil {
			w.WriteHeader(404)
			return
		}

		query := "DELETE FROM follower WHERE who_id=? and WHOM_ID=?"
		minitwit.Db.Exec(query, userId, unfollowsUserId)

		w.WriteHeader(204)
		w.Write([]byte(""))
	} else if r.Method == "GET" {
		noFollowers = 100
		if arg, found := r.URL.Query()["no"]; found {
			noFollowers, _ = strconv.Atoi(arg[0])
		}
		query := "SELECT user.username FROM user INNER JOIN follower ON follower.whom_id=user.user_id WHERE follower.who_id=? LIMIT ?"
		followers := minitwit.QueryDb(query, false, userId, noFollowers)

		followerNames := make([]string, 0)
		for _, f := range followers {
			followerNames = append(followerNames, f["username"].(string))
		}

		followersResponse, _ := json.Marshal(followerNames)
		w.Write(followersResponse)
	}
}

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/fllws/{username}", FollowsHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/msgs", MessagesHandler)
	r.HandleFunc("/msgs/{username}", MessagesPerUsernameHandler)
	r.HandleFunc("/latest", LatestHandler)
}
