package dataLayer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConnection(t *testing.T) {
	client := InitDataLayer()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	usersCollection := client.Database("appointy").Collection("users")
	if _, err := usersCollection.Find(ctx, bson.D{}); err != nil {
		fmt.Println(err)
		t.Fatalf(`Basic find request was not executed without errors %v`, err)
	}

}
