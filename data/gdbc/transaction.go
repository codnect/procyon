package gdbc

import (
	"codnect.io/procyon/data/tx"
	"context"
	"errors"
	"time"
)

const txDefaultTimeout = time.Duration(-1)

type txAttributeTimeout struct{}
type txAttributePropagation struct{}
type txAttributeIsolation struct{}
type txAttributeReadOnly struct{}

var (
	txAttributeTimeoutKey     = &txAttributeTimeout{}
	txAttributePropagationKey = &txAttributePropagation{}
	txAttributeIsolationKey   = &txAttributeIsolation{}
	txAttributeReadOnlyKey    = &txAttributeReadOnly{}
)

func DefaultTxAttributes() tx.Attributes {
	return tx.Attributes{
		txAttributeIsolationKey:   IsolationDefault,
		txAttributePropagationKey: PropagationRequired,
		txAttributeTimeoutKey:     txDefaultTimeout,
		txAttributeReadOnlyKey:    false,
	}
}

func WithIsolation(isolation Isolation) tx.FuncOption {
	return func(attributes tx.Attributes) {
		attributes[txAttributeIsolationKey] = isolation
	}
}

func WithPropagation(propagation Propagation) tx.FuncOption {
	return func(attributes tx.Attributes) {
		attributes[txAttributePropagationKey] = propagation
	}
}

func WithTimeout(timeout time.Duration) tx.FuncOption {
	return func(attributes tx.Attributes) {
		attributes[txAttributeTimeoutKey] = timeout
	}
}

func WithReadOnly(readOnly bool) tx.FuncOption {
	return func(attributes tx.Attributes) {
		attributes[txAttributeReadOnlyKey] = readOnly
	}
}

type txOptions struct {
	Isolation   Isolation
	Propagation Propagation
	Timeout     time.Duration
	ReadOnly    bool
}

func defaultTxOpts() *txOptions {
	return &txOptions{
		Isolation:   IsolationDefault,
		Propagation: PropagationRequired,
		Timeout:     txDefaultTimeout,
		ReadOnly:    false,
	}
}

func txOptsFromAttributes(attributes tx.Attributes) *txOptions {
	if len(attributes) == 0 {
		return defaultTxOpts()
	}

	txOpts := defaultTxOpts()
	if timeout, ok := attributes[txAttributeTimeoutKey]; ok {
		txOpts.Timeout = timeout.(time.Duration)
	}

	if isolation, ok := attributes[txAttributeIsolation{}]; ok {
		txOpts.Isolation = isolation.(Isolation)
	}

	if propagation, ok := attributes[txAttributePropagationKey]; ok {
		txOpts.Propagation = propagation.(Propagation)
	}

	if readOnly, ok := attributes[txAttributeReadOnlyKey]; ok {
		txOpts.ReadOnly = readOnly.(bool)
	}

	return txOpts
}

type TxManager struct {
	dataSource DataSource
}

func NewTxManager(dataSource DataSource) *TxManager {
	return &TxManager{
		dataSource: dataSource,
	}
}

func (t *TxManager) ContextWithTx(ctx context.Context, attributes tx.Attributes) (context.Context, error) {
	// TODO: complete implementation
	txOpts := txOptsFromAttributes(attributes)

	if txOpts.Timeout < txDefaultTimeout {
		return nil, errors.New("invalid tx timeout")
	} else if txOpts.Timeout != txDefaultTimeout {

	}

	if txOpts.Propagation == PropagationMandatory {
		return nil, errors.New("no existing tx found for tx marked with propagation with 'mandatory'")
	} else if txOpts.Propagation == PropagationRequired ||
		txOpts.Propagation == PropagationRequiredNew ||
		txOpts.Propagation == PropagationNested {

	}

	return nil, nil
}

func (t *TxManager) Commit(ctx context.Context) error {
	return nil
}

func (t *TxManager) Rollback(ctx context.Context) error {
	return nil
}
