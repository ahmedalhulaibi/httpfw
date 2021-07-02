package httpfw

import "net/http"

const XUserUUID = "X-User-UUID"
const XRequestID = "X-Request-ID"

func GetStringFromHeader(r *http.Request, header string) string {
	return r.Header.Get(header)
}
