package signin


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
        
          "signin.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "signin.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\xc1\x8e\xd3\x30\x10\x86\xef\x48\xbc\x83\xe5\xeb\x6e\x6a\x71\x43\x55\x0d\x42\x9c\xe1\x00\x4f\x30\x6b\x4f\x92\x29\xb6\xc7\xf2\x4c\xdb\x8d\xb4\x0f\x8f\xba\x49\x20\x50\x55\x70\xca\x64\xe6\xff\xbf\xfc\x13\xd9\x87\x51\x73\xfa\xf0\xf6\x8d\x31\x87\x11\x21\xce\x55\x46\x05\x13\x46\x68\x82\xea\xed\x49\xfb\xee\xbd\xdd\x4c\x0a\x64\xf4\xf6\x4c\x78\xa9\xdc\xd4\x9a\xc0\x45\xb1\xa8\xb7\x17\x8a\x3a\xfa\x88\x67\x0a\xd8\xbd\xbe\x3c\x1a\x2a\xa4\x04\xa9\x93\x00\x09\xfd\xbb\x47\x93\xe1\x99\xf2\x29\xaf\x8d\x2d\xb8\x36\xae\xd8\x74\xf2\x96\x87\xbd\x92\x26\xdc\xc0\x85\x86\x42\xe5\xae\x3c\xa2\x84\x46\x55\x89\xcb\xc6\xf4\xa9\xd6\xfd\xad\x71\x5e\x40\x2f\xa4\x8a\x6d\x1f\xa0\xc5\xed\x77\x4e\x39\x43\x9b\xee\xeb\xff\x1d\xec\x4f\xfd\x7f\x26\x4b\x54\x7e\x98\x86\xc9\x5b\xa8\x35\x61\xa7\x7c\x0a\x63\x47\x81\x4b\x57\x1b\x06\xce\x95\x05\xa3\x35\x63\xc3\xde\x5b\x10\x41\x15\x47\x79\x70\x57\x89\xb8\xbf\x4d\xbb\x5a\x86\x1b\xf2\x75\x60\x8d\x4e\x15\xbd\xa5\x0c\x03\xba\xab\xea\x1e\xb2\x87\xf3\x2f\x92\xd9\xa0\x66\xf9\xa8\x5a\x65\xef\x5c\xcf\x45\x65\x37\x30\x0f\x09\xa1\x92\xec\x02\x67\x17\x44\x3e\xf6\x90\x29\x4d\xfe\x0b\x17\x56\x2e\x2f\x5f\x59\xf9\xe1\x3b\x14\x59\x2a\x6c\xd4\xbf\x7c\xe3\x27\x56\x5e\x1e\x0f\x9f\xb9\x44\x2c\x82\x71\x6d\x5c\xbd\x76\x8e\x2e\x3a\x25\x94\x11\x51\xad\x71\x73\x18\xf7\xfb\xc4\x3e\x71\x9c\x5e\x2b\x63\x0e\xf3\xdf\x5e\xb6\x54\x7c\x56\x77\x84\x33\xcc\x5d\x6b\xa4\x05\x6f\x8f\xe2\x32\x50\xd9\x1d\xc5\xae\x2e\x37\x0b\x16\xf2\xc2\x5b\x6f\xc7\xcf\x00\x00\x00\xff\xff\x25\xe3\xaa\x32\x26\x03\x00\x00"),
          path: "signin.html",
          root: "signin.html",
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

// FindFile returns a io.Reader by seeking the giving file path if it exists.
func FindFile(path string, doGzip bool) (io.Reader, int64, error){
  item, ok := assetFiles[path]
  if !ok {
    return nil,0, fmt.Errorf("File %q not found in file system", path)
  }

	datalen := int64(len(item.data))
  if !doGzip {
    return bytes.NewReader(item.data), datalen, nil
  }

  gzr, err := gzip.NewReader(bytes.NewReader(item.data))
	return gzr, datalen, err
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
