package ipfs

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aperturerobotics/pbobject"
	"github.com/golang/protobuf/proto"
	api "github.com/ipfs/go-ipfs-api"
)

// FileShell manages reading/writing pbobject from an IPFS API shell as file objects.
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
	objWrapper, _, err := s.objTable.Encode(ctx, obj)
	if err != nil {
		return "", err
	}

	dat, err := proto.Marshal(objWrapper)
	if err != nil {
		return "", err
	}

	return s.Add(bytes.NewReader(dat))
}

// GetProtobufObject gets a protobuf object from IPFS.
func (s *FileShell) GetProtobufObject(ctx context.Context, hash string) (pbobject.Object, error) {
	rc, err := s.Cat(hash)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rc.Close()
	}()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	wrapper := &pbobject.ObjectWrapper{}
	if err := proto.Unmarshal(data, wrapper); err != nil {
		return nil, err
	}

	return s.objTable.DecodeWrapper(ctx, wrapper)
}
