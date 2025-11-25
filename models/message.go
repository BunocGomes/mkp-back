package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message representa a estrutura de uma mensagem de chat no MongoDB.
type Message struct {
    // Adicionamos as tags json:"..." para o Frontend entender
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RemetenteID    uint               `bson:"remetente_id" json:"remetente_id"`
	DestinatarioID uint               `bson:"destinatario_id" json:"destinatario_id"`
	Conteudo       string             `bson:"conteudo" json:"conteudo"`
	Timestamp      time.Time          `bson:"timestamp" json:"timestamp"`
}