// player.go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	min = []float64{0, 0}
	max = []float64{0, 0}
)

func initPlayerTank(p *Player) {
	p.Name = p.Name + "_tank"
	p.defense = 4.0
	p.viewRadius = 200.0
	p.attackStrength = 4.0
	p.attackRadius = 120.0
	p.maxAvailableShots = 10
	p.availableShots = 10
}

// getFoodTarget returns the food that is closest within the view radius of the player.
func getFoodTarget(p *Player, data PlayerInput) (mgl32.Vec2, bool) {
	targetAcquired := false
	ok := false
	var target mgl32.Vec2
	var closestFood []int
	tmpPos := p.Pos[0]
	min[0], min[1] = float64(tmpPos[0]-p.viewRadius), float64(tmpPos[1]-p.viewRadius)
	max[0], max[1] = float64(tmpPos[0]+p.viewRadius), float64(tmpPos[1]+p.viewRadius)

	data.Food.Search(min, max,
		func(min, max []float64, value interface{}) bool {
			if data.FoodDict[value.(int)].P.Sub(tmpPos).Len() < p.viewRadius {
				closestFood = append(closestFood, value.(int))
			}
			return true
		},
	)

	for _, f := range closestFood {
		if !targetAcquired || tmpPos.Sub(data.FoodDict[f].P).Len() < tmpPos.Sub(target).Len() {
			target = data.FoodDict[f].P
			targetAcquired = true
			ok = true
		}
	}

	return target, ok
}

func getEnemyTarget(p *Player, data PlayerInput) (int, bool) {
	targetAcquired := false
	ok := false
	var target int
	var closestEnemy []int
	tmpPos := p.Pos[0]
	min[0], min[1] = float64(tmpPos[0]-p.viewRadius), float64(tmpPos[1]-p.viewRadius)
	max[0], max[1] = float64(tmpPos[0]+p.viewRadius), float64(tmpPos[1]+p.viewRadius)

	data.Players.Search(min, max,
		func(min, max []float64, value interface{}) bool {
			if value.(int) != p.Id && data.PlayerDict[value.(int)].Player.Pos[0].Sub(tmpPos).Len() < p.viewRadius {
				closestEnemy = append(closestEnemy, value.(int))
			}
			return true
		},
	)

	var currentEnemy *Player
	for _, f := range closestEnemy {
		e := data.PlayerDict[f].Player
		if !targetAcquired || tmpPos.Sub(e.Pos[0]).Len() < tmpPos.Sub(currentEnemy.Pos[0]).Len() {
			currentEnemy = e
			target = f
			targetAcquired = true
			ok = true
		}
	}

	return target, ok
}

func runPlayerTank(p *Player, output chan PlayerOutput) {

	id := p.Id
	input := p.inputChannel

	target := mgl32.Vec2{rand.Float32() * FIELD_SIZE[0], rand.Float32() * FIELD_SIZE[1]}
	shootingTarget := target

	lastRandomTarget := time.Now()
	lastTimedTarget := time.Now().Add(-1 * time.Minute)
	timedTargetAcquired := false

	for {
		// Only react, if there is new data incoming. Blocking read.
		// We don't need the position, but might need the filled kd-tree.
		data := <-input
		targetAcquired := false
		action := DEFAULT

		// Determine any enemies within our attack zone and attack them. Highest priority!

		if timedTargetAcquired && time.Now().Sub(lastTimedTarget).Seconds() >= 1 {
			//fmt.Printf("Old: %v\nNew: %v\n", lastTimedTarget, time.Now())
			timedTargetAcquired = false
		}

		if !timedTargetAcquired {

			enemyTarget, ok := getEnemyTarget(p, data)

			if ok {
				targetAcquired = true
				enemy := data.PlayerDict[enemyTarget].Player
				enemyDistance := enemy.Pos[0].Sub(p.Pos[0]).Len()

				if enemy.Size > (p.Size*2.0) || p.availableShots < 1 {
					// Flee in the opposite direction
					target = p.Pos[0].Sub(enemy.Pos[0]).Add(p.Pos[0])
					timedTargetAcquired = true
					lastTimedTarget = time.Now()
					fmt.Printf("flee: enemy: %v, me: %v, shots: %v\n", enemy.Size, p.Size, p.availableShots)
				} else {
					// Circle around target
					target = data.PlayerDict[enemyTarget].Player.Pos[0]
					targetDir := enemy.Pos[0].Sub(p.Pos[0])
					// This is the tangential/perpendicular vector. We try to circle the enemy!
					target[0], target[1] = targetDir[1], -targetDir[0]

					// With the tendency to merge onto his position
					target = target.Add(p.Pos[0]).Add(targetDir.Mul(0.1))

					if p.availableShots >= 1 && enemyDistance <= p.attackRadius {
						fmt.Printf("attack\n")
						action = NEW_BULLET
						shootingTarget = enemy.Pos[0]
					}
				}
			}

			// If no enemy is present, try to get some food
			if !targetAcquired {
				foodTarget, ok := getFoodTarget(p, data)
				if ok {
					targetAcquired = true
					target = foodTarget
				}
			}
		}

		// If we still have no target, just go somewhere for a second...
		if !targetAcquired && time.Now().Sub(lastRandomTarget).Seconds() >= 1 {
			target = mgl32.Vec2{rand.Float32() * FIELD_SIZE[0], rand.Float32() * FIELD_SIZE[1]}
			lastRandomTarget = time.Now()
		}

		select {
		case output <- PlayerOutput{id, action, shootingTarget, target}:

		default:
			fmt.Println("Writing to channel failed.")
		}
	}

}
