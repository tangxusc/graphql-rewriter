package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSjsonSetError(t *testing.T) {
	err := fmt.Errorf("test err")
	want := []byte(`{"errors":[{"message":"test err"}]}`)
	buffer := &bytes.Buffer{}
	writeError(buffer, err)
	value := buffer.Bytes()
	if !reflect.DeepEqual(want, value) {
		t.Errorf("excepted:%s,got:%s", want, value)
	}

}

func TestSjsonSetData(t *testing.T) {
	message := json.RawMessage(`{"queryProduct":[{"productID":"0x2","name":"test"}]}`)
	want := []byte(`{"data":{"queryProduct":[{"productID":"0x2","name":"test"}]}}`)
	buffer := &bytes.Buffer{}
	writeData(buffer, message)
	got := buffer.Bytes()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("excepted:%s,got:%s", want, got)
	}

}
