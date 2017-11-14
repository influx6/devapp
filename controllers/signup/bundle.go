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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xcd\x8a\xdc\x30\x0c\xc7\xef\x03\xf3\x0e\x42\xf7\x59\xbf\x80\xed\xcb\x5e\x7b\x28\x94\x3d\x2f\x1e\x5b\x99\x98\x3a\x76\xf0\xc7\x4e\x4b\xc8\xbb\x97\xc4\x71\xb2\xb0\x53\xfa\xc1\xde\x2c\x4b\xff\xbf\x7e\x02\x69\x9a\x0c\x75\xd6\x13\xa0\x0e\x3e\x93\xcf\x38\xcf\xe7\x13\x4f\xa4\xb3\x0d\x1e\xb4\x53\x29\x09\x4c\xf6\xe6\xcb\x08\xb7\x68\xcd\xa5\x2b\xce\xd5\x97\x26\x9f\x29\x42\xf0\x74\xb9\xf7\xc1\xd1\xf1\xba\xf4\x64\x6f\x7d\x46\x79\x3e\x01\xf0\x2e\xc4\xa1\x39\x95\x44\xf1\x75\xfd\xa8\x9e\xf5\xbd\xda\x5d\x5d\xd0\xdf\x8f\x1e\x08\x6a\x65\x10\xc8\x16\x51\x62\x9e\xee\x08\x03\xe5\x3e\x18\x81\x63\x48\x9b\xfd\xd2\xc0\x92\x33\x89\xf2\xfb\x26\x5e\x0d\xf4\xda\x12\xad\x12\x80\x5b\x3f\x96\xbd\x70\x07\x46\xc8\x3f\x47\x12\x98\xe9\x47\x46\x58\xb4\x87\x0b\xc2\x9b\x72\x85\x04\x4e\xd3\xd3\xb3\xd2\x3d\x99\x97\x44\xf1\xe9\x65\x4b\xcf\x33\xc2\xe8\x94\xa6\x3e\x38\x43\x51\x60\x4b\xec\x78\xac\x61\xfc\x8e\x77\x54\x29\xdd\x43\x34\xff\xce\xdb\x94\x8d\xf9\x88\x1f\x32\x7f\xdd\xd2\x1f\x98\x5b\xe2\x3f\x98\x75\xf0\x9d\x8d\xc3\xe7\xb1\x37\xc7\x3f\xcc\xf0\x5c\xab\x3e\x8c\xb2\xfd\xc3\xdf\x8f\x74\x00\x5f\x4b\xce\xc1\x6f\x7c\xa9\x5c\x07\xbb\x6f\xc3\x16\xc9\x6f\xeb\xda\x72\x56\x4b\x1f\x7b\x73\xb6\x6c\xb5\x3c\x9f\x38\xdb\x0e\x49\x9e\x4f\xd3\x44\xde\x2c\xc7\x75\x9c\x9c\x57\x6f\xf5\xdc\x8a\xab\x3a\x67\x25\x57\xd0\x47\xea\x04\x32\x94\x5f\xc2\xcd\x7a\xce\x94\xe4\xcc\xd9\xd5\x6e\x2d\x7c\xe0\x94\x74\xb4\x63\x4e\xab\x1b\x00\xaf\x21\xa4\xa8\x05\xa2\xe4\xac\xc6\xef\xa5\xbf\x02\x00\x00\xff\xff\xb8\x71\xbc\x46\xf8\x03\x00\x00"),
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
