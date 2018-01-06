package pbobject

import (
	"fmt"
	"sync"
)

// ObjectTable looks up known object types.
type ObjectTable struct {
	typeMap sync.Map // map[uint32]func() Object
}

// NewObjectTable builds a new object table.
func NewObjectTable() *ObjectTable {
	return &ObjectTable{}
}

// RegisterType registers a new object constructor.
func (o *ObjectTable) RegisterType(overwrite bool, ctor func() Object) error {
	sample := ctor()
	key := sample.GetObjectTypeID().GetCrc32()
	var loaded bool
	if overwrite {
		o.typeMap.Store(key, ctor)
	} else {
		_, loaded = o.typeMap.LoadOrStore(key, ctor)
	}

	if loaded {
		return fmt.Errorf("type already registered: %s", sample.GetObjectTypeID().GetTypeUuid())
	}

	return nil
}

// RegisterTypes registers one or more object types.
func (o *ObjectTable) RegisterTypes(overwrite bool, ctors ...func() Object) error {
	for _, ctor := range ctors {
		if err := o.RegisterType(overwrite, ctor); err != nil {
			return err
		}
	}

	return nil
}

// DecodeWrapper attempts to decode the wrapped object.
func (o *ObjectTable) DecodeWrapper(wrapper *ObjectWrapper, encConf EncryptionConfig) (Object, error) {
	key := wrapper.GetObjectTypeCrc()
	ctorInter, ok := o.typeMap.Load(key)
	if !ok {
		return nil, fmt.Errorf("type identifier not known: %d", key)
	}

	ctor := ctorInter.(func() Object)
	obj := ctor()
	err := wrapper.DecodeToObject(obj, encConf)
	if err != nil {
		return nil, err
	}
	if objts, ok := obj.(TimestampAwareObject); ok {
		t := wrapper.GetTimestamp().ToTime()
		objts.SetEncodedTimestamp(t)
	}
	return obj, nil
}
