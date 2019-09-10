// player.go
package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	rbang "github.com/tidwall/rbang-go"
)

const (
	PLAYER_SEGMENT_SIZE      = 15.0
	PLAYER_MAX_SEGMENT_COUNT = 5
	PLAYER_INITIAL_SIZE      = 5
	PLAYER_MAX_SIZE          = 35.0
	PLAYER_MIN_SIZE          = 2.0
	PLAYER_SPEED             = 1.5
	PLAYER_SIZE_REDUCTION    = 0.01
	BULLET_RANGE             = 300
	BULLET_MIN_DAMAGE        = 0.3
	BULLET_DAMAGE            = 1.0
)

type PlayerType int

const (
	PLAYER_TYPE_TANK PlayerType = iota
	PLAYER_TYPE_VANGUARD
	PLAYER_TYPE_MAGE
)

var (
	g_playerIDs      int = 0
	c_availableNames     = [...]string{"Rodrigo", "Zayne", "Reyna", "Avah", "Campbell", "Tianna", "Cameron", "Nayeli", "Ricky", "Myles", "Aiden", "Cristian", "Edward", "Meredith", "Waylon", "Kolby", "Lorelai", "Brooks", "Pranav", "Gage", "Jaydin", "Braelyn", "Damarion", "Miracle", "Drew", "Harley", "Carley", "Rachael", "Dillan", "Mckenzie", "Emmalee", "Kayla", "Joe", "Lucia", "Ivy", "Cayden", "Dominik", "Kristian", "Emiliano", "Jasmine", "Luciana", "Valery", "Savanna", "Natalee", "Matteo", "Kenyon", "Bridger", "Litzy", "Soren", "Blaine", "Maribel", "London", "Scott", "Santos", "Theresa", "Delaney", "Keshawn", "Sharon", "Misael", "Benjamin", "Leanna", "Fatima", "Rory", "Susan", "Larry", "Sabrina", "Camden", "Javier", "Chaya", "Harper", "Elliana", "Izayah", "Matilda", "Cherish", "Beau", "Zoey", "King", "Aryanna", "Maximilian", "Dashawn", "Kamren", "Caden", "Ashlynn", "Emma", "Gary", "Julia", "Alessandro", "Eden", "Margaret", "Jewel", "Tania", "Myah", "Brodie", "Nasir", "Kaylen", "Jadiel", "Trystan", "Roselyn", "Kyleigh", "Lana", "Andrea", "Dexter", "Emery", "Elian", "Alexis", "Angela", "Sergio", "Kailyn", "Jett", "Emmanuel", "Keegan", "Ryker", "Ibrahim", "Cohen", "Dalia", "Anahi", "Jordan", "Raphael", "Janelle", "Mireya", "Sydnee", "Abram", "Ari", "Amaris", "Sasha", "Harry", "Jerry", "Eliezer", "Cash", "Landyn", "Maren", "Braedon", "Maia", "Lea", "August", "Valeria", "Orlando", "Liana", "Darion", "Angel", "Max", "Damion", "Shiloh", "Alden", "Annabella", "Danny", "Davian", "Immanuel", "Albert", "Jayce", "Shawn", "Jazlynn", "Patricia", "Alijah", "Destinee", "Stacy", "Kendal", "Jordan", "Aracely", "Raiden", "Kendall", "Cooper", "Jayson", "Jamir", "Makayla", "Harold", "Anabelle", "Talia", "Joaquin", "Bryson", "Carleigh", "Alyson", "Izabella", "Alison", "Connor", "Fisher", "Micah", "Colby", "Paityn", "Desiree", "Hadassah", "Isabel", "Nash", "Sloane", "Jack", "Lilliana", "Toby", "Zander", "Tatum", "Jamya", "Braylen", "Sanai", "Norah", "Tamia", "Mattie", "Seamus", "Mina", "Alicia", "Jordon", "Anton", "Zion", "Amanda", "Nathalie", "Adolfo", "Maximillian", "Damian", "Elisabeth", "Ignacio", "Roland", "Abagail", "Darren", "Anika", "Shamar", "Joy", "Mercedes", "Kamron", "Harper", "Yazmin", "Ramiro", "Aliza", "Malakai", "Aleena", "Ryann", "Faith", "Clay", "April", "Lydia", "Emmy", "Leticia", "Ayden", "Martin", "Jolie", "Janet", "Aubrie", "Ronin", "Angie", "Moshe", "Damaris", "Charlie", "Nola", "Gabriella", "Marshall", "Gideon", "Kaitlyn", "Carlee", "Alyvia", "Abel", "Alfonso", "Lindsey", "Cassius", "Damari", "Rhett", "Kaiden", "Tucker", "Meadow", "Eileen", "Maria", "Brayden", "Yasmine", "Ali", "Lucy", "Bradley", "Brett", "Mariam", "Shaniya", "Jaime", "Tiara", "Ayana", "Maliyah", "Pedro", "Elena", "Noah", "Makenzie", "Gretchen", "Jessie", "Emily", "Anthony", "Alisson", "Kaleigh", "Marley", "Lennon", "Brady", "Finley", "Marcel", "Madilynn", "Brandon", "Simeon", "Landen", "Sara", "Kaitlin", "Esmeralda", "Quinten", "Jaelynn", "Sammy", "Kobe", "Marisa", "Reginald", "Marley", "Lauren", "Roberto", "Yadiel", "Leonel", "Dax", "Antwan", "Trey", "Chad", "Jane", "Derek", "Camila", "Dylan", "Azul", "Sarahi", "Nataly", "Yahir", "Bo", "Ellie", "Nathanael", "Valentino", "Mathew", "Teagan", "Jamie", "Ronald", "Sheldon", "Danica", "Phillip", "Baylee", "Janessa", "Kaylin", "Kaylie", "Elsa", "Owen", "Mikaela", "Rose", "Devon", "Rylee", "Alessandra", "Rolando", "Rocco", "Jenny", "Lamar", "Juliet", "Jerimiah", "Adalynn", "Elisa", "Jairo", "Elise", "Sarah", "Ricardo", "Alaina", "Giovanni", "Jordyn", "Elias", "Harley", "Rylee", "Maeve", "Helen", "Tanner", "Ashlee", "Brianna", "Maggie", "Cyrus", "Boston", "Kiera", "Lamont", "Gabriela", "Ava", "Dorian", "Kiara", "Deja", "Itzel", "Ryan", "Ian", "Zaid", "Leah", "Lacey", "Zane", "Aileen", "Luz", "Osvaldo", "Judah", "Jared", "Kristen", "Tristian", "Briana", "Bailee", "Milton", "Audrey", "Noelle", "Tyrone", "Braeden", "Cara", "Zaiden", "Addison", "Avery", "Janiah", "Micah", "Tara", "Taryn", "Reese", "Tristen", "Douglas", "Amelie", "Joselyn", "Brisa", "Kymani", "Zechariah", "Beckham", "Arianna", "Hadley", "Kallie", "Sanaa", "Kole", "Cruz", "Marely", "Adrianna", "Jose", "Calvin", "Rey", "Dwayne", "Kirsten", "Andreas", "Melina", "Josephine", "Jaylene", "Amira", "Carsen", "Fernando", "Alexandra", "Jaylynn", "Makhi", "Andrew", "Lorenzo", "Samson", "Alejandro", "Santino", "Angelina", "Cortez", "Alejandra", "Brendon", "Maximo", "Giada", "Darwin", "Molly", "Zoe", "Heath", "Allie", "Angel", "Magdalena", "Kenna", "Greyson", "Giovanna", "Lance", "Trevin", "Zain", "Elizabeth", "Byron", "Nicolas", "Eddie", "Alexzander", "Elliott", "Rosa", "Jagger", "Ernest", "Arjun", "Johanna", "Zion", "Terrell", "Maverick", "Dominic", "Karen", "Karina", "Yoselin", "Logan", "Sophia", "Vicente", "Tyshawn", "Jonathan", "Katherine", "Josiah", "Rishi", "Sofia", "Kaylah", "Francisco", "Todd", "Dakota", "Camryn", "London", "Moriah", "Charlotte", "Adison", "Jessie", "Griffin", "Zackary", "Dalton", "Clarissa", "Lorena", "Nathalia", "Braydon", "Jamarcus", "Giancarlo"}
)

type PlayerInput struct {
	Players    rbang.RTree
	PlayerDict map[int]*PlayerAction
	Food       rbang.RTree
	FoodDict   map[int]Food
}

type PlayerOutput struct {
	Id             int
	Action         Shooting
	ShootingTarget mgl32.Vec2
	TargetPos      mgl32.Vec2
}

type Bullet struct {
	originPos mgl32.Vec2
	Pos       mgl32.Vec2 `json:"pos"`
	dir       mgl32.Vec2
}

type BulletHit struct {
	target  int
	shooter int
}

type Player struct {
	inputChannel chan PlayerInput

	// Only exported fields are included in the Json
	Id      int          `json:"id"`
	Pos     []mgl32.Vec2 `json:"pos"`
	Bullets []Bullet     `json:"bullets"`
	Color   mgl32.Vec3   `json:"color"`
	Name    string       `json:"name"`
	Size    float32      `json:"size"`

	// Player stats
	playerType        PlayerType
	defense           float32
	viewRadius        float32
	attackStrength    float32
	attackRadius      float32
	maxAvailableShots int
	availableShots    float32
}

func calcSegment(size float32, positions []mgl32.Vec2, i int, xin, yin float32) {
	dx := float64(xin - positions[i].X())
	dy := float64(yin - positions[i].Y())
	angle := math.Atan2(dy, dx)
	positions[i][0] = xin - float32(math.Cos(angle))*size*2
	positions[i][1] = yin - float32(math.Sin(angle))*size*2
}

func (p *Player) calcUpdate(app *Application, output PlayerOutput) []BulletHit {

	min := []float64{0, 0}
	max := []float64{0, 0}

	// Calc the very first position manually
	dir := output.TargetPos.Sub(p.Pos[0]).Normalize()

	speedFactor := (1.0 - (p.Size-PLAYER_MIN_SIZE)/(PLAYER_MAX_SIZE-PLAYER_MIN_SIZE))
	speedFactor = float32(math.Pow(float64(speedFactor), 2.0))

	dir = dir.Mul(speedFactor)

	if dir.Len() < 0.5 {
		if dir.Len() <= EPS {
			dir[0] = EPS
			dir[1] = EPS
		}
		dir = dir.Normalize()
	}
	dir = dir.Mul(PLAYER_SPEED)

	p.Pos[0] = p.Pos[0].Add(dir)

	p.Size -= PLAYER_SIZE_REDUCTION + p.Size*0.0001

	if p.Size > PLAYER_MAX_SIZE {
		p.Size = PLAYER_MAX_SIZE
	}

	for i := 0; i < len(p.Pos)-1; i++ {
		calcSegment(p.Size, p.Pos, i+1, p.Pos[i].X(), p.Pos[i].Y())
	}

	bulletHit := make([]BulletHit, 0)
	removeIndices := make([]int, 0)
	for i, b := range p.Bullets {

		if b.Pos.Sub(b.originPos).Len() > BULLET_RANGE {
			removeIndices = append(removeIndices, i)
		} else {
			p.Bullets[i].Pos = b.Pos.Add(b.dir.Mul(8.2))

			tmpPos := p.Bullets[i].Pos
			possibleBulletRange := float64(PLAYER_SEGMENT_SIZE * PLAYER_MAX_SEGMENT_COUNT)
			min[0], min[1] = float64(tmpPos[0])-possibleBulletRange, float64(tmpPos[1])-possibleBulletRange
			max[0], max[1] = float64(tmpPos[0])+possibleBulletRange, float64(tmpPos[1])+possibleBulletRange

			app.playerTree.Search(min, max,
				func(min, max []float64, value interface{}) bool {

					pId := value.(int)

					if pId == p.Id {
						return true
					}

					for _, tail := range app.Players[pId].Player.Pos {

						if tail.Sub(tmpPos).Len() < app.Players[pId].Player.Size {
							// Remove the bullet. It found its target.
							removeIndices = append(removeIndices, i)

							bulletHit = append(bulletHit, BulletHit{pId, p.Id})

							return false
						}
					}

					return true
				},
			)

		}
	}

	for i := len(removeIndices) - 1; i >= 0; i-- {
		p.Bullets[len(p.Bullets)-1], p.Bullets[removeIndices[i]] = p.Bullets[removeIndices[i]], p.Bullets[len(p.Bullets)-1]
		p.Bullets = p.Bullets[:len(p.Bullets)-1]
	}

	if output.Action == NEW_BULLET {
		if p.availableShots >= 1.0 {
			p.Bullets = append(p.Bullets, Bullet{p.Pos[0], p.Pos[0], output.ShootingTarget.Sub(p.Pos[0]).Normalize()})
			p.availableShots -= 1.0
		}
	}

	return bulletHit

}

func RegisterPlayer(outputChannel chan PlayerOutput, playerType PlayerType) *Player {

	newPlayer := &Player{
		inputChannel: make(chan PlayerInput, 10),
		Id:           g_playerIDs,
		Pos:          []mgl32.Vec2{mgl32.Vec2{rand.Float32() * FIELD_SIZE[0], rand.Float32() * FIELD_SIZE[1]}},
		Bullets:      []Bullet{},
		Color:        mgl32.Vec3{rand.Float32() * 255., rand.Float32() * 255., rand.Float32() * 255.},
		Name:         c_availableNames[g_playerIDs%len(c_availableNames)],
		Size:         PLAYER_INITIAL_SIZE,
		playerType:   playerType,
	}
	g_playerIDs++

	segments := rand.Intn(PLAYER_MAX_SEGMENT_COUNT - 1)
	for i := 0; i < segments; i++ {
		newPlayer.Pos = append(newPlayer.Pos, mgl32.Vec2{newPlayer.Pos[0].X() + float32(i)*PLAYER_SEGMENT_SIZE, newPlayer.Pos[0].Y()})
	}

	switch newPlayer.playerType {
	case PLAYER_TYPE_TANK:
		initPlayerTank(newPlayer)
		go runPlayerTank(newPlayer, outputChannel)
	case PLAYER_TYPE_MAGE:
		initPlayerMage(newPlayer)
		go runPlayerMage(newPlayer, outputChannel)
	case PLAYER_TYPE_VANGUARD:
		initPlayerVanguard(newPlayer)
		go runPlayerVanguard(newPlayer, outputChannel)
	}

	return newPlayer
}
