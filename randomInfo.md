# Conocimiento random.
Aqui voy a reunir cosas que he ido aprendiendo que por lo que sea he ido borrando del codigo.

##### Sort
- Problema: No poder ordenar una array de in32.
Podemos usar sort.Sort() con cualquier tipo si definimos less, len, y swap.
PERO lo que de verdad mola es que podemos conseguir lo mismo asi (sorted in place y se conserva el type):
```go
    sort.Sort(arr, func(i,j int)bool{arr[i]<arr[j]})
    fmt.Printf("%T", arr[0]) // prints: int32
```




