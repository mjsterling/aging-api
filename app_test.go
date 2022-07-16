package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateUser(t *testing.T) {
	postData, _ := json.Marshal(map[string]string{
		"email": "test@test.test",
	})
	body := bytes.NewBuffer(postData)
	response, err := http.Post("http://localhost:6000/users", "application/json", body)
	if response.Status != (string)(http.StatusOK) {
		t.Errorf("Create User = %v, want: %v, error: %v", response.Status, 200, err)
	}
}
