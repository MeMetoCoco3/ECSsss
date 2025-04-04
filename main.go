package main

import (
	"log"
	"os"

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
	renderSys := *NewSystem(world, &DrawSystem{})
	movementSys := *NewSystem(world, &MovementSystem{})
	collisionSys := *NewSystem(world, &CollisionSystem{})
	pjTexture := rl.LoadTexture("assets/player/fishy.png")
	defer rl.UnloadTexture(pjTexture)

	player := make(map[ComponentID]any)
	player[positionID] = Position{
		X: 200,
		Y: 200,
	}
	player[movementID] = Movement{VelocityX: 0, VelocityY: 0, Speed: 500, Grounded: true}
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
	// squareCenter := make(map[ComponentID]any)
	// squareCenter[positionID] = Position{X: SCREENWIDTH / 2, Y: SCREENHEIGHT / 2}
	// squareCenter[spriteID] = Sprite{Width: 200, Height: 200, Color: rl.Red}
	// squareCenter[collidesID] = Collides{X: SCREENWIDTH / 2, Y: SCREENHEIGHT / 2, Width: 200, Height: 200}
	border1 := make(map[ComponentID]any)
	border2 := make(map[ComponentID]any)
	border3 := make(map[ComponentID]any)
	border4 := make(map[ComponentID]any)
	border1[positionID] = Position{X: 0, Y: 0}
	border1[spriteID] = Sprite{Width: SCREENWIDTH, Height: 100, Color: rl.Red}
	border1[collidesID] = Collides{X: 0, Y: 0, Width: SCREENWIDTH, Height: 100}
	border2[positionID] = Position{X: SCREENWIDTH - 20, Y: 0}
	border2[spriteID] = Sprite{Width: 100, Height: SCREENHEIGHT, Color: rl.Red}
	border2[collidesID] = Collides{X: SCREENWIDTH - 20, Y: 0, Width: 100, Height: SCREENHEIGHT}
	border3[positionID] = Position{X: 0, Y: 0}
	border3[spriteID] = Sprite{Width: 100, Height: SCREENHEIGHT, Color: rl.Red}
	border3[collidesID] = Collides{X: 0, Y: 0, Width: 100, Height: SCREENHEIGHT}
	border4[positionID] = Position{X: 0, Y: SCREENHEIGHT - 20}
	border4[spriteID] = Sprite{Width: SCREENWIDTH, Height: 100, Color: rl.Red}
	border4[collidesID] = Collides{X: 0, Y: SCREENHEIGHT - 20, Width: SCREENWIDTH, Height: 100}
	world.CreateEntity(player)
	// world.CreateEntity(squareCenter)
	world.CreateEntity(border1)
	world.CreateEntity(border2)
	world.CreateEntity(border3)
	world.CreateEntity(border4)
	//
	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		movementSys.Update(dt)
		collisionSys.Update(dt)
		log.Println(world.nextEntityID)
		rl.BeginDrawing()
		rl.ClearBackground(VICOLOR)
		renderSys.Update(dt)
		rl.EndDrawing()
	}
}

func exit() {
	os.Exit(1)
}
