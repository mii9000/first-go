package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

type inputs []formInput

func postArrange(values ...string) *http.Request {
	data := url.Values{}
	data.Set("first_name", values[0])
	data.Set("last_name", values[1])
	data.Set("email", values[2])
	data.Set("phone_number", values[3])
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestHandleFunc_POST_Success(t *testing.T) {
	w := httptest.NewRecorder()
	req := postArrange("John", "Doe", "john@email.com", "096675432")
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
	//arrange
	w := httptest.NewRecorder()
	req := postArrange("", "", "jane@email.com", "")

	//act
	handleFunc(w, req)

	//assert
	if w.Code != http.StatusBadRequest {
		t.Fail()
		t.Logf("Expected status code to be 400")
	}

	//cleanup
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
}

//FIXED:adds proper asserts
func TestHandleFunc_GET_Success(t *testing.T) {
	//arrange
	handleFunc(httptest.NewRecorder(), postArrange("Get", "User", "user@email.com", "123456789"))
	w := httptest.NewRecorder()

	//act
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	handleFunc(w, req)
	inputs := make([]formInput, 0)
	err := json.NewDecoder(w.Body).Decode(&inputs)

	//assert
	if w.Code != http.StatusOK {
		t.Fail()
		t.Logf("Expected status code to be 200")
	}
	if err != nil {
		t.Fail()
		t.Logf("could not parse to list of form inputs")
	}
	if len(inputs) != 1 {
		t.Fail()
		t.Logf("expected 1 record")
	}

	//cleanup
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
}

//FIXED:adds test to check data gets appended and not replaced
func TestHandleFunc_POST_Append_Success(t *testing.T) {
	//arrange
	handleFunc(httptest.NewRecorder(), postArrange("User", "One", "user1@email.com", "123456789"))
	handleFunc(httptest.NewRecorder(), postArrange("Two", "User", "user2@email.com", "987654321"))
	w := httptest.NewRecorder()

	//act
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	handleFunc(w, req)
	inputs := make([]formInput, 0)
	json.NewDecoder(w.Body).Decode(&inputs)

	//assert
	if len(inputs) != 2 {
		t.Fail()
		t.Logf("expected 2 records")
	}
	if inputs[0].Email != "user1@email.com" {
		t.Fail()
		t.Logf("expected first user email to be 'user1@rmail.com'")
	}
	if inputs[1].Email != "user2@email.com" {
		t.Fail()
		t.Logf("expected second user email to be 'user1@rmail.com'")
	}

	//cleanup
	ioutil.WriteFile(dataFile, []byte("[]"), os.ModeAppend)
}

//FIXED:rename test name according to intention of code block
func TestHandleFunc_PUT_NotFound(t *testing.T) {
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
