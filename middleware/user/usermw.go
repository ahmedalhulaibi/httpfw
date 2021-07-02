package user

import (
	"context"
	"net/http"

	httputil "github.com/ahmedalhulaibi/httpfw"
)

const ContextKey = "user"

type UserMiddleware struct {
	h http.Handler
}

func New() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return UserMiddleware{
			h: h,
		}
	}
}

func (u UserMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if user := httputil.GetStringFromHeader(r, httputil.XUserUUID); user != "" {
		r = r.WithContext(context.WithValue(r.Context(), ContextKey, user))
	}

	u.h.ServeHTTP(w, r)
}
