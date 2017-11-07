package twoauth


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
        
          "twoauth.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "twoauth.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\x4d\x8e\xdb\x30\x0c\x85\xf7\x05\x7a\x07\x41\xdb\x19\x47\xe8\xae\x08\xa2\x16\x45\xd7\xed\xa2\x3d\x01\x47\xa6\x2d\xa6\x92\x28\x88\x74\x32\x06\xe6\xf0\x45\xc6\x4e\xeb\xfe\x04\xc8\xca\x34\xf9\xde\x87\x47\x41\x3a\x44\xcd\xe9\xc3\xdb\x37\xc6\x1c\x22\x42\xbf\x54\x19\x15\x4c\x88\xd0\x04\xd5\xdb\x49\x87\xee\xbd\xdd\x4c\x0a\x64\xf4\xf6\x44\x78\xae\xdc\xd4\x9a\xc0\x45\xb1\xa8\xb7\x67\xea\x35\xfa\x1e\x4f\x14\xb0\x7b\xfd\x79\x34\x54\x48\x09\x52\x27\x01\x12\xfa\x77\x8f\x26\xc3\x33\xe5\x29\x5f\x1b\x5b\x70\x6d\x5c\xb1\xe9\xec\x2d\x8f\x7b\x25\x4d\xb8\x81\xeb\x99\x61\xd2\x78\x53\xdf\xa3\x84\x46\x55\x89\xcb\xc6\xf5\xa9\xd6\xfd\x7f\x9c\xcb\x0a\x7a\x26\x55\x6c\xfb\x00\xad\xdf\x78\x64\xca\x19\xda\x7c\x5b\x7f\x47\xb4\x3f\x0d\xf7\x66\x4b\x54\x7e\x98\x86\xc9\x5b\xa8\x35\x61\xa7\x3c\x85\xd8\x51\xe0\xd2\xd5\x86\x81\x73\x65\xc1\xde\x9a\xd8\x70\xf0\x16\x44\x50\xc5\x51\x1e\xdd\x45\x22\xee\x6f\xd3\xae\x96\xf1\x1f\xf2\x65\x60\x8d\xce\x15\xbd\xa5\x0c\x23\xba\x8b\xea\x16\x72\x80\xd3\x2f\x92\xd9\xa0\x16\x79\x54\xad\xb2\x77\x6e\xe0\xa2\xb2\x1b\x99\xc7\x84\x50\x49\x76\x81\xb3\x0b\x22\x1f\x07\xc8\x94\x66\xff\x85\x0b\x2b\x97\x97\xaf\xac\xfc\xf0\x1d\x8a\xac\x15\x36\x1a\x5e\xbe\xf1\x13\x2b\xaf\x9f\x87\xcf\x5c\x7a\x2c\x82\xfd\xb5\x71\xf1\xda\x25\xba\xe8\x9c\x50\x22\xa2\x5a\xe3\x96\x30\xee\xf7\xad\x7d\xe2\x7e\x7e\xad\x8c\x39\x2c\xc7\xbd\x6e\xa9\xf8\xac\xee\x08\x27\x58\xba\xd6\x48\x0b\xde\x1e\xc5\x65\xa0\xb2\x3b\x8a\xbd\xba\xdc\x22\x58\xc9\x2b\xef\xfa\x42\x7e\x06\x00\x00\xff\xff\x40\x20\x99\x57\x2a\x03\x00\x00"),
          path: "twoauth.html",
          root: "twoauth.html",
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
