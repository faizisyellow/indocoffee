package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/faizisyellow/indocoffee/internal/logger"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"go.uber.org/zap"
)

type Envelope struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

func ResponseSuccess(w http.ResponseWriter, r *http.Request, data any, status int) {

	logger.Logger.Infow(
		"Success Response",
		zap.Any("Path", r.URL),
		zap.String("Method", r.Method),
		zap.Int("Status", status),
	)

	err := WriteHttpJson(w, Envelope{Data: data, Error: nil}, status)
	if err != nil {
		fallbackServerError(w)
	}

}

func validErrorService(err error) string {
	if err == nil {
		return "error service is nil"
	}

	v, ok := err.(*errorService.ErrorService)
	if !ok {
		return ""
	}

	if v.Internal == nil {
		return "error is nill"
	}

	return v.InternalError()
}

func ResponseClientError(w http.ResponseWriter, r *http.Request, rErr error, status int) {
	logger.Logger.Warnw(
		"Client Error Response",
		zap.Any("Path", r.URL),
		zap.String("Method", r.Method),
		zap.Int("Status", status),
		zap.String("Internal Error", validErrorService(rErr)),
	)

	err := WriteHttpJson(w, Envelope{Data: nil, Error: rErr}, status)
	if err != nil {
		fallbackServerError(w)
	}

}

func ResponseServerError(w http.ResponseWriter, r *http.Request, rErr error, status int) {
	logger.Logger.Errorw(
		"Server Error",
		zap.Any("Path", r.URL),
		zap.String("Method", r.Method),
		zap.Int("Status", status),
		zap.Error(rErr),
		zap.String("Internal Error", validErrorService(rErr)),
	)

	err := WriteHttpJson(w, Envelope{Data: nil, Error: errors.New("server encountered an internal error")}, status)
	if err != nil {
		fallbackServerError(w)
	}
}

func fallbackServerError(w http.ResponseWriter) {

	logger.Logger.Errorw(
		"Server Error",
		zap.Error(fmt.Errorf("error sending json response")),
	)

	w.WriteHeader(http.StatusInternalServerError)

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(Envelope{
		Data:  nil,
		Error: fmt.Errorf("server encounter error"),
	})

	logger.Logger.Errorw("Server Error", err)
}
