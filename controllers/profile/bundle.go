package profile


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
        
          "profile.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "profile.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xcc\xdb\x09\xc3\x30\x0c\x85\xe1\x77\x83\x77\x30\x5a\x40\x0b\xa8\x9a\xa0\x4b\x98\x5a\x69\x0d\x46\x82\x58\x29\x14\xe1\xdd\x4b\x2f\x90\xd7\x73\xf8\xbf\x88\x26\x5b\x57\x29\xa0\xf5\x09\x6b\xe5\x44\xc7\xe0\x9c\x4a\xa1\xd1\x99\x6a\x79\xec\xb2\x5d\x00\xa7\xcc\xd9\x4d\xb1\xc9\xf4\xdd\x5e\xc0\x57\xbb\xdb\xe1\x84\x95\x09\x47\xe7\x9c\x08\xbf\x61\x84\x68\xfb\x38\xa7\x7c\x33\x75\x51\x87\xdf\xfa\xbf\xdf\x01\x00\x00\xff\xff\x6e\xa0\x29\x72\x79\x00\x00\x00"),
          path: "profile.html",
          root: "profile.html",
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
