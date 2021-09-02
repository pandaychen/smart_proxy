package http

//https://pkg.go.dev/net/http

import "net/http"

type SmartProxyError struct {
	StatusCode int
	Errmsg     string
}

func (e *SmartProxyError) Error() string {
	return e.Errmsg
}

//内置错误
var (
	ErrorHostNotMatch = &SmartProxyError{http.StatusBadRequest, "Request Host Not Match with Proxy Host"}
)

func SmartProxyResponse(w http.ResponseWriter, err *SmartProxyError) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(err.StatusCode)
	w.Write([]byte(err.Errmsg))
}

func SmartProxyResponseWithError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
