package resolvers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"io"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_ResolveRequest2(t *testing.T) {
	r, err := http.NewRequest("GET", "http://example.com", nil)
	assert.Nil(t, err)
	r.Header.Set("x-trace-id", "1234")

	out, err := ResolveRequest(r)
	assert.Nil(t, err)

	expected := logrus.Fields{
		"request": logrus.Fields{
			"method":  "GET",
			"url":     "http://example.com",
			"headers": http.Header{"X-Trace-Id": []string{"1234"}},
		},
	}
	assert.Equal(t, expected, out)
}

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

func Test_readRequestBody(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    io.Reader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRequestBody(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readRequestBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readRequestBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readResponseBody(t *testing.T) {
	res := http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader("Hello World!")),
	}

	out, err := ResolveResponse(&res)
	assert.Nil(t, err)

	expected := logrus.Fields{
		"request": logrus.Fields{
			"statusCode": 200,
			"headers":    http.Header(nil),
			"body":       "Hello World!",
		},
	}
	assert.Equal(t, expected, out)

	// assert that the request body is still full
	expectedBody := "Hello World!"
	b, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, string(b))
}
