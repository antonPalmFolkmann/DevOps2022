package minitwit

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
	Db, _ = sql.Open("sqlite3", ":memory:")
	DATABASE = ":memory:"
	InitDb()
}

// Helper functions Login, Logout, and RegisterAndLogin

func login(username string, password string) string {
	var jsonData = []byte(`{
		"username": username,
		"password": password
	}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
	w := httptest.NewRecorder()
	return Login(w, req)
}

func logout() string {
	req, err := http.NewRequest("POST", "/logout", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("fatal")
	}
	w := httptest.NewRecorder()
	return Logout(w, req)
}

func register(username string, password string, password2 string, email string) string {
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
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
	w := httptest.NewRecorder()

	return Register(w, req)
}

func addMessage(text string) string {
	req, err := http.NewRequest("POST", "/add_message", bytes.NewBuffer([]byte(text)))
	if err != nil {
		log.Fatal("fatal")
	}
	w := httptest.NewRecorder()
	return AddMessage(w, req)
}

func registerAndLogin(username string, password string) []string {
	var results []string
	results[0] = register(username, password, "", "")
	results[1] = login(username, password)
	return results
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Tell the client that the API version is 1.3
	w.Header().Add("API-VERSION", "1.3")
	w.Write([]byte("ok"))
}

func TestHandler(t *testing.T) {
	setUp()
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

func TestRegister(t *testing.T) {
	setUp()
	var tests = []struct {
		name      string
		username  string
		password  string
		password2 string
		email     string
		want      string
	}{
		{"registerTest_succesful", "user1", "default", "", "", "You were successfully registered and can login now"},
		{"registerTest_usernameTaken", "user1", "default", "", "", "The username is already taken"},
		{"registerTest_noUsername", "", "default", "", "", "You have to enter a username"},
		{"registerTest_noPassword", "meh", "", "", "", "You have to enter a password"},
		{"registerTest_diffPasswords", "meh", "x", "y", "", "The two passwords do not match"},
		{"registerTest_invalidEmail", "meh", "foo", "", "broken", "You have to enter a valid email address"},
	}

	// Make sure registering works
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := register(tt.username, tt.password, tt.password2, tt.email)
			if got != tt.want {
				t.Errorf("register(%s, %s, %s, %s) got %v, want %v", tt.username, tt.password, tt.password2, tt.email, got, tt.want)
			}
		})
	}
}

func TestLoginLogout(t *testing.T) {
	setUp()
	// Make sure logging in and logging out works

	rv1 := registerAndLogin("user1", "default")
	assert.Equal(t, "You were successfully registered and can login now", rv1[0])
	assert.Equal(t, "You were logged in", rv1[1])

	rv2 := logout()
	assert.Equal(t, "You were logged out", rv2)

	rv3 := login("user1", "wrongpassword")
	assert.Equal(t, "Invalid password", rv3)

	rv4 := login("user2", "wrongpassword")
	assert.Equal(t, "Invalid username", rv4)
}

func TestMessageRecording(t *testing.T) {
	setUp()
	// check if adding messages works
	_ = registerAndLogin("foo", "default")
	rv1 := addMessage("test message 1")
	rv2 := addMessage("<test message 2>")
	// TODO: GET ALL TWEET MESSAGES IN TEXT TO CHECK THEY WERE UPLOADED rv = get('/')
	assert.Equal(t, "test message 1", rv1)
	assert.Equal(t, "&lt;test message 2&gt;", rv2)
}

func TestTimelines(t *testing.T) {
	setUp()
	// Make sure that timelines work
	_ = registerAndLogin("foo", "default")
	rv1 := addMessage("the message by foo")
	_ = logout()
	_ = registerAndLogin("bar", "default")
	rv2 := addMessage("the message by bar")
	// TODO: GET ALL TWEET MESSAGES FROM PUBLIC TIMELINE IN TEXT TO CHECK THEY WERE UPLOADED rv = get("/public") (#67)

	assert.Equal(t, "the message by foo", rv1)
	assert.Equal(t, "the message by bar", rv2)

	// bar's timeline should just show bar's message
	// TODO: GET TWEETS FROM "bars"'s timeline rv = get("/")
	// assert.NotContains(t, "the message by foo", rv.data)
	// assert.Equal(t, "the message by bar", rv2)

	// now let's follow foo
	// TODO: FOLLOW SPECIFIC PERSON FROM USERNAME rv = get("/foo/follow", follow_redirects == True)
	// TODO: GET RESPONE MESSAGE FROM FOLLOW PERSON assert.Contains(t, "You are now following &#34;foo&#34;", rv.data)

	/*
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
	*/
}
