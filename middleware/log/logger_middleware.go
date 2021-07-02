package log

import (
	"net/http"
	"time"

	"github.com/ahmedalhulaibi/loggy"
)

// Logger middleware.
type Logger struct {
	h      http.Handler
	logger loggy.Logger
}

// SetLogger sets the logger to `log`. If you have used logger.New(), you can use this to set your
// logger. Alternatively, if you already have your log.Logger, then you can just call logger.NewLogger() directly.
func (l *Logger) SetLogger(logger loggy.Logger) {
	l.logger = logger
}

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// capture status.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// capture written bytes.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// NewLogger logger middleware with the given log.Logger.
func New(logger loggy.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			logger: logger,
			h:      h,
		}
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &wrapper{w, 0, 200}

	// get the context since we'll use it a few times
	ctx := r.Context()

	// Get the RequestID (should have been set in previous middleware using gomiddleware/reqid)
	// and put it into the logger which goes in the context.
	// rid := reqid.ReqIdFromContext(ctx)
	// logger := log.With(l.logger, "rid", rid)

	defer l.logger.Sync()

	// log the request.start
	l.logger.Log(ctx,
		"request started",
		"method", r.Method,
		"uri", r.RequestURI,
		"evt", "request.start",
	)

	// continue to the next middleware
	l.h.ServeHTTP(res, r.WithContext(ctx))

	// log the request.end
	l.logger.Log(ctx,
		"request ended",
		"status", res.status,
		"size", res.written,
		"duration", time.Since(start),
		"evt", "request.end",
	)
}
