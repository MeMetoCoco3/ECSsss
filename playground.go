package main

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	/*
		for i := ArchetypeID(1); i <= 2048; i = i << 1 {

			if i&512 == 512 {
				fmt.Printf("The number %d is represented %016b\n", i, i)
			} else {
				fmt.Println(i)
			}
		}
	*/
	fmt.Printf("Size of rl.Color: %d bytes\n", unsafe.Sizeof(rl.Color{}))
	fmt.Printf("Size of rl.Texture2D: %d bytes\n", unsafe.Sizeof(rl.Texture2D{}))
	fmt.Printf("Size of Sprite: %d bytes\n", unsafe.Sizeof(Sprite{}))

	fmt.Printf("Value of a not defined texture2d: %v\n", Sprite{}.Texture)

	fmt.Println(NewArchetype([]ComponentID{movementID, healthID}))
	fmt.Println(NewArchetype([]ComponentID{positionID, spriteID, movementID, healthID}))

	e1 := GetMaskFromComponents([]ComponentID{positionID, spriteID, movementID, healthID})
	fmt.Println(uint32(e1))
	fmt.Println(GetComponentsFromMask(uint32(e1)))
}
