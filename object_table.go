package pbobject

import (
	"context"
	"fmt"
	"sync"
)

// ObjectTable looks up known object types.
type ObjectTable struct {
	typeMap           sync.Map // map[uint32]func() Object
	typeEncryptionMap sync.Map // map[uint32]EncryptionConfig
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

// RegisterTypeEncryption registers a default encryption config for a type.
func (o *ObjectTable) RegisterTypeEncryption(typeIDCrc uint32, encConf EncryptionConfig) {
	o.typeEncryptionMap.Store(typeIDCrc, encConf)
}

// DecodeWrapper attempts to decode the wrapped object using the registered encryption conf for the type.
func (o *ObjectTable) DecodeWrapper(ctx context.Context, wrapper *ObjectWrapper) (Object, error) {
	var encConf EncryptionConfig
	if encConfInter, ok := o.typeEncryptionMap.Load(wrapper.GetObjectTypeCrc()); ok {
		encConf = encConfInter.(EncryptionConfig)
	}
	encConf.Context = ctx
	return o.DecodeWrapperWithEncConf(wrapper, encConf)
}

// Encode attempts to encode the object using the registered encryption conf for the type.
// Returns the unencrypted data as well.
func (o *ObjectTable) Encode(ctx context.Context, obj Object) (*ObjectWrapper, []byte, error) {
	var encConf EncryptionConfig
	objType := obj.GetObjectTypeID().GetCrc32()
	if encConfInter, ok := o.typeEncryptionMap.Load(objType); ok {
		encConf = encConfInter.(EncryptionConfig)
	}

	encConf.Context = ctx
	return NewObjectWrapper(obj, encConf)
}

// DecodeWrapperWithEncConf attempts to decode the wrapped object with an encryption config.
func (o *ObjectTable) DecodeWrapperWithEncConf(wrapper *ObjectWrapper, encConf EncryptionConfig) (Object, error) {
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
	return obj, nil
}
