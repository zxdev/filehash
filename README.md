# filehash
create and manage files with an embedded checksum header

The filehash header holds the file checksum value based on a SHA256 of the data. The header is available immediatley after a successful filehash.Open when reading from file, and immediately after a successful filehash.Close when writing to a file.

```
                       0       2     34       36 ... [n]
header/data layout : [[prefix][hash][suffix]][data...]
```

Sample Usage

```
var test = "test.fh"

func TestWriter() {
	fh, err := filehash.NewWriter(test)
	if err != nil {
		log.Println(err)
	}
	defer fh.Close()
	fh.Write([]byte("hello filehash"))
}

func TestReader() {
	fh, err := filehash.NewReader(test)
	if err != nil {
		log.Prinln(err)
	}
	defer fh.Close()
	var buf = make([]byte, 512)
	fh.Read(buf[:])
	log.Println(fh.Hex())
	log.Println(string(buf))
}
```
