package iocketsdk

import "encoding/json"

type Payload struct {
	E string          `json:"e"`
	M json.RawMessage `json:"m"`
}

type Message struct {
	ChatID    string `json:"chat_external_id"`
	ClientID  string `json:"client_external_id"`
	MessageID string `json:"message_external_id"`
	Content   string `json:"content"`
}

type MessageCreate struct {
	From      string `json:"from"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
	ChatID    string `json:"external_id"`
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
	ExternalID string `json:"external_id"`
	CategoryID string `json:"category_id"`
}

type TicketClient struct {
	ExternalID string          `json:"external_id"`
	ExtraData  json.RawMessage `json:"extra_data"`
	Platform   string          `json:"platform"`
}

type TicketClose struct {
	ExternalID       string `json:"external_id"`
	ClientExternalID string `json:"client_external_id"`
}

