package iocketsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

//======================================================================
// ESTRUTURA PRINCIPAL DO SDK
//======================================================================

type IocketURL string

const (
	LOCAL                 IocketURL = "localhost:8080"
	PROD                  IocketURL = "api.iocket.com"
	endpointCreateTicket            = "/bot/ticket"
	endpointGetCategories           = "/bot/categories"
	endpointSendMessage             = "/ticket/message"
)

// Bot é a estrutura principal do cliente SDK.
type Bot struct {
	token      string
	baseURL    url.URL
	httpClient *http.Client
	conn       *websocket.Conn
	logger     *Logger
	ack        bool

	// Handlers de evento com segurança de tipos.
	OnConnect       func(b *Bot, channel *Channel)
	OnMessageCreate func(b *Bot, event MessageCreateEvent)
	OnTicketClaimed func(b *Bot, event TicketClaimedEvent)
	OnTicketClosed  func(b *Bot, event TicketClosedEvent)
	OnDisconnect    func(err error)

	wsURL url.URL
}

// Channel é a estrutura recebida na conexão.
type Channel struct {
	ID   Nanoid `json:"id"`
	Name string `json:"name"`
}

// New cria uma nova instância do Bot, já com o logger configurado.
func New(token string) *Bot {
	bot := &Bot{
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     NewLogger(), // Logger integrado!
	}
	bot.SetEnvironment(PROD)
	return bot
}

// SetEnvironment configura o ambiente (LOCAL ou PROD).
func (b *Bot) SetEnvironment(env IocketURL) {
	schemeRest := "https"
	if env == LOCAL {
		schemeRest = "http"
	}
	b.baseURL = url.URL{Scheme: schemeRest, Host: string(env)}
}

// Run inicia a conexão WebSocket e o loop de escuta.
func (b *Bot) Run() error {
	wsURL := b.baseURL
	wsURL.Scheme = "wss"
	if b.baseURL.Scheme == "http" {
		wsURL.Scheme = "ws"
	}
	wsURL.Path = "/gateway"
	wsURL.RawQuery = "token=" + b.token

	b.logger.Info("Connecting to", wsURL.String())
	b.wsURL = wsURL
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		b.logger.Error("Failed to connect:", err)
		return fmt.Errorf("failed to connect: %w", err)
	}
	b.conn = c

	c.SetPingHandler(func(appData string) error {
		deadline := time.Now().Add(10 * time.Second)
		return c.WriteControl(websocket.PingMessage, []byte(appData), deadline)
	})

	var channelInfo Channel
	if err := c.ReadJSON(&channelInfo); err != nil {
		b.logger.Error("Failed to read channel info:", err)
		return fmt.Errorf("failed to read channel info: %w", err)
	}

	if b.OnConnect != nil {
		go b.OnConnect(b, &channelInfo)
	}

	b.logger.Info("Connected to channel:", channelInfo.Name)

	go b.listen()

	return nil
}

// listen é o loop principal que lê mensagens do WebSocket.
func (b *Bot) listen() {
	defer b.conn.Close()
	for {
		_, data, err := b.conn.ReadMessage()
		if err != nil {
			b.logger.Error("Connection read error:", err)
			if b.OnDisconnect != nil {
				go b.OnDisconnect(err)
			}
			b.logger.Warn("Trying to reconnect")
			go b.reconnect()
			return
		}
		b.dispatch(data)
	}
}

func (b *Bot) reconnect() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	times := 0
	maxTimes := 30
	for range ticker.C {
		if times == maxTimes {
			return
		}

		times += 1
		b.logger.Warn("Attempt", times)
		if err := b.Run(); err != nil && maxTimes > times {
			continue
		}
		return
	}
}

// --- Métodos de API REST ---

func (b *Bot) SendMessage(req MessageBotRequest) error {
	resp, err := b.post(endpointSendMessage, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return b.newAPIError(resp)
	}
	b.logger.Info("Message sent successfully:", req.MessageExternalID)
	return nil
}

func (b *Bot) CreateTicket(req CreateTicketRequest) (*Ticket, error) {
	var ticket Ticket
	resp, err := b.post(endpointCreateTicket, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, b.newAPIError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&ticket); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &ticket, nil
}

func (b *Bot) GetCategories() ([]Category, error) {
	var categories []Category
	resp, err := b.get(endpointGetCategories)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, b.newAPIError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return categories, nil
}

// --- Funções de Ajuda Internas ---

func (b *Bot) post(endpoint string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		b.logger.Error("Failed to marshal POST payload:", err)
		return nil, err
	}

	url := b.baseURL.String() + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+b.token)

	return b.httpClient.Do(req)
}

func (b *Bot) get(endpoint string) (*http.Response, error) {
	url := b.baseURL.String() + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bot "+b.token)

	return b.httpClient.Do(req)
}

func (b *Bot) dispatch(data []byte) {
	var p Payload
	if err := json.Unmarshal(data, &p); err != nil {
		b.logger.Warn("Failed to unmarshal websocket payload:", err)
		return
	}

	switch p.Event {
	case "MESSAGE_CREATE":
		if b.OnMessageCreate != nil {
			var event MessageCreateEvent
			if err := json.Unmarshal(p.Data, &event); err == nil {
				go b.OnMessageCreate(b, event)
			} else {
				b.logger.Warn("Failed to unmarshal MESSAGE_CREATE data:", err)
			}
		}
	case "TICKET_CLAIMED":
		if b.OnTicketClaimed != nil {
			var event TicketClaimedEvent
			if err := json.Unmarshal(p.Data, &event); err == nil {
				go b.OnTicketClaimed(b, event)
			} else {
				b.logger.Warn("Failed to unmarshal TICKET_CLAIMED data:", err)
			}
		}
	case "TICKET_CLOSED":
		if b.OnTicketClosed != nil {
			var event TicketClosedEvent
			if err := json.Unmarshal(p.Data, &event); err == nil {
				go b.OnTicketClosed(b, event)
			} else {
				b.logger.Warn("Failed to unmarshal TICKET_CLOSED data:", err)
			}
		}
	case "HEARTBEAT_ACK":
		b.logger.Info("Pong")
	default:
		b.logger.Warn("Received unknown event type:", p.Event)
	}
}

func (b *Bot) newAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	err := fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	b.logger.Error(err.Error()) // Loga o erro da API antes de retorná-lo.
	return err
}
