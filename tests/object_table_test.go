package tests

import (
	"testing"

	"github.com/aperturerobotics/pbobject"
	"github.com/aperturerobotics/pbobject/mock"
)

// TestObjectTable tests the object table.
func TestObjectTable(t *testing.T) {
	table := pbobject.NewObjectTable()

	testMsg := "testing 1234"
	ctor := func() pbobject.Object {
		return &mock.MockObject{}
	}
	err := table.RegisterType(false, ctor)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = table.RegisterType(false, ctor)
	if err == nil {
		t.Fatal("expected error when registering duplicate type")
	}

	inObj := &mock.MockObject{Message: testMsg}
	objw, err := pbobject.NewObjectWrapper(inObj, pbobject.EncryptionConfig{})
	if err != nil {
		t.Fatal(err.Error())
	}

	obj, err := table.DecodeWrapper(objw, pbobject.EncryptionConfig{})
	if err != nil {
		t.Fatal(err.Error())
	}

	gotMsg := obj.(*mock.MockObject).GetMessage()
	if gotMsg != testMsg {
		t.Fatalf("got message %s expected %s", gotMsg, testMsg)
	}
}
