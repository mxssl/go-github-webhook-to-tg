package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// test status code for method != POST
func Test_webhookHandlerWrongStatusCode(t *testing.T) {
	var methods = []string{
		"GET",
		"HEAD",
		"PUT",
		"DELETE",
		"CONNECT",
		"OPTIONS",
		"TRACE",
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "localhost", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		webhookHandler(res, req)

		exp := 403
		act := res.Code

		if exp != act {
			t.Fatalf("Expected %v got %v", exp, act)
		}
	}
}
