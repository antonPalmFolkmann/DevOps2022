package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setUp() (r *mux.Router) {
	log := logrus.New()
	db, _ := gorm.Open("sqlite3", ":memory:")
	storage.Migrate(db)

	userService := services.NewUserService(db, log)
	_ = userService.CreateUser("rnsk", "rnsk@rnsk.com", "rnsk")
	_ = userService.CreateUser("siu", "uwu@uwu.mail", "o_o")
	messageService := services.NewMessageService(db, log)
	_ = messageService.CreateMessage("rnsk", "ITS A RNSK EAT RNSK WORLD!")
	_ = messageService.CreateMessage("siu", "SIIIIIIIIIIIIIIIIIIIIIIIIUUUUUUUUUUUU")
	_ = messageService.CreateMessage("rnsk", "rnsking is the newing sagging")

	store := sessions.NewCookieStore([]byte("supersecret1234"))
	userController := controllers.NewUserController(userService, messageService, store)
	messageController := controllers.NewMessage(store, messageService, userService)

	r = mux.NewRouter()
	userController.SetupRoutes(r)
	messageController.SetupRoutes(r)
	return r
}

func login(r *mux.Router) *http.Cookie {
	user := &controllers.UserReq{
		Username: "rnsk",
		Password: "rnsk",
	}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	session := resp.Result().Cookies()[0]
	return session
}

func TestRegisterGivenValidUserReturnsStatusCreated(t *testing.T) {
	r := setUp()

	user := &controllers.RegisterReq{
		UserReq: controllers.UserReq{
			Username: "faker",
			Password: "what was that!",
		},
		Email: "hideonbush@gangdem.com",
	}

	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonUser))

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
}

func TestRegisterGivenTakenUsernameReturnsStatusConflict(t *testing.T) {
	r := setUp()

	user := &controllers.RegisterReq{
		UserReq: controllers.UserReq{
			Username: "rnsk",
			Password: "rnsk",
		},
		Email: "hideonbush@gangdem.com",
	}

	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonUser))

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Result().StatusCode)
}

func TestLoginGivenExistingUserReturnsStatusOK(t *testing.T) {
	r := setUp()

	login := &controllers.UserReq{
		Username: "rnsk",
		Password: "rnsk",
	}

	jsonLogin, _ := json.Marshal(login)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Len(t, resp.Result().Cookies(), 1)
	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestLoginWhileLoggedInReturnsStatusBadRequest(t *testing.T) {
	r := setUp()

	session := login(r)

	user := &controllers.UserReq{
		Username: "rnsk",
		Password: "rnsk",
	}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
}

func TestLoginGivenWrongPasswordReturnsStatusForbidden(t *testing.T) {
	r := setUp()

	login := &controllers.UserReq{
		Username: "rnsk",
		Password: "superstrongwhatareyougoingtodoaboutit!",
	}

	jsonLogin, _ := json.Marshal(login)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Result().StatusCode)
}

func TestLoginGivenNonexistentUserReturnsStatusForbidden(t *testing.T) {
	r := setUp()

	login := &controllers.UserReq{
		Username: "yo!",
		Password: "yoyo!",
	}

	jsonLogin, _ := json.Marshal(login)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Result().StatusCode)
}

func TestLogoutWhileLoggedinReturnsStatusOk(t *testing.T) {
	r := setUp()
	session := login(r)

	user := &controllers.UserReq{
		Username: "rnsk",
		Password: "rnsk",
	}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/logout", bytes.NewBuffer(jsonUser))
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestLogoutWhileNotLoggedinReturnsStatusBadrequest(t *testing.T) {
	r := setUp()

	user := &controllers.UserReq{
		Username: "rnsk",
		Password: "rnsk",
	}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/logout", bytes.NewBuffer(jsonUser))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Result().StatusCode)
}

func TestFollowsGivenExistentUserReturnsStatusOK(t *testing.T) {
	r := setUp()
	session := login(r)

	req, _ := http.NewRequest("GET", "/fllw/siu", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestFollowsGivenNonExistentUserReturnsStatusNotFound(t *testing.T) {
	r := setUp()
	session := login(r)

	req, _ := http.NewRequest("GET", "/fllw/whywouldilivealie", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
}

func TestUnfollowsGivenExistentUserReturnsStatusOK(t *testing.T) {
	r := setUp()
	session := login(r)

	req, _ := http.NewRequest("GET", "/unfllw/siu", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestUnfollowsGivenNonExistentUserReturnsStatusNotFound(t *testing.T) {
	r := setUp()
	session := login(r)

	req, _ := http.NewRequest("GET", "/unfllw/whatupinthislife", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
}

func TestTimelineWhenLoggedInReturnsMessages(t *testing.T) {
	r := setUp()
	session := login(r)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}
