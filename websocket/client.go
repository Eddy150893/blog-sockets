package websocket

import "github.com/gorilla/websocket"

type Client struct { //Se encargara de manejar conexiones
	hub      *Hub
	id       string          //identificar al cliente
	socket   *websocket.Conn // una conexion
	outbound chan []byte     //enviar mensajes
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
