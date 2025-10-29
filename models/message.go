package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message representa a estrutura de uma mensagem de chat no MongoDB.
type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`  
	RemetenteID    uint               `bson:"remetente_id"`   
	DestinatarioID uint               `bson:"destinatario_id"`  
	Conteudo       string             `bson:"conteudo"`        
	Timestamp      time.Time          `bson:"timestamp"`        
}