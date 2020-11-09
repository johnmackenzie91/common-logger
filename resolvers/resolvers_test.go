package resolvers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_ResolveRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", strings.NewReader("Hello World!"))
	assert.Nil(t, err)
	r.Header.Set("x-trace-id", "1234")

	out, err := ResolveRequest(r)
	assert.Nil(t, err)

	expected := logrus.Fields{
		"request": logrus.Fields{
			"method":  "GET",
			"url":     "http://example.com",
			"body":    "Hello World!",
			"headers": http.Header{"X-Trace-Id": []string{"1234"}},
		},
	}
	assert.Equal(t, expected, out)

	// assert that the request body is still full
	expectedBody := "Hello World!"
	b, err := ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(b))
}

func Test_ResolveJSONRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", strings.NewReader("{\"msg\":\"hello world\"}"))
	assert.Nil(t, err)
	r.Header.Set("x-trace-id", "1234")

	out, err := ResolveJSONRequest(r)
	assert.Nil(t, err)

	expectedBody := make(map[string]interface{})
	expectedBody["msg"] = "hello world"
	expected := logrus.Fields{
		"request": logrus.Fields{
			"method":  "GET",
			"url":     "http://example.com",
			"body":    expectedBody,
			"headers": http.Header{"X-Trace-Id": []string{"1234"}},
		},
	}
	assert.Equal(t, expected, out)

	// assert that the request body is still full
	b, err := ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"msg\":\"hello world\"}", string(b))
}

func Test_ResolveJSONRequest_Falls_Back_To_Test_ResolveRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", strings.NewReader("{bad json}"))
	assert.Nil(t, err)
	r.Header.Set("x-trace-id", "1234")

	out, err := ResolveJSONRequest(r)
	assert.Nil(t, err)

	expectedBody := make(map[string]interface{})
	expectedBody["msg"] = "hello world"
	expected := logrus.Fields{
		"request": logrus.Fields{
			"method":  "GET",
			"url":     "http://example.com",
			"body":    "{bad json}",
			"headers": http.Header{"X-Trace-Id": []string{"1234"}},
		},
	}
	assert.Equal(t, expected, out)

	// assert that the request body is still full
	b, err := ioutil.ReadAll(r.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{bad json}", string(b))
}
