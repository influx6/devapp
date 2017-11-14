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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xcf\x6a\xc3\x30\x0c\xc6\xef\x81\xbc\x83\xd0\xbd\xf8\x05\x9c\x3c\xc1\x0e\x83\xb1\x73\x71\x1d\x25\x31\x73\xe4\x60\x39\xeb\x46\xe8\xbb\x8f\xfc\x71\xd3\x43\xc7\xd8\xa9\x52\x25\x7d\xdf\x8f\xf8\x9b\xe7\x86\x5a\xc7\x04\x68\x03\x27\xe2\x84\xb7\x5b\x59\x68\x21\x9b\x5c\x60\xb0\xde\x88\x54\xe8\x43\xe7\x18\xba\xe8\x9a\x53\x3b\x79\xbf\x55\x96\x38\x51\x84\xc0\x74\xba\xf6\xc1\xd3\x51\x9d\x7a\x72\x5d\x9f\xb0\x2e\x0b\x00\xdd\x86\x38\x64\xa1\x49\x28\x9e\xd7\x3f\x56\xc9\xad\x5c\xd5\x2e\x3e\xd8\x8f\xc3\x02\xc1\xac\x04\x15\x2a\x21\x11\x17\x58\x31\x5d\x11\x06\x4a\x7d\x68\x2a\x1c\x83\xec\xfa\x8b\x83\x23\xdf\x08\xa5\x47\x17\x36\x03\x9d\xf3\x20\x6f\x02\x68\xc7\xe3\x74\x5f\xbc\x13\x23\xa4\xef\x91\x2a\x4c\xf4\x95\x10\x96\xdb\x43\x05\xe1\xd3\xf8\x89\x2a\x44\x18\xbd\xb1\xd4\x07\xdf\x50\xac\xf0\x3d\xcf\x33\x86\xca\x76\xbf\x71\x8d\x46\xe4\x1a\x62\xf3\x7f\xae\x7c\x99\xd9\x8e\xfe\x39\xdb\x6b\x9e\xff\xc9\x76\x10\x5c\xa6\x94\x02\xef\x86\x32\x5d\x06\x77\xff\x14\x7b\x57\xbf\x2c\x8f\xa6\xd5\xb6\x79\x5c\x9a\xcc\x2d\xae\xe3\x69\x44\xe8\x23\xb5\xcb\xcb\x6d\x6d\xfd\xb6\xfe\x6a\x65\x9e\xd3\x68\xb5\xc4\xa0\x2e\x0b\xad\xf6\xdc\xd5\x65\x31\xcf\xc4\xcd\x92\xc5\x23\xa1\x62\xa3\x1b\x93\xac\x09\x05\xd0\x5b\x0b\x12\x6d\x85\x58\x6b\xb5\xf5\x8f\xa7\x3f\x01\x00\x00\xff\xff\xcc\xc8\xb0\x51\xdd\x02\x00\x00"),
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
