package pbobject

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/aperturerobotics/objectenc"
	"github.com/aperturerobotics/objectsig"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-crypto"
	lpeer "github.com/libp2p/go-libp2p-peer"
	"github.com/pkg/errors"
)

// Object is a protobuf-encoded object.
type Object interface {
	proto.Message

	// GetObjectTypeID returns the object type string, used to identify types.
	GetObjectTypeID() *ObjectTypeID
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
	// SignerKeys are the keys to sign the buffer with when making an object wrapper.
	SignerKeys []crypto.PrivKey
	// VerifyKeys requires the set of public keys to have signed the object.
	VerifyKeys []crypto.PubKey
	// CompressionType is the compression type to use.
	// Defaults to uncompressed
	CompressionType objectenc.CompressionType
}

// GetContext returns the context.
func (c *EncryptionConfig) GetContext() context.Context {
	if c.Context == nil {
		return context.TODO()
	}
	return c.Context
}

// NewObjectWrapper builds a new object wrapper.
// The unencrypted data is also returned for convenience.
func NewObjectWrapper(obj Object, econf EncryptionConfig) (*ObjectWrapper, []byte, error) {
	ctx := econf.GetContext()
	data, err := proto.Marshal(obj)
	if err != nil {
		return nil, nil, err
	}

	// Build the signatures.
	var sigs []*objectsig.Signature
	for _, signer := range econf.SignerKeys {
		sig, err := objectsig.NewSignature(signer, data)
		if err != nil {
			return nil, nil, err
		}
		sigs = append(sigs, sig)
	}

	encBlob, err := objectenc.EncryptWithResolver(
		ctx,
		econf.ResourceLookup,
		econf.EncryptionType,
		econf.CompressionType,
		data,
	)
	if err != nil {
		return nil, nil, err
	}

	return &ObjectWrapper{
		ObjectTypeCrc: obj.GetObjectTypeID().GetCrc32(),
		EncBlob:       encBlob,
		Signatures:    sigs,
	}, data, nil
}

// DecodeToObject decodes the object wrapper to a pre-identified object.
func (w *ObjectWrapper) DecodeToObject(obj Object, encConf EncryptionConfig) error {
	expectedID := obj.GetObjectTypeID().GetCrc32()
	if expectedID != w.ObjectTypeCrc {
		return fmt.Errorf("object type mismatch: expected %d != actual %d", expectedID, w.ObjectTypeCrc)
	}

	ctx := encConf.GetContext()
	objData, err := w.GetEncBlob().DecryptWithResolver(ctx, encConf.ResourceLookup)
	if err != nil {
		return err
	}

	// TODO: optimize
VerifyLoop:
	for _, ver := range encConf.VerifyKeys {
		var verErr error
		for _, sig := range w.Signatures {
			if err := sig.MatchesPublicKey(ver); err != nil {
				continue
			}

			verErr = sig.Verify(ver, objData)
			if verErr == nil {
				continue VerifyLoop
			}
		}

		var peerIDStr string
		peerID, err := lpeer.IDFromPublicKey(ver)
		if err != nil {
			peerIDStr = err.Error()
		} else {
			peerIDStr = peerID.Pretty()
		}

		if verErr == nil {
			return errors.Errorf("object not signed by key: %s", peerIDStr)
		}

		return errors.WithMessage(verErr, fmt.Sprintf("key %s verify error", peerIDStr))
	}

	if err := proto.Unmarshal(objData, obj); err != nil {
		return err
	}

	return nil
}

// objectTableKey is a pointer used as the key for object tables.
var objectTableKey = struct{ objectTableKey string }{}

// WithObjectTable attaches an object table to a context.
func WithObjectTable(parent context.Context, table *ObjectTable) context.Context {
	return context.WithValue(parent, &objectTableKey, table)
}

// GetObjectTable returns the object table in the context.
func GetObjectTable(ctx context.Context) *ObjectTable {
	v := ctx.Value(&objectTableKey)
	if v == nil {
		return nil
	}

	return v.(*ObjectTable)
}
