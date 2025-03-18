package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// An entity is a index of a object in the whole world.
type Entity uint32
type Entities []Entity
type ComponentID uint32
type State uint32

const (
	positionID ComponentID = 1 << iota
	spriteID
	movementID
	healthID
	aliveID
)

/*
	func GetComponent(id ComponentID) Component {
		switch id {
		case positionID:
			return &Position{}
		case spriteID:
			return &Sprite{}
		case movementID:
			return &Movement{}
		case healthID:
			return &Health{}
		case aliveID:
			return &Alive{}
		default:
			return nil
		}
	}
*/
const (
	PAUSE State = iota
	PLAY
	MENU
)

type Component interface {
	Type() ComponentID
}

// ===COMPONENTS===
type Position struct {
	X float32
	Y float32
}

func (c *Position) Type() ComponentID {
	return positionID
}

// +++++++++++

type Sprite struct {
	Width   float32
	Height  float32
	Texture rl.Texture2D
	Color   rl.Color
}

func (c *Sprite) Type() ComponentID {
	return spriteID
}
func (c *Sprite) Draw(x, y float32) {
	// Texture is not declared, we draw a rectangle
	if c.Texture.ID == 0 {
		rl.DrawRectangle(int32(x), int32(y), int32(c.Width), int32(c.Height), c.Color)
		return
	}
	rl.DrawTexture(c.Texture, int32(x), int32(y), c.Color)
}

// +++++++++++
type Movement struct {
	VelocityX float32
	VelocityY float32
	Speed     float32
}

func (c *Movement) Type() ComponentID {
	return movementID
}

// +++++++++++
type Health struct {
	Max     int32
	Current int32
}

func (c *Health) Type() ComponentID {
	return healthID
}

// +++++++++++
type Alive struct {
	IsAlive bool
}

func (c *Alive) Type() ComponentID {
	return aliveID
}

// ===ARCHETYPE===
type Archetype struct {
	Mask          ComponentID
	Entities      Entities
	Components    map[ComponentID][]interface{}
	EntityToIndex map[Entity]int
}

func NewArchetype(componentsID ...ComponentID) *Archetype {
	componentMap := make(map[ComponentID][]interface{})
	for _, currComp := range componentsID {
		componentMap[currComp] = make([]interface{}, 0)
	}
	mask := GetMaskFromComponents(componentsID...)
	return &Archetype{
		Mask:       mask,
		Entities:   Entities{},
		Components: componentMap,
	}
}

func (a *Archetype) AddEntity(entity Entity) (idx int) {
	idx = len(a.Entities) - 1
	a.Entities[idx] = entity
	return
}

func (a *Archetype) RemoveEntity(entity Entity) {
	if entity < 0 || int(entity) >= len(a.Entities) {
		return
	}
	lastIndex := len(a.Entities) - 1
	a.Entities[entity] = a.Entities[lastIndex]
	a.Entities = a.Entities[:lastIndex]
}

// ===WORLD===
type World struct {
	nextEntityID Entity
	state        State
	entityMask   map[Entity]ComponentID
	archetypes   map[ComponentID]*Archetype
}

func NewWorld() *World {
	return &World{
		nextEntityID: 1,
		state:        PAUSE,
		entityMask:   make(map[Entity]ComponentID),
		archetypes:   make(map[ComponentID]*Archetype),
	}
}

func (w *World) CreateEntity(components ...ComponentID) (entity Entity) {
	entity = w.nextEntityID
	w.nextEntityID++

	// Build archetype if not exists
	mask := GetMaskFromComponents(components...)
	archetype, exists := w.archetypes[mask]
	if !exists {
		newArchetype := NewArchetype(components...)
		w.archetypes[mask] = newArchetype
	}

	w.entityMask[entity] = mask
	archetype.AddEntity(entity)
	return
}

func (w *World) AddComponent(entity Entity, component ComponentID) {
	mask, ok := w.entityMask[entity]
	if !ok {
		// Si no existe, la creamos.
		w.CreateEntity(component)
		return
	}
	oldArchetype := w.archetypes[mask]
	oldArchetype.RemoveEntity(entity)

	components := GetComponentsFromMask(mask & component)
	w.CreateEntity(components...)

	w.nextEntityID--
}

func (w *World) RemoveComponent(entity Entity, component ComponentID) {
	mask, ok := w.entityMask[entity]
	if !ok || mask&component == 0 {
		return
	}

	oldArchetype := w.archetypes[mask]
	oldArchetype.RemoveEntity(entity)

	mask = mask ^ component
	w.entityMask[entity] = mask
	w.nextEntityID--
	w.CreateEntity(mask)
}

func (w *World) RemoveEntity(entity Entity) {
	mask, ok := w.entityMask[entity]
	if !ok {
		return
	}

	entityArchetype := w.archetypes[mask]
	entityArchetype.RemoveEntity(entity)

	delete(w.entityMask, entity)
	w.nextEntityID--

}

func (w *World) HasComponent(entity Entity, component ComponentID) bool {
	mask, ok := w.entityMask[entity]
	if !ok {
		return false
	}

	return hasComponents(mask, component)
}

func (w *World) HasComponents(entity Entity, components ...ComponentID) bool {
	mask, ok := w.entityMask[entity]
	if !ok {
		return false
	}
	var componentsMask ComponentID

	for i := range components {
		componentsMask |= components[i]
	}

	return hasComponents(mask, componentsMask)
}

func (w *World) Query(components ...ComponentID) []*Archetype {
	var result []*Archetype
	mask := GetMaskFromComponents(components...)

	for k, v := range w.archetypes {
		if hasComponents(mask, k) {
			result = append(result, v)
		}
	}

	return result
}

// ===BASE SYSTEM===

type BaseSystem struct {
	World *World
}

func (s *BaseSystem) SetWorld(w World) {
	s.World = &w
}

type System interface {
	Update(dt float32)
	SetWorld(w *World)
}

func NewSystem[T System](w *World, s T) *T {
	s.SetWorld(w)
	return &s
}

// ===SYSTEM===

type MovementSystem struct {
	BaseSystem
	world *World
}

func (s *MovementSystem) Update(dt float32) {
	//...
}

type RenderingSystem struct {
	BaseSystem
	world *World
}

func (s *RenderingSystem) Update(dt float32) {
	s.Draw()
}

func (s *RenderingSystem) Draw() {
	archetypes := s.World.Query(positionID, spriteID)
	for archIdx := range archetypes {
		entities := archetypes[archIdx].Entities
		for entIdx := range entities {
			entities[entIdx].Draw()

		}
	}
}
