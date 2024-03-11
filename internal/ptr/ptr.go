package ptr

// ToPtr returns a pointer to the given value.
func ToPtr[T any](v T) *T {
	return &v
}

// ToVal returns the value of the given pointer.
// If the given pointer is nil, it safely returns the zero value of the type.
func ToVal[T any](v *T) T {
	if v == nil {
		var zeroValue T
		return zeroValue
	}
	return *v
}
