package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/BunocGomes/mkp-back/dto"
	"github.com/BunocGomes/mkp-back/models"
	"github.com/BunocGomes/mkp-back/service"
	"github.com/gorilla/websocket"
)

// ... (Constantes continuam iguais) ...
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// --- MUDANÇA: Campos capitalizados ---
// Client é a ponte entre a conexão websocket e o hub.
type Client struct {
	Hub *Hub // <--- Capitalizado
	// A conexão websocket.
	Conn *websocket.Conn // <--- Capitalizado
	// O ID do usuário (vindo do nosso banco de dados, extraído do JWT).
	UserID uint // <--- Capitalizado
	// Um canal de buffer para mensagens de saída.
	Send chan []byte // <--- Capitalizado
}

// --- MUDANÇA: Campos capitalizados ---
// Hub mantém o conjunto de clientes ativos e transmite mensagens para eles.
type Hub struct {
	// clients (minúsculo) - Ninguém de fora precisa acessar o map diretamente.
	clients map[uint]*Client

	// Mensagens de entrada dos clientes.
	Broadcast chan *models.Message // <--- Capitalizado

	// Canal para registrar clientes.
	Register chan *Client // <--- Capitalizado

	// Canal para remover o registro de clientes.
	Unregister chan *Client // <--- Capitalizado
}

// --- MUDANÇA: Método capitalizado (ReadPump) e referências internas atualizadas ---
// ReadPump bombeia mensagens do websocket para o hub (O que o usuário DIGITA).
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c // <--- Usa campos capitalizados
		c.Conn.Close()        // <--- Usa campos capitalizados
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msgDTO dto.WebSocketMessageDTO
		if err := json.Unmarshal(messageBytes, &msgDTO); err != nil {
			log.Printf("erro ao decodificar json do websocket: %v", err)
			continue
		}

		// c.UserID (capitalizado)
		savedMessage, err := service.SaveMessage(c.UserID, msgDTO.DestinatarioID, msgDTO.Conteudo)
		if err != nil {
			log.Printf("erro ao salvar mensagem no mongo: %v", err)
			continue
		}

		// c.Hub.Broadcast (capitalizado)
		c.Hub.Broadcast <- savedMessage
	}
}

// --- MUDANÇA: Método capitalizado (WritePump) e referências internas atualizadas ---
// WritePump bombeia mensagens do hub para o websocket (O que o usuário RECEBE).
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close() // <--- Usa campos capitalizados
	}()

	for {
		select {
		// c.Send (capitalizado)
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// NewHub cria um novo Hub.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *models.Message), // <--- Capitalizado
		Register:   make(chan *Client),         // <--- Capitalizado
		Unregister: make(chan *Client),         // <--- Capitalizado
		clients:    make(map[uint]*Client),     // (continua minúsculo, é interno)
	}
}

// Run inicia o Hub. Esta função deve ser executada como uma goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client.UserID] = client
			log.Printf("🟢 CONECTOU: UserID %d (Total online: %d)", client.UserID, len(h.clients))

		case client := <-h.Unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				log.Printf("🔴 DESCONECTOU: UserID %d (Total online: %d)", client.UserID, len(h.clients))
			}

		case message := <-h.Broadcast:
			// --- BLOCO DE DEBUG NOVO ---
			onlineIDs := []uint{}
			for id := range h.clients {
				onlineIDs = append(onlineIDs, id)
			}
			log.Printf("🔍 DEBUG: Tentando enviar de %d para %d", message.RemetenteID, message.DestinatarioID)
			log.Printf("📋 LISTA DE ONLINE: %v", onlineIDs)
			// ---------------------------

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("❌ Erro JSON: %v", err)
				continue
			}

			// 1. Envia para o DESTINATÁRIO
			if client, ok := h.clients[message.DestinatarioID]; ok {
				select {
				case client.Send <- messageBytes:
					log.Printf("✅ ENTREGUE para %d", message.DestinatarioID)
				default:
					close(client.Send)
					delete(h.clients, message.DestinatarioID)
				}
			} else {
				log.Printf("⚠️ FALHA: %d não está no mapa %v", message.DestinatarioID, onlineIDs)
			}

			// 2. Envia para o REMETENTE
			if client, ok := h.clients[message.RemetenteID]; ok {
				select {
				case client.Send <- messageBytes:
					// Sucesso silencioso para o remetente
				default:
					close(client.Send)
					delete(h.clients, message.RemetenteID)
				}
			}
		}
	}
}
