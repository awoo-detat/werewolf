package player

import (
	"log/slog"

	"github.com/awoo-detat/werewolf/gamechannel"
	"github.com/awoo-detat/werewolf/gamechannel/client"
	"github.com/awoo-detat/werewolf/gamechannel/server"
	"github.com/awoo-detat/werewolf/role"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Role        *role.Role `json:"-"`
	Views       []*View    `json:"-"`
	socket      Communicator
	gameChannel gamechannel.GameChannel
}

func NewPlayer(socket Communicator) *Player {
	p := &Player{
		ID:     uuid.New(),
		Views:  []*View{},
		socket: socket,
	}
	p.Message(server.IDSet, p.ID)
	return p
}

func (p *Player) String() string {
	if len(p.Name) != 0 {
		return p.Name
	}
	return p.ID.String()
}

func (p *Player) SetName(name string) {
	p.Name = name
	slog.Info("setting player name", "ID", p.ID, "Name", p.Name)
}

func (p *Player) SetRole(r *role.Role) {
	p.Role = r
	p.Message(server.RoleAssigned, r)
	slog.Info("setting player role", "player", p, "role", r)
}

func (p *Player) SetGameChannel(gc gamechannel.GameChannel) {
	p.gameChannel = gc
}

func (p *Player) AddView(v *View) {
	p.Views = append(p.Views, v)
	p.Message(server.View, v)
	slog.Info("adding view", "view", v, "player", p)
}

// Message handles sending a message to the client, wrapping it in error handling
func (p *Player) Message(t server.MessageType, payload interface{}) error {
	m, err := server.NewMessage(t, payload)
	if err != nil {
		return err
	}
	slog.Info("sending message to player", "message", m, "player", p)
	return p.socket.WriteMessage(1, m)
}

func (p *Player) Reconnect(c Communicator) {
	slog.Info("player reconnecting", "player", p)
	p.socket = c
	go p.Play()
	p.gameChannel <- &gamechannel.Activity{Type: gamechannel.Reconnect, From: p.ID}
}

func (p *Player) Play() {
	defer p.socket.Close()

	for {
		_, c, err := p.socket.ReadMessage()
		if err != nil {
			// TODO!!
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Info("closing connection", "player", p)
				break
			}
			slog.Error("error reading message", "player", p, "error", err)
			break
		}

		m, err := client.Decode(c)
		if err != nil {
			continue
		}

		switch m.Type {
		case client.Awoo:
			p.Message(server.Awoo, "awooooooooo")
		case client.SetName:
			p.SetName(m.PlayerName)
			p.gameChannel <- &gamechannel.Activity{Type: gamechannel.SetName, From: p.ID, Value: p.Name}
		case client.SetRoleset:
			p.gameChannel <- &gamechannel.Activity{Type: gamechannel.SetRoleset, Value: m.Roleset}
		case client.Vote:
			p.gameChannel <- &gamechannel.Activity{Type: gamechannel.Vote, From: p.ID, Value: m.Target}
		case client.Quit:
			slog.Info("player is quitting", "player", p)
			p.gameChannel <- &gamechannel.Activity{Type: gamechannel.Quit, From: p.ID}
			p.socket.Close()
			break
		default:
			slog.Warn("unknown message ", "message", m)
		}
	}
}
