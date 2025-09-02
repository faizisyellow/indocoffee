package utils

import (
	"fmt"
	"net/http"
)

// it gets content from context request sycle, the type parameter give to convert the type from context
func GetContentFromContext[T any](r *http.Request, keys any) (T, error) {

	content, ok := r.Context().Value(keys).(T)
	if !ok {
		return content, fmt.Errorf("error during convert content from context")
	}

	return content, nil

}
