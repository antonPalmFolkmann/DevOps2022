package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func SetUp() {
	db, _ = sql.Open("sqlite3", ":memory:")
	// Use httptest package instead of minitwit.app.test_client()
	DATABASE = ":memory:"
	InitDb()
}

func TestSomething(t *testing.T) {
	assert.True(t, true, "True is true!")
}

/*

func setUp(self):
        """Before each test, set up a blank database"""
        self.db = tempfile.NamedTemporaryFile()
        self.app = minitwit.app.test_client()
        minitwit.DATABASE = self.db.name
        minitwit.init_db()

# helper functions

def register(self, username, password, password2=None, email=None):
        """Helper function to register a user"""
        if password2 is None:
            password2 = password
        if email is None:
            email = username + '@example.com'
        return self.app.post('/register', data={
            'username':     username,
            'password':     password,
            'password2':    password2,
            'email':        email,
        }, follow_redirects=True)

def login(self, username, password):
        """Helper function to login"""
        return self.app.post('/login', data={
            'username': username,
            'password': password
        }, follow_redirects=True)

def register_and_login(self, username, password):
        """Registers and logs in in one go"""
        self.register(username, password)
        return self.login(username, password)

def logout(self):
        """Helper function to logout"""
        return self.app.get('/logout', follow_redirects=True)

def add_message(self, text):
        """Records a message"""
        rv = self.app.post('/add_message', data={'text': text},
                                    follow_redirects=True)
        if text:
            assert 'Your message was recorded' in rv.data
        return rv

*/
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
