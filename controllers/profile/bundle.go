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
    
      ".css": []string{  // all .css assets.
        
          "profile.css",
        
      },
    
      ".html": []string{  // all .html assets.
        
          "profile.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "profile.css": { // all .css assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\x2b\x28\xca\x4f\xcb\xcc\x49\xd5\xcd\xcc\x2b\x29\xca\xaf\x2e\x48\x4c\x49\xc9\xcc\x4b\xb7\x32\x34\x28\xa8\xb0\x2e\xc9\x2f\xb0\x32\x32\xa8\x05\x04\x00\x00\xff\xff\x94\x60\x03\x3d\x23\x00\x00\x00"),
          path: "profile.css",
          root: "profile.css",
        },
      
    
      
        "profile.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xcf\x8e\xd4\x30\x0c\xc6\xef\x95\xfa\x0e\x56\x2e\x0b\x87\x6e\xa4\x3d\xb7\x95\x90\x80\x13\x47\x10\xe7\x4c\xea\x4e\xa3\xf1\xd8\x10\x27\x54\x55\xd5\x77\x47\xfd\xc3\x0c\x8c\x10\x62\x2f\x89\x63\xe5\xb3\x7f\xf1\x97\xb2\x98\xe7\x0e\xfb\xc0\x08\x86\xdd\x0f\xb3\x2c\x65\x51\x67\x6a\xcb\x02\xa0\xa6\xd0\xd6\x0e\x86\x88\x7d\x63\xac\xa2\x6a\x10\xb6\x1d\x6a\x8a\x32\x99\xf6\x93\x9c\x25\xa7\xda\xba\xb6\xb6\x14\x36\xc5\x3c\x87\x1e\x9e\xbf\x28\xc6\x75\xf9\x3c\xca\x47\xe7\x93\x44\x58\xab\x3e\xd6\xcb\x8a\x51\x6d\x1a\xa5\xdf\xee\xd8\x2e\xa8\x3b\x11\x1a\x10\xf6\x14\xfc\xa5\x31\x5e\xb8\x0f\xf1\xfa\xe6\xe9\x5d\x44\x98\x24\x83\xe6\x23\x18\x83\x0e\x90\x04\x0e\x0d\xdc\x3a\x3d\xbd\x35\xed\xfb\xc7\xe4\x03\x22\x92\xe2\xff\x00\x21\xbf\x9a\x67\x97\xfc\x89\xf3\x81\xff\x4d\xc3\xdd\x36\x73\xbb\x0d\xfd\x76\xbe\xbb\xe2\x85\x13\x72\xda\x9d\x51\xf4\x29\x08\x83\x27\xa7\xda\x98\x6f\x51\xfa\x40\x58\x05\x4e\x51\xe0\x1c\x43\x57\x9d\x48\xfc\x65\x0f\xfb\x4c\xb4\x47\x1e\x39\x61\x04\x61\xac\xc6\x41\x08\x4d\x5b\x16\xf5\xf0\xd2\x7e\x45\xf2\x72\xc5\x15\x7d\x92\x1c\xe1\xa8\x07\xf3\xfc\x3d\x4b\xc2\xbb\x95\x91\xdd\x15\x97\xa5\xb6\xc3\xcb\x2a\xb5\x07\xc7\xdf\x89\x35\x4d\x84\xba\x03\x53\xe0\x0b\x44\xa4\xe6\xc8\x0e\x88\xc9\xdc\xbe\x54\x72\x29\xf8\x6a\x7d\x61\x14\xa2\x75\xfc\x07\xc1\xaf\xfd\xd9\xab\x9a\xdf\xbb\xfc\x0c\x00\x00\xff\xff\xcb\x81\xbd\xf9\xb0\x02\x00\x00"),
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
