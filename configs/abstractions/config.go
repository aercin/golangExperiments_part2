package abstractions

type Config interface {
	GetValue(key string) any
}
