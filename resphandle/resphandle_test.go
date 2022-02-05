package resphandle

import (
	"net/http"
	"testing"
)

func CheckStstusCode(t *testing.T) {
	var w http.ResponseWriter

	HandleResp(w, 200, true, "sucesss")
}

func CheckStatusCode2(t *testing.T) {
	var w http.ResponseWriter

	HandleResp(w, 200, true, "sucesss")
}
