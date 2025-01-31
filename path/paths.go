package path

import "strings"

// Paths is a collection of exact attribute paths.
type Paths []Path

// Contains returns true if the collection of paths includes the given path.
func (p Paths) Contains(checkPath Path) bool {
	for _, path := range p {
		if path.Equal(checkPath) {
			return true
		}
	}

	return false
}

// String returns the human-readable representation of the path collection.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
//
// Empty paths are skipped.
func (p Paths) String() string {
	var result strings.Builder

	result.WriteString("[")

	for pathIndex, path := range p {
		if path.Equal(Empty()) {
			continue
		}

		if pathIndex != 0 {
			result.WriteString(",")
		}

		result.WriteString(path.String())
	}

	result.WriteString("]")

	return result.String()
}
