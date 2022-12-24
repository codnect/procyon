package tx

import "time"

type Option func(options *Options)

type Options struct {
	Propagation    Propagation
	Timeout        time.Duration
	ReadOnly       bool
	IsolationLevel IsolationLevel
}

func DefaultOptions() *Options {
	return &Options{
		Propagation:    PropagationRequired,
		Timeout:        -1,
		ReadOnly:       false,
		IsolationLevel: LevelDefault,
	}
}

func WithOptions(opts *Options, options ...Option) *Options {
	if opts == nil {
		opts = DefaultOptions()
	} else {
		opts = &Options{
			Propagation:    opts.Propagation,
			Timeout:        opts.Timeout,
			ReadOnly:       opts.ReadOnly,
			IsolationLevel: opts.IsolationLevel,
		}
	}

	for _, option := range options {
		option(opts)
	}

	return opts
}

func NewOptions(opts ...Option) *Options {
	return WithOptions(DefaultOptions(), opts...)
}

func WithPropagation(propagation Propagation) Option {
	return func(options *Options) {
		options.Propagation = propagation
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithReadOnly(readOnly bool) Option {
	return func(options *Options) {
		options.ReadOnly = readOnly
	}
}

func WithIsolation(isolationLevel IsolationLevel) Option {
	return func(options *Options) {
		options.IsolationLevel = isolationLevel
	}
}
