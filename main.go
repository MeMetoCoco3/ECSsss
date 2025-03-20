package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREENWIDTH = 600
	SCREENHEIGHT
)

func main() {
	rl.InitWindow(SCREENWIDTH, SCREENHEIGHT, "Snake")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	world := NewWorld()
	renderSys := *NewSystem(world, &RenderingSystem{})

	// WARN: OJO CON ESTE ANY!!
	player := make(map[ComponentID]any)
	player[spriteID] = Sprite{
		Width:  32,
		Height: 32,
		Color:  rl.Black,
	}
	player[positionID] = Position{
		X: SCREENWIDTH / 2,
		Y: SCREENHEIGHT / 2,
	}
	world.CreateEntity(player)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		rl.BeginDrawing()
		rl.ClearBackground(VICOLOR)
		renderSys.Update(dt)
		rl.EndDrawing()
	}
}
