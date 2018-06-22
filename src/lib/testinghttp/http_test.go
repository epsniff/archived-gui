package testinghttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	for i := 0; i < 100; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		//st = time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		//gou.Debugf("finished testcase get: %v rpc-time: %v", tc.name, time.Since(st))
		//st = time.Now()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := string(b)
		if body == "" {
			t.Fatal("no body")
		}
	}
}
