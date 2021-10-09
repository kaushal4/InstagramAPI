package user

import (
	"appointy/InstagramAPI/dataLayer"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func deleteTestUser() (*mongo.DeleteResult, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("users")
	result, err := usersCollection.DeleteMany(ctx, bson.M{"id": 100})
	fmt.Println(result.DeletedCount)
	return result, err
}

func TestCheckUserFields(t *testing.T) {
	test := User{1, "", "", ""}
	if checkUserFields(test) == true {
		t.Fatalf(`returned true instead of false`)
	}

	test = User{1, "v", "v", "v"}
	if checkUserFields(test) == false {
		t.Fatalf(`returned false instead of true`)
	}
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{
		"id": 1,
		"name": "kaushal",
		"email": "kdv@gmail.com",
		`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserWithWrongMethod(t *testing.T) {
	req, err := http.NewRequest("PUT", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostUser(t *testing.T) {
	var jsonStr = []byte(`{"name":"kaushal","id":100,"password":"test","email":"dks@gmail.com"}`)
	req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "operation successfull"
	if strings.Compare(rr.Body.String(), expected) != 0 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	deleteTestUser()

}

func TestPostUserMissingValues(t *testing.T) {
	var jsonStr = []byte(`{"name":"kaushal"}`)
	req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUsersById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
