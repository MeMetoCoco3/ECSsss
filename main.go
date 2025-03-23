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
	renderSys := *NewSystem(world, &AnimationSystem{})

	player := make(map[ComponentID]any)
	player[positionID] = Position{
		X: SCREENWIDTH / 2,
		Y: SCREENHEIGHT / 2,
	}
	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)

	player[animationID] = Animation{
		First:         0,
		Last:          13,
		Current:       0,
		Speed:         0.1,
		Duration_left: 0.1,
		Type:          REPEATING,

		NumFramesPerRow: 14,
		SizeTile:        64,
		XPad:            32, YPad: 0, XOffset: 16, YOffset: 8,
		Sprite: Sprite{
			Texture: pjTexture,
			Width:   32,
			Height:  32,
			Color:   rl.White,
		},
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
