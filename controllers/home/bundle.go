package home


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type fileData struct{
  path string
  root string
  data []byte
}

var (
  assets = map[string][]string{
    
      ".html": []string{  // all .html assets.
        
          "home.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "home.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\x4d\x8a\xe4\x30\x0c\x85\xf7\x81\xdc\x41\x68\x5f\xf8\x02\x4e\x4e\x30\x8b\x81\x61\xd6\x85\xcb\x51\x12\x33\x8e\x1c\x2c\x7b\xaa\x9b\x50\x77\x6f\xf2\xe3\x4a\x2d\xba\xe9\x55\xaf\x22\x45\xd2\x7b\x1f\xf8\x2d\x4b\x47\xbd\x63\x02\xb4\x81\x13\x71\xc2\xc7\xa3\xae\xb4\x90\x4d\x2e\x30\x58\x6f\x44\x1a\xf4\x61\x70\x0c\x43\x74\xdd\xa5\xcf\xde\xef\x95\x25\x4e\x14\x21\x30\x5d\xee\x63\xf0\x74\x56\x97\x91\xdc\x30\x26\x6c\xeb\x0a\x40\xf7\x21\x4e\x45\x28\x0b\xc5\xeb\xf6\x63\x93\xdc\xcb\x4d\xed\xe6\x83\xfd\x77\x5a\x20\x98\x8d\xa0\x41\x25\x24\xe2\x02\x2b\xa6\x3b\xc2\x44\x69\x0c\x5d\x83\x73\x90\x43\x7f\x75\x70\xe4\x3b\xa1\xf4\xea\xc2\x66\xa2\x6b\x19\x94\x4d\x00\xed\x78\xce\xcf\xc5\x27\x31\x42\x7a\x9f\xa9\xc1\x44\x6f\x09\x61\xbd\x3d\x55\x10\xfe\x1b\x9f\xa9\x41\x84\xd9\x1b\x4b\x63\xf0\x1d\xc5\x06\xff\x96\x79\xc1\x50\xc5\xee\x67\xb9\x66\x23\x72\x0f\xb1\xfb\x8a\xeb\x77\x99\x7f\xcb\x75\xba\xdf\x72\x4a\x81\x0f\x33\xc9\xb7\xc9\x3d\xed\x8e\xae\xfd\xb5\x3e\x98\x56\xfb\xe6\x79\x69\x0a\xb3\xb8\x81\xf3\x8c\x30\x46\xea\xd7\x57\xdb\xdb\xf6\xcf\xf6\xd5\xca\x7c\x4e\xa3\xd5\x1a\x81\xb6\xae\xb4\x3a\x32\xd7\xd6\xd5\xb2\x10\x77\x6b\x0e\xcf\x74\x8a\x8d\x6e\x4e\xb2\xa5\x13\x40\xef\x2d\x48\xb4\x0d\x62\xab\xd5\xde\xbf\x9e\x7e\x04\x00\x00\xff\xff\xd0\x3c\x6b\x59\xd9\x02\x00\x00"),
          path: "home.html",
          root: "home.html",
        },
      
    
  }
)

//==============================================================================

// FilesFor returns all files that use the provided extension, returning a
// empty/nil slice if none is found.
func FilesFor(ext string) []string {
  return assets[ext]
}

// MustFindFile calls FindFile to retrieve file reader with path else panics.
func MustFindFile(path string, doGzip bool) (io.Reader, int64) {
  reader, size, err := FindFile(path, doGzip)
  if err != nil {
    panic(err)
  }

  return reader, size
}

// FindDecompressedGzippedFile returns a io.Reader by seeking the giving file path if it exists.
// It returns an uncompressed file.
func FindDecompressedGzippedFile(path string) (io.Reader, int64, error){
	return FindFile(path, true)
}

// MustFindDecompressedGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindDecompressedGzippedFile(path string) (io.Reader, int64){
	reader, size, err := FindDecompressedGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindGzippedFile returns a io.Reader by seeking the giving file path if it exists.
// It returns an uncompressed file.
func FindGzippedFile(path string) (io.Reader, int64, error){
	return FindFile(path, false)
}

// MustFindGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindGzippedFile(path string) (io.Reader, int64){
	reader, size, err := FindGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFile returns a io.Reader by seeking the giving file path if it exists.
func FindFile(path string, doGzip bool) (io.Reader, int64, error){
	reader, size, err := FindFileReader(path)
	if err != nil {
		return nil, size, err
	}

	if !doGzip {
		return reader, size, nil
	}

  gzr, err := gzip.NewReader(reader)
	return gzr, size, err
}

// MustFindFileReader returns bytes.Reader for path else panics.
func MustFindFileReader(path string) (*bytes.Reader, int64){
	reader, size, err := FindFileReader(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFileReader returns a io.Reader by seeking the giving file path if it exists.
func FindFileReader(path string) (*bytes.Reader, int64, error){
  item, ok := assetFiles[path]
  if !ok {
    return nil,0, fmt.Errorf("File %q not found in file system", path)
  }

  return bytes.NewReader(item.data), int64(len(item.data)), nil
}

// MustReadFile calls ReadFile to retrieve file content with path else panics.
func MustReadFile(path string, doGzip bool) string {
  body, err := ReadFile(path, doGzip)
  if err != nil {
    panic(err)
  }

  return body
}

// ReadFile attempts to return the underline data associated with the given path
// if it exists else returns an error.
func ReadFile(path string, doGzip bool) (string, error){
  body, err := ReadFileByte(path, doGzip)
  return string(body), err
}

// MustReadFileByte calls ReadFile to retrieve file content with path else panics.
func MustReadFileByte(path string, doGzip bool) []byte {
  body, err := ReadFileByte(path, doGzip)
  if err != nil {
    panic(err)
  }

  return body
}

// ReadFileByte attempts to return the underline data associated with the given path
// if it exists else returns an error.
func ReadFileByte(path string, doGzip bool) ([]byte, error){
  reader, _, err := FindFile(path, doGzip)
  if err != nil {
    return nil, err
  }

  if closer, ok := reader.(io.Closer); ok {
    defer closer.Close()
  }

  var bu bytes.Buffer

  _, err = io.Copy(&bu, reader);
  if err != nil && err != io.EOF {
   return nil, fmt.Errorf("File %q failed to be read: %+q", path, err)
  }

  return bu.Bytes(), nil
}
