// Copyright 2025 Codnect
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"codnect.io/tag"
)

// PropertyBinder interface defines the method to bind configuration properties to a target.
type PropertyBinder interface {
	// Bind method binds the properties with the given name to the target.
	Bind(name string, target any) error
}

// DefaultPropertyBinder is the default implementation of the PropertyBinder interface.
type DefaultPropertyBinder struct {
	propSources  *PropertySources
	propResolver PropertyResolver
}

// NewDefaultPropertyBinder function creates a new DefaultPropertyBinder with the provided property sources.
func NewDefaultPropertyBinder(propSources *PropertySources) *DefaultPropertyBinder {
	if propSources == nil {
		panic("nil property sources")
	}

	return &DefaultPropertyBinder{
		propSources:  propSources,
		propResolver: NewDefaultPropertyResolver(propSources),
	}
}

// Bind method binds the properties with the given name to the target.
func (b *DefaultPropertyBinder) Bind(name string, target any) error {
	if target == nil {
		return errors.New("nil target")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("empty or blank property name")
	}

	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return errors.New("target must be a non-nil pointer")
	}

	return b.bindValue(name, targetVal.Elem())
}

// findProperty method searches for the property with the given name in the property sources.
func (b *DefaultPropertyBinder) findProperty(name string) (any, bool) {
	for _, source := range b.propSources.Slice() {
		if value, ok := source.Value(name); ok {
			return value, true
		}
	}

	return nil, false
}

// bindValue method binds the value to the target based on its kind.
func (b *DefaultPropertyBinder) bindValue(name string, targetVal reflect.Value) error {
	switch targetVal.Kind() {
	case reflect.Map:
		return b.bindMap(name, targetVal)
	case reflect.Slice:
		return b.bindSlice(name, targetVal)
	case reflect.Struct:
		return b.bindStruct(name, targetVal)
	case reflect.Interface:
		return b.bindAny(name, targetVal)
	default:
		return b.bindScalar(name, targetVal)
	}
}

// bindAny method binds any property to the target.
func (b *DefaultPropertyBinder) bindAny(name string, targetVal reflect.Value) error {
	propVal, ok := b.findProperty(name)
	if !ok {
		return ErrNoPropertyFound
	}

	targetVal.Set(reflect.ValueOf(propVal))
	return nil
}

// bindScalar method binds a scalar property to the target.
func (b *DefaultPropertyBinder) bindScalar(name string, target reflect.Value) error {
	val, ok := b.findProperty(name)
	if !ok {
		return ErrNoPropertyFound
	}

	return b.setScalar(target, val)
}

// bindSlice method binds a slice property to the target.
func (b *DefaultPropertyBinder) bindSlice(name string, targetVal reflect.Value) error {
	if v, ok := b.findProperty(name); ok {
		err := b.setValue(targetVal, v)
		if err != nil {
			return err
		}

		return nil
	}

	elemType := targetVal.Type().Elem()

	for i := 0; i < math.MaxInt64; i++ {
		indexedName := fmt.Sprintf("%s[%d]", name, i)
		elemVal := reflect.New(elemType)

		err := b.Bind(indexedName, elemVal.Interface())
		if err != nil {
			if errors.Is(err, ErrNoPropertyFound) {
				break
			}

			return fmt.Errorf("cannot append property %q to slice []%s: %w", name, elemType, err)
		}

		targetVal.Set(reflect.Append(targetVal, elemVal.Elem()))
	}

	if targetVal.Len() == 0 {
		return ErrNoPropertyFound
	}

	return nil
}

// bindMap method binds a map property to the target.
func (b *DefaultPropertyBinder) bindMap(name string, target reflect.Value) error {
	if target.IsNil() {
		target.Set(reflect.MakeMap(target.Type()))
	}

	rootKey := fmt.Sprintf("%s.", name)
	subKeys := make(map[string]struct{})
	for _, propSrc := range b.propSources.Slice() {
		for _, propName := range propSrc.PropertyNames() {
			if strings.HasPrefix(propName, rootKey) {
				fieldKey := strings.TrimPrefix(propName, rootKey)
				fieldKey = strings.TrimSpace(fieldKey)
				parts := strings.SplitN(fieldKey, ".", 2)
				subKeys[parts[0]] = struct{}{}
			}
		}
	}

	for subKey := range subKeys {
		valType := target.Type().Elem()
		val := reflect.New(valType)
		err := b.Bind(name+"."+subKey, val.Interface())
		if err != nil {
			return fmt.Errorf("cannot bind map property %q: %w", name, err)
		}

		key := reflect.ValueOf(subKey)
		target.SetMapIndex(key, val.Elem())
	}

	return nil
}

// bindStruct method binds struct properties to the target struct.
func (b *DefaultPropertyBinder) bindStruct(name string, targetVal reflect.Value) error {
	targetType := targetVal.Type()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldVal := targetVal.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		propTag := &PropertyTag{}
		tags := string(field.Tag)
		if len(tags) == 0 {
			continue
		}

		err := tag.Parse(tags, propTag)
		if err != nil {
			return fmt.Errorf("failed to parse tags for field %s: %w", field.Name, err)
		}

		propName := fmt.Sprintf("%s.%s", name, propTag.Name)

		err = b.bindValue(propName, fieldVal)

		if errors.Is(err, ErrNoPropertyFound) {
			if propTag.Default != nil {
				err = b.setValue(fieldVal, propTag.Default)
				if err != nil {
					return fmt.Errorf("failed to set default value for property %q: %w", propName, err)
				}
			} else if !propTag.Optional {
				return fmt.Errorf("missing required property: %q", propName)
			}
		} else if err != nil {
			return fmt.Errorf("failed to bind property %q: %w", propName, err)
		}

		switch field.Type.Kind() {
		case reflect.Slice, reflect.Map:
			if fieldVal.Len() == 0 && propTag.Default != nil {
				err = b.setValue(fieldVal, propTag.Default)
				if err != nil {
					return fmt.Errorf("failed to set default value for property %q: %w", propName, err)
				}
			}
		default:
			continue
		}
	}

	return nil
}

func (b *DefaultPropertyBinder) setValue(target reflect.Value, val any) error {
	source := reflect.ValueOf(val)

	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}

	if source.Type().ConvertibleTo(target.Type()) {
		target.Set(source.Convert(target.Type()))
		return nil
	}

	switch target.Kind() {
	case reflect.Slice:
		return b.copySlice(target, source)
	case reflect.Map:
		return b.copyMap(target, source)
	default:
		return b.setScalar(target, val)
	}
}

// setScalar method sets a scalar value to the target based on its kind.
func (b *DefaultPropertyBinder) setScalar(target reflect.Value, val any) error {
	switch target.Kind() {
	case reflect.Bool:
		bVal, err := toBool(val)
		if err != nil {
			return err
		}
		target.SetBool(bVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := toInt64(val, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := toUint64(val, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := toFloat64(val, target.Type().Bits())
		if err != nil {
			return err
		}
		target.SetFloat(f)
	case reflect.String:
		target.SetString(fmt.Sprint(val))
	default:
		return fmt.Errorf("unsupported target type: %s", target.Kind())
	}

	return nil
}

func (b *DefaultPropertyBinder) copySlice(target reflect.Value, source reflect.Value) error {
	if source.Kind() == reflect.String {
		parts := strings.Split(source.String(), ",")
		out := reflect.MakeSlice(target.Type(), 0, len(parts))
		elemType := target.Type().Elem()

		for _, p := range parts {
			p = strings.TrimSpace(p)
			elem := reflect.New(elemType).Elem()

			if err := b.setValue(elem, p); err != nil {
				return fmt.Errorf("cannot append element %q to slice []%s: %w", p, elemType, err)
			}

			out = reflect.Append(out, elem)
		}

		target.Set(out)
		return nil
	}

	out := reflect.MakeSlice(target.Type(), 0, source.Len())
	elemType := target.Type().Elem()

	for i := 0; i < source.Len(); i++ {
		sv := source.Index(i)
		if sv.Kind() == reflect.Interface && !sv.IsNil() {
			sv = sv.Elem()
		}

		elem := reflect.New(elemType).Elem()

		if err := b.setValue(elem, sv.Interface()); err != nil {
			return fmt.Errorf("cannot append element %q to slice []%s: %w", sv, elemType, err)
		}

		out = reflect.Append(out, elem)
	}

	target.Set(out)
	return nil
}

func (b *DefaultPropertyBinder) copyMap(target reflect.Value, source reflect.Value) error {
	keyType := target.Type().Key()
	valType := target.Type().Elem()

	iter := source.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		if key.Kind() == reflect.Interface && !key.IsNil() {
			key = key.Elem()
		}

		if val.Kind() == reflect.Interface && !val.IsNil() {
			val = val.Elem()
		}

		k := reflect.New(keyType).Elem()
		if err := b.setValue(k, key.Interface()); err != nil {
			return fmt.Errorf("cannot convert map key (%s) to %s: %w", key.Type(), keyType, err)
		}

		v := reflect.New(valType).Elem()
		if err := b.setValue(v, val.Interface()); err != nil {
			return fmt.Errorf("cannot convert map value (%s) to %s: %w", val.Type(), valType, err)
		}

		target.SetMapIndex(k, v)
	}

	return nil
}

// toBool function converts a value to a boolean.
func toBool(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(strings.TrimSpace(v))
	default:
		return strconv.ParseBool(fmt.Sprint(val))
	}
}

// toInt64 function converts a value to an int64.
func toInt64(val any, bits int) (int64, error) {
	switch v := val.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(v).Uint()), nil
	case float32, float64:
		return int64(reflect.ValueOf(v).Float()), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, bits)
	default:
		return strconv.ParseInt(fmt.Sprint(val), 10, bits)
	}
}

// toUint64 function converts a value to a uint64.
func toUint64(val any, bits int) (uint64, error) {
	switch v := val.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint(), nil
	case int, int8, int16, int32, int64:
		i := reflect.ValueOf(v).Int()
		if i < 0 {
			return 0, fmt.Errorf("cannot convert negative integer %d to uint", i)
		}
		return uint64(i), nil
	case float32, float64:
		f := reflect.ValueOf(v).Float()
		if f < 0 {
			return 0, fmt.Errorf("cannot convert negative float %f to uint", f)
		}
		return uint64(f), nil
	case string:
		return strconv.ParseUint(strings.TrimSpace(v), 10, bits)
	default:
		return strconv.ParseUint(fmt.Sprint(val), 10, bits)
	}
}

// toFloat64 function converts a value to a float64.
func toFloat64(val any, bits int) (float64, error) {
	switch v := val.(type) {
	case float32, float64:
		return reflect.ValueOf(v).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int()), nil
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return float64(reflect.ValueOf(v).Uint()), nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(v), bits)
	default:
		return strconv.ParseFloat(fmt.Sprint(val), bits)
	}
}
