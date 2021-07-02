package requestid

import (
	"context"
	"net/http"

	httputil "github.com/ahmedalhulaibi/httpfw"
	"github.com/google/uuid"
)

const ContextKey = "request_id"

type requestIDExtractor interface {
	GetRequestID(r *http.Request) string
}

type RequestIDMiddleware struct {
	h     http.Handler
	ridex requestIDExtractor
}

func New(ridex requestIDExtractor) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return RequestIDMiddleware{
			h:     h,
			ridex: ridex,
		}
	}
}

func (ridmw RequestIDMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if requestID := ridmw.ridex.GetRequestID(r); requestID != "" {
		r = r.WithContext(context.WithValue(r.Context(), ContextKey, requestID))
	}

	ridmw.h.ServeHTTP(w, r)
}

type RequestIDGenerator struct{}

func NewRequestIDGenerator() RequestIDGenerator {
	return RequestIDGenerator{}
}

func (r RequestIDGenerator) GetRequestID(_ *http.Request) string {
	return uuid.NewString()
}

type RequestIDExtractor struct {
	header string
}

func NewRequestIDExtractor(header string) RequestIDExtractor {
	return RequestIDExtractor{header: header}
}

func (ridex RequestIDExtractor) GetRequestID(r *http.Request) string {
	return httputil.GetStringFromHeader(r, ridex.header)
}

type RequestIDChain struct {
	extractors []requestIDExtractor
}

func NewRequestIDChain(extractors ...requestIDExtractor) RequestIDChain {
	return RequestIDChain{
		extractors: extractors,
	}
}

func (ridex RequestIDChain) GetRequestID(r *http.Request) string {
	for _, extractor := range ridex.extractors {
		if reqID := extractor.GetRequestID(r); reqID != "" {
			return reqID
		}
	}

	return ""
}
