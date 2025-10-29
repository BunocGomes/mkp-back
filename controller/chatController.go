package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/BunocGomes/mkp-back/chat"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader armazena as configurações para "promover" uma conexão HTTP para WebSocket.

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWs(hub *chat.Hub, c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		log.Println("Erro: UserID não encontrado no contexto. O AuthMiddleware está sendo usado?")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Falha ao promover para websocket: %v", err)
		return
	}

	client := &chat.Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID.(uint),
		Send:   make(chan []byte, 256),
	}

	// 4. Registrar o novo cliente no Hub
	client.Hub.Register <- client

	// 5. Iniciar as "bombas" (leitores e escritores) em goroutines
	// Isso permite que o cliente leia e escreva mensagens em paralelo,
	// e libera o handler HTTP para terminar.
	go client.WritePump()
	go client.ReadPump()
}

// GetMessageHistory é um handler HTTP normal para buscar o histórico do chat.
func GetMessageHistory(c *gin.Context) {
	// 1. Pega o ID do usuário logado (do token)
	userID1, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	// 2. Pega o ID do *outro* usuário (da URL)
	userID2Str := c.Param("userID")
	userID2, err := strconv.ParseUint(userID2Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// 3. Chama o service para buscar o histórico no MongoDB
	messages, err := service.GetMessageHistory(userID1.(uint), uint(userID2))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
