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

func GetComponentsFromMask(mask ComponentID) (components []ComponentID) {
	components = make([]ComponentID, 32)
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

func hasComponents(mask, componentsMask ComponentID) bool {
	return (uint32(mask) & uint32(componentsMask)) == uint32(componentsMask)
}
