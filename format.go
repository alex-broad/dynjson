package dynjson

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type formatter interface {
	typ() reflect.Type
	format(src reflect.Value) (reflect.Value, error)
}

// Formatter is a dynamic API format formatter.
type Formatter struct {
	mu         sync.Mutex
	builders   map[reflect.Type]builder
	formatters map[reflect.Type]map[string]formatter
}

// NewFormatter creates a new formatter.
func NewFormatter() *Formatter {
	return &Formatter{
		builders:   map[reflect.Type]builder{},
		formatters: map[reflect.Type]map[string]formatter{},
	}
}

// Format formats either a struct or a slice, returning only the selected fields (or all if none specified).
func (f *Formatter) Format(o interface{}, fields []string) (interface{}, error) {
	if len(fields) == 0 {
		return o, nil
	}
	err := detectDuplicateFields(fields)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(o)
	t := v.Type()
	f.mu.Lock()
	defer f.mu.Unlock()
	b := f.builders[t]
	if b == nil {
		var err error
		b, err = makeBuilder(t)
		if err != nil {
			return nil, err
		}
		f.builders[t] = b
		f.formatters[t] = map[string]formatter{}
	}
	key := strings.Join(fields, ",")
	ff := f.formatters[t][key]
	if ff == nil {
		var err error
		ff, err = b.build(fields, "")
		if err != nil {
			return nil, err
		}
		f.formatters[t][key] = ff
	}
	v, err = ff.format(v)
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

// detectDuplicateFields returns an error if passed the same field more than once.
func detectDuplicateFields(fields []string) error {
	h := make(map[string]int)
	for _, f := range fields {
		h[f]++
	}
	e := []string{}
	for f, count := range h {
		if count > 1 {
			e = append(e, f)
		}
	}
	if len(e) > 0 {
		return fmt.Errorf("duplicate fields detected: %s", strings.Join(e, ", "))
	}
	return nil
}
