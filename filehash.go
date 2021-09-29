package filehash

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

// unique header signature
var (
	prefix = [2]byte{0xCC, 0xFF}
	suffix = [2]byte{0xFF, 0xCC}
)

// Header type holds the file checksum value based on a SHA256
//
//                        0       2     34       36 ... [n]
// 	header+data layout : [[prefix][hash][suffix]][data]
type Header struct {
	checksum       [32]byte
	prefix, suffix [2]byte
}

// Hex value of the header checksum
func (h *Header) Hex() string { return fmt.Sprintf("%064x", h.checksum[:]) }

// Reader type for filehash files
type Reader struct {
	Header
	f *os.File
}

// Open a filehash file
func NewReader(path string) (*Reader, error) {
	z := new(Reader)
	return z, z.Open(path)
}

// Open file and read the filehash header
func (r *Reader) Open(path string) error {
	var err error
	if r.f, err = os.Open(path); err != nil {
		return err
	}
	return r.readHeader()
}

// Close file
func (r *Reader) Close() {
	r.f.Sync()
	r.f.Close()
}

// Read supports io.Reader interface
func (r *Reader) Read(p []byte) (int, error) { return r.f.Read(p) }

func (r *Reader) readHeader() error {
	r.f.Read(r.prefix[:])
	r.f.Read(r.checksum[:])
	r.f.Read(r.suffix[:])
	if r.prefix != prefix && r.suffix != suffix {
		return errors.New("filehash: invalid header")
	}
	return nil
}

// Writer type for filehash files
type Writer struct {
	Header
	h hash.Hash
	f *os.File
	t io.Writer
}

// Create a filehash file
func NewWriter(path string) (*Writer, error) {
	z := new(Writer)
	return z, z.Create(path)
}

// Create file and write a blank FileHash header
func (w *Writer) Create(path string) error {
	var err error
	if w.f, err = os.Create(path); err != nil {
		return err
	}
	w.writeHeader()
	w.prefix = prefix
	w.suffix = suffix
	w.h = sha256.New()
	w.t = io.MultiWriter(w.f, w.h)
	return nil
}

// Writer supports io.Writer interface
func (w *Writer) Write(p []byte) (int, error) { return w.t.Write(p) }

func (w *Writer) writeHeader() {
	w.f.Seek(0, 0)
	w.f.Write(w.prefix[:])
	w.f.Write(w.checksum[:])
	w.f.Write(w.suffix[:])
}

// Close and finalize the filehash header
func (w *Writer) Close() {
	copy(w.checksum[:], w.h.Sum(nil))
	w.writeHeader()
	w.f.Sync()
	w.f.Close()
}
