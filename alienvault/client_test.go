package alienvault

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestClientAuth(t *testing.T) {

	actualToken := ""
	var postedData []byte

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualToken = r.Header.Get("X-XSRF-TOKEN")
		w.Header().Add("Set-Cookie", "XSRF-TOKEN=abc123")
		w.Header().Add("Set-Cookie", "SESSION=mysession")

		if strings.HasSuffix(r.RequestURI, "/login") {
			var err error
			postedData, err = ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}))
	defer ts.Close()

	creds := Credentials{
		Username: "something",
		Password: "something",
	}

	client := NewClient(strings.Replace(ts.URL, "https://", "", -1), creds)

	err := client.Authenticate()
	require.Nil(t, err)

	expectedCreds, err := json.Marshal(creds)
	require.Nil(t, err)

	assert.Equal(t, string(expectedCreds), string(postedData))

	assert.Equal(t, "abc123", actualToken)

}
