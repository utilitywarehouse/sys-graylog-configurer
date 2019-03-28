package graylog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCallAPI(t *testing.T) {
	t.Run("detects invalid request", func(t *testing.T) {
		cl := NewClient("", "")

		err := cl.callAPI("//", "", nil, nil)
		if err == nil {
			t.Error("invalid request not detected")
		}
	})
	t.Run("uses client url", func(t *testing.T) {
		hitServer := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hitServer = true

		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI("", "", nil, nil)
		if hitServer == false {
			t.Error("call is not using client URL")
		}
	})
	t.Run("proper http method", func(t *testing.T) {
		method := "METHOD"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				t.Error("http method not sent")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI(method, "", nil, nil)
	})
	t.Run("proper path", func(t *testing.T) {
		path := "/path"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() != path {
				t.Error("http path not sent")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI("", path, nil, nil)
	})
	t.Run("proper headers", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Header.Get("bob"))
			if r.Header.Get("Accept") != "application/json" {
				t.Error("'Accept' header not sent")
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("'Content-Type' header not sent")
			}
			if r.Header.Get("X-Requested-By") != "graylog-configurer" {
				t.Error("'X-Requested-By' header not sent")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI("", "", nil, nil)
	})
	t.Run("detects unauthorized", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		err := cl.callAPI("", "", nil, nil)
		if !strings.Contains(err.Error(), "401 Unauthorized") {
			t.Error("unauthorized request not detected")
		}

	})
	t.Run("basic auth", func(t *testing.T) {
		password := "password"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, p, _ := r.BasicAuth()
			if p != password {
				t.Error("basic auth password not sent")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, password)

		cl.callAPI("", "", nil, nil)
	})
	t.Run("ignore output if not requested", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `notJson`)
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		err := cl.callAPI("", "", nil, nil)
		if err != nil {
			t.Error("output parsed without need")
		}
	})
	t.Run("detects invalid json reply", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "notJson")
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		var out interface{}
		err := cl.callAPI("", "", nil, &out)
		if !strings.Contains(err.Error(), "decoding body") {
			t.Error("invalid json reply not detected")
		}
	})
	t.Run("returns output", func(t *testing.T) {
		type Testable struct {
			Test string
		}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"test":"text"}`)
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		var out Testable
		err := cl.callAPI("", "", nil, &out)
		if out.Test != "text" || err != nil {
			t.Error("output not recovered")
		}
	})
	t.Run("detect unmarshable input", func(t *testing.T) {
		cl := NewClient("", "")

		err := cl.callAPI("", "", make(chan int), nil)
		if err == nil {
			t.Error("non json payload not detected")
		}
	})
	t.Run("posting json", func(t *testing.T) {
		type Testable struct {
			Test string
		}
		payload := Testable{Test: "testable"}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := new(Testable)
			json.NewDecoder(r.Body).Decode(value)
			if value.Test != payload.Test {
				t.Error("json payload not sent")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI("", "", payload, nil)
	})
	t.Run("posting json respects tags", func(t *testing.T) {
		type Testable struct {
			Test string `json:"test_tag"`
		}
		payload := Testable{Test: "testable"}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := ioutil.ReadAll(r.Body)
			if !strings.Contains(string(body), "test_tag") {
				t.Error("json tag not respected")
			}
		}))
		defer ts.Close()
		cl := NewClient(ts.URL, "")

		cl.callAPI("", "", payload, nil)
	})
}
