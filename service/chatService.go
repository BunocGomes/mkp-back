package service

import (
	"context"
	"time"

	"github.com/BunocGomes/mkp-back/database"
	"github.com/BunocGomes/mkp-back/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveMessage(remetenteID, destinatarioID uint, conteudo string) (*models.Message, error) {
	collection := database.MongoDatabase.Collection("messages")

	message := &models.Message{
		RemetenteID:    remetenteID,
		DestinatarioID: destinatarioID,
		Conteudo:       conteudo,
		Timestamp:      time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func GetMessageHistory(userID1, userID2 uint) ([]models.Message, error) {
	collection := database.MongoDatabase.Collection("messages")

	filter := bson.M{
		"$or": []bson.M{
			{"remetente_id": userID1, "destinatario_id": userID2},
			{"remetente_id": userID2, "destinatario_id": userID1},
		},
	}

	// Ordena pela mais antiga primeiro (timestamp: 1)
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var messages []models.Message
	if err = cursor.All(context.TODO(), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
