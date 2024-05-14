package player

import (
	"log/slog"

	"github.com/awoo-detat/werewolf/gamechannel"
	"github.com/awoo-detat/werewolf/gamechannel/client"
	"github.com/awoo-detat/werewolf/gamechannel/server"
	"github.com/awoo-detat/werewolf/role"

	"github.com/google/uuid"
)

type Player struct {
	ID          uuid.UUID
	Name        string
	Role        *role.Role
	Views       []*View
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
	slog.Info("setting player role", "player", p, "role", r)
}

func (p *Player) AddView(v *View) {
	p.Views = append(p.Views, v)
	slog.Info("adding view", "view", v, "player", p)
}

// Message handles sending a message to the client, wrapping it in error handling
func (p *Player) Message(t server.MessageType, payload interface{}) error {
	m, err := server.NewMessage(t, payload)
	if err != nil {
		return err
	}
	return p.socket.WriteMessage(1, m)
}

func (p *Player) Play() {
	defer p.socket.Close()

	for {
		_, c, err := p.socket.ReadMessage()
		if err != nil {
			slog.Error("error reading message", "player", p, "error", err)
		}

		m := client.Decode(c)

		switch m.Type {
		case client.Awoo:
			p.Message(server.Awoo, "awooooooooo")
		default:
			slog.Warn("unknown message ", "message", m)
		}
	}
}
