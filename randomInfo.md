# Conocimiento random.
Aqui voy a reunir cosas que he ido aprendiendo que por lo que sea he ido borrando del codigo.

## Sort
- Problema: No poder ordenar una array de in32.
Podemos usar sort.Sort() con cualquier tipo si definimos less, len, y swap.
PERO lo que de verdad mola es que podemos conseguir lo mismo asi (sorted in place y se conserva el type):
```go
    sort.Sort(arr, func(i,j int)bool{arr[i]<arr[j]})
    fmt.Printf("%T", arr[0]) // prints: int32
```

## Interfaces
- Problema: Como inicializar arrays de Componentes con una funcion.
El problema aparece cuando decimos que queremos devolver una interfaz y lo que vamos a devolver es un slice.
Golang no convierte automaticamente Slices de un datatype a slices de interfaces.
La solucion es decir que devolvemos una unica interfaz.
```go
// Bad!!
func GetComponentFromID(id ComponentID) []interface{} {
	return make([]Position, 0)

// Good!!
func GetComponentFromID(id ComponentID) interface{} {
	return make([]Position, 0)

```

## Definir sistemas
- Problema: Tenemos sistemas que hacen cosas similares O dependen de otros sistemas. En este caso la
  idea era que el sistema ANIMACION cambiara el rectangulo que tenia el SPRITE para animar. Lo cual
  es mucho mas complejo y complicado que hacer dos sistemas diferentes, uno para entidades
  estaticas, y otro para entidades animadas.

## Epsilon
Es la unidad de computacion mas pequena, varia entre ordenadores. En este caso he usado una variable
epsilon para definir una distancia minuscula entre dos entidades en mi sistema de colisiones. 
