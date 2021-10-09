package user

import (
	"appointy/InstagramAPI/cryptoPass"
	"appointy/InstagramAPI/dataLayer"
	"appointy/InstagramAPI/responses"
	"appointy/InstagramAPI/utility"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/* Mutex locks are applied to post requests so that we can aviod threads
reading outdated data or multiple threads trying to add the same data */
var lock sync.Mutex

/*
the User struct defines the structure of user
This will be used to send data to the database and read from it
password -> string contains the hash of the password sent by the user
*/
type User struct {
	ID       int    `json:"id" bson:"id"`
	NAME     string `json:"name" bson:"name"`
	EMAIL    string `json:"email" bson:"email"`
	PASSWORD string `json:"password" bson:"password"`
}

/*
The below function is used to make sure that the post request contains the adequte fields
before quering it to database
*/
func checkUserFields(user User) bool {
	v := reflect.ValueOf(user)
	for i := 0; i < v.NumField(); i++ {
		temp := v.Field(i).Interface()
		if temp == nil || temp == "" {
			return false
		}
	}
	return true
}

/*
The below function is used to find the user from the database using the user id
*/
func findUser(id int) (User, error) {
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("users")
	var result User
	if err := usersCollection.FindOne(ctx, bson.M{"id": id}).Decode(&result); err == nil {
		return result, nil
	} else {
		return User{}, errors.New("failed to find user with provided ID")
	}
}

/*
The below function is used to create a user in the database
The object recieved from the post request is checked and then added to the database
*/

func createUser(user User) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	client := dataLayer.InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("users")
	user.PASSWORD = string(cryptoPass.Encrypt([]byte(user.PASSWORD), os.Getenv("HASH_KEY")))
	if _, err := usersCollection.InsertOne(ctx, user); err == nil {
		return "success", nil
	} else {
		return "", err
	}
}

/*
The below function is the primary function used to serve requests
It checks if the method is post or get
for get requests it extracts the user id from the url and calls the findUser() function
for post requests it extracts post data from the request body and calls the createUser() function
*/
func GetUsersById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var err error
			var idParsed int
			if idParsed, err = utility.ExtractID(r.URL.Path, "/users/"); err != nil {
				responses.SetError(w, "Invalid ID")
				return
			}
			if result, err := findUser(idParsed); err == nil {
				if userJson, err := json.MarshalIndent(result, "", "   "); err == nil {
					responses.SetResponse(w)
					w.Write(userJson)
					return
				} else {
					responses.SetError(w, err.Error())
				}

			} else {
				responses.SetError(w, "Could not fetch user :(, The user might not exist")
			}
		}
	case "POST":
		{

			decoder := json.NewDecoder(r.Body)
			var user User
			err := decoder.Decode(&user)
			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}
			fmt.Println(user)
			if !checkUserFields(user) {
				responses.SetError(w, "Request body missing fields.")
				return
			}
			if _, err := createUser(user); err == nil {
				responses.SetResponse(w)
				w.Write([]byte("operation successfull"))
				return
			}
			responses.SetError(w, err.Error())

		}
	default:
		{
			responses.SetError(w, "only post and get requests allowed!")
		}
	}
}
