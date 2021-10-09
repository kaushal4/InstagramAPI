# Instagram API

## How to run

- Install go in your device
- Install mongodb and start mongod
- Clone the repository
- Execute `go run tidy` in root folder
- Execute `go run` in root folder
- Execute `go test ./... -cover` to run unit tests with coverage

## API documentation

### Create a user

- It is a **POST** request
- JSON request body must be provided
- URL is '/users'

### Get a user using id

- This is a **GET** request
- Id is a url parameter
- URL is '/users/<id here>'

### Create a Post

- This is a **POST** request
- JSON request body must be provided
- URL is '/posts'

### Get a post using id

- It is a **GET** request
- Id is a url parameter
- URL is '/posts/<id here>'

### List all posts of a user

- It is a **GET** request
- Id is a url parameter
- Pagination Offset must be provided as JSON body
- URL should be '/posts/users/<Id here>'

## Folder Structure

- Root folder contains main.go

- User folder contains the user package with user struct and user based functions

- Post folder contains the post package with post struct and post based functions

- dataLayer folder contains dataLayer package with contains mongodb client function

- responses folder contains some standard http responses which are used in serving GET and POST request

- utility folder contains utility package which has some basic utility functions used in other packages

- CryptoPass folder contains the functions used to encrypt,decrypt and compare passwords

## API endpoints and screenshots

### Create a User

![This is an image](screenshots/apis/postUser.png)

### Get a user using id

![This is an image](screenshots/apis/postUserProof.png)

### Create a Post

![This is an image](screenshots/apis/postPost.png)

### Get post using ID

![This is an image](screenshots/apis/postProof.png)

### Get post using ID (pagination is implemented using offset=1)

![This is an image](screenshots/apis/getPostOfUser.png)

### setting a different offset=2

![This is an image](<screenshots/apis/getPostOfUser(1).png>)

## Password Encryption

![This is an image](screenshots/testing/passwordEncryption.png)

## Unit testing with overall coverage

![This is an image](screenshots/testing/all_tests.png)
