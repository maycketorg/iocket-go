package iocketsdk

import "encoding/json"

type Payload struct {
	E string          `json:"e"`
	M json.RawMessage `json:"m"`
}

type MessagePayload struct {
	ID             string `json:"id"`                    // TicketID
	Message        *Message        `json:"message"`               // Message
	Status         bool            `json:"status,omitempty"`      // Status
	ChatExternalID *string         `json:"external_id,omitempty"` // ChatID External
}

type Message struct {
	ID        string `json:"id,omitempty"` // MessageID
	From      string          `json:"from"`         // From Name
	Timestamp int64           `json:"timestamp"`    // Timestamp
	Content   string          `json:"content"`      // Content
}
type CreateTicket struct {
	CategoryID string               `json:"category_id"`
	Name       string               `json:"name"`
	Platform   CreateTicketPlatform `json:"platform"`
}

type CreateTicketPlatform struct {
	ExternalID        string          `json:"external_id"`
	Username          string          `json:"username"`
	ExtraData         json.RawMessage `json:"extra_data"`
	ChannelExternalID string          `json:"channel_external_id"`
}

type Channel struct {
	ID         string     `json:"id"`
	OrgID      string     `json:"org_id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClaimTicket struct {
	ID           string        `json:"id"`
	AgentName    string        `json:"agent_name"`
	Ticket       *Ticket       `json:"ticket"`
	TicketClient *TicketClient `json:"client"`
}

type Ticket struct {
	ID         string    `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Category   *Category `json:"category,omitempty"`
	CategoryID string    `json:"category_id,omitempty"`
}

type TicketClient struct {
	ExternalID string          `json:"external_id"`
	ExtraData  json.RawMessage `json:"extra_data"`
	Platform   string          `json:"platform"`
}

type TicketClose struct {
	ID           string        `json:"id"`
	AgentName    string        `json:"agent_name"`
	Ticket       *Ticket       `json:"ticket"`
	TicketClient *TicketClient `json:"client"`
}
