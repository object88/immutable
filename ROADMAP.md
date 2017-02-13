## Goals
Create a strongly typed immutable collections package without using generics.

Example:
``` Go
x := NewStringToIntegerMap()
x, _ = x.Set('a', 1)
x, _ = x.Set('b', 2)
y, ok := x.Get('a')
if !ok {
  // Error!
}
```


`Hashmap<K, V>`

Common methods:
* `New[K]To[V]Hashmap() (result *[K]To[V]Hashmap)`
* `Get(key K) (result V)`
* `Set(key K, value V)`
* `Insert(key K, value V)`
* `Remove(key K)`
* `ForEach(func (key K, value V))`
* `Map(func (key K, value V) value V2) (result Hashmap<K, V2>)`
* `Reduce(func (key K, value V), acc V2) (result V2)`

### Example:
_`StringToIntegerHashmap` is a string to integer hashmap_

``` Go
type InternalHashmap

type StringToIntHashmap struct {
  GetHash func(key String) uint64

  i InternalHashmap
}
```

`NewStringToIntegerMap() (hashmap *StringToIntegerHashmap)`

``` Go
(h *StringToIntHashmap) Get(key string) (result int) {
  hash := h.GetHash(key)
}
```

`(h *StringToIntHashmap) Set(key string, value int)`

`(h *StringToIntHashmap) Insert(key string, value int)`

`(h *StringToIntHashmap) Remove(key string)`

`(h *StringToIntHashmap) ForEach(func (key string, value int))`

`(h *StringToIntHashmap) Map(func (key string, value int) value string) (result StringToStringHashmap)`

`(h *StringToIntHashmap) Reduce(func (key string, value int), acc bool) (result bool)`


### Methodology
Implement a composable core hashmap library which operates on `Element`.


Concrete implementations will live in own namespaces e.g.:
`https://github.com/object88/immutable/hashmap/stringtointeger`
`https://github.com/object88/immutable/hashmap/stringtostring`
