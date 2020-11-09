# The Common Logger

This is an attempt to build a minimal, slim and flexible interface for logging withing go apps.

The interface looks like;

```go
type ErrorInfoDebugger interface {
	Error(...interface{})
	Info(...interface{})
	Debug(...interface{})
}
```

The idea being that you can put `almost` whatever you want into the 
above functions and this library will attempt to transform them into logrus.Fields and supply them to the log line.

For complex types such as context (where you may want to parse a request-id),
or a request object (where you may want to remove PII data). Custom call back functions can be supplied.

### Why only three methods?

This interface has come around because on the back of Dave Cheney's [post](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) but with the addition of the `Error(...interface)` method.

### But it is tied to logrus!
Logrus is the logging library I am most confident with. Also creating an abstraction layer would just add more complexity.

### Any examples?

Example can be found [here]()