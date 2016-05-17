package interpol

import "errors"

// ErrMapKeyNotFound is returned by Exec() when the given map doesn't have the
// key needed by the template string.
var ErrMapKeyNotFound = errors.New(`Key was not found in the map!`)

// ErrSpawnerNotFound is returned by New() if there's no lookup function for the
// given type.
var ErrSpawnerNotFound = errors.New(`Spawner for the type was not found!`)
