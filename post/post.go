package post

import (
	"appointy/InstagramAPI/dataLayer"
	"appointy/InstagramAPI/responses"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Post struct {
	ID       int    `json:"id"`
	IMGURL   string `json:"imgurl"`
	CAPTION  string `json:"caption"`
	TMESTAMP string `json:"timestamp"`
	USERID   int    `json:"userid"`
}

func checkPostFields(post Post) bool {
	v := reflect.ValueOf(post)
	for i := 0; i < v.NumField(); i++ {
		temp := v.Field(i).Interface()
		if temp == nil || temp == "" {
			return false
		}
	}
	return true
}

func createPost(post Post) (string, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("posts")

	if _, err := usersCollection.InsertOne(ctx, post); err == nil {
		return "success", nil
	} else {
		return "", err
	}
}

func findPost(id int) (Post, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	postsCollection := client.Database("appointy").Collection("posts")
	var result Post
	if err := postsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&result); err == nil {
		return result, nil
	} else {
		return Post{}, errors.New("failed to find post with provided ID")
	}
}

func GetPostsById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var id string
			idParam := r.URL.Query()["id"]
			if len(idParam) == 0 {
				responses.SetError(w, "no parameters passed")
				return
			} else {
				id = idParam[0]
			}
			var idParsed int
			if i, err := strconv.Atoi(id); err == nil {
				idParsed = i
			} else {
				responses.SetError(w, "Invalid ID")
				return
			}
			if result, err := findPost(idParsed); err == nil {
				//var message string
				if postJson, err := json.MarshalIndent(result, "", "   "); err == nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(postJson)
					return
				} else {
					responses.SetError(w, err.Error())
				}
			} else {
				responses.SetError(w, "Could not fetch post :(, The post might not exist")
			}
		}
	case "POST":
		{

			decoder := json.NewDecoder(r.Body)
			var post Post
			err := decoder.Decode(&post)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			post.TMESTAMP = time.Now().String()
			fmt.Println(post)
			if !checkPostFields(post) {
				responses.SetError(w, "Request body missing fields.")
				return
			}
			if _, err := createPost(post); err == nil {
				responses.SetResponse(w)
				w.Write([]byte("operation successfull"))
				return
			}
			responses.SetError(w, "could not post :(")
		}
	default:
		{
			responses.SetError(w, "Only get and post requests allowed")
		}
	}
}

func fetchPostsByUser(id int) ([]Post, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	postsCollection := client.Database("appointy").Collection("posts")
	var posts []Post
	cur, err := postsCollection.Find(ctx, bson.M{"userid": id})
	if err != nil {
		return make([]Post, 0), errors.New("error fetching posts")
	}
	for cur.Next(ctx) {
		var post Post
		err := cur.Decode(&post)
		if err != nil {
			return make([]Post, 0), errors.New("error decoding posts")
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var user struct {
				USERID int `json:"userid"`
			}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&user)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			if posts, err := fetchPostsByUser(user.USERID); err == nil {
				if len(posts) == 0 {
					responses.SetError(w, "could not fetch posts :(, The user might not have any posts")
					return
				}
				postsJson, _ := json.MarshalIndent(posts, "", "   ")
				responses.SetResponse(w)
				w.Write(postsJson)
				return
			} else {
				responses.SetError(w, err.Error())
			}
		}
	default:
		{
			responses.SetError(w, "Only GET requests allowed!")
		}
	}
}
