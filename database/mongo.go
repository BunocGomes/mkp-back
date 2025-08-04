package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongo() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("A variável de ambiente MONGO_URI não foi definida.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Falha ao conectar ao MongoDB: ", err)
	}
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Não foi possível pingar o servidor MongoDB: ", err)
	}
	fmt.Println("Conexão com o MongoDB estabelecida com sucesso.")
}