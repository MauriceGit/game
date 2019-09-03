// player.go
package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type GuiPlayerUpdate struct {
	Id        int          `json:"id"`
	Size      float32      `json:"size"`
	Positions []mgl32.Vec2 `json:"positions"`
	// If bullets exist, they are moving fast, so need constant update
	Bullets []Bullet `json:"bullets"`
}

type GuiNewPlayer struct {
	Id    int        `json:"id"`
	Color mgl32.Vec3 `json:"color"`
	Name  string     `json:"name"`
}

type GuiUpdate struct {
	// All ids from Players that can be deleted from the Gui
	RemovedPlayer []int `json:"removedPlayer"`
	// Complete information for every player that moved/changed
	UpdatedPlayer []GuiPlayerUpdate `json:"updatedPlayer"`
	// Basic information for new players. No positions or bullets yet
	NewPlayer []GuiNewPlayer `json:"newPlayer"`
	// All ids of food that can be removed from the Gui
	RemovedFood []int `json:"removedFood"`
	// Static information for new food
	NewFood []Food `json:"newFood"`
}

type GuiConnection struct {
	connection net.Conn
	isNew      bool
}

func handleGuiConnectionInput(conn net.Conn) {
	for {
		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			fmt.Printf("Websocket read client error: %v\n", err)
			// We don't close the connection when a read fails. This is only done in the sender routine.
			// Which then also handles the organisational stuff around connections.
			return
		}
		fmt.Printf("<--- %v\n", string(msg))
	}
}

func HandleGuiConnections(app *Application) {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Websocket listen error: %v\n", err)
	}
	u := ws.Upgrader{
		OnHeader: func(key, value []byte) (err error) {
			//fmt.Printf("non-websocket header: %q=%q\n", key, value)
			return
		},
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Websocket accept error: %v\n", err)
		}

		fmt.Printf("New connection\n")

		_, err = u.Upgrade(conn)
		if err != nil {
			fmt.Printf("Websocket upgrade error: %v\n", err)
		}

		app.guiConnectionMutex.Lock()
		app.guiConnections = append(app.guiConnections, GuiConnection{conn, true})
		app.guiConnectionMutex.Unlock()
		go handleGuiConnectionInput(conn)
	}
}

func removeFromSlice(s []GuiConnection, i int) []GuiConnection {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func BroadcastGuiInfo(app *Application, update GuiUpdate) {

	var removeGuiIndices []int

	message, err := json.Marshal(update)
	if err != nil {
		fmt.Printf("json Marshal error: %v\n", err)
		return
	}

	app.guiConnectionMutex.Lock()

	allExistingPlayers := make([]GuiNewPlayer, 0)
	allExistingFood := make([]Food, 0)

	for i, gui := range app.guiConnections {

		specificMessage := message
		if gui.isNew {

			if len(allExistingPlayers) == 0 {
				// Just create this list once.
				allExistingPlayers = make([]GuiNewPlayer, 0, len(app.Players))
				for _, v := range app.Players {
					allExistingPlayers = append(allExistingPlayers, GuiNewPlayer{v.Player.Id, v.Player.Color, v.Player.Name})
				}
				allExistingFood = make([]Food, 0, len(app.Food))
				for _, f := range app.Food {
					allExistingFood = append(allExistingFood, f)
				}
			}

			// Just once overwrite the update message.
			thisUpdate := update
			thisUpdate.NewPlayer = allExistingPlayers
			thisUpdate.NewFood = allExistingFood

			specificMessage, err = json.Marshal(thisUpdate)
			if err != nil {
				fmt.Printf("Specific json Marshal error: %v\n", err)
				continue
			}
			app.guiConnections[i].isNew = false
		}

		err = wsutil.WriteServerMessage(gui.connection, ws.OpBinary, []byte(string(specificMessage)))
		if err != nil {
			fmt.Printf("Websocket write message error: %v\n", err)
			removeGuiIndices = append(removeGuiIndices, i)
			gui.connection.Close()
		}

	}

	// Removes all Guis from the connection list, where a write failed.
	// We iterate in reverse so all the indices stay correct!
	for i := len(removeGuiIndices) - 1; i >= 0; i-- {
		app.guiConnections = removeFromSlice(app.guiConnections, removeGuiIndices[i])
	}

	app.guiConnectionMutex.Unlock()
}
