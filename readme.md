# Instagram API

## how to run

- install go in your device
- install mongodb and start mongod
- clone the repository
- execute `go run tidy` in root folder
- execute `go run` in root folder

## API endpoints and screenshots

- Create a User
  ![This is an image](screenshots/apis/postUser.png){:height="50%" width="50%"}
- Get a user using id
  ![This is an image](screenshots/apis/postUserProof.png)
- Create a Post
  ![This is an image](screenshots/apis/postPost.png | width=250)
- Get post using ID
  ![This is an image](screenshots/apis/postProof.png | width=250)
- Get post using ID (pagination is implemented using offset=1)
  ![This is an image](screenshots/apis/getPostOfUser.png | width=250)
- setting a different offset=2
  ![This is an image](screenshots/apis/getPostOfUser(1).png | width=250)
