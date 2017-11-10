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
        
          "twoauth.html",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "twoauth.css": { // all .css assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00"),
          path: "twoauth.css",
          root: "twoauth.css",
        },
      
    
      
        "twoauth.html": { // all .html assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\x4f\x6b\xdc\x30\x10\xc5\xef\x81\x7c\x87\x61\xce\xdd\xf5\xa5\xf4\xd0\xae\x02\xa5\x2d\xa5\xb7\xfe\xa1\xe7\x30\x96\xc6\xb6\x88\xac\x31\xd2\x38\xdb\xd4\xf8\xbb\x17\xd9\xeb\xb0\x74\x37\xd0\x5e\xcc\x93\x1f\xf3\x7b\xc3\x93\xa6\xc9\x71\xe3\x23\x03\x46\x7a\xc4\x79\xbe\xbd\x99\x26\x8e\x6e\x15\x9b\x65\x25\x2a\x47\x5d\xec\x43\x66\xab\x5e\x22\xd8\x40\x39\x1b\xd4\xa3\xec\x1a\xb2\x2a\x09\xda\xe4\xdd\xae\x0e\x62\x1f\x56\xd9\x8c\x21\xe0\xdd\xed\x0d\xc0\xdf\x43\x35\xa5\x9d\x15\xc7\x17\x23\x1d\x85\x06\x24\xf2\x22\xd6\x59\x80\x83\xef\x5b\xc8\xc9\x1a\x74\xa4\xf4\xd6\xf7\xd4\x72\x35\xc4\xf6\x5d\x4d\x99\xdf\xbc\x7e\x35\x4d\xfb\x9f\x99\xd3\xfe\xdb\xf7\x79\x46\xa8\xd6\xc4\xea\x14\x79\x35\x5f\x8f\xb2\xee\xfc\x3f\x5b\x0c\x77\x9f\xa2\x72\x82\xf7\xa3\x76\x92\xfc\x6f\x5a\x88\x1f\x0a\xa0\x49\xd2\xc3\x67\x91\x36\xf0\x62\x73\x54\x6f\x49\x25\x1d\xaa\x61\x1b\x6f\x24\xf5\xd7\xf3\x77\xc5\x42\xa0\x65\x45\x83\x08\x3d\x6b\x27\xce\xe0\x16\x5d\x2a\x88\xc3\xa8\xa0\x4f\x03\x1b\xec\xbc\x73\x1c\x11\x22\xf5\x6c\x70\xcc\x9c\xee\x8b\x44\x78\xa4\x30\xb2\xc1\xad\x8f\xf2\x29\xc6\x3c\xff\x3b\x68\x18\xeb\xe0\xed\xbd\x77\x17\xb4\xaf\x8b\xf3\xe5\xe3\x8b\x34\xe5\x5f\x7a\xce\x2a\xcd\x3c\x53\x10\x86\x40\x96\x3b\x09\x8e\x93\xc1\xb5\xab\x52\xd5\x52\xe0\x0b\xc4\x3c\xd6\xbd\x7f\x66\x6e\xa7\x13\xf1\xc7\x7a\xdc\xea\xad\x4a\x89\x17\x57\x7f\xae\xaf\x3c\xec\xac\x4f\x81\xf3\xfa\xae\x83\x8f\x0f\x90\x38\x98\xd3\xdf\x8e\x59\x11\xba\xc4\x8d\xc1\x2a\x2b\xa9\xb7\x95\x1e\x85\x46\xed\xaa\x9e\xb2\x72\xda\xdb\x9c\xf1\x1c\xfc\x27\x00\x00\xff\xff\x3a\xd1\x7a\x22\x4d\x03\x00\x00"),
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
