// game.go
package main

import (
	"fmt"
	"math"
	"math/rand"

	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	//"github.com/pkg/profile"
	rbang "github.com/tidwall/rbang-go"
)

type Shooting int

const (
	EPS                       = 0.0000001
	MAX_PLAYER_COUNT          = 13
	MAX_FOOD_COUNT            = 5000
	MAX_FOOD_SIZE             = 2.5
	DEFAULT          Shooting = iota
	NEW_BULLET                = iota
)

var (
	maxFoodId = 0
	// Constant, but needs to be declared as var because of non-constant array type
	FIELD_SIZE = [...]float32{3000, 2000}
)

type PlayerAction struct {
	Player *Player `json:"player"`
	action PlayerOutput
}

// Food is always stationary.
type Food struct {
	Id int        `json:"id"`
	P  mgl32.Vec2 `json:"p"`
	S  float32    `json:"s"`
}

type Application struct {
	playerOutput chan PlayerOutput

	Players    map[int]*PlayerAction `json:"players"`
	playerTree rbang.RTree
	Food       map[int]Food `json:"food"`
	foodTree   rbang.RTree

	guiConnectionMutex *sync.Mutex
	guiConnections     []GuiConnection
}

// generateNewPlayers should be remodelled to read from a channel (non-blocking), if a new player
// should be created. If so, create it. This channel could be filled from some browser or internally...
func (app *Application) generateNewPlayers() []GuiNewPlayer {

	guiNewPlayer := make([]GuiNewPlayer, 0)

	if len(app.Players) < MAX_PLAYER_COUNT {
		for i := 0; i < MAX_PLAYER_COUNT-len(app.Players); i++ {
			pa := &PlayerAction{RegisterPlayer(app.playerOutput, PLAYER_TYPE_TANK), PlayerOutput{}}
			app.Players[pa.Player.Id] = pa

			pPos := []float64{float64(pa.Player.Pos[0].X()), float64(pa.Player.Pos[0].Y())}
			app.playerTree.Insert(pPos, nil, pa.Player.Id)

			guiNewPlayer = append(guiNewPlayer, GuiNewPlayer{pa.Player.Id, pa.Player.Color, pa.Player.Name})
		}
	}

	return guiNewPlayer
}

func collectPlayerActions(app *Application) {
	allUpdates := false

	// Get all action updates from Players/Bots
	for !allUpdates {
		select {
		case playerData := <-app.playerOutput:
			player := app.Players[playerData.Id]
			player.action = playerData
			app.Players[playerData.Id] = player
		default:
			allUpdates = true
		}
	}
}

func (app *Application) generateNewFood(count int) []Food {
	newFood := make([]Food, 0, count)
	for i := 0; i < count; i++ {
		f := Food{maxFoodId, mgl32.Vec2{rand.Float32() * FIELD_SIZE[0], rand.Float32() * FIELD_SIZE[1]}, rand.Float32() * MAX_FOOD_SIZE}
		app.Food[maxFoodId] = f
		newFood = append(newFood, f)

		app.foodTree.Insert([]float64{float64(f.P.X()), float64(f.P.Y())}, nil, f.Id)

		maxFoodId++

	}
	return newFood
}

func (app *Application) manageFood() []int {
	removedFood := make([]int, 0)

	// Cache the array.
	min := []float64{0, 0}
	max := []float64{0, 0}

	// Check how much food was eaten this round.
	// For efficiency reasons, we only check within a square, not the circle!
	for _, p := range app.Players {

		var eatenFood []int
		tmpPos := p.Player.Pos[0]
		size := float64(p.Player.Size)
		min[0], min[1] = float64(tmpPos[0])-size, float64(tmpPos[1])-size
		max[0], max[1] = float64(tmpPos[0])+size, float64(tmpPos[1])+size

		app.foodTree.Search(min, max,
			func(min, max []float64, value interface{}) bool {
				eatenFood = append(eatenFood, value.(int))
				return true
			},
		)

		for _, f := range eatenFood {

			fPos := app.Food[f].P
			if fPos.Sub(p.Player.Pos[0]).Len() <= float32(size) {

				p.Player.Size += app.Food[f].S
				// Player can never exceed the maximum food size.
				// This avoids situations where a player grows so quick, that he eats the entire field
				// within one single frame before being sized down again in the  next, triggering the entire food list again...
				if p.Player.Size > PLAYER_MAX_SIZE {
					p.Player.Size = PLAYER_MAX_SIZE
				}

				removedFood = append(removedFood, f)
				delete(app.Food, f)

				app.foodTree.Delete([]float64{float64(fPos[0]), float64(fPos[1])}, nil, f)
			}
		}
	}
	return removedFood
}

func runSimulation(app *Application) {

	var guiNewPlayer []GuiNewPlayer

	var guiRemovedFood []int
	var guiNewFood []Food
	playerTreePos := []float64{0, 0}

	framecount := 0

	for {

		time.Sleep(16 * time.Millisecond)

		guiNewPlayer = app.generateNewPlayers()

		collectPlayerActions(app)

		guiUpdatedPlayer := make([]GuiPlayerUpdate, 0, len(app.Players))
		guiRemovedPlayer := make([]int, 0)

		bulletHits := make([]BulletHit, 0)

		// Update world order
		for k, v := range app.Players {

			// Save the old position
			playerTreePos[0], playerTreePos[1] = float64(v.Player.Pos[0].X()), float64(v.Player.Pos[0].Y())

			bulletHits = append(bulletHits, v.Player.calcUpdate(app, v.action)...)

			guiUpdatedPlayer = append(guiUpdatedPlayer, GuiPlayerUpdate{v.Player.Id, v.Player.Size, v.Player.Pos, v.Player.Bullets})

			guiRemovedFood = append(guiRemovedFood, app.manageFood()...)

			app.Players[k].Player.availableShots += 0.01

			// Re-insert the updated player position and ID into the tree for the next round
			app.playerTree.Delete(playerTreePos, nil, v.Player.Id)
			playerTreePos[0], playerTreePos[1] = float64(v.Player.Pos[0].X()), float64(v.Player.Pos[0].Y())
			app.playerTree.Insert(playerTreePos, nil, v.Player.Id)

			if v.Player.Size <= PLAYER_MIN_SIZE {
				guiRemovedPlayer = append(guiRemovedPlayer, k)
			}

		}

		guiNewFood = app.generateNewFood(MAX_FOOD_COUNT - len(app.Food))

		// Calculating the mass reduction here means, that it will be one frame too late in the Gui and
		// will be taken into account for other players in the next frame as well... Not bad per se but good to know.
		for _, v := range bulletHits {
			// Cheap way of verifying, that we don't accidentally try to remove a player twice because he was hit by a bullet twice in a frame
			// and added twice by several players to be removed...
			if app.Players[v.target].Player.Size > PLAYER_MIN_SIZE {

				s := app.Players[v.shooter].Player
				t := app.Players[v.target].Player
				//sizeBeforeHit := app.Players[v.target].Player.Size
				app.Players[v.target].Player.Size -= float32(math.Min(float64(s.attackStrength+BULLET_DAMAGE-t.defense), BULLET_MIN_DAMAGE))

				if app.Players[v.target].Player.Size <= PLAYER_MIN_SIZE {
					guiRemovedPlayer = append(guiRemovedPlayer, v.target)
				}
			}

		}

		for _, index := range guiRemovedPlayer {
			fmt.Printf("Player %v (%v) died.\n", index, app.Players[index].Player.Name)
			p := app.Players[index].Player

			playerTreePos[0], playerTreePos[1] = float64(p.Pos[0].X()), float64(p.Pos[0].Y())
			app.playerTree.Delete(playerTreePos, nil, p.Id)

			delete(app.Players, index)
		}

		if framecount%6 == 0 {
			for _, v := range app.Players {
				select {
				case v.Player.inputChannel <- PlayerInput{app.playerTree, app.Players, app.foodTree, app.Food}:
				default:
					fmt.Println("Write to Bot input channel failed.")
				}
			}
		}

		guiUpdate := GuiUpdate{guiRemovedPlayer, guiUpdatedPlayer, guiNewPlayer, guiRemovedFood, guiNewFood}

		// Broadcast new simulation state to all Guis.
		BroadcastGuiInfo(app, guiUpdate)

		guiRemovedFood = make([]int, 0)
		guiNewFood = make([]Food, 0)
		framecount++
	}
}

func main() {

	//defer profile.Start().Stop()

	app := &Application{
		Players:            make(map[int]*PlayerAction),
		playerOutput:       make(chan PlayerOutput, MAX_PLAYER_COUNT),
		guiConnections:     make([]GuiConnection, 0),
		guiConnectionMutex: &sync.Mutex{},
	}

	app.playerTree = rbang.New(2)

	//app.generateNewPlayers()

	app.Food = make(map[int]Food, MAX_FOOD_COUNT)
	app.foodTree = rbang.New(2)
	app.generateNewFood(MAX_FOOD_COUNT)

	go HandleGuiConnections(app)

	runSimulation(app)

}
