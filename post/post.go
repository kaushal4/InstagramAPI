package post

import (
	"appointy/InstagramAPI/dataLayer"
	"appointy/InstagramAPI/responses"
	"appointy/InstagramAPI/utility"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
The below struct is used to encode the post request json and also to push document to mongodb collection
timestamp is calulated using time.Now() which corresponds with the time the request was fullfilled
*/
type Post struct {
	ID       int    `json:"id"`
	IMGURL   string `json:"imgurl"`
	CAPTION  string `json:"caption"`
	TMESTAMP string `json:"timestamp"`
	USERID   int    `json:"userid"`
}

/* Mutex locks are applied to post requests so that we can aviod threads
reading outdated data or multiple threads trying to add the same data */
var lock sync.Mutex

/*
The below function is used to make sure that the post request contains the adequte fields
before quering it to database
*/
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

/*
The below fuction is used to add post document to the database
*/
func createPost(post Post) (string, error) {
	lock.Lock()
	defer lock.Unlock()
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

/*
The below function will find a post corresponding to a particular post id
*/
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

/*
The below function is the primary function used to serve requests
It checks if the method is post or get
for get requests it extracts the user id from the url and calls the findPost() function
for post requests it extracts post data from the request body and calls the createPost() function
*/
func GetPostsById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			idParsed, err := utility.ExtractID(r.URL.Path, "/posts/")
			if err != nil {
				responses.SetError(w, "Invalid ID")
				return
			}
			if result, err := findPost(idParsed); err == nil {
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

/*
The below function is used to fetch posts using user id
*/
func fetchPostsByUser(id int, offset int) ([]Post, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	postsCollection := client.Database("appointy").Collection("posts")
	var posts []Post

	opts := options.Find()
	opts.SetSort(bson.D{{"id", 1}})
	opts.SetSkip(int64(offset))
	opts.SetLimit(2)
	cur, err := postsCollection.Find(ctx, bson.M{"userid": id}, opts)
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

/*
The below function handles get request made to /posts/users/<userid>
This checks for offset value in request body and the user id value in the url
It then calls the fetchPostByUser function
*/

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var user struct {
				OFFSET string `json:"offset"`
			}
			var offset int
			var id int
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&user)
			if err != nil {
				responses.SetError(w, err.Error())
				return
			}
			if user.OFFSET == "" {
				offset = 0
			} else if offset, err = utility.GetIDFromString(user.OFFSET); err != nil {
				responses.SetError(w, err.Error())
				return
			}
			if id, err = utility.ExtractID(r.URL.Path, "/posts/users/"); err != nil {
				responses.SetError(w, "Invalid ID")
				return
			}
			if posts, err := fetchPostsByUser(id, offset); err == nil {
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
