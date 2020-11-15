package commonlogger

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

// logLine is a struct that hold values to be used in constructing the log line
type logLine struct {
	ctx    context.Context
	msg    string
	fields logrus.Fields
	time   time.Time
}

// buildLogLine builds a basic structure, these fields we will be putting into the logger call
func (l Logger) buildLogLine(opts ...interface{}) logLine {
	line := logLine{
		fields: logrus.Fields{},
		time:   time.Now(),
	}
	for _, o := range opts {
		switch v := o.(type) {
		case string:
			line.msg = v
		case context.Context:
			// set the context to be used when calling log method
			line.ctx = v
			if l.strategies.ContextResolver != nil {
				fields, err := l.strategies.ContextResolver(v)
				if err != nil {
					// do something.....
					continue
				}

				for key, val := range fields {
					line.fields[key] = val
				}
			}
		case *http.Request:
			if l.strategies.RequestResolver != nil {
				fields, err := l.strategies.RequestResolver(v)
				if err != nil {
					// do something.....
					continue
				}

				for key, val := range fields {
					line.fields[key] = val
				}
			}
		case *http.Response:
			if l.strategies.ResponseResolver != nil {
				fields, err := l.strategies.ResponseResolver(v)
				if err != nil {
					// do something.....
					continue
				}

				for key, val := range fields {
					line.fields[key] = val
				}
			}
		case url.URL:
			line.fields["url"] = v
		case error:
			line.fields["error"] = v
		case logrus.Fields:
			line.fields = mergeFields(line.fields, v)
		case time.Time:
			line.time = v
		}
	}
	return line
}

func mergeFields(fieldsOne, fieldsTwo logrus.Fields) logrus.Fields {
	f := logrus.Fields{}
	for key, value := range fieldsOne {
		f[key] = value
	}
	for key, value := range fieldsTwo {
		f[key] = value
	}
	return f
}
