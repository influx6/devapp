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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\x3d\x6e\xe4\x30\x0c\x85\x7b\x9f\x82\x60\x3f\xd0\x05\x64\x9f\x60\x8b\x05\x16\x5b\x0f\x34\x32\x6d\x0b\x91\x29\x43\x94\x32\x09\x0c\xdf\x3d\xb0\xe5\x9f\x14\x09\x52\xa5\x12\x29\x92\xef\x7d\x00\x39\xcf\x2d\x75\x8e\x09\xd0\x06\x4e\xc4\x09\x97\xa5\xd2\x42\x36\xb9\xc0\x60\xbd\x11\xa9\xd1\x87\xde\x31\xf4\xd1\xb5\xb7\x2e\x7b\x5f\x22\x4b\x9c\x28\x42\x60\xba\x3d\x87\xe0\xe9\x8a\x6e\x03\xb9\x7e\x48\xd8\x54\x00\xba\x0b\x71\x3c\x74\xb2\x50\xbc\x6f\x1f\x9b\x62\x09\x37\xb1\x87\x0f\xf6\xe5\x72\x40\x30\x1b\x40\x8d\x4a\x48\xc4\x05\x56\x4c\x4f\x84\x91\xd2\x10\xda\x1a\xa7\x20\x45\x7e\x35\x70\xe4\x5b\xa1\xf4\xd9\x84\xcd\x48\xf7\xa3\xb0\x37\x02\x68\xc7\x53\x3e\xfb\x4e\x5c\x84\xf4\x3e\x51\x8d\x89\xde\x12\xc2\x3a\x7a\x89\x20\xbc\x1a\x9f\xa9\x46\x84\xc9\x1b\x4b\x43\xf0\x2d\xc5\x1a\xff\x1f\xf5\x1d\x42\x1d\x66\xbf\x09\x35\x19\x91\x67\x88\xed\x77\x50\x7f\x8f\xfa\x0f\x50\xa7\xf5\x23\xa7\x14\x78\x77\x92\xfc\x18\xdd\xe9\xb5\x67\xcd\x9f\x75\x51\x5a\x95\xce\x73\xd0\x1c\xbc\xe2\x7a\xce\x13\xc2\x10\xa9\x5b\x97\x55\xd2\xe6\xdf\xf6\x6a\x65\xbe\x22\xd1\x6a\xdd\x7b\x53\x69\xb5\x9f\x59\x53\xcd\x33\x71\xbb\x2c\xd5\x75\x8d\x62\xa3\x9b\x92\xac\xd7\x08\xa0\x4b\x06\x12\x6d\x8d\xd8\x68\x55\xf2\x6b\xee\x23\x00\x00\xff\xff\x4d\x3a\xb5\x09\xc6\x02\x00\x00"),
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
