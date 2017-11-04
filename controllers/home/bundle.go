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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\xc1\x8e\x14\x21\x10\x86\xef\x26\xbe\x03\xe1\xba\xdb\x43\xbc\x99\xc9\xa2\x31\x9e\xf5\xa0\x4f\x50\x0b\xd5\x4d\x8d\x40\x11\xaa\x66\x66\x3b\xd9\x87\x37\xb3\xdd\xa3\xed\x9a\x89\x9e\x28\xaa\xfe\xff\xcb\x0f\x81\x87\xa4\x25\x7f\x78\xfb\xc6\x98\x87\x84\x10\x97\xaa\xa0\x82\x09\x09\xba\xa0\x7a\x7b\xd4\x71\x78\x6f\x37\x93\x0a\x05\xbd\x3d\x11\x9e\x1b\x77\xb5\x26\x70\x55\xac\xea\xed\x99\xa2\x26\x1f\xf1\x44\x01\x87\x97\xcd\xbd\xa1\x4a\x4a\x90\x07\x09\x90\xd1\xbf\xbb\x37\x05\x9e\xa8\x1c\xcb\xb5\xb1\x05\xb7\xce\x0d\xbb\xce\xde\xf2\xb4\x57\xd2\x8c\x1b\x78\xe2\x82\x37\xc5\x11\x25\x74\x6a\x4a\x5c\x37\x96\x4f\xad\xed\x5f\xdb\x96\xf0\x7a\x26\x55\xec\xfb\x00\x3d\x6e\x0c\x72\x2c\x05\xfa\x7c\x5b\xff\xaf\x50\x7f\xaa\xff\x2b\x55\xa6\xfa\xc3\x74\xcc\xde\x42\x6b\x19\x07\xe5\x63\x48\x03\x05\xae\x43\xeb\x18\xb8\x34\x16\x8c\xd6\xa4\x8e\xa3\xb7\x20\x82\x2a\x8e\xca\xe4\x2e\x12\x71\xaf\x4d\xbb\x56\xa7\xbf\xc8\x97\x81\x35\x3a\x37\xf4\x96\x0a\x4c\xe8\x2e\xaa\x5b\xc8\x11\x4e\xbf\x48\x66\x83\x5a\xe4\x49\xb5\xc9\xde\xb9\x91\xab\xca\x6e\x62\x9e\x32\x42\x23\xd9\x05\x2e\x2e\x88\x7c\x1c\xa1\x50\x9e\xfd\x17\xae\xac\x5c\x9f\xbf\xb2\xf2\xdd\x77\xa8\xb2\x56\xd8\x69\x7c\xfe\xc6\x8f\xac\xbc\x2e\x77\x9f\xb9\x46\xac\x82\xf1\xda\xb8\x78\xed\x12\x5d\x74\xce\x28\x09\x51\xad\x71\x4b\x18\xf7\xfb\xa5\x3e\x72\x9c\x5f\x2a\x63\x1e\x96\xbb\x5e\x4f\xa9\xf8\xa4\xee\x00\x27\x58\xba\xd6\x48\x0f\xde\x1e\xc4\x15\xa0\xba\x3b\x88\xbd\xba\xdc\x22\x58\xc9\x2b\xef\xfa\x2b\x7e\x06\x00\x00\xff\xff\xbe\x0d\x34\x95\x1e\x03\x00\x00"),
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
func MustFindFile(path string, doGzip bool) io.Reader {
  reader, err := FindFile(path, doGzip)
  if err != nil {
    panic(err)
  }

  return reader
}

// FindFile returns a io.Reader by seeking the giving file path if it exists.
func FindFile(path string, doGzip bool) (io.Reader, error){
  item, ok := assetFiles[path]
  if !ok {
    return nil, fmt.Errorf("File %q not found in file system", path)
  }

  if !doGzip {
    return bytes.NewReader(item.data), nil
  }

  return gzip.NewReader(bytes.NewReader(item.data))
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
  reader, err := FindFile(path, doGzip)
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
