package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"os"
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
	renderSys := *NewSystem(world, &DrawSystem{})
	movementSys := *NewSystem(world, &MovementSystem{})

	player := make(map[ComponentID]any)
	player[positionID] = Position{
		X: SCREENWIDTH / 2,
		Y: SCREENHEIGHT / 2,
	}
	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)
	player[movementID] = Movement{VelocityX: 0, VelocityY: 0, Speed: 50}
	player[playerControlledID] = PlayerControlled{}
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

	enemy1 := make(map[ComponentID]any)
	enemy1[positionID] = Position{X: SCREENWIDTH / 2, Y: SCREENHEIGHT / 2}
	enemy1[spriteID] = Sprite{Width: 20, Height: 20, Color: rl.Black}
	enemy1[IAControlledID] = IAControlled{}
	enemy1[movementID] = Movement{VelocityX: 0, VelocityY: 0, Speed: 50}

	world.CreateEntity(player)
	world.CreateEntity(enemy1)

	for i := range world.Query(playerControlledID, movementID, positionID) {
		log.Printf("Entity:%d\n ", i)
	}

	for i := range world.Query(IAControlledID, movementID, positionID) {
		log.Printf("Entity:%d\n ", i)
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		rl.BeginDrawing()
		rl.ClearBackground(VICOLOR)
		movementSys.Update(dt)
		renderSys.Update(dt)
		rl.EndDrawing()
	}
}

func exit() {
	os.Exit(1)
}
