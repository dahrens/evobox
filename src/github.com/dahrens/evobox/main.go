package main

import (
	evo "github.com/dahrens/evobox/evolution"
	"github.com/gin-gonic/gin"
	"time"
	"strconv"
	"log"
	"strings"
)

type Environment struct {
	Evolvers  []evo.Evolver
	Creatures evo.Creatures
	Clock     *time.Ticker
	Tick      int
	Speed	  time.Duration
}

func NewEnvironment() *Environment {
	e := new(Environment)
	e.Evolvers = make([]evo.Evolver, 0)
	e.Creatures = make(evo.Creatures, 0)
	e.Clock = time.NewTicker(time.Second)
	e.Tick = 0
	e.Speed = time.Second
	return e
}

func (env *Environment) Init() {
	go env.SpawnMany(100, evo.GENDER_MALE)
	go env.SpawnMany(100, evo.GENDER_FEMALE)
}

func (env *Environment) SpawnMany(n int, gender evo.Gender) {
	for i := 0; i < n; i++ {
		e := evo.NewCreature(0, 0, float32(i + 10), gender)
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
		for i, e := range env.Evolvers {
			if e.Alive() {
				e.Pulse() <- env.Tick
			} else {
				env.RemoveEvolver(i)
			}
		}
	}
}

func (env *Environment) RemoveEvolver(i int) {
	go env.removeRelated(env.Evolvers[i])
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

func (env *Environment) SortCreatures(c *gin.Context) {
	column_id := c.Request.Form.Get("order[0][column]")
	direction := c.Request.Form.Get("order[0][dir]")
	param_parts := make([]string, 3)
	param_parts = append(param_parts, "columns[" ,column_id, "][data]")
	param := strings.Join(param_parts, "")
	column := c.Request.Form.Get(param)
	env.Creatures.Sort(column, direction)
}

func (env *Environment) ListCreatures(c *gin.Context) {
	result := make(map[string]interface{})

	c.Request.ParseForm()
	env.SortCreatures(c)

	draw, err := strconv.Atoi(c.Request.Form.Get("draw"))
    start, err := strconv.Atoi(c.Request.Form.Get("start"))
    length, err := strconv.Atoi(c.Request.Form.Get("length"))
    if err != nil { log.Println(err)}

	end := start + length
	if end > len(env.Creatures) {
		end = len(env.Creatures)
	}

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
	r.LoadHTMLGlob("templates/*")

	r.GET("/creatures", env.ListCreatures)
	r.GET("/index", func(c *gin.Context) {
		obj := gin.H{"title": "Main website"}
		c.HTML(200, "index.tpl", obj)
	})
	r.GET("/pause", env.Pause)
	r.GET("/start", env.Start)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
