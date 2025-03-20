package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// An entity is a index of a object in the whole world.
type Entity uint32
type ComponentID uint32
type State uint32

const (
	positionID ComponentID = 1 << iota
	spriteID
	movementID
	healthID
	aliveID
)

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

// +++++++++++/
type Alive struct {
	IsAlive bool
}

func (c *Alive) Type() ComponentID {
	return aliveID
}

// ===ARCHETYPE===
type Archetype struct {
	Mask          ComponentID
	Entities      []Entity
	Components    map[ComponentID]interface{}
	EntityToIndex map[Entity]int
}

func NewArchetype(componentsID ...ComponentID) *Archetype {
	archetype := &Archetype{
		Mask:       GetMaskFromComponents(componentsID...),
		Entities:   []Entity{},
		Components: make(map[ComponentID]interface{}),
	}
	for _, currComp := range componentsID {
		archetype.Components[currComp] = GetComponentFromID(currComp)
	}
	return archetype
}

func (a *Archetype) AddEntity(entity Entity, components map[ComponentID]interface{}) (idx int) {
	idx = len(a.Entities)
	// INFO: Maybe better use commented line instead of append()
	//	    a.Entities[idx] = entity
	a.Entities = append(a.Entities, entity)

	for k, v := range components {
		switch k {

		// HACK: Add components here as needed.
		case positionID:
			positions := a.Components[k].([]Position)
			a.Components[k] = append(positions, v.(Position))
		case spriteID:
			sprites := a.Components[k].([]Sprite)
			a.Components[k] = append(sprites, v.(Sprite))
		case movementID:
			movements := a.Components[k].([]Movement)
			a.Components[k] = append(movements, v.(Movement))
		case healthID:
			health := a.Components[k].([]Health)
			a.Components[k] = append(health, v.(Health))
		case aliveID:
			alives := a.Components[k].([]Alive)
			a.Components[k] = append(alives, v.(Alive))
		default:
			continue
		}
	}

	a.EntityToIndex[entity] = idx
	return idx
}

func (a *Archetype) RemoveEntity(entity Entity) {
	idx, exists := a.EntityToIndex[entity]
	if !exists || idx < 0 || idx >= len(a.Entities) {
		return
	}

	lastIdx := len(a.Entities) - 1

	// Si no es la ultima, cogemos la ultima entidad y la swapeamos con la que
	// queremos eliminar.
	if idx != lastIdx {
		lastEntity := a.Entities[lastIdx]
		a.Entities[idx] = lastEntity
		a.EntityToIndex[lastEntity] = idx

		for k, v := range a.Components {
			// HACK: Add components here as needed.
			switch k {
			case positionID:
				components := v.([]Position)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case spriteID:
				components := v.([]Sprite)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case movementID:
				components := v.([]Movement)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case healthID:
				components := v.([]Health)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case aliveID:
				components := v.([]Alive)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			default:
				continue
			}
		}
	} else {
		for k, v := range a.Components {

			// HACK: Add components here as needed.
			switch k {
			case positionID:
				components := v.([]Position)
				a.Components[k] = components[:lastIdx]
			case spriteID:
				components := v.([]Sprite)
				a.Components[k] = components[:lastIdx]
			case movementID:
				components := v.([]Movement)
				a.Components[k] = components[:lastIdx]
			case healthID:
				components := v.([]Health)
				a.Components[k] = components[:lastIdx]
			case aliveID:
				components := v.([]Alive)
				a.Components[k] = components[:lastIdx]
			default:
				continue
			}
		}
	}
	a.Entities = a.Entities[:lastIdx]
	delete(a.EntityToIndex, entity)
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

func (w *World) CreateEntity(components map[ComponentID]interface{}) (entity Entity) {
	entity = w.nextEntityID
	w.nextEntityID++

	var mask ComponentID
	for k, _ := range components {
		mask |= GetMaskFromComponents(k)
	}

	// Build archetype if not exists
	archetype, exists := w.archetypes[mask]
	if !exists {
		newArchetype := NewArchetype(mask)
		w.archetypes[mask] = newArchetype
	}

	w.entityMask[entity] = mask
	archetype.AddEntity(entity, components)
	return
}

func (w *World) AddComponent(entity Entity, components map[ComponentID]interface{}) {
	mask, ok := w.entityMask[entity]
	if !ok {
		// Si no existe, la creamos.
		w.CreateEntity(components)
		return
	}
	oldArchetype := w.archetypes[mask]
	idx := oldArchetype.EntityToIndex[entity]

	for k, v := range oldArchetype.Components {
		// HACK: Add components here as needed.
		switch k {
		case positionID:
			component := v.([]Position)[idx]
			components[k] = component
		case spriteID:
			component := v.([]Sprite)[idx]
			components[k] = component
		case movementID:
			component := v.([]Movement)[idx]
			components[k] = component
		case healthID:
			component := v.([]Health)[idx]
			components[k] = component
		case aliveID:
			component := v.([]Alive)[idx]
			components[k] = component
		default:
			continue
		}
	}
	w.CreateEntity(components)
	oldArchetype.RemoveEntity(entity)
	w.nextEntityID--
}

func (w *World) RemoveComponent(entity Entity, component ComponentID) {
	mask, ok := w.entityMask[entity]
	if !ok || mask&component == 0 {
		return
	}

	var components map[ComponentID]interface{}
	oldArchetype := w.archetypes[mask]
	idx := oldArchetype.EntityToIndex[entity]

	for k, v := range oldArchetype.Components {
		if k == component {
			continue
		}
		// HACK: Add components here as needed.
		switch k {
		case positionID:
			component := v.([]Position)[idx]
			components[k] = component
		case spriteID:
			component := v.([]Sprite)[idx]
			components[k] = component
		case movementID:
			component := v.([]Movement)[idx]
			components[k] = component
		case healthID:
			component := v.([]Health)[idx]
			components[k] = component
		case aliveID:
			component := v.([]Alive)[idx]
			components[k] = component
		default:
			continue
		}
	}
	oldArchetype.RemoveEntity(entity)

	mask = mask ^ component
	w.entityMask[entity] = mask
	w.nextEntityID--
	w.CreateEntity(components)
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
// Los sistemas son la funcionalidad de los componentes , cada arquetipo guardara en una array componentes independientemente
// estos componentes estaran alineados con su correspondiente entidas gracias a la array de entidades del arquetipo.
// WARN: Cuando trabajamos con el field mask, lo lamamos componentID, pero realmente es el conjunto de varios formando una mascara, cambiar el tipo.

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
