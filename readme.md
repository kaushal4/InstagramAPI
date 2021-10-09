# Instagram API

## How to run

- Install go in your device
- Install mongodb and start mongod
- Clone the repository
- Execute `go run tidy` in root folder
- Execute `go run` in root folder
- Execute `go test ./... -cover` to run unit tests with coverage

## API documentation

#### create a user

- Should be a POST request
- Use JSON request body
- URL should be '/users'

#### Get a user using id

- This is a **GET** request
- Id is a url parameter
- URL is '/users/<id here>'

#### Create a Post

- This is a **POST** request
- JSON request body must be provided
- URL is '/posts'

#### Get a post using id

- It is a **GET** request
- Id is a url parameter
- URL is '/posts/<id here>'

#### List all posts of a user

- Should be a **GET** request
- Id is a url parameter
- Pagination Offset must be provided as JSON body
- URL should be '/posts/users/<Id here>'

## API endpoints and screenshots

- Create a User
  ![This is an image](screenshots/apis/postUser.png)
- Get a user using id
  ![This is an image](screenshots/apis/postUserProof.png)
- Create a Post
  ![This is an image](screenshots/apis/postPost.png)
- Get post using ID
  ![This is an image](screenshots/apis/postProof.png)
- Get post using ID (pagination is implemented using offset=1)
  ![This is an image](screenshots/apis/getPostOfUser.png)
- setting a different offset=2
  ![This is an image](<screenshots/apis/getPostOfUser(1).png>)

## Password Encryption

![This is an image](screenshots/testing/passwordEncryption.png)

## Unit testing with overall coverage

![This is an image](screenshots/testing/all_tests.png)
