package web

import "time"

type ValueBinder struct {
	err error
}

func NewValueBinder(source any) *ValueBinder {
	return &ValueBinder{}
}

func (b *ValueBinder) Bool(name string, dest *bool) *ValueBinder {
	return nil
}

func (b *ValueBinder) BoolOrDefault(name string, dest *bool, defaultValue bool) *ValueBinder {
	return nil
}

func (b *ValueBinder) String(name string, dest *string) *ValueBinder {
	return nil
}

func (b *ValueBinder) StringOrDefault(name string, dest *string) *ValueBinder {
	return nil
}

func (b *ValueBinder) Float32(name string, dest *float32) *ValueBinder {
	return nil
}

func (b *ValueBinder) Float32OrDefault(name string, dest *float32) *ValueBinder {
	return nil
}

func (b *ValueBinder) Float64(name string, dest *float64) *ValueBinder {
	return nil
}

func (b *ValueBinder) Float64OrDefault(name string, dest *float64) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int8(name string, dest *int8) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int8OrDefault(name string, dest *int8) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int16(name string, dest *int16) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int16OrDefault(name string, dest *int16) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int32(name string, dest *int32) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int32OrDefault(name string, dest *int32) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int64(name string, dest *int64) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int64OrDefault(name string, dest *int64) *ValueBinder {
	return nil
}

func (b *ValueBinder) Int(name string, dest *int) *ValueBinder {
	return nil
}

func (b *ValueBinder) IntOrDefault(name string, dest *int) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt8(name string, dest *uint8) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt8OrDefault(name string, dest *uint8) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt16(name string, dest *uint16) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt16OrDefault(name string, dest *uint16) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt32(name string, dest *uint32) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt32OrDefault(name string, dest *uint32) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt64(name string, dest *uint64) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt64OrDefault(name string, dest *uint64) *ValueBinder {
	return nil
}

func (b *ValueBinder) UInt(name string, dest *uint) *ValueBinder {
	return nil
}

func (b *ValueBinder) UIntOrDefault(name string, dest *uint) *ValueBinder {
	return nil
}

func (b *ValueBinder) Time(name string, dest *time.Time) *ValueBinder {
	return nil
}

func (b *ValueBinder) TimeOrDefault(name string, dest *time.Time, defaultValue time.Time) *ValueBinder {
	return nil
}

func (b *ValueBinder) Duration(name string, dest *time.Duration) *ValueBinder {
	return nil
}

func (b *ValueBinder) DurationOrDefault(name string, dest *time.Duration, defaultValue time.Duration) *ValueBinder {
	return nil
}

func (b *ValueBinder) Err() error {
	return b.err
}
