package middleware

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"go-middleware/handlers"
	"time"
	"fmt"
)

func TestCacheMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/withCaching", nil)
	if err != nil {
		t.Fatal(err)
	}
	withHeaders := Compose(handlers.Basehandler(), CacheMiddleware())

	// Create handler and call with recorder
	rr := httptest.NewRecorder()
	withHeaders.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	rb1 := rr.Body.String()

	rr = httptest.NewRecorder()
	withHeaders.ServeHTTP(rr, req)

	rb2 := rr.Body.String()

	fmt.Printf("result 1: %s, result 2: %s", rb1, rb2)

	if rb1 != rb2 {
		t.Errorf("expecting cached result, handler returned wrong response body: got %v want %v",
			rb2, rb1)
	}

	// Now pause 10 seconds
	time.Sleep(time.Second * 12)

	rr = httptest.NewRecorder()
	withHeaders.ServeHTTP(rr, req)

	rb3 := rr.Body.String()
	fmt.Printf("result 1: %s, result 3: %s", rb1, rb3)
	if rb1 == rb3 {
		t.Errorf("expecting uncached result, handler returned wrong response body: got %v wanted to be different %v",
			rb3, rb1)
	}
}
