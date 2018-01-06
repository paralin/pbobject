package ipfs

import (
	"bytes"
	"context"

	"github.com/aperturerobotics/pbobject"
	"github.com/golang/protobuf/proto"
	api "github.com/ipfs/go-ipfs-api"
)

// FileShell manages reading/writing pbobject from an IPFS API shell.
type FileShell struct {
	*api.Shell
	objTable *pbobject.ObjectTable
}

// NewFileShell builds a new file shell.
func NewFileShell(sh *api.Shell, objTable *pbobject.ObjectTable) *FileShell {
	return &FileShell{Shell: sh, objTable: objTable}
}

// AddProtobufObject adds a protobuf object to IPFS and returns the hash.
func (s *FileShell) AddProtobufObject(ctx context.Context, obj pbobject.Object) (string, error) {
	objWrapper, err := s.objTable.Encode(obj)
	if err != nil {
		return "", err
	}

	dat, err := proto.Marshal(objWrapper)
	if err != nil {
		return "", err
	}

	return s.Add(bytes.NewReader(dat))
}
