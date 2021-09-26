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
	ErrorHostNotMatch            = &SmartProxyError{http.StatusBadRequest, "Request Host Not Match with Proxy Host"}
	ErrorNoneProperlyBackendNode = &SmartProxyError{http.StatusBadRequest, "Can not found Online Backend Node"}
	ErrorCreateReverseProxy      = &SmartProxyError{http.StatusBadRequest, "Create Reverse Proxy Error"}
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

// 封装http.ResponseWriter，metrics需要HTTP CODE
type SmartProxyResponseWriter struct {
	http.ResponseWriter
	HttpRetCode int
	Bytes       int
}

func NewSmartProxyResponseWriter(resp http.ResponseWriter, code int) *SmartProxyResponseWriter {
	return &SmartProxyResponseWriter{
		ResponseWriter: resp,
		HttpRetCode:    code,
	}
}

//
func (w *SmartProxyResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.Bytes += size
	return size, err
}

func (w *SmartProxyResponseWriter) WriteHeader(retcode int) {
	w.HttpRetCode = retcode
	w.ResponseWriter.WriteHeader(retcode)
}
