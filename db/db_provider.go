package db_provider

import (
    "context"
	"fmt"
	"time"

	"github.com/jonathanhamel4/goscrape/types"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db_client *mongo.Client
var db_ctx context.Context

func verifyError(err error) {
	if err != nil {
		panic(err)
	}
}

func getCollection(collection string) mongo.Collection {
	return db_client.Database("imdb").Collection("movies")
}

func ConnectDB(connectionString string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, conErr := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	verifyError(conErr)

	pingErr := client.Ping(ctx, readpref.Primary())
	verifyError(pingErr)

	db_client = client
	db_ctx = ctx

	fmt.Println("Connected to MongoDB!")
}


func InsertMovies(movies []*types.Movie) {
	collection := getCollection("movies")

	for _, v := range movies {
		collection.InsertOne(db_ctx, v)
	}
}

// func GetMovie() {
// 	var result *types.Movie = {}
// 	collection := getCollection("movies")

// 	collection.FindOne(db_ctx, bson.M{"Title": "The Wrong Missy (2020)"}).Decode(&result)
// 	fmt.Println(result.Title)
// }