package post

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

func deleteTestPost() (*mongo.DeleteResult, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("posts")
	result, err := usersCollection.DeleteMany(ctx, bson.M{"id": 100})
	fmt.Println(result.DeletedCount)
	return result, err
}

func TestCheckPostFields(t *testing.T) {
	test := Post{1, "", "", "", 0}
	if checkPostFields(test) == true {
		t.Fatalf(`returned true instead of false`)
	}

	test = Post{1, "v", "v", "v", 0}
	if checkPostFields(test) == false {
		t.Fatalf(`returned false instead of true`)
	}
}

func TestGetPost(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{
		"id": 1,
    "imgurl": "url",
		`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostWithWrongMethod(t *testing.T) {
	req, err := http.NewRequest("PUT", "/posts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostPost(t *testing.T) {
	var jsonStr = []byte(`{"caption":"love will come set me free","imgurl":"url2","id":100,"userid":1}`)
	req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsById)
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
	deleteTestPost()

}

func TestPostUserMissingValues(t *testing.T) {
	var jsonStr = []byte(`{"id":1}`)
	req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetPostByUserID(t *testing.T) {
	var jsonStr = []byte(`{"offset":"1"}`)
	req, err := http.NewRequest("GET", "/posts/users/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsByUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{
        "id": 2,
        "imgurl": "url1",
        "caption": "Nature has a calming quality",
        "timestamp": "2021-10-09 06:52:21.6279256 +0530 IST m=+284.468516101",
        "userid": 1
    }`
	if strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	deleteTestPost()

}

func TestGetPostByUserIDWrongMethod(t *testing.T) {
	var jsonStr = []byte(`{"offset":"1"}`)
	req, err := http.NewRequest("POST", "/posts/users/1", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPostsByUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}
