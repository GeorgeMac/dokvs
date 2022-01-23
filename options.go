package etcdoc

type Options[T any] []func(*T)

func (o Options[T]) Apply(t *T) {
	for _, opt := range o {
		opt(t)
	}
}
