package dto

// trocadas através do WebSocket.
type WebSocketMessageDTO struct {
	Tipo           string `json:"tipo"`           // "nova_mensagem", "erro", "notificacao"
	DestinatarioID uint   `json:"destinatario_id"` 
	Conteudo       string `json:"conteudo"`       
}