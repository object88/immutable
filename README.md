# immutable
#### A Go / Golang library for immutable collections

_Note:_ This library is under development and is not meant for public consumption at this time.

In general, this library seeks to avoid errors due to nil pointers or missing entries.  For example, if the `Size` method is invoked on a nil pointer, the value `0` is returned, rather than returning an error.  It is assumed that if the caller needs to know that the pointer receiver is not nil, the caller should perform that check.

### Hashmap
The hashmap is a key-value pair collection based on Go's native map implementation.

#### Methods

`Get(key Key) Value`

The get method searches the collection for a key-value pair with the matching key, and returns the value.  In the case where there is no matching key, or the pointer receiver is nil, then `nil` is returned.

`Insert(key Key, value Value) (*HashMap, error)`

The insert method creates a copy of the provided hashmap collection with the provided key-value pair added.

`Remove(key Key) (*HashMap, error)`

The remove method creates a copy of the provided hashmap collection, with the entry at the specified key removed.  If the method would result in no change, the same reference is returned. 

### Common functions

Collections have some common functions for operations across the entire collection

`Filter`

`ForEach`

`Map`

`Reduce`

`Size() uint32`

This size method returns the number of key-value pairs in a collection.
