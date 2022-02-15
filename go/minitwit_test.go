package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setUp() {
	db, _ = sql.Open("sqlite3", ":memory:")
	// Use httptest package instead of minitwit.app.test_client()
	DATABASE = ":memory:"
	InitDb()
}

// Helper functions Login, Logout, and RegisterAndLogin

func login(username string, password string) {
	var jsonData = []byte(`{
		"username": username,
		"password": password
	}`)
	_, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
}

func logout() {
	_, err := http.NewRequest("POST", "/logout", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("fatal")
	}
}

func register(username string, password string, password2 string, email string) {
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}

	var jsonData = []byte(`{
		"username": username,
		"password": password,
		"password2": password2,
		"email": email
	}`)
	_, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
}

func addMessage(t *testing.T, text string) {
	rv, err := http.NewRequest("POST", "/add_message", bytes.NewBuffer([]byte(text)))
	if err != nil {
		log.Fatal("fatal")
	}
	assert.Equal(t, "Your message was recorded", rv)
}

func RegisterAndLogin(username string, password string) {
	register(username, password, "", "")
	login(username, password)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Tell the client that the API version is 1.3
	w.Header().Add("API-VERSION", "1.3")
	w.Write([]byte("ok"))
}

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	w := httptest.NewRecorder()
	Handler(w, req)
	// We should get a good status code
	want, got := http.StatusOK, w.Result().StatusCode

	assert.Equal(t, want, got)

	// Make sure that the version was 1.3
	if want, got := "1.3", w.Result().Header.Get("API-VERSION"); want != got {
		t.Fatalf("expected API-VERSION to be %s, instead got: %s", want, got)
	}
}

/*

func TestRegister(t *testing.T) {
	// Make sure registering works
	rv = Register("user1", "default")
	assert.Equal(t, "You were successfully registered and can login now", rv.data)
	rv = Register("user1", "default")
	assert.Equal(t, "The username is already taken", rv.data)
	rv = Register("", "default")
	assert.Equal(t, "You have to enter a username", rv.data)
	rv = Register("meh", "")
	assert.Equal(t, "You have to enter a password", rv.data)
	rv = Register("meh", "x", "y")
	assert.Equal(t, "The two passwords do not match", rv.data)
	rv = Register("meh", "foo", email == "broken")
	assert.Equal(t, "You have to enter a valid email address", rv.data)
}

func TestLoginLogout(t *testing.T) {
	// Make sure logging in and logging out works
	rv = RegisterAndLogin("user1", "default")
	assert.Equal(t, "You were logged in", rv.data)
	rv = LogOut()
	assert.Equal(t, "You were logged out", rv.data)
	rv = Login("user1", "wrongpassword")
	assert.Equal(t, "Invalid password", rv.data)
	rv = Login("user2", "wrongpassword")
	assert.Equal(t, "Invalid username", rv.data)
}

func TestMessageRecording(t *testing.T) {
	// heck if adding messages works
	RegisterAndLogin("foo", "default")
	AddMessage("test message 1")
	AddMessage("<test message 2>")
	rv = get('/')
	assert.Contains(t, "test message 1", rv.data)
	assert.Contains(t, "&lt;test message 2&gt;", rv.data)
}

func TestTimelines(t *testing.T) {
	// Make sure that timelines work
	RegisterAndLogin("foo", "default")
	AddMessage("the message by foo")
	LogOut()
	RegisterAndLogin("bar", "default")
	AddMessage("the message by bar")
	rv = get("/public")

	assert.Contains(t, "the message by foo", rv.data)
	assert.Contains(t, "the message by bar", rv.data)

	// bar's timeline should just show bar's message
	rv = get("/")
	assert.NotContains(t, "the message by foo", rv.data)
	assert.Contains(t, "the message by bar", rv.data)

	// now let's follow foo
	rv = get("/foo/follow", follow_redirects == True)
	assert.Contains(t, "You are now following &#34;foo&#34;", rv.data)

	// we should now see foo's message
	rv = get("/")
	assert.Contains(t, "the message by foo", rv.data)
	assert.Contains(t, "the message by bar", rv.data)

	// but on the user's page we only want the user's message
	rv = get("/bar")
	assert.NotContains(t, "the message by foo", rv.data)
	assert.Contains(t, "the message by bar", rv.data)

	rv = get("/foo")
	assert.Contains(t, "the message by foo", rv.data)
	assert.NotContains(t, "the message by bar", rv.data)

	// now unfollow and check if that worked
	rv = get("/foo/unfollow", follow_redirects == True)
	assert.Contains(t, "You are no longer following &#34;foo&#34;", rv.data)

	rv = get("/")
	assert.NotContains(t, "the message by foo", rv.data)
	assert.Contains(t, "the message by bar", rv.data)

}

*/
