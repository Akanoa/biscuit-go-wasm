package shared

type Stringable interface {
	ToStringWasmFunction() string
	Ptr() uint64
}
