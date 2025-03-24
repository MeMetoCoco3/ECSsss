package main

import (
	"log"

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
	animationID
	playerControlledID
	IAControlledID
	collidesID
	enemyID
	staticID
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

func (c *Movement) Update(dt float32) {

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

// +++++++++++/
type AnimationType int
type AnimationDirection int

const (
	REPEATING AnimationType = iota
	ONESHOT
)
const (
	LEFT  AnimationDirection = -1
	RIGHT AnimationDirection = 1
)

type Animation struct {
	Sprite          Sprite
	First           int
	Last            int
	Current         int
	NumFramesPerRow int
	SizeTile        int
	XPad            int
	YPad            int
	XOffset         int
	YOffset         int
	Direction       AnimationDirection
	Type            AnimationType
	Speed           float32
	Duration_left   float32
}

func (c *Animation) Draw(x, y float32) {
	if c.Duration_left <= 0 {
		c.Duration_left = c.Speed
		c.Current++

		if c.Current > c.Last {
			switch c.Type {
			case REPEATING:
				c.Current = c.First
				break
			case ONESHOT:
				c.Current = c.Last
				break
			}
		}
	}

	rl.DrawTexturePro(
		c.Sprite.Texture,
		c.AnimationFrame(c.NumFramesPerRow, c.SizeTile, c.XPad, c.YPad, c.XOffset, c.YOffset),
		rl.Rectangle{x, y, 128.0, 128.0},
		rl.Vector2{0.0, 0.0}, 0.0, rl.White)

}

func (c *Animation) AnimationFrame(numFramesPerRow, sizeTile, xPad, yPad, xOffset, yOffset int) rl.Rectangle {
	x := int((c.Current % numFramesPerRow) * (sizeTile + xPad))
	y := int((c.Current / numFramesPerRow) * (sizeTile + yPad))
	return rl.Rectangle{
		X:      float32(x + xOffset),
		Y:      float32(y + yOffset),
		Width:  float32(sizeTile),
		Height: float32(sizeTile),
	}
}

// TODO:Booth of this could be a struct with a state, but for now lets leave it like this.
// +++++++++++
type PlayerControlled struct{}

func (c *PlayerControlled) Type() ComponentID { return playerControlledID }

// +++++++++++
type IAControlled struct{}

func (c *IAControlled) Type() ComponentID { return IAControlledID }

// +++++++++++
type Collides struct {
	Width  float32
	Height float32
}

func (c *Collides) Type() ComponentID { return collidesID }

// +++++++++++
type Enemy struct{}

func (c *Enemy) Type() ComponentID { return enemyID }

// +++++++++++
type Static struct{}

func (c *Static) Type() ComponentID { return staticID }

/*
// +++++++++++
type inputReaction uint8

// WARN: Add here extra behaviors for input
const (
	Player inputReaction = iota
)

type Listen struct {
	Behavior inputReaction
}

func (c *Listen) Type() ComponentID {
	return listenID
}

func (c *Listen)Update(dt float32){
  switch c.Behavior{
  case Player:

  }
}
*/

// ===ARCHETYPE===
type Archetype struct {
	Mask          ComponentID
	Entities      []Entity
	Components    map[ComponentID]any
	EntityToIndex map[Entity]int
}

func NewArchetype(componentsID ...ComponentID) *Archetype {
	archetype := &Archetype{
		Mask:          GetMaskFromComponents(componentsID...),
		Entities:      make([]Entity, 0),
		Components:    make(map[ComponentID]any),
		EntityToIndex: make(map[Entity]int),
	}
	for _, currComp := range componentsID {
		archetype.Components[currComp] = GetArrayComponentsFromID(currComp)
	}
	log.Printf("(-)COMPONENTS: %v \n", archetype.Components)
	return archetype
}

func (a *Archetype) AddEntity(entity Entity, components map[ComponentID]any) (idx int) {
	log.Printf("(-)Adding entity n: %d\n", entity)
	idx = len(a.Entities)
	// INFO: Maybe better use commented line instead of append()
	//	    a.Entities[idx] = entity
	a.Entities = append(a.Entities, entity)
	log.Println("Entity added")
	for k, v := range components {
		switch k {

		// HACK: Add components here as needed.
		case positionID:
			positions := a.Components[k].([]Position)
			a.Components[k] = append(positions, v.(Position))
			log.Println("Component added")
		case spriteID:
			sprites := a.Components[k].([]Sprite)
			a.Components[k] = append(sprites, v.(Sprite))
			log.Println("Component added")
		case movementID:
			movements := a.Components[k].([]Movement)
			a.Components[k] = append(movements, v.(Movement))
			log.Println("Component added")
		case healthID:
			health := a.Components[k].([]Health)
			a.Components[k] = append(health, v.(Health))
			log.Println("Component added")
		case aliveID:
			alives := a.Components[k].([]Alive)
			a.Components[k] = append(alives, v.(Alive))
			log.Println("Component added")
		case animationID:
			animations := a.Components[k].([]Animation)
			a.Components[k] = append(animations, v.(Animation))
			log.Println("Component added")
		case playerControlledID:
			pjControlled := a.Components[k].([]PlayerControlled)
			a.Components[k] = append(pjControlled, v.(PlayerControlled))
			log.Println("Component added")
		case IAControlledID:
			iaControlled := a.Components[k].([]IAControlled)
			a.Components[k] = append(iaControlled, v.(IAControlled))
			log.Println("Component added")
		case collidesID:
			collides := a.Components[k].([]Collides)
			a.Components[k] = append(collides, v.(Collides))
			log.Println("Component added")
		case enemyID:
			iaControlled := a.Components[k].([]Enemy)
			a.Components[k] = append(iaControlled, v.(Enemy))
			log.Println("Component added")
		case staticID:
			static := a.Components[k].([]Static)
			a.Components[k] = append(static, v.(Static))
			log.Println("Component added")
		default:
			continue
		}
	}
	a.EntityToIndex[entity] = idx
	return
}

func (a *Archetype) RemoveEntity(entity Entity) {
	idx, exists := a.EntityToIndex[entity]
	if !exists || idx < 0 || idx >= len(a.Entities) {
		return
	}

	lastIdx := len(a.Entities) - 1

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
			case playerControlledID:
				components := v.([]PlayerControlled)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case IAControlledID:
				components := v.([]IAControlled)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case enemyID:
				components := v.([]Enemy)
				components[idx] = components[lastIdx]
				a.Components[k] = components[:lastIdx]
			case staticID:
				components := v.([]Static)
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
			case animationID:
				components := v.([]Animation)
				a.Components[k] = components[:lastIdx]
			case playerControlledID:
				components := v.([]PlayerControlled)
				a.Components[k] = components[:lastIdx]
			case IAControlledID:
				components := v.([]IAControlled)
				a.Components[k] = components[:lastIdx]
			case collidesID:
				components := v.([]Collides)
				a.Components[k] = components[:lastIdx]
			case enemyID:
				components := v.([]Enemy)
				a.Components[k] = components[:lastIdx]
			case staticID:
				components := v.([]Static)
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
		nextEntityID: 0,
		state:        PAUSE,
		entityMask:   make(map[Entity]ComponentID),
		archetypes:   make(map[ComponentID]*Archetype),
	}
}

func (w *World) CreateEntity(components map[ComponentID]any) (entity Entity) {
	entity = w.nextEntityID
	w.nextEntityID++

	var mask ComponentID
	var compList []ComponentID
	for k := range components {
		mask |= GetMaskFromComponents(k)
		compList = append(compList, k)
	}

	// Build archetype if not exists
	archetype, exists := w.archetypes[mask]
	if !exists {
		log.Println("Archetype Created")
		newArchetype := NewArchetype(compList...)
		w.archetypes[mask] = newArchetype
		archetype = newArchetype
	}

	w.entityMask[entity] = mask
	archetype.AddEntity(entity, components)
	return
}

// ...
func (w *World) AddComponent(entity Entity, components map[ComponentID]any) {
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
		case playerControlledID:
			component := v.([]PlayerControlled)[idx]
			components[k] = component
		case IAControlledID:
			component := v.([]IAControlled)[idx]
			components[k] = component
		case collidesID:
			component := v.([]Collides)[idx]
			components[k] = component
		case enemyID:
			component := v.([]Enemy)[idx]
			components[k] = component
		case staticID:
			component := v.([]Enemy)[idx]
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

	components := make(map[ComponentID]any)
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
		case animationID:
			component := v.([]Animation)[idx]
			components[k] = component
		case playerControlledID:
			component := v.([]PlayerControlled)[idx]
			components[k] = component
		case IAControlledID:
			component := v.([]IAControlled)[idx]
			components[k] = component
		case collidesID:
			component := v.([]Collides)[idx]
			components[k] = component
		case enemyID:
			component := v.([]Enemy)[idx]
			components[k] = component
		case staticID:
			component := v.([]Static)[idx]
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

	return hasComponent(mask, component)
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

	return hasComponent(mask, componentsMask)
}

func (w *World) Query(components ...ComponentID) []*Archetype {
	var result []*Archetype
	mask := GetMaskFromComponents(components...)
	log.Println(components)
	log.Println(mask)
	for k, v := range w.archetypes {

		if hasComponent(mask, k) {
			result = append(result, v)
		}
	}

	return result
}

// ===BASE SYSTEM===

type BaseSystem struct {
	World *World
}

func (s *BaseSystem) setWorld(w *World) {
	s.World = w
}

type System interface {
	Update(dt float32)
	setWorld(w *World)
}

func NewSystem[T System](w *World, s T) *T {
	s.setWorld(w)
	return &s
}

// ===SYSTEM===

type MovementSystem struct {
	BaseSystem
}

func (s *MovementSystem) Update(dt float32) {
	archetypes := s.World.Query(positionID, movementID, playerControlledID)
	for archIdx := range archetypes {
		entities := archetypes[archIdx].Entities
		mover := archetypes[archIdx].Components[movementID].([]Movement)

		log.Println("On player")
		for idx := range entities {
			inputX, inputY := GetInput()
			mover[idx].VelocityX = inputX * mover[idx].Speed
			mover[idx].VelocityY = inputY * mover[idx].Speed
		}
	}

	archetypes = s.World.Query(positionID, movementID, IAControlledID)
	for archIdx := range archetypes {
		log.Println("On Enemy")
		entities := archetypes[archIdx].Entities
		mover := archetypes[archIdx].Components[movementID].([]Movement)

		for idx := range entities {
			// TODO: Define AI Behavior
			mover[idx].VelocityY = 1 * mover[idx].Speed
		}
	}

	archetypes = s.World.Query(positionID, movementID)
	for archIdx := range archetypes {
		entities := archetypes[archIdx].Entities
		position := archetypes[archIdx].Components[positionID].([]Position)
		mover := archetypes[archIdx].Components[movementID].([]Movement)

		for idx := range entities {
			position[idx].X += mover[idx].VelocityX * dt
			position[idx].Y += mover[idx].VelocityY * dt
		}
	}
}

// +++++++++++
type DrawSystem struct {
	BaseSystem
}

func (s *DrawSystem) Update(dt float32) {
	// Sprite
	archetypes := s.World.Query(positionID, spriteID)
	for archIdx := range archetypes {
		entities := archetypes[archIdx].Entities
		position := archetypes[archIdx].Components[positionID].([]Position)
		sprite := archetypes[archIdx].Components[spriteID].([]Sprite)

		for idx := range entities {
			// WARN: IM not sure it is necessary to deal with entIdx, i believe that id is important if you want to
			// deal with an exact entitie, but in our case, we want to update everything
			// So i belive this could be substituted by a plain range
			// entityID := entities[entIdx]
			// log.Printf("(-)Drawing Entity: %d\n", entityID)
			// sprite[entityID].Draw(position[entityID].X, position[entityID].Y)
			sprite[idx].Draw(position[idx].X, position[idx].Y)
		}
	}

	// Animation
	archetypes = s.World.Query(positionID, animationID)
	for archIdx := range archetypes {
		entities := archetypes[archIdx].Entities
		position := archetypes[archIdx].Components[positionID].([]Position)
		animation := archetypes[archIdx].Components[animationID].([]Animation)

		for idx := range entities {
			animation[idx].Duration_left -= dt
			animation[idx].Draw(position[idx].X, position[idx].Y)
		}
	}
}

// +++++++++++
type CollisionSystem struct {
	BaseSystem
}

func (s *CollisionSystem) Update(dt float32) {
	archetypes := s.World.Query(positionID, collidesID)
	for i := range archetypes {
		entitiesA := archetypes[i].Entities
		positionA := archetypes[i].Components[positionID].([]Position)
		colliderA := archetypes[i].Components[collidesID].([]Collides)

		movementA, isMovingA := archetypes[i].Components[movementID].([]Movement)
		for j := range archetypes {
			entitiesB := archetypes[j].Entities
			positionB := archetypes[j].Components[positionID].([]Position)
			colliderB := archetypes[j].Components[collidesID].([]Collides)

			movementB, isMovingB := archetypes[j].Components[movementID].([]Movement)
			for idxA := range entitiesA {
				for idxB := range entitiesB {
					if archetypes[i] == archetypes[j] && idxA == idxB {
						continue
					}
					switch CheckRectCollision(positionA[idxA], colliderA[idxA], positionB[idxB], colliderB[idxA]) {
					case noC:
						continue
					case topC:
						if isMovingA {
							movementA[idxA].VelocityY = 0
						}
						if isMovingB {
							movementB[idxB].VelocityY = 0
						}
						log.Println("Custom collision for movement objects")
					case bottomC:
						if isMovingA {
							movementA[idxA].VelocityY = 0
						}
						if isMovingB {
							movementB[idxB].VelocityY = 0
						}
						log.Println("Custom collision for movement objects")
					case leftC:
						if isMovingA {
							movementA[idxA].VelocityX = 0
						}
						if isMovingB {
							movementB[idxB].VelocityX = 0
						}
						log.Println("Custom collision for movement objects")
					case rightC:
						if isMovingA {
							movementA[idxA].VelocityX = 0
						}
						if isMovingB {
							movementB[idxB].VelocityX = 0
						}
						log.Println("Custom collision for movement objects")
					}
				}

			}
		}
	}
}
