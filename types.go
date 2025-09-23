package iocketsdk

import (
	"encoding/json"
)

// Nanoid define um alias para o tipo de ID usado, melhorando a legibilidade.
type Nanoid string

//======================================================================
// ESTRUTURA DE ENVELOPE WEBSOCKET
//======================================================================

// Payload é o envelope genérico para todos os eventos recebidos via WebSocket.
// 'Event' contém o nome do evento (ex: "TICKET_CREATE").
// 'Data' contém o payload JSON específico do evento.
type Payload struct {
	Event string          `json:"e"`
	Data  json.RawMessage `json:"m"`
}

//======================================================================
// PAYLOADS DE REQUISIÇÃO (Dados enviados PARA o servidor)
//======================================================================

// CreateTicketRequest é a estrutura para criar um novo ticket.
type CreateTicketRequest struct {
	CategoryID Nanoid                      `json:"category_id"`
	Name       string                      `json:"name"`
	Platform   CreateTicketPlatformRequest `json:"platform"`
}

// CreateTicketPlatformRequest contém os dados da plataforma para um novo ticket.
type CreateTicketPlatformRequest struct {
	ExternalID        string          `json:"external_id"`
	Username          string          `json:"username"`
	ExtraData         json.RawMessage `json:"extra_data"`
	ChannelExternalID string          `json:"channel_external_id"`
}

// MessageBotRequest é a estrutura para enviar uma mensagem vinda de um bot.
type MessageBotRequest struct {
	ChatExternalID    string `json:"chat_external_id"`
	ClientExternalID  string `json:"client_external_id"`
	MessageExternalID string `json:"message_external_id"`
	Content           string `json:"content"`
}

//======================================================================
// PAYLOADS DE EVENTOS (Dados recebidos DO servidor via WebSocket)
//======================================================================

// MessageCreateEvent é o payload para o evento 'MESSAGE_CREATE'.
type MessageCreateEvent struct {
	TicketID       Nanoid  `json:"id"`
	Message        Message `json:"message"`
	ChatExternalID *string `json:"external_id,omitempty"`
}

// TicketClaimedEvent é o payload para o evento 'TICKET_CLAIMED'.
type TicketClaimedEvent struct {
	ID        Nanoid       `json:"id"`
	AgentName string       `json:"agent_name"`
	Ticket    Ticket       `json:"ticket"`
	Client    TicketClient `json:"client"`
}

// TicketClosedEvent é o payload para o evento 'TICKET_CLOSED'.
type TicketClosedEvent struct {
	ID        Nanoid       `json:"id"`
	AgentName string       `json:"agent_name"`
	Ticket    Ticket       `json:"ticket"`
	Client    TicketClient `json:"client"`
}

//======================================================================
// ENTIDADES PRINCIPAIS E COMPONENTES (Usados nos payloads acima)
//======================================================================

// Message representa uma única mensagem, alinhada com a estrutura do servidor.
type Message struct {
	ID        Nanoid `json:"id,omitempty"`
	From      From   `json:"from"`      // Objeto polimórfico para o remetente
	Timestamp int64  `json:"timestamp"` // Unix timestamp em segundos
	Content   string `json:"content"`
}

// From é uma estrutura especial para lidar com o remetente polimórfico (Client ou Employer).
// Apenas um dos campos será preenchido.
type From struct {
	Client   *Client
	Employer *Employer
}

// UnmarshalJSON implementa a lógica para desserializar o campo 'from' corretamente.
func (f *From) UnmarshalJSON(data []byte) error {
	// Primeiro, tentamos desserializar como Employer verificando a presença de um campo único como 'role'.
	var peek struct {
		Role string `json:"role"`
	}
	if err := json.Unmarshal(data, &peek); err == nil && peek.Role != "" {
		var emp Employer
		if err := json.Unmarshal(data, &emp); err != nil {
			return err
		}
		f.Employer = &emp
		f.Client = nil
		return nil
	}

	// Se não for um Employer, assumimos que é um Client.
	var cli Client
	if err := json.Unmarshal(data, &cli); err != nil {
		return err
	}
	f.Client = &cli
	f.Employer = nil
	return nil
}

// Client representa um cliente final nos payloads.
type Client struct {
	ID   Nanoid `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Employer representa um agente/funcionário nos payloads.
type Employer struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// Ticket representa um ticket dentro de um evento.
type Ticket struct {
	ID         Nanoid   `json:"id,omitempty"`
	ExternalID string   `json:"external_id,omitempty"`
	Category   Category `json:"category,omitempty"`
}

// TicketClient representa dados de um cliente dentro de um evento de ticket.
type TicketClient struct {
	ExternalID string          `json:"external_id"`
	ExtraData  json.RawMessage `json:"extra_data"`
	Platform   string          `json:"platform"`
}

// Category representa uma categoria. Usado tanto em tickets quanto na resposta do endpoint de categorias.
type Category struct {
	ID   Nanoid `json:"id"`
	Name string `json:"name"`
}
