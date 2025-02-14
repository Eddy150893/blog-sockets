package websocket

//Se encargara de manejar,centralizar y distribuir los clientes
import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		http.Error(w, "Coul not open websocket connection", http.StatusBadRequest)
	}
	client := NewClient(hub, socket)
	hub.register <- client

	go client.Write()
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("Client Connected", client.socket.RemoteAddr())
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client Disconnected", client.socket.RemoteAddr())
	client.socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}
	//Desplazamiento de elementos!!!!
	//copia las posiciones i+1 hasta :(final)
	//en las posiciones i hasta : (final)
	copy(hub.clients[i:], hub.clients[i+1:])
	//Deja la ultima posicion vacia
	hub.clients[len(hub.clients)-1] = nil
	//Redimensiona desde :(cero) hasta el final -1(nil)
	hub.clients = hub.clients[:len(hub.clients)-1]
}

func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}
