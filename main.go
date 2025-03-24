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
	collisionSys := *NewSystem(world, &CollisionSystem{})
	player := make(map[ComponentID]any)
	player[positionID] = Position{
		X: 40,
		Y: 40,
	}

	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)
	player[movementID] = Movement{VelocityX: 0, VelocityY: 0, Speed: 500}
	player[collidesID] = Collides{X: player[positionID].(Position).X, Y: player[positionID].(Position).Y, Width: 128, Height: 128}
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
	enemy1[spriteID] = Sprite{Width: 200, Height: 200, Color: rl.Black}
	enemy1[IAControlledID] = IAControlled{}
	enemy1[movementID] = Movement{VelocityX: 0, VelocityY: 0, Speed: 50}
	enemy1[collidesID] = Collides{Width: 200, Height: 200}
	world.CreateEntity(player)
	world.CreateEntity(enemy1)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		collisionSys.Update(dt)
		movementSys.Update(dt)

		log.Println(world.nextEntityID)
		os.Exit(1)
		rl.BeginDrawing()
		rl.ClearBackground(VICOLOR)
		renderSys.Update(dt)
		rl.EndDrawing()
	}
}

func exit() {
	os.Exit(1)
}
