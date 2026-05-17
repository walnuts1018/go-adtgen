package resolvepkg

type Customer struct {
	Name string
}

type Address struct {
	City string
}

type Envelope[T any] struct {
	Value T
}
