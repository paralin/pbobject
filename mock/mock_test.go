package mock

import (
	"testing"
)

func TestCrc32(t *testing.T) {
	obj := &MockObject{}
	typeID := obj.GetObjectTypeID()
	crc := typeID.GetCrc32()
	t.Logf("%s crc32: %d", typeID.TypeUuid, crc)
	if crc != 3957319027 {
		t.Fail()
	}
}
