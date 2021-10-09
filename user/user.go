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
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID       int    `json:"id" bson:"id"`
	NAME     string `json:"name" bson:"name"`
	EMAIL    string `json:"email" bson:"email"`
	PASSWORD string `json:"password" bson:"password"`
}

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

func createUser(user User) (string, error) {
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

func GetUsersById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var id string
			var err error
			idParam := r.URL.Query()["id"]
			if len(idParam) == 0 {
				responses.SetError(w, "no parameter passed")
				return
			} else {
				id = idParam[0]
			}
			var idParsed int
			if idParsed, err = utility.GetIDFromString(id); err != nil {
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
