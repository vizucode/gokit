package logger

import (
	"context"
)

// Values create contract store value into context
type Values interface {
	Set(key Flags, value interface{})
	Load(key Flags) (interface{}, bool)
	LoadAndDelete(key Flags) (interface{}, bool)
	Delete(key Flags)
}

// Set value to keys
func (l *Locker) Set(key Flags, value interface{}) {
	l.data.Store(key, value)
}

// Delete value from sync
func (l *Locker) Delete(key Flags) {
	l.data.Delete(key)
}

// Load value from key
func (l *Locker) Load(key Flags) (interface{}, bool) {
	return l.data.Load(key)
}

// LoadAndDelete from key
func (l *Locker) LoadAndDelete(key Flags) (interface{}, bool) {
	return l.data.LoadAndDelete(key)
}

func extract(ctx context.Context) (Values, bool) {
	var (
		lock = new(Locker)
		ok   bool
	)

	if ctx == nil {
		return lock, false
	}

	lock, ok = ctx.Value(LogKey).(*Locker)
	return lock, ok
}
