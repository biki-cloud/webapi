package os_test

import (
	"os"
	"testing"
	os2 "webapi/pkg/os"
)

func TestMove(t *testing.T) {
	b := []byte("aaaaabbbbbccccceeeeedddddfffffgggggiiiii")
	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
	s, err := f.Stat()
	if err != nil {
		panic(err)
	}
	dstFile := "dst"
	mover := os2.NewMover()
	err = mover.Move(testFile, dstFile)
	if err != nil {
		panic(err)
	}

	if !os2.FileExists(dstFile) {
		t.Errorf("got: false, want: true")
	}

	d, err := os.Stat(dstFile)
	if err != nil {
		panic(err)
	}

	if s.Size() != d.Size() {
		t.Errorf("got: %v, want: %v", s.Size(), d.Size())
	}

	os.Remove(dstFile)

}
