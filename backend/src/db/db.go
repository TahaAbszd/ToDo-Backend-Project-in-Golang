package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatal((".env file not found"))
	}
	mongoURI := os.Getenv("MONGO_URI")
	db_Name := os.Getenv("DB_NAME")

	if mongoURI == "" || db_Name == "" {
		log.Fatal("db credintials are missing")
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//free the memmory after finishing the function
	defer cancel()
	/*if the server is down and client try to connect the db the Connect() will be succesfull but the the error wont
	 */
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}
	//check if the db is pinging
	ping_err := client.Ping(ctx, nil)
	if ping_err != nil {
		log.Fatal("db ping failed: ", ping_err)
	}
	fmt.Println("connected to db")
	DB = client.Database(db_Name)
}
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
