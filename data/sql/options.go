package sql

import "time"

type TransactionOption func(options *TransactionOptions)

type TransactionOptions struct {
	propagation    Propagation
	timeout        time.Duration
	readOnly       bool
	isolationLevel IsolationLevel
}

func (o *TransactionOptions) Propagation() Propagation {
	return o.propagation
}

func (o *TransactionOptions) Timeout() time.Duration {
	return o.timeout
}

func (o *TransactionOptions) IsReadOnly() bool {
	return o.readOnly
}

func (o *TransactionOptions) IsolationLevel() IsolationLevel {
	return o.isolationLevel
}

func (o *TransactionOptions) Merge(options *TransactionOptions) *TransactionOptions {
	copyOptions := new(TransactionOptions)
	*copyOptions = *o

	copyOptions.propagation = options.propagation
	copyOptions.timeout = options.timeout
	copyOptions.readOnly = options.readOnly
	copyOptions.isolationLevel = options.isolationLevel

	return copyOptions
}

func (o *TransactionOptions) Override(options ...TransactionOption) *TransactionOptions {
	copyOptions := new(TransactionOptions)
	*copyOptions = *o

	for _, option := range options {
		option(copyOptions)
	}

	return copyOptions
}

func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		propagation:    PropagationRequired,
		timeout:        -1,
		readOnly:       false,
		isolationLevel: LevelDefault,
	}
}

func NewTransactionOptions(opts ...TransactionOption) *TransactionOptions {
	return DefaultTransactionOptions().Override(opts...)
}

func WithPropagation(propagation Propagation) TransactionOption {
	return func(options *TransactionOptions) {
		options.propagation = propagation
	}
}

func WithTimeout(timeout time.Duration) TransactionOption {
	return func(options *TransactionOptions) {
		options.timeout = timeout
	}
}

func WithReadOnly(readOnly bool) TransactionOption {
	return func(options *TransactionOptions) {
		options.readOnly = readOnly
	}
}

func WithIsolation(isolationLevel IsolationLevel) TransactionOption {
	return func(options *TransactionOptions) {
		options.isolationLevel = isolationLevel
	}
}
