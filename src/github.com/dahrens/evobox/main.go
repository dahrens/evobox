package main

import (
	evo "github.com/dahrens/evobox/evolution"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"time"
	"strconv"
	"strings"
	"math/rand"
)

type Environment struct {
	Evolvers  []evo.Evolver
	Creatures evo.Creatures
	Clock     *time.Ticker
	Tick      int
	Speed	  time.Duration
	Rand *rand.Rand
}

func NewEnvironment() *Environment {
	env := new(Environment)
	env.Evolvers = make([]evo.Evolver, 0)
	env.Creatures = make(evo.Creatures, 0)
	env.Clock = time.NewTicker(time.Second)
	env.Tick = 0
	env.Speed = time.Second
	env.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	return env
}

func (env *Environment) Init() {
	go env.SpawnMany(15, evo.GENDER_MALE)
	go env.SpawnMany(15, evo.GENDER_FEMALE)
}

func (env *Environment) SpawnMany(n int, gender evo.Gender) {
	for i := 0; i < n; i++ {
		e := evo.NewCreature(env.Rand.Intn(32), env.Rand.Intn(32), float32(i + 10), gender, env.Rand)
		env.Spawn(e)
	}
}

func (env *Environment) Spawn(e evo.Evolver) {
	go e.Evolve()
	env.Evolvers = append(env.Evolvers, e)
	switch obj := e.(type) {
	case *evo.Creature:
		env.Creatures = append(env.Creatures, obj)
	}
}

func (env *Environment) Run() {
	for range env.Clock.C {
		env.Tick++
		t := make([]evo.Evolver, len(env.Evolvers), cap(env.Evolvers))
		copy(t, env.Evolvers)
		removed := 0
		for i, e := range t {
			if e.Alive() {
				e.Pulse() <- env.Tick
			} else {
				env.RemoveEvolver(i-removed)
				removed++
			}
		}
	}
}

func (env *Environment) RemoveEvolver(i int) {
	env.removeRelated(env.Evolvers[i])
	env.Evolvers[i] = nil // or the zero value of T
	env.Evolvers = append(env.Evolvers[:i], env.Evolvers[i+1:]...)
}

func (env *Environment) removeRelated(e evo.Evolver) {
	switch o := e.(type) {
	case *evo.Creature:
		var i int
		var c *evo.Creature
		for i, c = range env.Creatures {
			if c == o {
				break
			}
		}
		env.Creatures[i] = nil // or the zero value of T
		env.Creatures = append(env.Creatures[:i], env.Creatures[i+1:]...)
	}
}

func (env *Environment) Start(c *gin.Context) {
	go env.Run()
}

func (env *Environment) Pause(c *gin.Context) {
	env.Clock.Stop()
	env.Clock = time.NewTicker(env.Speed)
}

func (env *Environment) Reset(c *gin.Context) {
	env.Evolvers = make([]evo.Evolver, 0)
	env.Creatures = make(evo.Creatures, 0)
	env.Tick = 0
	env.Init()
}

func (env *Environment) listCreaturesParseRequest(c *gin.Context) (int, int, int, string, string) {
	c.Request.ParseForm()
	draw, _ := strconv.Atoi(c.Request.Form.Get("draw"))
    start, _ := strconv.Atoi(c.Request.Form.Get("start"))
    length, _ := strconv.Atoi(c.Request.Form.Get("length"))
    column_id := c.Request.Form.Get("order[0][column]")
	direction := c.Request.Form.Get("order[0][dir]")
	param_parts := make([]string, 3)
	param_parts = append(param_parts, "columns[" ,column_id, "][data]")
	param := strings.Join(param_parts, "")
	column := c.Request.Form.Get(param)
	end := start + length
	if end > len(env.Creatures) {
		end = len(env.Creatures)
	}
	return draw, start, end, column, direction
}

func (env *Environment) ListCreatures(c *gin.Context) {
	draw, start, end, column, direction := env.listCreaturesParseRequest(c)
	env.Creatures.Sort(column, direction)

	result := make(map[string]interface{})
	result["data"] = env.Creatures.ToMap(start, end)
	result["status"] = 200
	result["recordsTotal"] = len(env.Creatures)
	result["recordsFiltered"] = len(env.Creatures)
	result["draw"] = draw
	c.JSON(200, result)
}

func main() {
	env := NewEnvironment()
	env.Init()

	r := gin.Default()

	r.GET("/creatures", env.ListCreatures)
	r.GET("/pause", env.Pause)
	r.GET("/start", env.Start)
	r.GET("/reset", env.Reset)

	r.Use(static.Serve("/", static.LocalFile("public", false)))

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
