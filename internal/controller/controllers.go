package controller

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/propane-server/internal/service"
)

type Controllers interface {
	Measurement() MeasurementController
	Info() InfoController

	Route(e *echo.Echo)
	Loop()
}

type controllers struct {
	measurementController MeasurementController
	infoController        InfoController

	connections map[*connection]bool
	broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

func NewControllers(services service.Services) Controllers {
	c := &controllers{}
	c.measurementController = newMeasurementController(services.Measurement(), c.publish)
	c.infoController = newInfoController()
	return c
}

func (c *controllers) Measurement() MeasurementController {
	return c.measurementController
}

func (c *controllers) Info() InfoController {
	return c.infoController
}

func (c *controllers) Route(e *echo.Echo) {
	e.GET("/", c.infoController.Info)

	e.POST("/insert", c.measurementController.Insert)
	e.GET("/ws", c.handleWebsockets)
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

type Broadcaster func([]byte)

func (c *controllers) Loop() {
	for {
		select {
		case conn := <-c.register:
			c.connections[conn] = true
		case conn := <-c.unregister:
			delete(c.connections, conn)
			close(conn.send)
		case m := <-c.broadcast:
			for conn := range c.connections {
				select {
				case conn.send <- m:
				default:
					delete(c.connections, conn)
					close(conn.send)
				}
			}
		}
	}
}

func (c *controllers) publish(b []byte) {
	c.broadcast <- b
}

var (
	upgrade = websocket.Upgrader{}
)

func (c *controllers) handleWebsockets(ct echo.Context) error {
	ws, err := upgrade.Upgrade(ct.Response(), ct.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	conn := &connection{send: make(chan []byte, 256), ws: ws}
	c.register <- conn
	defer func() { c.unregister <- conn }()
	go conn.writePump()
	conn.readPump()
	return nil
}

func (c *connection) readPump() {
	defer func() { c.ws.Close() }()
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *connection) writePump() {
	defer func() { c.ws.Close() }()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			err := c.ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}
