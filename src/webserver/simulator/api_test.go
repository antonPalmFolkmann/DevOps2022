package simulator

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	BASE_URL = "http://127.0.0.1:8081"
	DATABASE = "../minitwit.db"
	USERNAME = "simulator"
	PWD      = "super_safe!"
)

var (
	CREDENTIALS         = strings.Join([]string{USERNAME, PWD}, "")
	ENCODED_CREDENTIALS = encodeCredentials()
	HEADERS             = http.Header{
		"Connection":    []string{"close"},
		"Content-type":  []string{"application/json"},
		"Authorization": []string{fmt.Sprintf("Basic %x", ENCODED_CREDENTIALS)},
	}
)

// Helper method for ENCODED_CREDENTIALS
func encodeCredentials() []byte {
	encoded := base64.StdEncoding.EncodeToString([]byte(CREDENTIALS))
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return decoded
}

type latestResponse struct {
	Latest int `json:"latest"`
}

func TestLatestReturnsLatest(t *testing.T) {
	// Post something to update LATEST
	target := fmt.Sprintf("%s/register", BASE_URL)

	params := url.Values{}
	params.Set("latest", "1337")

	data := url.Values{
		"username": []string{"test"},
		"email":    []string{"test@test"},
		"pwd":      []string{"foo"},
	}

	target = target + "?" + params.Encode()
	resp, err := http.PostForm(target, data)
	if err != nil {
		log.Fatalf("api_test.go:60 Failed to PostFrom: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify that the latest was updated
	target = fmt.Sprintf("%s/register", BASE_URL)
	req, _ := http.NewRequest(http.MethodGet, target, nil)
	req.Header = HEADERS

	resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var respData latestResponse
	json.Unmarshal(body, &respData)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, respData.Latest, 1337)
}
