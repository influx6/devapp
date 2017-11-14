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
    
      ".css": []string{  // all .css assets.
        
          "twoauth.css",
        
      },
    
      ".html": []string{  // all .html assets.
        
          "twoauth-qr.html",
        
          "twoauth.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "twoauth.css": { // all .css assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x8f\xdd\x6a\xc4\x20\x10\x85\x5f\xc5\x9b\x42\x0b\x35\xb8\x91\x85\x45\x9f\x66\xfc\xa9\x4a\xe3\x8e\x98\x59\xd6\x22\x79\xf7\x52\x93\x40\x61\x2f\x06\x66\x0e\x1f\xe7\x63\x26\x7a\x22\xff\x02\x4b\x58\x3b\x61\x51\x57\xb1\x4d\x06\x2a\xb7\xe8\x3c\x4b\x39\xf4\x67\x72\x14\xd5\x55\xbc\xe9\xe8\x53\x88\x34\x56\x97\xd6\xb2\xc0\x8f\x32\x0b\xda\x6f\x9d\xa1\x86\x74\x57\x82\xc1\x83\x50\x1b\xac\xce\x57\x75\x11\xa5\xb1\x15\x97\xe4\x58\x0d\x06\xde\xe7\x79\xfe\x3c\xe7\xf2\xa1\x33\x34\xbe\x57\x4b\x21\x4a\x1b\xf7\x21\x18\xc1\x36\xad\x05\xac\xe7\x06\x89\x30\x73\x29\xfa\x6e\x39\x02\x25\xff\x31\x84\xe5\x85\x1b\x00\xfb\x7b\xe6\x41\x84\x77\x2e\x6f\xb1\x9f\xfd\xb7\xd2\xb6\xdf\x00\x00\x00\xff\xff\x3e\x4a\x0c\x93\xf8\x00\x00\x00"),
          path: "twoauth.css",
          root: "twoauth.css",
        },
      
    
      
        "twoauth-qr.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x92\xc1\x8a\xdb\x30\x10\x86\xef\x0b\x79\x87\x41\xe7\x3a\x2a\x24\xf4\xb0\x55\x0c\x7b\x69\xaf\xdd\xdd\xf6\x01\x26\xf2\xd8\x12\x3b\xd1\x18\x69\x5c\x13\x8c\xdf\xbd\xb8\x4e\xb6\xdb\x92\x42\x4f\xbd\x48\xc3\xc0\xff\xcf\xa7\xd1\x3f\x4d\x0d\xb5\x31\x11\x98\x84\xdf\xcd\x3c\x6f\xee\xdc\xc0\xf5\xe6\x0e\xc0\x71\xac\x1d\x42\xc8\xd4\x1e\x8c\xed\xb3\xb4\x91\xc9\xd4\x5f\xd6\xc2\x59\xac\x9d\xe5\x58\x6f\xee\x9c\xfd\xa9\x98\x26\x4a\xcd\x62\xf0\xcb\xd2\x4b\x52\x4a\xba\xda\x16\xf2\x1a\x25\x81\x67\x2c\xe5\x60\x74\x94\xaa\x45\xaf\x92\xa1\xcb\xb1\xa9\x8e\x2c\xfe\x65\x2d\xdb\x81\x19\x24\x51\x35\x06\x59\x46\x2e\x34\xd3\x14\x5b\xd8\x7e\x2b\x94\xb7\xcf\x44\xe9\xeb\x28\x9f\x56\xf1\xe2\x0d\xf0\xa7\xfb\x11\x73\xe5\xa5\xa1\x7f\xf2\x06\x70\x61\x77\x55\x96\x1e\x3d\x55\x47\x51\x95\x53\xb5\x7b\x0f\x9e\x92\x52\xae\x90\x63\x97\x4c\xfd\x1c\x64\xe0\x06\x90\x33\x61\x73\x86\x23\x41\x4c\x70\x96\x21\xc3\x67\x91\x8e\x09\x1e\x06\x0d\x94\x34\x7a\x5c\xe0\x1e\xfa\xde\xd9\xb0\x5b\xf7\x69\x2f\x8c\x97\xf7\x10\x17\xfa\xcf\xf0\x1e\x13\x3c\x3e\xc1\x18\x35\xdc\xe4\x7d\x65\x05\x70\xf1\xd4\x41\xc9\xfe\x60\x1a\x54\xbc\x8f\x27\xec\xc8\xf6\xa9\xfb\x78\xc4\x42\x1f\xf6\xef\xa6\x69\xfd\x8c\xc7\xa7\x79\x36\x60\x5f\x51\xf6\xbf\xa3\xa8\xf4\x7f\xc5\x59\xf4\xf7\x70\x35\x5a\x8e\x84\x27\x9a\x67\x67\xc3\xfe\xf6\xc6\xd6\x7c\xbd\x6d\xdf\x08\x5d\xd1\x33\x53\x59\x33\xc7\x31\xbd\x40\x26\x3e\x5c\xba\x81\x48\xcd\x35\xd2\x45\x51\xa3\xaf\x96\x90\x66\x61\xa6\x5c\xac\x8e\x82\x83\x86\xeb\xbd\xf5\xa5\x98\xb7\x53\x7e\x04\x00\x00\xff\xff\xa1\xf5\xae\xe8\x2e\x03\x00\x00"),
          path: "twoauth-qr.html",
          root: "twoauth-qr.html",
        },
      
        "twoauth.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\x31\x8f\xd4\x30\x10\x85\xfb\x93\xf6\x3f\x8c\xa6\xcf\xa5\xa0\xa1\x48\x56\x42\x80\x10\x1d\x0d\xf5\xc9\x6b\x4f\xce\xd6\x3a\x9e\xc8\x1e\xdf\x71\x44\xfb\xdf\x91\xed\x04\x56\x90\x82\x26\x7e\xc9\xb3\xbf\x71\xe6\xcd\xba\x1a\x9a\x5c\x20\xc0\xa0\x5e\xf0\x76\x3b\x3d\xac\x2b\x05\xd3\xc4\x6e\x69\x0e\x42\x41\xaa\x3d\x24\xd2\xe2\x38\x80\xf6\x2a\xa5\x11\xe5\x95\xbb\x49\x69\xe1\x08\xcf\xd1\x99\xee\xe2\x59\x5f\x9b\x9c\xb2\xf7\xc0\x81\xba\x57\xcb\x9e\xf0\x7c\x7a\x00\x38\x38\xde\x4e\x77\x9a\x0d\x1d\x23\xaa\xd2\x14\x84\xe2\xdf\x38\x80\x61\x39\x7f\xae\xce\x87\x2c\x96\xa3\xfb\xa9\x2a\xfe\x63\xa1\x4d\x91\x67\xf8\xc2\xfc\xec\xa9\xda\x14\xc4\x69\x25\x1c\x87\x7e\xd9\x8f\x4f\x1c\xe7\xe3\xcb\x74\xc5\x42\x50\xf5\xbe\x23\xf6\x89\x52\x72\x1c\xfa\xdf\xbb\xfa\xc9\x05\x97\x2c\xc2\x4c\x62\xd9\x8c\xb8\x70\x92\xfd\x5e\x00\x83\x0b\x4b\x16\x90\xb7\x85\x46\xb4\xce\x18\x0a\x08\x41\xcd\x34\x62\x4e\x14\x9f\x8a\x44\x78\x51\x3e\xd3\x88\xeb\xfa\xf8\x3d\x51\xac\x8f\x62\xdc\x6e\xff\x0f\x5a\xf2\xc5\x3b\xfd\xe4\xcc\x3f\xb4\x6f\xd5\xf9\xfa\xe9\x80\xb6\xfd\xf2\x86\xa6\x1f\x82\x5b\x81\xa6\x1b\x5e\xf8\x5a\x6a\x6d\x54\x84\xc5\x2b\x4d\x96\xbd\xa1\x38\x62\x6b\x6c\xe9\x6b\xed\xf6\x5d\x85\x4b\x16\xf9\x13\xf1\xf6\xd6\x96\xee\xdd\x7b\xbb\xd3\x53\xbe\xcc\x4e\xf0\x7c\x17\x0d\x0d\x7d\xdb\xb7\xc7\xd3\x97\x10\xda\xe4\xf4\xdb\xe8\x9c\x4f\x0f\xf7\xfa\x60\x5e\x93\xbc\x79\x4a\x6d\x5c\xbd\x0b\x57\x88\xe4\xc7\xed\xab\x25\x12\x04\x1b\x69\x2a\x91\x8a\x12\xa7\xbb\x32\xdf\x91\xbd\xa7\x98\x4a\xba\x2a\x8b\xdd\xd7\x47\x9d\x12\xde\x57\xf9\x15\x00\x00\xff\xff\x1e\x4e\x48\xb3\x31\x03\x00\x00"),
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
