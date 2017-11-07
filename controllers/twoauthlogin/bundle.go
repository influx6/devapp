package twoauthlogin


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
        
          "twoauthlogin.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "twoauthlogin.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x93\xc1\x8e\xd3\x30\x10\x86\xef\x48\xbc\x83\xe5\xeb\x6e\x6a\x71\x43\x55\x0d\x42\x9c\xe1\x00\x4f\x30\xeb\x4c\xe2\x29\xb6\xc7\xf2\x4c\xda\x8d\xb4\x0f\x8f\xda\xa4\x10\x40\x95\xd8\x53\x26\xe3\xff\xfb\x34\x13\xc5\x87\xa8\x39\x7d\x78\xfb\xc6\x98\x43\x44\xe8\x97\x2a\xa3\x82\x09\x11\x9a\xa0\x7a\x3b\xe9\xd0\xbd\xb7\x9b\x93\x02\x19\xbd\x3d\x11\x9e\x2b\x37\xb5\x26\x70\x51\x2c\xea\xed\x99\x7a\x8d\xbe\xc7\x13\x05\xec\xae\x2f\x8f\x86\x0a\x29\x41\xea\x24\x40\x42\xff\xee\xd1\x64\x78\xa6\x3c\xe5\x5b\x63\x2b\xae\x8d\x2b\x36\x9d\xbd\xe5\x71\xaf\xa4\x09\x37\x72\x3d\x33\x4c\x1a\x13\x8f\x54\xee\x42\x3d\x4a\x68\x54\x95\xb8\x6c\xd0\x4f\xb5\xee\xef\xe1\xcb\x32\x7a\x26\x55\x6c\xfb\x00\xad\xdf\x80\x32\xe5\x0c\x6d\xbe\x9f\xff\xdf\x21\xff\xa4\x5e\x35\x65\xa2\xf2\xc3\x34\x4c\xde\x42\xad\x09\x3b\xe5\x29\xc4\x8e\x02\x97\xae\x36\x0c\x9c\x2b\x0b\xf6\xd6\xc4\x86\x83\xb7\x20\x82\x2a\x8e\xf2\xe8\x2e\x11\x71\x7f\x43\xbb\x5a\xc6\x7f\xcc\x97\x03\x6b\x74\xae\xe8\x2d\x65\x18\xd1\x5d\x52\xf7\x94\x03\x9c\x7e\x99\xcc\x46\xb5\xc4\xa3\x6a\x95\xbd\x73\x03\x17\x95\xdd\xc8\x3c\x26\x84\x4a\xb2\x0b\x9c\x5d\x10\xf9\x38\x40\xa6\x34\xfb\x2f\x5c\x58\xb9\xbc\x7c\x65\xe5\x87\xef\x50\x64\xad\xb0\xd1\xf0\xf2\x8d\x9f\x58\x79\x7d\x3c\x7c\xe6\xd2\x63\x11\xec\x6f\x8d\x0b\x6b\x97\xd1\x45\xe7\x84\x12\x11\xd5\x1a\xb7\x0c\xe3\x7e\xff\xc9\x4f\xdc\xcf\xd7\xca\x98\xc3\xf2\xcd\xd7\x2d\x15\x9f\xd5\x1d\xe1\x04\x4b\xd7\x1a\x69\xc1\xdb\xa3\xb8\x0c\x54\x76\x47\xb1\x37\xca\x2d\x81\xd5\xbc\xfa\x6e\xb7\xe6\x67\x00\x00\x00\xff\xff\x11\x8b\x60\x31\x3e\x03\x00\x00"),
          path: "twoauthlogin.html",
          root: "twoauthlogin.html",
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
