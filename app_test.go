package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	postData, _ := json.Marshal(map[string]string{
		"email":    "test@test.test",
		"password": "123456",
	})
	body := bytes.NewBuffer(postData)
	response, err := http.Post("http://localhost:6000/user", "application/json", body)
	if response.Status != "201 Created" {
		t.Errorf("Create User: response: %v, want: %v, error: %v", response.Status, "201 Created", err)
	}
}

func TestLogin(t *testing.T) {
	postData, _ := json.Marshal(map[string]string{
		"email":    "test@test.test",
		"password": "123456",
	})
	body := bytes.NewBuffer(postData)
	response, err := http.Post("http://localhost:6000/login", "application/json", body)
	if err != nil {
		t.Errorf("Login: error: %v", err)
	}
	defer response.Body.Close()
	responseBody, _ := io.ReadAll(response.Body)
	want := "200 OK"
	if response.Status != want {
		t.Errorf("Login User: response: %v, want: %v, error: %v", response.Status, want, (string)(responseBody))
	}
}
