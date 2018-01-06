# Protobuf Object

> Standardized protobuf object wrapper for blob storage.

## Introduction

pbobject is a standard wrapper for protobuf objects, intended to safely store things in IPFS.

## Go Usage

```go
import (
	"github.com/aperturerobotics/pbobject"
)

// GetObjectTypeID returns the object type ID.
func (o *MyObject) GetObjectTypeID() *pbobject.ObjectTypeID {
	return pbobject.NewObjectTypeID("/my/object/0.0.1")
}

table := pbobject.NewObjectTable()
table.RegisterType(false, func() { return &MyObject{} })

// Encrypt / store the object
inObj := &MyObject{}
objw, err := pbobject.NewObjectWrapper(inObj, pbobject.EncryptionConfig{})

// Identify and decode the wrapper.
obj, err := table.DecodeWrapper(objw, pbobject.EncryptionConfig{})
if err != nil {
    t.Fatal(err.Error())
}
```
