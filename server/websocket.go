package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lxzan/gws"
	"github.com/mjishu/pokeDate/database"
)

// ? How to make other connectiosn listen for just this message?

const (
	PingInternval = 5 * time.Second
	PingWait      = 10 * time.Second
)

type Handler struct {
	pool *pgxpool.Pool
}

func (c *Handler) OnOpen(socket *gws.Conn) {
	fmt.Println("open conenction")
	_ = socket.SetDeadline(time.Now().Add(PingInternval + PingWait))
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	fmt.Println("closed connection")
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	fmt.Println("ping")
	_ = socket.SetDeadline(time.Now().Add(PingInternval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {
	fmt.Println("pong")
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) { // how to get pool here
	fmt.Printf("message is %v\n", message)
	var messageInfo database.WSMessage

	body, err := io.ReadAll(message.Data)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &messageInfo)
	if err != nil {
		return
	}

	// insert into database here with messageInfo,
	err = database.CreateWSMessage(c.pool, messageInfo)
	if err != nil {
		return
	}

	defer message.Close()
	socket.WriteMessage(message.Opcode, message.Bytes())
}
