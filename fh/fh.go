package main

import (
	"girhub.com/zxdev/filehash"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*

	fh conversion utility

	This utility will add the filehash header to a file and remove
	the filehash header from a file; operates like gzip

*/

func main() {
	switch len(os.Args) {

	case 2: // add fh header

		// fh {file}
		if !strings.HasPrefix(os.Args[1], "-") {

			r, err := os.Open(os.Args[1])
			if err != nil {
				fmt.Println("fh: bad source file")
				return
			}

			tmp := "/tmp/fh-" + filepath.Base(os.Args[1])
			if fh, err := filehash.NewWriter(tmp); err == nil {
				io.Copy(fh, r)
				fh.Close()
			}

			r.Close()
			os.Remove(os.Args[1])
			os.Rename(tmp, os.Args[1])
		}

	case 3: // delete fh header

		// fh -d {file}
		fh, err := filehash.NewReader(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		tmp := "/tmp/fh-" + filepath.Base(os.Args[2])
		if w, err := os.Create(tmp); err == nil {
			io.Copy(w, fh)
			w.Close()
			fh.Close()
			os.Remove(os.Args[2])
			os.Rename(tmp, os.Args[2])
			return
		}
		fh.Close()

	default:
		fmt.Println("fh - filehash utility")
		fmt.Println("usage: fh [-d] {file}")
	}

}
