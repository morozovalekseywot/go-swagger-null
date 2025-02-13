package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Byte is a nullable int.
type Byte struct {
	Byte  byte
	Valid bool
	Set   bool
}

// NewByte creates a new Byte
func NewByte(value byte, valid bool) Byte {
	return Byte{
		Byte:  value,
		Valid: valid,
		Set:   true,
	}
}

// ByteFrom creates a new Byte that will always be valid.
func ByteFrom(value byte) Byte {
	return NewByte(value, true)
}

// ByteFromPtr creates a new Byte that be null if i is nil.
func ByteFromPtr(ptr *byte) Byte {
	if ptr == nil {
		return NewByte(0, false)
	}

	return NewByte(*ptr, true)
}

// IsValid returns true if this carries and explicit value and
// is not null.
func (b Byte) IsValid() bool {
	return b.Set && b.Valid
}

// IsSet returns true if this carries an explicit value (null inclusive)
func (b Byte) IsSet() bool {
	return b.Set
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Byte) UnmarshalJSON(data []byte) error {
	b.Set = true

	if len(data) == 0 || bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Byte = 0
		return nil
	}

	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if len(x) > 1 {
		return errors.New("json: cannot convert to byte, text len is greater than one")
	}

	b.Byte = x[0]
	b.Valid = true

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Byte) UnmarshalText(text []byte) error {
	b.Set = true
	if len(text) == 0 {
		b.Valid = false
		return nil
	}

	if len(text) > 1 {
		return errors.New("text: cannot convert to byte, text len is greater than one")
	}

	b.Valid = true
	b.Byte = text[0]

	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Byte) MarshalJSON() ([]byte, error) {
	if !b.IsValid() {
		return NullBytes, nil
	}

	return []byte{'"', b.Byte, '"'}, nil
}

// MarshalText implements encoding.TextMarshaler.
func (b Byte) MarshalText() ([]byte, error) {
	if !b.IsValid() {
		return []byte{}, nil
	}

	return []byte{b.Byte}, nil
}

// SetValue changes this Byte's value and also sets it to be non-null.
func (b *Byte) SetValue(value byte) {
	b.Byte = value
	b.Valid = true
	b.Set = true
}

// Ptr returns a pointer to this Byte's value, or a nil pointer if this Byte is null.
func (b Byte) Ptr() *byte {
	if !b.IsValid() {
		return nil
	}

	return &b.Byte
}

// IsZero returns true for invalid Bytes, for future omitempty support (Go 1.4?)
func (b Byte) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Byte) Scan(value interface{}) error {
	if value == nil {
		b.Byte, b.Valid, b.Set = 0, false, false
		return nil
	}

	val := value.(string)
	if len(val) == 0 {
		b.Byte, b.Valid, b.Set = 0, false, false
		return nil
	}

	b.Byte, b.Valid, b.Set = val[0], true, true

	return nil
}

// Value implements the driver Valuer interface.
func (b Byte) Value() (driver.Value, error) {
	if !b.IsValid() {
		return nil, nil
	}

	return []byte{b.Byte}, nil
}
