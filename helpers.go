package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var VICOLOR = rl.Color{252, 163, 17, 255}

func GetMaskFromComponents(componentsID ...ComponentID) ComponentID {
	var mask ComponentID
	for i := range len(componentsID) {
		mask |= componentsID[i]
	}
	return mask
}

// HACK: Add components here as needed.
func GetArrayComponentsFromID(id ComponentID) any {
	switch id {
	case positionID:
		return make([]Position, 0)
	case spriteID:
		return make([]Sprite, 0)
	case movementID:
		return make([]Movement, 0)
	case healthID:
		return make([]Health, 0)
	case aliveID:
		return make([]Alive, 0)
	case animationID:
		return make([]Animation, 0)
	case playerControlledID:
		return make([]PlayerControlled, 0)
	case IAControlledID:
		return make([]IAControlled, 0)
	case collidesID:
		return make([]Collides, 0)
	case enemyID:
		return make([]Enemy, 0)
	default:
		return nil
	}
}

func GetComponentsFromMask(mask ComponentID) []ComponentID {
	components := make([]ComponentID, 32)
	count := 0
	bitValue := 1
	for mask != 0 {
		bit := mask & 1
		if bit != 0 {
			components[count] = ComponentID(bitValue)
			count++
		}
		bitValue = bitValue << 1
		mask = mask >> 1
	}
	return components[:count]
}

func hasComponent(mask, componentsMask ComponentID) bool {
	newMask := (uint32(mask) & uint32(componentsMask))
	result := newMask == uint32(mask)
	return result
}

func GetInput() (x float32, y float32) {
	if rl.IsKeyDown(rl.KeyLeft) {
		x = -1
	} else if rl.IsKeyDown(rl.KeyRight) {
		x = 1
	}
	if rl.IsKeyDown(rl.KeyUp) {
		y = -1
	} else if rl.IsKeyDown(rl.KeyDown) {
		y = 1
	}
	return
}

type collisionType uint8

const (
	noC collisionType = iota
	topC
	bottomC
	rightC
	leftC
)

func CheckRectCollision(aPos Position, aSize Collides, bPos Position, bSize Collides) collisionType {
	aTop := aPos.Y
	aBottom := aPos.Y + aSize.Height
	aRight := aPos.X + aSize.Width
	aLeft := aPos.X
	bTop := bPos.Y
	bBottom := bPos.Y + bSize.Height
	bRight := bPos.X + bSize.Width
	bLeft := bPos.X

	// A Bottom collision means, the bottom of A hits the top of B.
	// Bottom collision
	if aBottom > bTop && aTop < bTop &&
		aRight > bLeft && aLeft < bRight {
		return bottomC
	}

	// Top collision
	if aTop < bBottom && aBottom > bBottom &&
		aRight > bLeft && aLeft < bRight {
		return topC
	}

	// Right collision
	if aRight > bLeft && aLeft < bLeft &&
		aBottom > bTop && aTop < bBottom {
		return rightC
	}

	// Left collision
	if aLeft < bRight && aRight > bRight &&
		aBottom > bTop && aTop < bBottom {
		return leftC
	}

	return noC
	/*
		// Additional check for full overlap
		if aLeft <= bLeft && aRight >= bRight &&
			aTop <= bTop && aBottom >= bBottom {
			log.Println("A contains B")
		}

		if bLeft <= aLeft && bRight >= aRight &&
			bTop <= aTop && bBottom >= aBottom {
			log.Println("B contains A")
		}
	*/

}
