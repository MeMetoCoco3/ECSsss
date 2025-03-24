# SNAKE INVADERS
Game created using a ECS with archetypes to define de different types of entities.
It is written in Golang using Raylib.

## Archetypes ID
The archetype ID is a uint64 object, which we will use to mask and check which kind 
of entities we are working with. This are the values that are asigned to each component.
- 1 = Position
- 2 = Sprite
- 4 = Movement
- 8 = Health
- 16 = Alive
- 32 = Animation
- 64 = PlayerControlled
- 128 = IAControlled
- 256 = Collides
- 512 = Enemy
- 1028 = Static

## Systems
I do not know if this is intrinsec from the ECS, but in this case, systems are isolated, so the functionality
of each cannot depend from the other. I realized this while trying to do an animate system and a sprite system, being the first
one just something that changes the square that is going to be changed, but for that, we need to query 3 components instead of 2.
Which would not be a problem, except for the fact that then i need to query the rest of objects with 2 components to draw everything
and then, i will also redraw the things that where drawn with 3 components.

We can use tagstructs, to define custom behavior on some of this systems. Some of this systems can
query multiple times to define different behavior for similar functionality(Ej drawing animation and sprite)
