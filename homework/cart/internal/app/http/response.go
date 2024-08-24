package http

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"route256/cart/pkg/prometheus"
)

func GetErrorResponse(ctx context.Context, w http.ResponseWriter, handlerName string, err error, statusCode int) {
	w.WriteHeader(statusCode)
	prometheus.IncHttpResponseStatusTotalCounter(handlerName, strconv.Itoa(statusCode))

	buf := bytes.NewBufferString(handlerName)
	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteString("\n")
	_, _ = w.Write(buf.Bytes())
}

func GetSuccessResponse(ctx context.Context, w http.ResponseWriter, handlerName string) {
	w.WriteHeader(http.StatusOK)
	prometheus.IncHttpResponseStatusTotalCounter(handlerName, strconv.Itoa(http.StatusOK))
}

func GetNoContentResponse(ctx context.Context, w http.ResponseWriter, handlerName string) {
	w.WriteHeader(http.StatusNoContent)
	prometheus.IncHttpResponseStatusTotalCounter(handlerName, strconv.Itoa(http.StatusNoContent))
}

func GetSuccessResponseWithBody(ctx context.Context, w http.ResponseWriter, body []byte, handlerName string) {
	w.WriteHeader(http.StatusOK)
	prometheus.IncHttpResponseStatusTotalCounter(handlerName, strconv.Itoa(http.StatusOK))

	_, _ = w.Write(body)
}
