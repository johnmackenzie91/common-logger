package example_test

import (
	"context"
	"os"
	"time"

	"net/http"
	"strings"

	"errors"

	"github.com/johnmackenzie91/commonlogger"
	"github.com/johnmackenzie91/commonlogger/resolvers"
	"github.com/sirupsen/logrus"
)

func Example_Test() {
	// Configure logrus the way YOU want to
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)

	// Instansiate a logger with the commonlogger.ErrorInfoDebugger interface
	sut := commonlogger.New(l, commonlogger.Config{
		// for every context.Context received attempt to get the "x-request-id" value
		ContextResolver: func(ctx context.Context) (logrus.Fields, error) {
			f := logrus.Fields{}
			if reqID := ctx.Value("x-request-id"); reqID != "" {
				f["x-request-id"] = reqID
			}
			return f, nil
		},
		// we are "hopefully" receiving JSON requests, attempt to log out request using predefined helper func
		RequestResolver: resolvers.ResolveJSONRequest,
	})

	// A time is set for example test assertion
	nowTime := time.Date(2006, 01, 02, 4, 5, 6, 7, time.UTC)

	// A request comes into our app
	ctx := context.WithValue(context.Background(), "x-request-id", "0000-1111-2222-3333")
	r, _ := http.NewRequest("GET", "http://example.com", strings.NewReader("{\"key\":\"value\"}"))

	// Log an info
	sut.Info(ctx, r, "we receive a request", nowTime)
	// {"level":"info","x-request-id":"0000-1111-2222-3333","msg":"we receive a request","request":{"body":{"key":"value"},"headers":{},"method":"GET","url":"http://example.com"},"time":"2006-01-02T04:05:06Z"}

	sut.Debug(ctx, "some debug line", nowTime)
	// {"level":"debug","x-request-id":"0000-1111-2222-3333","msg":"some debug line","time":"2006-01-02T04:05:06Z"}

	// An error has been received
	err := errors.New("some error")

	sut.Error(ctx, err, nowTime)
	// {"level":"error", "x-request-id":"0000-1111-2222-3333", "error":"some error","msg":"","time":"2006-01-02T04:05:06Z"}
}
