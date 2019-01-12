// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/reports/html/commons.gohtml
// templates/reports/html/default.gohtml
// templates/reports/html/timelog.gohtml

package tom


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}


type asset struct {
	bytes []byte
	info  fileInfoEx
}

type fileInfoEx interface {
	os.FileInfo
	MD5Checksum() string
}

type bindataFileInfo struct {
	name        string
	size        int64
	mode        os.FileMode
	modTime     time.Time
	md5checksum string
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) MD5Checksum() string {
	return fi.md5checksum
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _bindataReportsHtmlCommonsgohtml = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x58\xeb\x6e\xdb\xb8\x12\xfe\x9f\xa7\x20\x74\x50\xa0\x2d\x4c\xd9\x4e" +
	"\x4f\x4e\x5b\x35\x35\x4e\xaf\xd8\xfe\xe8\x62\x11\x77\x1f\x80\x12\x69\x99\x35\x2f\x02\x39\x4e\xea\x15\xfc\xee\x0b" +
	"\xea\x62\x4b\x14\xa5\xb8\x68\xf3\x27\x14\x39\x9c\xf9\xe6\xf6\x91\x74\x59\x52\xb6\xe1\x8a\xa1\x68\xbd\x97\x92\x98" +
	"\x43\x74\x3c\x5e\x21\x84\x50\x59\x62\x34\x7f\x9e\x6b\x38\x14\x2c\x41\x39\x87\xed\x3e\x8d\x33\x2d\xe7\xdf\x89\xb2" +
	"\xda\xe4\xf3\x5c\x03\x97\xac\xfd\x67\x58\xa1\x0d\xc4\x77\xcc\xee\x05\xbc\xdf\x67\x3b\x06\xcf\xe7\x08\x1f\x8f\x57" +
	"\x95\xb6\x5b\x20\xa9\x60\x28\x13\xc4\xda\xb7\x91\x6d\x4c\xad\xaa\xb5\x7a\x3d\xd5\xf4\x70\xfe\x2e\x4b\xbe\x89\xbf" +
	"\xd8\xbf\x8c\xfe\xce\xb2\x46\x5f\x03\xec\xbc\xc5\xac\x7a\x13\xf5\x24\x5d\x95\x25\x5f\xbe\x52\x28\x6a\x36\x47\xc7" +
	"\xe3\xed\x1c\x68\x50\xb6\x05\x04\x1c\x04\x8b\x56\x65\x19\x7f\x73\xa3\xd0\x8e\xdb\x79\xd7\x5e\x59\x32\x45\x3b\x80" +
	"\x1c\x5e\xa4\x34\xa0\xf8\x23\x01\x76\x47\x54\xce\xe2\x4f\xb2\x80\xc3\xcf\x82\xfe\xc6\x25\x43\xc6\xed\x4f\x2e\x03" +
	"\x2e\x19\xae\xc4\x2b\xf4\x67\xe3\x5f\xb9\xe2\x92\x88\x35\x18\xae\xf2\x9f\xf7\x67\x80\xb3\x87\xd1\x90\x6c\xc7\x28" +
	"\x72\xc6\xc3\x28\x7d\x84\x82\x6d\xc0\x01\x94\x5c\x7d\xdc\x1b\x02\x5c\x2b\x14\xb7\x23\x5f\x81\x8f\xac\xaa\x84\x3b" +
	"\xbd\x57\x94\xd1\x9f\x0d\xe7\xa7\x1f\x24\x03\x04\x8f\x02\xbe\x18\x74\xa5\x70\x0c\xf9\x2f\xc6\xf5\xb3\x21\x92\xd9" +
	"\xd1\x88\xae\xca\x72\xa3\x8d\x24\xf0\xe7\x5e\xa6\xcc\xa0\xb8\x92\xff\xa0\xf7\x0a\xa6\x43\x78\x3b\xef\x74\xd7\xed" +
	"\xbc\x6a\xc5\xd5\x55\x0b\xed\xea\x4c\x01\x99\x96\x52\xab\x0f\xeb\x75\x4b\x02\xb7\x16\x0e\x82\x21\xc7\x01\x6f\x23" +
	"\x60\x3f\x60\x9e\x59\xdb\x69\xdb\xc4\x68\x0d\xa8\xec\x21\xc5\x78\xa3\x15\x60\xcb\xff\x61\x09\x5a\x5e\x17\xf0\x26" +
	"\xb4\xbc\x21\x92\x8b\x43\x82\x22\x7b\xb0\xc0\x24\xde\xf3\x68\x86\x30\x29\x0a\xc1\x70\x3d\x35\x43\xef\x05\x57\xbb" +
	"\xaf\x24\x5b\x57\xdf\x9f\xb5\x82\x19\x8a\xd6\x2c\xd7\x0c\xfd\xfd\x25\x9a\xa1\x3b\x9d\x6a\xd0\x33\xf4\x07\x13\xf7" +
	"\x0c\x78\x46\x66\xe8\x9d\xe1\x44\xcc\x90\x25\xca\x62\xcb\x0c\xdf\xcc\x50\xf4\xce\x29\x45\x1f\xb4\xd0\x06\x7d\x92" +
	"\xfa\xbb\xb3\x74\x52\x13\x98\x59\x1f\x64\xaa\x45\xe4\xc3\xae\x6a\xa2\x87\x5d\x6a\xa5\x6d\x41\x32\x36\x2e\xfa\xc0" +
	"\x78\xbe\x85\x04\x29\x97\x38\xf1\xe6\xca\x13\xcc\x1c\xa8\x04\xa5\x82\x64\x3b\x5f\x49\x9a\xb7\xcb\x0f\x5b\x0e\x01" +
	"\x1b\x20\x58\x2b\xf1\x9f\x17\xcb\xff\xdd\xa4\xff\xf5\x65\x8c\x7e\xc0\x9a\xd2\x93\x14\xab\xfe\x06\x28\xaa\x82\xc0" +
	"\x0f\x9c\xc2\x36\x41\x2f\x17\x4f\x06\xa6\xaa\x75\x4a\x80\xb4\x42\xaf\x6f\x46\x84\x60\xdb\x1a\xbb\x27\xe6\x69\xe3" +
	"\xe0\xb3\x51\xd9\x34\xef\x49\xb7\x2e\x8f\x6f\x68\xe3\x99\x6a\x41\x47\x85\x5c\xe4\xeb\xf2\x5b\xc4\xaf\x0c\x93\x03" +
	"\x87\x9b\xf3\x07\xa7\xda\x50\x66\xfa\x88\x7b\x31\xeb\x00\x39\x9e\x95\x6c\x41\x0a\xaf\xe8\x3b\x25\x5f\xab\x39\x4d" +
	"\x78\xbe\x3c\x1a\x9d\x94\x64\xbb\xdc\x38\xb6\xc3\x97\xc4\xa6\x57\x90\x1d\xd3\xf5\x54\x18\x7f\x7d\x16\xf7\x1d\x68" +
	"\x12\x5b\x6b\xe8\x14\x84\x8f\xee\x14\x31\x41\x0a\xcb\x12\xd4\x8e\x7a\x86\xce\x96\xb6\x33\x04\xd4\x33\x25\xb8\x62" +
	"\x78\xdb\x24\x72\x19\x5f\xdf\x54\x29\xea\x4a\x14\x84\x52\xae\x72\x97\xbf\x6a\x15\x2d\x7b\x22\xc7\x8e\xfe\xb8\xea" +
	"\x03\x67\xa5\x1e\x79\xb6\x06\x9a\x16\xe1\x88\x6c\xfd\x70\xb8\x96\xc3\x55\x73\xbb\xe6\x7d\x30\xa4\xe8\x43\x74\x64" +
	"\x88\x89\xe0\xb9\x4a\x90\x3b\x28\xfa\xab\xf7\xcc\x38\x3e\x12\xad\x44\xaa\x01\xb4\x0c\x9b\xf6\xc3\xe3\xef\x05\x5d" +
	"\x8c\x60\x66\x84\x0e\x91\x8f\xd5\x8f\xdf\x73\x53\x75\xd9\xef\xe5\x50\xc5\x75\x4b\x7d\xd0\x77\xa1\x0d\x6d\xe3\x7a" +
	"\x5b\xea\xe9\x70\x99\xc6\x67\xde\x79\xbc\x58\xcf\xec\x34\xa9\x4c\x53\x8a\xc0\x24\xca\xf9\xb6\xe5\x82\x3e\xbd\x56" +
	"\x78\xf9\x6c\x98\x84\xb1\x28\x5e\xc0\x0e\xb1\x73\x8b\xd3\x41\x83\xf5\x2b\xaa\x3e\x0e\x02\x4e\xbd\x7e\xfd\x24\xac" +
	"\xb6\x2a\x07\x6c\x5c\xbc\x3c\xd5\xdd\x5a\xac\xd6\x47\x70\x05\xfa\xa3\x9f\xf6\xf3\x91\x32\x95\xc2\x21\xf7\x76\x0f" +
	"\xfc\xd8\xef\xd5\x33\x80\x86\x76\xc3\xc9\x24\x7b\xd0\x21\xae\x49\xd0\x8b\xe2\x07\xb2\x5a\x70\xda\xc0\x0c\xb1\xb7" +
	"\x87\x57\x12\x93\x73\x95\xa0\x05\x5a\xa0\x17\xe3\x6d\x7f\x42\x34\xc8\xff\x89\x37\xe2\x9b\x09\x77\x52\xd0\x40\xfc" +
	"\xa3\xa0\xb6\x8d\x41\x17\xc9\x80\xb7\x02\x9b\x81\xce\xba\x5f\x7e\x3b\x3f\x12\xf9\x26\x06\x95\xb1\xeb\x53\x9c\xbc" +
	"\xfb\x44\x8f\x34\x62\x77\x39\xb9\xbc\x80\xd0\x63\x64\xd8\xa3\xb2\x90\xf2\xc0\x09\xe5\x5f\xa5\x2e\x20\x0c\xef\x46" +
	"\xf5\x6c\x9c\x8d\x7f\xc5\x81\x8d\xbb\xce\xfe\x1e\x27\x7e\x13\xa4\x38\xad\xde\xbe\xfe\x7f\x1b\x2e\x3a\x77\x12\x25" +
	"\x55\xc5\x4f\x2a\x0b\xed\xad\x9b\xe5\x7a\xaa\x59\x80\x18\xc0\x94\x1c\x66\xed\xd8\x85\xa0\xfa\xd0\x05\x0e\x04\xed" +
	"\x62\x1f\x95\x06\xe6\x7b\x84\xa5\xc5\xdb\x43\xb1\x65\xca\x86\xc8\x01\x3f\xb0\x74\xc7\x61\x4a\x64\x64\xa9\x79\xd9" +
	"\xcc\xab\xa7\xcd\xea\x6a\xe2\x9d\x83\x24\xa3\x9c\xbc\x8d\x0a\xc3\x15\x74\x5e\x3d\xff\x2f\x48\xee\xfb\x5a\xb3\xdf" +
	"\x10\x44\x1b\xdc\x97\xf1\x8d\x0c\x27\xc5\x3d\xcd\x26\x9f\x50\x8b\x22\x4c\xe8\x13\xd7\xb8\xe5\xc2\xbf\xc5\x3b\xc8" +
	"\x38\x35\x8c\xec\x30\x57\x96\x53\x07\xf6\x5e\x73\x1a\xd6\x6c\x06\x74\xf8\xd8\x6e\x4f\x8a\x6c\xc0\x51\xb7\x1f\x79" +
	"\xef\xf2\x32\x6a\x24\x65\x1b\x6d\x2e\x36\x32\x29\x34\x81\x97\x72\x5b\x08\x72\x48\xea\x48\x62\x07\x89\x19\xec\x8e" +
	"\xfd\x11\x7a\x08\xa4\xea\x37\x81\x1e\x90\xd0\x2f\x5a\xb9\xdc\x6b\x77\xa3\x19\xb8\x8c\xba\x1d\x12\xfa\xa1\x60\x6f" +
	"\x41\xcb\xfa\x87\x82\x66\xf9\xdf\x00\x00\x00\xff\xff\x12\x5c\xce\xb4\x49\x14\x00\x00")

func bindataReportsHtmlCommonsgohtmlBytes() ([]byte, error) {
	return bindataRead(
		_bindataReportsHtmlCommonsgohtml,
		"reports/html/commons.gohtml",
	)
}



func bindataReportsHtmlCommonsgohtml() (*asset, error) {
	bytes, err := bindataReportsHtmlCommonsgohtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "reports/html/commons.gohtml",
		size: 5193,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1546272929, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataReportsHtmlDefaultgohtml = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x56\x39\x6f\xe3\x38\x14\xee\xf5\x2b\xb8\xc4\x56\x41\x24\x6d\xba\x20" +
	"\xa0\x0c\x6c\x8e\xc5\x16\x53\xe5\x00\xa6\xa5\xc5\x67\x8b\x33\x3c\x3c\xe2\x93\x1d\x43\xd0\x7f\x1f\x90\xb2\x25\xd9" +
	"\xa6\x13\x17\x19\x37\x92\xde\x49\x7e\xdf\x3b\xdc\xb6\x29\xc9\xaf\x96\x16\xb7\x2b\xb8\x23\x4b\x89\x55\x33\xcf\x4a" +
	"\xab\xf3\x1f\xdc\x38\x5b\x2f\x73\xb4\x3a\x5f\xda\xd4\x3f\x6a\x58\xd9\x1a\xb3\x67\x70\x8d\xc2\xfb\xa6\xfc\x09\x78" +
	"\x95\x93\xb4\xeb\x92\xa4\x6d\x05\x2c\xa4\x01\x42\x7b\x39\xed\xba\x84\x10\x42\xbe\x26\x7c\x1f\xe9\xef\x79\x90\x91" +
	"\xbb\x82\x64\xa3\x50\x01\x5f\xb8\x20\xfb\x06\x7c\xf1\x50\x49\x25\x6a\x30\xa3\xde\x55\x76\xf3\xf4\xce\xcb\xde\xef" +
	"\x7f\xee\x9e\x6d\x63\x04\x88\x89\x65\x30\x65\x42\xae\x49\xa9\xb8\x73\x05\xed\x13\xd1\x59\x50\xf4\x71\xe4\x82\xf4" +
	"\xa9\x76\x91\xf7\x3f\x86\x7c\xae\x60\xef\x18\x3e\x52\x2b\xc4\xc4\x77\xb4\xac\x80\x8b\x98\xbc\x3e\x15\xee\x1c\x86" +
	"\xb8\x55\xba\x91\x02\xe8\x8c\xb9\x15\x37\x83\x54\xa2\x02\x3a\x6b\xdb\xec\xd5\xbf\x75\x1d\xcb\xbd\x7a\xc6\x72\xac" +
	"\x3e\x0f\x29\x75\xf0\x95\x37\xb7\x86\xd0\xc7\xa6\xe6\x28\xad\xa1\x3e\xc8\x39\xef\x1e\x85\x01\xd0\x23\x24\x2e\xca" +
	"\xd3\x33\x21\x2e\xcb\x06\x46\x44\x72\xb0\x3c\x06\x98\x8f\x73\x06\xde\xb9\x15\xdb\x53\x79\xdb\xd6\xdc\x2c\x21\xce" +
	"\xea\xe8\x7c\x86\x9b\x5e\x29\x0e\xb0\xc7\x48\xf6\x89\xed\x31\x20\x5a\x9a\x3d\xea\x24\xdb\xbf\x7d\x16\xe7\x72\x0a" +
	"\x2e\xc9\x1a\x62\x5c\x9e\x3a\xce\x07\x39\xcb\xc9\x07\x0c\xc6\x39\x61\xb8\xb0\x16\xf7\x47\x76\xcd\x1c\x2d\x72\x15" +
	"\xed\xa5\xf3\x3d\x33\xd4\xda\x6b\x70\xfe\xa8\xc4\x22\x65\x3a\xc5\x67\x37\x70\x0e\xc9\xf9\x13\xbd\x11\x4b\x7a\xca" +
	"\xcd\x57\xf5\x89\xc7\xf8\x50\xc1\xf2\x30\xb8\xa6\xf3\xae\x8f\x79\x38\x00\xb3\x27\xbd\xc2\xed\x51\xaa\xa0\x31\x16" +
	"\x49\xf6\xc8\x11\x9e\x7d\x4f\x45\xed\xc8\xd1\x8c\x0d\xc3\x8b\x84\x1e\x0c\x23\x6c\xf0\xf6\x77\x15\x72\x3d\x4b\x3e" +
	"\xba\x64\xdb\x82\x72\x70\x3c\x8b\x4f\x46\xb8\x8b\xd4\xce\xbe\xf1\x77\x8b\xe6\x5c\xeb\x8f\xf7\x3a\x77\x9b\xd1\x12" +
	"\x41\xaf\x14\xc7\x71\xfb\x0d\x0b\xea\xd4\x36\x4e\x56\x4c\x7e\x04\xc3\xd4\x64\xa7\x1a\x68\x62\x7f\x09\x5b\xfa\x1d" +
	"\x4b\x2a\xd4\x6a\x96\x30\xff\x20\x8a\x9b\x65\x41\xdb\x36\xbc\xdc\x73\x07\xa4\xeb\xa8\x57\x0e\x83\x92\x69\x40\x4e" +
	"\xca\x8a\xd7\x0e\xb0\xa0\x6f\xaf\xff\xa5\xb7\x74\xaa\x32\x5c\x43\x41\xd7\x12\x36\x7e\x37\x53\x52\x5a\x83\x60\xb0" +
	"\xa0\x1b\x29\xb0\x2a\x04\xac\x65\x09\x69\xf8\xb8\x26\x8d\x83\x3a\x75\x25\x57\xbe\x98\x0a\x63\xaf\x89\x34\x12\x25" +
	"\x57\x41\x08\xc5\x4d\xf6\xcf\x35\xd1\xfc\x5d\xea\x46\x1f\x88\xa4\x39\x14\x1d\x1c\xa1\x42\x5c\xa5\xf0\xab\x91\xeb" +
	"\x82\x7e\x4f\xdf\xfe\x4d\x1f\xac\x5e\x71\x94\x73\x05\x93\xf3\x48\x28\x40\xf8\x52\xda\x6d\xf3\x50\x5e\x7e\x13\x86" +
	"\x67\x92\x1c\x13\x55\x5a\xad\xad\x79\x78\x79\x19\xff\xa9\x4c\x94\x8d\x43\xab\x77\x4a\x96\xf7\x78\xb1\x7e\x66\x4d" +
	"\xed\x5e\x1a\xad\x79\xbd\xa5\xfb\x5a\xea\xba\x24\x5a\x0c\x83\x96\xe5\x7d\x10\x96\x07\x9e\x7e\x07\x00\x00\xff\xff" +
	"\xa9\xaf\x14\xc8\x7d\x09\x00\x00")

func bindataReportsHtmlDefaultgohtmlBytes() ([]byte, error) {
	return bindataRead(
		_bindataReportsHtmlDefaultgohtml,
		"reports/html/default.gohtml",
	)
}



func bindataReportsHtmlDefaultgohtml() (*asset, error) {
	bytes, err := bindataReportsHtmlDefaultgohtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "reports/html/default.gohtml",
		size: 2429,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1547324211, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

var _bindataReportsHtmlTimeloggohtml = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x56\xdd\x6b\xe4\x36\x10\x7f\xdf\xbf\x62\x2a\xfa\x14\x62\xbb\x47\x29" +
	"\x1c\x87\x1d\x68\xf3\x41\x0b\xa5\x94\x6c\x0e\xfa\xaa\xb5\x66\x6d\xf5\xf4\xe1\x5a\xe3\xcd\x2d\x66\xff\xf7\x22\xc9" +
	"\xde\xf5\xe6\xec\x24\x7b\x7e\x91\x35\xa3\xf9\xcd\x97\x66\x46\x7d\x9f\x40\x76\x55\x59\xda\x37\xf8\x09\x2a\x49\x75" +
	"\xb7\x49\x4b\xab\xb3\x7f\xb9\x71\xb6\xad\xb2\xca\x92\xd4\x38\x2e\x2d\x36\xb6\xa5\xf4\xb7\xae\xfc\x82\xf4\x18\x36" +
	"\x57\x19\x24\x87\xc3\x6a\xd5\xf7\x02\xb7\xd2\x20\xb0\xc8\x65\x87\xc3\x0a\x00\xe0\x4d\x0d\x64\x75\x56\xd9\xc4\x2f" +
	"\x03\xfc\x23\xba\x4e\x51\x84\x19\xe0\x23\xd2\x8f\x9b\x40\x83\x4f\x05\xa4\x27\xa2\x42\xbe\x75\x81\xf6\x27\xf2\xed" +
	"\x6d\x2d\x95\x68\xd1\x9c\xf8\xae\xb6\xcf\xf7\x5f\x79\x19\xe5\x7e\xe7\xee\xd1\x76\x46\xa0\x98\x9c\x0c\x47\x73\x21" +
	"\x77\x50\x2a\xee\x5c\xc1\xa2\x22\x76\x13\x18\x11\x47\x6e\x21\xbd\xd7\x0d\xed\x07\xe4\xf1\xcb\x89\x6f\x14\x8e\x82" +
	"\x61\x93\x08\x4e\x1c\xe2\xaf\x15\x62\x02\x73\x12\xaa\x91\x8b\x39\x7a\xfb\x2d\x31\xea\xcf\xae\x50\x6f\x50\x08\x14" +
	"\x40\x92\x14\x02\x59\xf8\x82\xd8\x0c\x3b\x6e\x44\xd4\xe8\x80\x6c\x85\x54\x63\x0b\xd2\x40\xd3\x4a\x43\x28\xe0\xef" +
	"\xbb\x07\x77\x95\xbd\x30\x7d\x62\x0d\x94\x56\xb9\x86\x9b\x82\xfd\xc2\x8e\xbe\x78\x60\x76\xd3\xf7\xe9\x93\xff\x3b" +
	"\x1c\xf2\x8c\xea\x19\x9b\xb3\x39\xa3\x17\x3d\x09\xca\xa2\x02\xc1\x29\xe0\xcb\x0f\x1f\x0d\xb0\x3b\xbf\x5b\x52\xf2" +
	"\xaa\xe4\x9a\x78\x4b\xdf\x27\x7a\x6f\xc4\x7b\x05\x7d\x05\x4c\xac\xed\x5a\x4e\xd2\x9a\x77\x4b\xd7\xc9\xb3\x14\x13" +
	"\x80\xbf\x2c\xa1\x5b\x94\x5e\x08\x6a\xb6\x78\x6f\x36\x56\xec\xbf\xa5\xf7\x7d\xcb\x4d\x85\x90\x3e\xb4\x5c\xa3\x1b" +
	"\x96\xc5\x6b\xb0\x90\xb2\xc8\x14\xa3\x2b\xce\x87\x3b\x11\x7c\xef\x9d\xd9\xda\x56\x73\xf2\xb9\x83\x34\xe4\x21\x78" +
	"\x34\x63\xe2\x22\xd0\x18\xd7\x88\xf4\x24\xf5\x77\x21\xd9\x66\x00\x5a\x3c\x0e\xc7\x32\xfe\xc3\xad\xc9\x36\x0d\x8a" +
	"\x85\x40\xcc\x4a\x48\x53\x29\xbc\xe3\x2f\xcb\x7f\x59\x6e\xea\xcf\xbd\x79\x9f\x2e\x54\x0e\x2f\x54\xe0\x43\x7f\xa1" +
	"\x92\x37\xcf\xbd\x75\xe6\xdd\x79\x19\x73\xab\xa5\x19\x0b\x06\xd2\xf1\xef\x82\xfc\x9a\x50\x2c\xbe\x15\x85\xb2\x79" +
	"\x4d\x72\xbe\x72\x96\x3c\xca\xb3\x85\xca\xc9\x69\x6b\x2d\x1d\xef\x57\xb7\x21\x4b\x5c\xcd\xb6\xf2\xd7\x1a\xdd\xd8" +
	"\x55\x7f\x3e\x15\xfe\x53\x00\xba\xb0\xe9\x2c\x07\x70\x19\xe3\xe6\xc2\xd6\xe2\x1d\x3e\x67\xe4\x59\x18\x29\xd3\x31" +
	"\x38\x73\x43\xa7\x83\x73\x66\x66\x08\xb9\x3b\x47\x3d\x36\xa5\x38\xe8\xe7\xda\x51\xdf\x13\xea\x46\xf9\xae\x32\xbe" +
	"\x27\x8e\x23\x7f\x62\xcb\x59\x4a\xa7\xfb\x41\xeb\x48\x5a\xe5\x3f\x08\x5b\xfa\x27\x08\xd4\xa4\xd5\xcd\x2a\xf7\x0b" +
	"\x28\x6e\xaa\x82\xa1\x61\x9e\x70\xec\xab\xb9\x46\xe2\x50\xd6\xbc\x75\x48\x05\xfb\xfc\xf4\x90\x7c\x64\x53\x96\xe1" +
	"\x1a\x0b\xb6\x93\xf8\xec\x9f\x2b\x6c\x62\x54\x69\x0d\xa1\xa1\x82\x3d\x4b\x41\x75\x21\x70\x27\x4b\x4c\xc2\xe6\x1a" +
	"\x3a\x87\x6d\xe2\x4a\xae\x7c\x48\x0b\x63\xaf\x41\x1a\x49\x92\xab\x40\xc4\xe2\x43\xfa\xd3\x35\x68\xfe\x55\xea\x4e" +
	"\x9f\x91\xa4\x39\x27\x9d\x19\x53\x13\x35\x09\xfe\xd7\xc9\x5d\xc1\xfe\x49\x3e\xff\x9a\xdc\x5a\xdd\x70\x92\x1b\x85" +
	"\xec\x64\x8f\xc4\x02\x45\x35\x76\xc7\x3c\x64\xc9\x5f\x8f\xb0\xae\x56\x2f\x43\x5e\x5a\xad\xad\xb9\x5d\xaf\x4f\xaf" +
	"\xb8\x09\xb3\x73\x64\xf5\xc0\xcc\xb3\x18\xb8\x3c\x56\xd1\xf4\xdc\xba\xd3\x9a\xb7\x7b\x36\xe6\xf9\x70\x58\xcd\xa6" +
	"\xf5\xc8\xcd\xb3\x08\x92\x67\x21\x49\xff\x07\x00\x00\xff\xff\x8c\xe9\x88\x64\x9c\x0a\x00\x00")

func bindataReportsHtmlTimeloggohtmlBytes() ([]byte, error) {
	return bindataRead(
		_bindataReportsHtmlTimeloggohtml,
		"reports/html/timelog.gohtml",
	)
}



func bindataReportsHtmlTimeloggohtml() (*asset, error) {
	bytes, err := bindataReportsHtmlTimeloggohtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "reports/html/timelog.gohtml",
		size: 2716,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1547324009, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}


//
// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
//
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
// nolint: deadcode
//
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

//
// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or could not be loaded.
//
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrNotExist}
}

//
// AssetNames returns the names of the assets.
// nolint: deadcode
//
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

//
// _bindata is a table, holding each asset generator, mapped to its name.
//
var _bindata = map[string]func() (*asset, error){
	"reports/html/commons.gohtml": bindataReportsHtmlCommonsgohtml,
	"reports/html/default.gohtml": bindataReportsHtmlDefaultgohtml,
	"reports/html/timelog.gohtml": bindataReportsHtmlTimeloggohtml,
}

//
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
//
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, &os.PathError{
					Op: "open",
					Path: name,
					Err: os.ErrNotExist,
				}
			}
		}
	}
	if node.Func != nil {
		return nil, &os.PathError{
			Op: "open",
			Path: name,
			Err: os.ErrNotExist,
		}
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}


type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{Func: nil, Children: map[string]*bintree{
	"reports": {Func: nil, Children: map[string]*bintree{
		"html": {Func: nil, Children: map[string]*bintree{
			"commons.gohtml": {Func: bindataReportsHtmlCommonsgohtml, Children: map[string]*bintree{}},
			"default.gohtml": {Func: bindataReportsHtmlDefaultgohtml, Children: map[string]*bintree{}},
			"timelog.gohtml": {Func: bindataReportsHtmlTimeloggohtml, Children: map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}