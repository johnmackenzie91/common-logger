package resolvers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Strategies struct {
	ContextResolver ContextResolver
	RequestResolver RequestResolver
}

// ContextResolver is a callback function that instructs commonlogger when to do when a context is logged.
// This can be useful when you are wanting to retrieve a request id that is stored in the context
type ContextResolver func(context.Context) (logrus.Fields, error)

// RequestResolver is a callback function that instructs commonlogger when to do when a *http.Request is logged.
// This can be useful when the request you received has PII data in it and
// you wish to set a custom formatter which excludes sensitive data.
type RequestResolver func(r *http.Request) (logrus.Fields, error)

// ResolveJSONRequest parses the request body as json.
// If successful the value for logrus.fields.request["body'] will be a json object,
// make it easier to query with something like jq. If json decode fails, this func falls back to ResolveRequest()
var ResolveJSONRequest = func(r *http.Request) (logrus.Fields, error) {
	f := logrus.Fields{}

	bodyContent, err := readRequestBody(r)

	if err != nil {
		return f, err
	}

	jsonMap := make(map[string]interface{})
	if err := json.NewDecoder(bodyContent).Decode(&jsonMap); err != nil {
		return ResolveRequest(r)
	}

	f["request"] = logrus.Fields{
		"method":  r.Method,
		"url":     r.URL.String(),
		"headers": r.Header,
		"body":    jsonMap,
	}

	return f, nil
}

// ResolveRequest writes request to logrus.Fields.request["request"]
var ResolveRequest = func(r *http.Request) (logrus.Fields, error) {
	f := logrus.Fields{}

	bodyContent, err := readRequestBody(r)

	if err != nil {
		return f, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(bodyContent); err != nil {
		return f, err
	}

	f["request"] = logrus.Fields{
		"method":  r.Method,
		"url":     r.URL.String(),
		"headers": r.Header,
		"body":    buf.String(),
	}

	return f, nil
}

// readRequestBody refills the request body after reading
// As the request body is a reader and can only be read once,
// we must refill for other processes further down the line who may wish to read it.
func readRequestBody(r *http.Request) (io.Reader, error) {
	// create a temp buffer
	b := bytes.NewBuffer(make([]byte, 0))

	// TeeReader returns a Reader that writes to b what it reads from r.Body.
	reader := io.TeeReader(r.Body, b)

	// NopCloser returns a ReadCloser with a no-op Close method wrapping the provided Reader r.
	r.Body = ioutil.NopCloser(b)

	if err := r.Body.Close(); err != nil {
		return reader, err
	}

	return reader, nil
}
