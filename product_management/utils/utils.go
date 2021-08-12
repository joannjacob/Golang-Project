package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Logger return log message
func Logger() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now(), r.Method, r.URL)
	})
}

func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ErrorsWrap(err error, s string) error {
	path := strings.Trim(getCallerString(), os.Getenv("GOPATH"))
	return errors.Wrap(err, path+s)
}

func getCallerString() string {
	_, f, l, _ := runtime.Caller(2)
	return fmt.Sprintf("%v:%v:", f, l)
}
