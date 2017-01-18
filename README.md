# immutable
#### A Go / Golang library for immutable collections

_Note:_ This library is under development and is not meant for public consumption at this time.

[![Build Status](https://travis-ci.org/object88/immutable.svg?branch=master)](https://travis-ci.org/object88/immutable) [![Stories in Ready](https://badge.waffle.io/object88/immutable.svg?label=ready&title=Ready)](http://waffle.io/object88/immutable)

In general, this library seeks to avoid errors due to nil pointers or missing entries.  For example, if the `Size` method is invoked on a nil pointer, the value `0` is returned, rather than returning an error.  It is assumed that if the caller needs to know that the pointer receiver is not nil, the caller should perform that check.

## Hashmap
The hashmap is a key-value pair collection based on Go's native map implementation.

### Methods

#### `Get(key Key) (result Value, ok bool)`

The get method searches the collection for a key-value pair with the matching key, and returns the value.  In the case where there is no matching key, or the pointer receiver is `nil`, then `ok` is `false`.  If a valid value is returned (including `nil`), then `ok` is `true`.

#### `GetKeys() (results []Key)`

The `GetKeys` method returns the collection of `Key` objects used to store values in the hashmap.  If the hashmap is unassigned, `nil` is returned, and if the hashmap does not have any contents, then a 0-length array is returned.

It is important to note that, just as the collection functions do not operate in any order, the `GetKeys` method will return the key collection in random order.

#### `Insert(key Key, value Value) (*HashMap, error)`

The insert method creates a copy of the provided hashmap collection with the provided key-value pair added.  The `key` may not be `nil`, but `value` may.

#### `Remove(key Key) (*HashMap, error)`

The remove method creates a copy of the provided hashmap collection, with the entry at the specified key removed.  If the method would result in no change, the same reference is returned.

## Common functions

Collections have some common functions for operations across the entire collection

#### `Filter`

#### `ForEach`

#### `Map`

#### `Reduce`

#### `Size() int`

The `size` method returns the number of key-value pairs in a collection.
