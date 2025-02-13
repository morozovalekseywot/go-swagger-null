package null

import (
	"bytes"
	"encoding/json"
	"testing"
)

var (
	hello        = []byte("hello")
	bytesJSON    = []byte(`"hello"`)
	b64BytesJSON = []byte(`"aGVsbG8="`)
)

func TestBytesFrom(t *testing.T) {
	i := BytesFrom(hello)
	assertBytes(t, i, "BytesFrom()")

	zero := BytesFrom(nil)
	if zero.Valid {
		t.Error("BytesFrom(nil)", "is valid, but should be invalid")
	}

	zero = BytesFrom([]byte{})
	if !zero.Valid {
		t.Error("BytesFrom([]byte{})", "is invalid, but should be valid")
	}
}

func TestBytesFromPtr(t *testing.T) {
	n := hello
	iptr := &n
	i := BytesFromPtr(iptr)
	assertBytes(t, i, "BytesFromPtr()")

	null := BytesFromPtr(nil)
	assertNullBytes(t, null, "BytesFromPtr(nil)")
}

func TestUnmarshalBytes(t *testing.T) {
	var i Bytes
	err := json.Unmarshal(b64BytesJSON, &i)
	maybePanic(err)
	assertBytes(t, i, "[]byte json")

	var ni Bytes
	err = ni.UnmarshalJSON([]byte{})
	if err == nil {
		t.Errorf("Expected error")
	}

	var null Bytes
	err = null.UnmarshalJSON(NullBytes)
	if err != nil {
		t.Error(err)
	}
	if null.Valid {
		t.Errorf("expected Valid to be false, got true")
	}
	if null.Bytes != nil {
		t.Errorf("Expected Bytes to be nil, but was not: %#v %#v", null.Bytes, []byte(`null`))
	}
	if !null.Set {
		t.Errorf("Expected Set to be true; got false")
	}
}

func TestTextUnmarshalBytes(t *testing.T) {
	var i Bytes
	err := i.UnmarshalText(hello)
	maybePanic(err)
	assertBytes(t, i, "UnmarshalText() []byte")

	var blank Bytes
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBytes(t, blank, "UnmarshalText() empty []byte")
}

func TestMarshalBytes(t *testing.T) {
	i := BytesFrom(hello)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, string(b64BytesJSON), "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewBytes(nil, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalBytesText(t *testing.T) {
	i := BytesFrom(bytesJSON)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewBytes(nil, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBytesPointer(t *testing.T) {
	i := BytesFrom([]byte(`"hello"`))
	ptr := i.Ptr()
	if !bytes.Equal(*ptr, bytesJSON) {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "pointer", ptr, `"hello"`)
	}

	null := NewBytes(nil, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBytesIsZero(t *testing.T) {
	i := BytesFrom(bytesJSON)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewBytes(nil, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewBytes(nil, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestBytesSetValid(t *testing.T) {
	change := NewBytes(nil, false)
	assertNullBytes(t, change, "SetValue()")
	change.SetValue(hello)
	assertBytes(t, change, "SetValue()")
}

func TestBytesScan(t *testing.T) {
	var i Bytes
	err := i.Scan(`hello`)
	maybePanic(err)
	assertBytes(t, i, "Scan() []byte")

	var null Bytes
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBytes(t, null, "scanned null")
}

func assertBytes(t *testing.T, i Bytes, from string) {
	if !bytes.Equal(i.Bytes, hello) {
		t.Errorf("bad %s []byte: %v ≠ %v\n", from, string(i.Bytes), "hello")
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBytes(t *testing.T, i Bytes, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
