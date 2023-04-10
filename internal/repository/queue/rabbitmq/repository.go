package rabbitmq

type Repository struct {
	options     Options
	isTracingOn bool
}

type Options struct{}

func New(options Options, isTracingOn bool) *Repository {
	return &Repository{options: options, isTracingOn: isTracingOn}
}
