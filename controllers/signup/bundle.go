package signup


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
        
          "signup.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "signup.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\xc1\x8e\xd3\x30\x10\x86\xef\x48\xbc\x83\xe5\xeb\x6e\x6a\x71\x43\x55\x0d\x42\x9c\xe1\x00\x4f\x30\x6b\x4f\x92\x29\xb6\xc7\xf2\x4c\xda\x8d\xb4\x0f\x8f\xba\x49\x21\x50\x55\x70\xca\x64\xe6\xff\xbf\xfc\x13\xd9\x87\x51\x73\xfa\xf0\xf6\x8d\x31\x87\x11\x21\x2e\x55\x46\x05\x13\x46\x68\x82\xea\xed\xa4\x7d\xf7\xde\x6e\x26\x05\x32\x7a\x7b\x22\x3c\x57\x6e\x6a\x4d\xe0\xa2\x58\xd4\xdb\x33\x45\x1d\x7d\xc4\x13\x05\xec\x5e\x5f\x1e\x0d\x15\x52\x82\xd4\x49\x80\x84\xfe\xdd\xa3\xc9\xf0\x4c\x79\xca\xd7\xc6\x16\x5c\x1b\x57\x6c\x3a\x7b\xcb\xc3\x5e\x49\x13\x6e\xe0\x42\x43\x99\xea\x5d\x79\x44\x09\x8d\xaa\x12\x97\x8d\xe9\x53\xad\xfb\x5b\xe3\xb2\x80\x9e\x49\x15\xdb\x3e\x40\x8b\xdb\xef\x4c\x39\x43\x9b\xef\xeb\xff\x1d\xec\x4f\xfd\x7f\x26\x4b\x54\x7e\x98\x86\xc9\x5b\xa8\x35\x61\xa7\x3c\x85\xb1\xa3\xc0\xa5\xab\x0d\x03\xe7\xca\x82\xd1\x9a\xb1\x61\xef\x2d\x88\xa0\x8a\xa3\x3c\xb8\x8b\x44\xdc\xdf\xa6\x5d\x2d\xc3\x0d\xf9\x32\xb0\x46\xe7\x8a\xde\x52\x86\x01\xdd\x45\x75\x0f\xd9\xc3\xe9\x17\xc9\x6c\x50\x8b\x7c\x54\xad\xb2\x77\xae\xe7\xa2\xb2\x1b\x98\x87\x84\x50\x49\x76\x81\xb3\x0b\x22\x1f\x7b\xc8\x94\x66\xff\x85\x0b\x2b\x97\x97\xaf\xac\xfc\xf0\x1d\x8a\xac\x15\x36\xea\x5f\xbe\xf1\x13\x2b\xaf\x8f\x87\xcf\x5c\x22\x16\xc1\x78\x6d\x5c\xbc\x76\x89\x2e\x3a\x27\x94\x11\x51\xad\x71\x4b\x18\xf7\xfb\xc4\x3e\x71\x9c\x5f\x2b\x63\x0e\xcb\xdf\x5e\xb7\x54\x7c\x56\x77\x84\x13\x2c\x5d\x6b\xa4\x05\x6f\x8f\xe2\x32\x50\xd9\x1d\xc5\x5e\x5d\x6e\x11\xac\xe4\x95\x77\xbd\x1d\x3f\x03\x00\x00\xff\xff\x98\x35\xfc\x93\x26\x03\x00\x00"),
          path: "signup.html",
          root: "signup.html",
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
