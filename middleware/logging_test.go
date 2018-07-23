package middleware

import (
	"go-middleware/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	withHeaders := Compose(handlers.Basehandler(), LoggingMiddleware())

	// Create handler and call with recorder
	withHeaders.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	/*	// Check for content-type header
		if h := rr.Header().Get("Content-Type"); h != "application/json" {
			t.Errorf("handler returned wrong status code: got %v want %v",
				h, http.StatusOK)
		}

		// Check for CORS header
		if h := rr.Header().Get("Access-Control-Allow-Origin"); h != "*" {
			t.Errorf("handler returned wrong status code: got %v want %v",
				h, http.StatusOK)
		}*/
}
