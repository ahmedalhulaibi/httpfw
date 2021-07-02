package instanceid

import (
	"context"
	"net/http"
)

const ContextKey = "instance_id"

type instanceIDInjector interface {
	GetInstanceID() string
}

type InstanceIDMiddleware struct {
	h http.Handler
	i instanceIDInjector
}

func New(i instanceIDInjector) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return InstanceIDMiddleware{
			h: h,
			i: i,
		}
	}
}

func (idmw InstanceIDMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if instanceID := idmw.i.GetInstanceID(); instanceID != "" {
		r = r.WithContext(context.WithValue(r.Context(), ContextKey, instanceID))
	}

	idmw.h.ServeHTTP(w, r)
}

type InstanceIDInjector struct {
	id string
}

func NewInstanceIDInjector(id string) InstanceIDInjector {
	return InstanceIDInjector{
		id: id,
	}
}

func (i InstanceIDInjector) GetInstanceID() string {
	return i.id
}
