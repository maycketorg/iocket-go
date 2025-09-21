package iocketsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
)

type IocketURL string

const (
	LOCAL IocketURL = "localhost:8080"
	URL   IocketURL = "api.iocket.com"

	RAW       = "https://api.iocket.com"
	RAW_LOCAL = "http://localhost:8080"

	MESSAGE        = "/ticket/message"
	CREATE_TICKET  = "/bot/ticket"
	GET_CATEGORIES = "/bot/categories"
)

type Bot struct {
	token    string
	Channel  Channel
	handlers map[reflect.Type][]reflect.Value
	route    IocketURL
}

func New(token string) *Bot {
	return &Bot{
		token:    token,
		handlers: make(map[reflect.Type][]reflect.Value),
		route:    URL,
	}
}

func (b *Bot) Set(url IocketURL) {
	b.route = url
}

func (b *Bot) Run() error {
	P("Starting Bot")
	
	var i string
	if b.route == LOCAL {
		i = "ws://"
	} else {
		i = "wss://"
	}
	c, _, err := websocket.DefaultDialer.Dial(i + string(b.route) + "/gateway" +"?token="+b.token, nil)
	if err != nil {
		return err
	}

	P("Getting channel informations")

	for {
		if c == nil {
			return errors.New("conection is closed")
		}
		_, data, err := c.ReadMessage()
		if err != nil {
			Perror(err)
			continue
		}
		var ch Channel
		if err := json.Unmarshal(data, &ch); err != nil {
			return err
		}
		b.Channel = ch
		break
	}

	P("Hello", b.Channel.Name)

	go func() {
		for {
			if c == nil {
				return
			}
			_, data, err := c.ReadMessage()
			if err != nil {
				Perror(err)
				continue
			}

			P(string(data))
			b.trigger(data)
		}
	}()

	return nil
}

func (b *Bot) Add(events ...interface{}) {
	for _, event := range events {
		v := reflect.TypeOf(event)
		if v.Kind() != reflect.Func {
			Perror("is not possible add other type")
			return
		}

		if v.NumIn() != 2 {
			Perror("invalid handler")
			return
		}

		param := v.In(1)
		b.handlers[param] = append(b.handlers[param], reflect.ValueOf(event))
	}
}

func (b *Bot) POST(m interface{}, ep string) (*http.Response, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", string(b.route)+ep, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bot "+b.token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Bot) GET(ep string) (*http.Response, error) {
	req, err := http.NewRequest("GET", string(b.route)+ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bot "+b.token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Bot) Send(m Message) error {
	r, err := b.POST(m, MESSAGE)
	if err != nil {
		return err
	}

	if r.StatusCode != 201 {
		return errors.New("invalid message")
	}

	return nil
}

func (b *Bot) CreateTicket(ct CreateTicket) error {
	r, err := b.POST(ct, CREATE_TICKET)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return errors.New("invalid to create ticket")
	}

	return nil
}

func (b *Bot) GetCategories() ([]Category, error) {
	r, err := b.GET(GET_CATEGORIES)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, errors.New("error for get categories")
	}
	defer r.Body.Close()
	var categories []Category
	if err := json.NewDecoder(r.Body).Decode(&categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (b *Bot) trigger(data []byte) error {
	var p Payload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	var m interface{}
	switch p.E {
	case "MESSAGE_CREATE":
		var mc MessageCreate
		if err := json.Unmarshal(p.M, &mc); err != nil {
			return err
		}
		m = mc
	case "CLAIM_TICKET":
		var ct ClaimTicket
		if err := json.Unmarshal(p.M, &ct); err != nil {
			return err
		}
		m = ct
	case "TICKET_CLOSE":
		var tc TicketClose
		if err := json.Unmarshal(p.M, &tc); err != nil {
			return err
		}
		m = tc
	default:
		Pwarn("update this package")
	}

	mType := reflect.TypeOf(m)
	for _, v := range b.handlers[mType] {
		v.Call([]reflect.Value{
			reflect.ValueOf(b),
			reflect.ValueOf(m),
		})
	}

	return nil
}
