package main

import (
	"appointy/InstagramAPI/cryptoPass"
	"appointy/InstagramAPI/post"
	"appointy/InstagramAPI/user"
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {
	//setting up routes
	cryptoPass.SetSecret()
	fmt.Println("Hello World!")
	http.HandleFunc("/users/", user.GetUsersById)
	http.HandleFunc("/posts/", post.GetPostsById)
	http.HandleFunc("/posts/users/", post.GetPostsByUser)

	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
