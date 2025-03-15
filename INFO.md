# SNAKE INVADERS
Game created using a ECS with archetypes to define de different types of entities.
It is written in Golang using Raylib.

#### Archetypes ID
The archetype ID is a uint64 object, which we will use to mask and check which kind 
of entities we are working with. This are the values that are asigned to each component.
- 1 = Position
- 2 = Sprite
- 4 = Movement
- 8 = Health
- 16 = Alive

