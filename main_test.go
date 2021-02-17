package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandleFunc_POST_Success(t *testing.T) {
	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("first_name", "John")
	data.Set("last_name", "Doe")
	data.Set("email", "email@example.com")
	data.Set("phone_number", "0819999999")
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handleFunc(w, req)
	//FIXED:change test to follow REST principles
	if w.Code != http.StatusCreated {
		t.Fail()
		t.Logf("Expected status code to be 201")
	}
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
}

//FIXED:add test to check bad requests
func TestHandleFunc_POST_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	data := url.Values{}
	data.Set("first_name", "John")
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handleFunc(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fail()
		t.Logf("Expected status code to be 400")
	}
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
}

func TestHandleFunc_GET_Success(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	handleFunc(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHandleFunc_GET_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/", nil)
	handleFunc(w, req)
	if w.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	oLoadEnv := loadEnv
	loadEnv = func(filename ...string) (err error) {
		os.Setenv("PORT", "8080")
		return
	}
	defer func() {
		loadEnv = oLoadEnv
		r := recover()
		if r != nil {
			t.Fail()
		}
	}()
	srv := run()
	//FIXED:add an assert to check there are handlers for request
	if srv.Handler == nil {
		t.Fail()
		t.Log("server does not have any handler attached")
	}
	time.Sleep(1 * time.Second)
	srv.Shutdown(context.TODO())
}
