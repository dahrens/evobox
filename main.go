package main

import (
	evo "github.com/dahrens/evobox/evolution"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Environment struct {
	World *evo.World
}

func NewEnvironment() *Environment {
	env := new(Environment)
	env.World = evo.NewWorld(32, 32)
	env.World.Init()
	return env
}

func (env *Environment) Start(c *gin.Context) {
	go env.World.Run()
}

func (env *Environment) Pause(c *gin.Context) {
	go env.World.Pause()
}

func (env *Environment) Reset(c *gin.Context) {
	env.World = evo.NewWorld(32, 32)
	env.World.Init()
}

func (env *Environment) listCreaturesParseRequest(c *gin.Context) (int, int, int, string, string) {
	c.Request.ParseForm()
	draw, _ := strconv.Atoi(c.Request.Form.Get("draw"))
	start, _ := strconv.Atoi(c.Request.Form.Get("start"))
	length, _ := strconv.Atoi(c.Request.Form.Get("length"))
	column_id := c.Request.Form.Get("order[0][column]")
	direction := c.Request.Form.Get("order[0][dir]")
	param_parts := make([]string, 3)
	param_parts = append(param_parts, "columns[", column_id, "][data]")
	param := strings.Join(param_parts, "")
	column := c.Request.Form.Get(param)
	end := start + length
	if end > len(env.World.Creatures) {
		end = len(env.World.Creatures)
	}
	return draw, start, end, column, direction
}

func (env *Environment) ListCreatures(c *gin.Context) {
	draw, start, end, column, direction := env.listCreaturesParseRequest(c)
	env.World.Creatures.Sort(column, direction)

	result := make(map[string]interface{})
	result["data"] = env.World.Creatures.ToMap(start, end)
	result["status"] = 200
	result["recordsTotal"] = len(env.World.Creatures)
	result["recordsFiltered"] = len(env.World.Creatures)
	result["draw"] = draw
	c.JSON(200, result)
}

func main() {
	env := NewEnvironment()

	r := gin.Default()

	r.GET("/creatures", env.ListCreatures)
	r.GET("/pause", env.Pause)
	r.GET("/start", env.Start)
	r.GET("/reset", env.Reset)

	r.Use(static.Serve("/", static.LocalFile("assets", false)))

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
