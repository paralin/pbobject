package ipfs

import (
	"context"
	"io"

	"github.com/aperturerobotics/pbobject"
	api "github.com/ipfs/go-ipfs-api"
)

// ObjectShell reads and writes blobs and objects from IPFS.
// Used in interface form so that caching layers can be added.
type ObjectShell interface {
	// AddProtobufObject adds a protobuf object to IPFS and returns the hash.
	AddProtobufObject(ctx context.Context, obj pbobject.Object) (string, error)
	// GetProtobufObject gets a protobuf object from IPFS.
	GetProtobufObject(ctx context.Context, hash string) (pbobject.Object, error)
	// Add a file to ipfs from the given reader, returns the hash of the added file
	Add(r io.Reader) (string, error)
	// Unpin the given path
	Unpin(path string) error
	// Pins returns a map of the pin hashes to their info (currently just the
	// pin type, one of DirectPin, RecursivePin, or IndirectPin. A map is returned
	// instead of a slice because it is easier to do existence lookup by map key
	// than unordered array searching.
	Pins() (map[string]api.PinInfo, error)
}

var ipfsObjectShellKey = &(struct{ ipfsObjectShellKey string }{})

// WithObjectShell attaches an IPFS object shell to a context.
func WithObjectShell(parent context.Context, sh ObjectShell) context.Context {
	return context.WithValue(parent, ipfsObjectShellKey, sh)
}

// GetObjectShell gets the IPFS file shell from a context.
func GetObjectShell(ctx context.Context) ObjectShell {
	v := ctx.Value(ipfsObjectShellKey)
	if v == nil {
		return nil
	}

	return v.(ObjectShell)
}
