package etcdoc

func ApplyAll[T any](t *T, opts ...func(*T)) {
	Options[T](opts).Apply(t)
}

type Options[T any] []func(*T)

func (o Options[T]) Apply(t *T) {
	for _, opt := range o {
		opt(t)
	}
}
