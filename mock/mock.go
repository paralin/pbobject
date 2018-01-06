package mock

import (
	"github.com/aperturerobotics/pbobject"
)

// NewMockObject builds a new mock object.
func NewMockObject(msg string) *MockObject {
	return &MockObject{Message: msg}
}

// GetObjectTypeID returns the object type ID.
func (o *MockObject) GetObjectTypeID() *pbobject.ObjectTypeID {
	return pbobject.NewObjectTypeID("/mock/object/0.0.1")
}
