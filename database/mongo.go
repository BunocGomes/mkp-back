package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabase *mongo.Database

func InitMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("A variável de ambiente MONGO_URI não foi definida.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Falha ao pingar o MongoDB: %v", err)
	}

	MongoDatabase = client.Database("mkp_chat_db")

	log.Println("Conectado ao MongoDB com sucesso!")
}