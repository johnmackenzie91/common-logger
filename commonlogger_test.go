package commonlogger

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/johnmackenzie91/commonlogger/resolvers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Info(t *testing.T) {
	output := bytes.Buffer{}

	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(&output)

	ctx := context.WithValue(context.Background(), "x-request-id", "0000-1111-2222-3333")
	timestamp := time.Date(2006, 01, 02, 4, 5, 6, 7, time.UTC)

	sut := New(l, Config{
		// ContextResolver is a callback for when a context.Context is passed into the logger.
		// This is useful for retrieving values from context such as a request id
		ContextResolver: func(ctx context.Context) (logrus.Fields, error) {
			f := logrus.Fields{}
			if reqID := ctx.Value("x-request-id"); reqID != "" {
				f["x_response_id"] = reqID
			}
			return f, nil
		},
		RequestResolver: resolvers.ResolveJSONRequest,
	})

	r, _ := http.NewRequest("POST", "http://example.com", strings.NewReader("{\"key\":\"value\"}"))
	r.Header.Set("x-trace-id", "1234")
	sut.Info(ctx, "some log line", timestamp, r)

	expected := `{"level":"info","msg":"some log line","request":{"body":{"key":"value"},"headers":{"X-Trace-Id":["1234"]},"method":"POST","url":"http://example.com"},"time":"2006-01-02T04:05:06Z","x_response_id":"0000-1111-2222-3333"}
`
	assert.Equal(t, expected, output.String())
}
