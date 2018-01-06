package pbobject

import (
	"context"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/aperturerobotics/objectenc"
	"github.com/aperturerobotics/timestamp"
	"github.com/golang/protobuf/proto"
)

// Object is a protobuf-encoded object.
type Object interface {
	proto.Message

	// GetObjectTypeID returns the object type string, used to identify types.
	GetObjectTypeID() *ObjectTypeID
}

// TimestampAwareObject is an object that is aware of timestamps.
type TimestampAwareObject interface {
	// SetEncodedTimestamp sets the encoded timestamp of the object.
	SetEncodedTimestamp(time.Time)
}

// NewObjectTypeID builds a new object type id.
func NewObjectTypeID(uniqueID string) *ObjectTypeID {
	return &ObjectTypeID{TypeUuid: uniqueID}
}

// GetCrc32 gets the crc32 of the id.
func (o *ObjectTypeID) GetCrc32() uint32 {
	dat, _ := proto.Marshal(o)
	return crc32.ChecksumIEEE(dat)
}

// EncryptionConfig sets the encryption settings, defaults are unencrypted.
type EncryptionConfig struct {
	// Context if set will limit how long resource resolution can continue.
	Context context.Context
	// EncryptionType is the kind of encryption to use, default is unencrypted.
	// Only used when encrypting - when decrypting the encoded method is respected.
	EncryptionType objectenc.EncryptionType
	// ResourceLookup if set will be used to look up necessary keys and other data.
	ResourceLookup objectenc.ResourceResolverFunc
}

// GetContext returns the context.
func (c *EncryptionConfig) GetContext() context.Context {
	if c.Context == nil {
		return context.TODO()
	}
	return c.Context
}

// NewObjectWrapper builds a new object wrapper.
func NewObjectWrapper(obj Object, econf EncryptionConfig) (*ObjectWrapper, error) {
	return NewObjectWrapperWithTimestamp(obj, econf, timestamp.Now())
}

// NewObjectWrapperWithTimestamp builds a new object wrapper with a preset timestamp.
func NewObjectWrapperWithTimestamp(obj Object, econf EncryptionConfig, ts timestamp.Timestamp) (*ObjectWrapper, error) {
	ctx := econf.GetContext()
	data, err := proto.Marshal(obj)
	if err != nil {
		return nil, err
	}

	encBlob, err := objectenc.EncryptWithResolver(ctx, econf.ResourceLookup, econf.EncryptionType, data)
	if err != nil {
		return nil, err
	}

	return &ObjectWrapper{
		ObjectTypeCrc: obj.GetObjectTypeID().GetCrc32(),
		Timestamp:     &ts,
		EncBlob:       encBlob,
	}, nil
}

// DecodeToObject decodes the object wrapper to a pre-identified object.
func (w *ObjectWrapper) DecodeToObject(obj Object, encConf EncryptionConfig) error {
	expectedID := obj.GetObjectTypeID().GetCrc32()
	if obj.GetObjectTypeID().GetCrc32() != w.ObjectTypeCrc {
		return fmt.Errorf("object type mismatch: expected %d != actual %d", expectedID, w.ObjectTypeCrc)
	}

	ctx := encConf.GetContext()
	objData, err := w.GetEncBlob().DecryptWithResolver(ctx, encConf.ResourceLookup)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(objData, obj); err != nil {
		return err
	}

	return nil
}
