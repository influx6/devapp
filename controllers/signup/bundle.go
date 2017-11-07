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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcf\xea\xdb\x30\x0c\x80\xef\x81\xbc\x83\xd0\xbd\xf8\x05\x6c\x5f\x76\xdd\x61\x30\x76\x2e\xae\xa3\x24\x66\x8e\x1d\xfc\xa7\xdd\x08\x7d\xf7\x91\x38\x6e\x7a\x68\x61\x85\xdf\xcd\x8a\xa4\x4f\x1f\x44\x5a\x96\x8e\x7a\xe3\x08\x50\x7b\x97\xc8\x25\xbc\xdf\xdb\x86\x47\xd2\xc9\x78\x07\xda\xaa\x18\x05\x46\x33\xb8\x3c\xc3\x10\x4c\x77\xea\xb3\xb5\xe5\xa5\xc9\x25\x0a\xe0\x1d\x9d\x6e\xa3\xb7\x74\xbc\x4e\x23\x99\x61\x4c\x28\xdb\x06\x80\xf7\x3e\x4c\x95\x94\x23\x85\xf3\xf6\xa1\x30\xcb\x7b\xc3\x5d\xac\xd7\xbf\x8f\x19\x08\x6a\x73\x10\xc8\xd6\xa6\xc8\x1c\xdd\x10\x26\x4a\xa3\xef\x04\xce\x3e\xee\xf8\x75\x80\x21\xdb\x45\x4a\xcf\x43\x9c\x9a\xe8\x5c\x13\xb5\x12\x80\x1b\x37\xe7\x47\xe1\x43\x18\x21\xfd\x9d\x49\x60\xa2\x3f\x09\x61\xed\x3d\x28\x08\x57\x65\x33\x09\x44\x98\xad\xd2\x34\x7a\xdb\x51\x10\xf8\xab\xe6\xab\x06\xab\xe3\xde\x79\xcd\x2a\xc6\x9b\x0f\xdd\xe7\x5e\xb5\xb3\xba\x1d\xf1\x6b\xb7\x1f\x35\xff\xb9\x9b\xf6\xae\x37\x61\xfa\x3a\xc7\x4a\x7c\xe7\xfa\xad\xa4\xe1\xff\x9d\x0f\xa3\x4b\x4e\xc9\xbb\x5d\x20\xe6\xcb\x64\x1e\xbf\x6f\x8f\xe4\xcf\x6d\xcf\x38\x2b\xa5\xaf\xd9\x9c\xad\x6b\x28\xdb\x86\xb3\x7d\xf3\x65\xdb\x2c\x0b\xb9\x6e\xbd\x86\xe3\x46\x9c\xba\x96\xfb\xc8\xb6\xf4\x59\x23\xb9\x82\x31\x50\x2f\x90\xa1\xfc\xee\x07\xe3\x38\x53\x92\x33\x6b\x36\xdc\x56\xf8\x82\x14\x75\x30\x73\x8a\x1b\x0d\x80\x97\x10\x62\xd0\x02\x51\x72\x56\xe2\xe7\xd6\x7f\x01\x00\x00\xff\xff\x19\x97\xa2\x91\xa9\x03\x00\x00"),
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
