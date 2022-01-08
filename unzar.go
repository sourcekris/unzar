package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/JoshVarga/blast"
)

const (
	name       = "Zip-Archiv" // Name of the archiver.
	ext        = "ZAR"        // File extension usually seen.
	id         = "PT&"        // File ID string.
	fFooterLen = 7            // fileID + 4 unknown bytes
)

var (
	fileID      = []byte(id)
	fset        = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	arcFile     = fset.String("e", "", fmt.Sprintf("The %s file to extract.", ext))
	dstPath     = fset.String("d", "", "Optional output directory to extract to.")
	bytesRead   int
	headerBytes int
)

func errpanic(e error) {
	if e != nil {
		panic(e)
	}
}

type header struct {
	id    []byte // Stores the file sig
	fSize int    // Total file size
}

type entry struct {
	cSize    int    // Compressed file size
	fnSize   int    // Filename length
	fnOff    int    // Filename offset
	fnLenOff int    // Filename Length byte offset
	fn       string // Filename for this header.
}

// Seeks backwards len bytes, reads forwards, then returns to original position in file.
func readBack(f *os.File, len int64) ([]byte, error) {
	buf := make([]byte, len)
	_, err := f.Seek(-len, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(-len, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func fTell(f *os.File) (int, error) {
	pos, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	_, err = f.Seek(pos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return int(pos), nil
}

func main() {
	var entries []*entry

	fset.Parse(os.Args[1:])

	if *arcFile == "" {
		fset.Usage()
		os.Exit(0)
	}

	if *dstPath != "" {
		if _, err := os.Stat(*dstPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "destination folder (%s) doesn't exist: %v", *dstPath, err)
			os.Exit(1)
		}
	}

	f, err := os.Open(*arcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v", *arcFile, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := f.Stat()
	errpanic(err)

	h := &header{
		id:    make([]byte, len(fileID)),
		fSize: int(s.Size()),
	}

	// Read the file id from the end.
	_, err = f.Seek(-3, io.SeekEnd)
	errpanic(err)
	_, err = f.Read(h.id)
	errpanic(err)

	bytesRead = bytesRead + 3

	if !reflect.DeepEqual(h.id, fileID) {
		fmt.Fprintf(os.Stderr, "file %s is not a %s file\n", *arcFile, ext)
		os.Exit(1)
	}

	// Skip backwards 4 bytes past the 2 unknown uint16s.
	_, err = f.Seek(-7, io.SeekCurrent)
	errpanic(err)

	bytesRead = bytesRead + 4

	// Read archive entries in a loop.
	for {
		// Read filesize first.
		buf, err := readBack(f, 4)
		errpanic(err)

		e := &entry{
			cSize: int(binary.LittleEndian.Uint32(buf)),
		}

		bytesRead = bytesRead + 4

		// Store current cursor and return to current pos.
		pos, err := fTell(f)
		errpanic(err)

		// Read the filename bytes until we reach the filelen indicator.
		for count := h.fSize - fFooterLen - headerBytes; ; count-- {
			buf, err := readBack(f, 1)
			errpanic(err)

			bv := int(buf[0]) - 0x80
			pc := int(pos) - count + 4

			if bv == pc {
				e.fnSize = bv
				pos, err = fTell(f)
				errpanic(err)
				e.fnLenOff = pos
				e.fnOff = pos + 1

				headerBytes = headerBytes + e.fnSize + 4 + 1
				bytesRead = bytesRead + e.fnSize + 1 // +1 for the filesize byte itself.
				break
			}
		}

		if e.fnSize > 12 {
			fmt.Fprintf(os.Stderr, "filename length is > 12: %d", e.fnSize)
			os.Exit(1)
		}

		fn := make([]byte, e.fnSize)
		_, err = f.Seek(int64(e.fnOff), io.SeekStart)
		errpanic(err)
		_, err = f.Read(fn)
		errpanic(err)
		e.fn = string(fn)

		_, err = f.Seek(int64(e.fnLenOff), io.SeekStart)
		errpanic(err)

		bytesRead = bytesRead + e.cSize
		entries = append(entries, e)

		if bytesRead == h.fSize {
			_, err := f.Seek(0, io.SeekStart)
			errpanic(err)
			break
		}
	}

	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		fmt.Printf("Extracting: %s (%d bytes compressed)\n", e.fn, e.cSize)

		// Read all the DCL compressed data.
		dcl := make([]byte, e.cSize)
		_, err = f.Read(dcl)
		errpanic(err)

		b := bytes.NewReader(dcl)
		r, err := blast.NewReader(b)
		errpanic(err)

		if *dstPath != "" {
			e.fn = *dstPath + "/" + e.fn
		}

		o, err := os.Create(e.fn)
		errpanic(err)

		_, err = io.Copy(o, r)
		errpanic(err)
		o.Close()
	}
}
