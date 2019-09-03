// player.go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

func initPlayerVanguard(p *Player) {
	p.Name = p.Name + "_van"
	p.defense = 0.0
	p.attackStrength = 1.0
	p.viewRadius = 120.0
	p.attackRadius = 100.0
	p.maxAvailableShots = 4
	p.availableShots = 4
}

func runPlayerVanguard(p *Player, output chan PlayerOutput) {

	id := p.Id
	input := p.inputChannel

	target := mgl32.Vec2{rand.Float32() * 1000., rand.Float32() * 1200.}

	lastRandomTarget := time.Now()

	for {
		// Only react, if there is new data incoming. Blocking read.
		// We don't need the position, but might need the filled kd-tree.
		data := <-input
		data = data
		targetAcquired := false

		var closestFood []int
		tmpPos := p.Pos[0]
		min := []float64{float64(tmpPos[0]) - 100, float64(tmpPos[1]) - 100}
		max := []float64{float64(tmpPos[0]) + 100, float64(tmpPos[1]) + 100}

		data.Food.Search(min, max,
			func(min, max []float64, value interface{}) bool {
				closestFood = append(closestFood, value.(int))
				return true
			},
		)

		for _, f := range closestFood {
			if !targetAcquired || tmpPos.Sub(data.FoodDict[f].P).Len() < tmpPos.Sub(target).Len() {
				target = data.FoodDict[f].P
				targetAcquired = true
			}
		}

		if len(closestFood) == 0 && time.Now().Sub(lastRandomTarget).Seconds() >= 1 {
			target = mgl32.Vec2{rand.Float32() * 1500., rand.Float32() * 1200.}
			lastRandomTarget = time.Now()
		}

		action := DEFAULT

		if rand.Int()%10 == 0 {
			action = NEW_BULLET
		}

		select {
		case output <- PlayerOutput{id, action, target, target}:

		default:
			fmt.Println("Writing to channel failed.")
		}
	}

}
