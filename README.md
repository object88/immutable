# immutable
#### A Go / Golang library for immutable collections

_Note:_ This library is under development and is not meant for public consumption at this time.

[![Build Status](https://travis-ci.org/object88/immutable.svg?branch=master)](https://travis-ci.org/object88/immutable) [![Stories in Ready](https://badge.waffle.io/object88/immutable.svg?label=ready&title=Ready)](http://waffle.io/object88/immutable)

## Hashmap
The hashmap is a key-value pair collection based on Go's native map implementation.

### Methods

#### `Get(key Element) (result Element, ok bool, err error)`

The get method searches the collection for a key-value pair with the matching key, and returns the value.

If a `nil` pointer receiver is used, then an error is returned.  Additionally, `result` is `nil` and `ok` is `false`.

If there is no matching key, then `result` is nil, `ok` is `false`, and `err` is `nil`.

If the key matches in the hash map, a valid value is returned (including `nil`), and `ok` is `true`.

#### `GetKeys() (results []Element, err error)`

The `GetKeys` method returns the collection of `Element` objects used to store values in the hashmap.  If the hashmap is unassigned, `err` is non-`nil`, and if the hashmap does not have any contents, then a 0-length array is returned.

It is important to note that, just as the collection functions do not operate in any order, the `GetKeys` method will return the key collection in random order.

#### `Insert(key Element, value Element) (*InternalHashmap, error)`

The insert method creates a copy of the provided hashmap collection with the provided key-value pair added.  The `key` may not be `nil`, but `value` may.

#### `Remove(key Element) (*InternalHashmap, error)`

The remove method creates a copy of the provided hashmap collection, with the entry at the specified key removed.  If the method would result in no change, the same reference is returned.

## Common functions

Collections have some common functions for operations across the entire collection

#### `Filter`

#### `ForEach`

#### `Map`

#### `Reduce`

#### `Size() int`

The `size` method returns the number of key-value pairs in a collection.
