// Code generated by vfsgen; DO NOT EDIT.

package fixtures

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Fixtures statically implements the virtual filesystem provided to vfsgen.
var Fixtures = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2024, 9, 4, 6, 23, 35, 127490235, time.UTC),
		},
		"/controlplane.yaml": &vfsgen۰CompressedFileInfo{
			name:             "controlplane.yaml",
			modTime:          time.Date(2024, 9, 4, 7, 20, 11, 578220700, time.UTC),
			uncompressedSize: 395,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x74\x8f\x4f\x8f\x82\x30\x10\xc5\xef\x7e\x8a\x09\xf7\xf5\x03\xf4\xd6\x54\x0f\x24\x22\x66\xad\x7b\x25\xa3\x0c\x48\x2c\x94\x4c\xa7\xfb\x27\xc6\xef\xbe\x01\xc2\xca\x1e\xbc\xb5\xaf\xbf\xf7\x5e\xdf\xad\xe9\x4a\x05\x69\x90\xc6\xe7\x3d\x31\x8a\xe7\x55\xe8\xe9\xa2\x56\x00\x2d\x85\xab\xf1\x5d\xd5\xd4\xc3\x0d\x80\x3a\x3c\x3b\x3a\xb0\x6f\x49\xae\x14\x43\x46\x5c\x93\x82\x0a\x5d\xa0\x11\x28\xa9\xc2\xe8\x64\xe9\x01\xe8\xd9\x7f\xff\x64\x24\x58\xa2\xe0\x2c\x02\xa4\x47\x9b\xe6\x45\xb6\xb5\xba\xd8\xec\x8f\x85\xd1\x07\x7b\x7a\xdf\x2a\x48\x84\x23\x25\xaf\x30\x7d\xb2\x79\xa1\x77\xbb\xdc\x68\xbb\x84\x3f\xd1\x45\x0a\x53\x7a\xed\xfc\x19\xdd\xdc\x34\x8c\x48\x37\x0a\x92\xfb\x1d\xd6\x1f\x23\xb6\xce\x46\x0d\x1e\x8f\xb9\xa7\x8d\x4e\x1a\xe3\x62\x10\xe2\xe7\x17\xa7\xbd\xa5\x82\xa1\xe5\x4f\xbd\x4c\xd8\x1e\x5b\xfa\x9f\x6a\x9e\x0f\x8b\xe8\x8e\xe4\xcb\xf3\x4d\xcd\x87\xb7\x97\x8e\xdf\x00\x00\x00\xff\xff\x8e\xbd\x91\x33\x8b\x01\x00\x00"),
		},
		"/controlplane1.yaml": &vfsgen۰CompressedFileInfo{
			name:             "controlplane1.yaml",
			modTime:          time.Date(2024, 8, 22, 9, 9, 25, 283382742, time.UTC),
			uncompressedSize: 726,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x84\x91\x41\x6f\x9b\x40\x10\x85\xef\xfc\x8a\x91\xef\x76\x95\xeb\xde\x10\x89\x5a\xa4\x38\xb6\x1a\xd2\x1e\xad\x01\xc6\x78\xeb\xdd\x9d\xd5\xec\x60\xd7\x8a\xf2\xdf\x2b\x03\x6e\x88\xaa\x2a\x9c\x98\x99\xef\x3d\xf4\x1e\x18\xed\x0f\x92\x64\x39\x18\xb0\x21\x29\x3a\xb7\xb2\x49\x2d\xaf\x2c\x7f\x39\xdd\xa1\x8b\x07\xbc\xcb\x3c\x29\xb6\xa8\x68\x32\x80\x80\x9e\x0c\x0c\xcc\xb2\xe1\xa0\xc2\x2e\x3a\x0c\x34\x9d\x52\xc4\xe6\xef\x3d\x5d\x92\x92\xcf\x00\x1c\xd6\xe4\xd2\x55\x0e\x10\x85\x4f\xb6\x25\x31\x20\x03\xd5\xa8\xcb\x8e\x36\xb4\x06\xca\xeb\xb8\x89\x24\xa8\x2c\x59\x8a\xd4\x5c\x15\x51\x78\x6f\x1d\x19\x68\xc9\x73\x06\xe0\x29\x1d\x0a\x0e\x7b\xdb\x8d\x7e\x14\xb0\x76\xb4\x15\xf6\xa4\x07\xea\xd3\x9a\xa4\x23\x03\x7b\x74\x89\x06\xa0\xa5\x3d\xf6\x4e\xe7\x9a\xc1\xf5\xf7\x65\x3d\xcb\x35\x3e\xe5\x73\x55\x6e\x76\xeb\x87\x2a\xdf\xdd\x3f\x3d\xef\x8a\x7c\x5b\xbd\x7c\x7f\x30\xb0\x50\xe9\x69\xf1\x3f\x2c\x7f\xa9\x36\xbb\xfc\xf1\x71\x53\xe4\xd5\x1c\x6e\xd8\x47\x0e\x14\x74\x8a\x4e\x9d\x50\x4a\x5f\x51\xe9\x8c\x97\x69\xb7\xfc\xd0\xe8\x48\x74\x23\x31\x7d\x6e\x0c\xd8\xce\x23\xd9\xf0\xa9\xd3\x84\x7c\x62\x75\x42\xd7\xd3\xa4\x4f\xb6\xa5\x06\xa5\x0c\xbf\xa8\x51\x96\x9f\x54\x1f\x98\x8f\xb7\x6e\x84\xce\x62\x95\xf2\x18\xbf\x55\xd5\x76\x2b\x5c\x7f\x28\xb9\x73\x5c\xa3\xbb\xc1\xd7\x7f\x54\xde\x1b\x58\xbc\xbe\xc2\x6a\x3d\x0c\xf0\xf6\x76\xeb\xcf\xf7\x4e\x6d\xe1\xfa\xa4\x24\xef\xd5\x37\xe3\xe2\x69\xc8\x30\x08\x8b\xf7\xcd\x4c\x1d\x48\xcf\x2c\x47\x73\x7b\x59\xfe\x8b\xfe\x09\x00\x00\xff\xff\x58\x24\x92\x8c\xd6\x02\x00\x00"),
		},
		"/eastwest-gateway.yaml": &vfsgen۰CompressedFileInfo{
			name:             "eastwest-gateway.yaml",
			modTime:          time.Date(2024, 9, 4, 7, 21, 3, 600292715, time.UTC),
			uncompressedSize: 1133,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x94\x93\x41\x8f\xd3\x3e\x10\xc5\xef\xf9\x14\x23\xfd\xcf\xe9\x9f\x2c\x8b\x40\xbe\x21\x88\x50\x85\xd8\xc2\x6e\xe8\x1e\xab\x49\x32\x4d\x4c\x1d\xdb\xf2\x4c\x5a\x55\xab\xfd\xee\x28\x6d\xd2\xa6\xad\x4a\xe1\x16\xcf\x3c\xff\xf4\xc6\x93\x87\x5e\xcf\x29\xb0\x76\x56\x81\xb6\x2c\x68\xcc\x44\xb3\x68\x37\xd1\xee\xff\x75\x82\xc6\xd7\x98\x44\x2b\x6d\x4b\x05\xd3\xae\x3e\xf3\x14\x50\x5c\x88\x1a\x12\x2c\x51\x50\x45\x00\x16\x1b\x52\x40\xc8\xb2\x21\x96\x88\x3d\x15\x5d\xd9\x07\xb7\xd4\xa6\xeb\x34\x5e\xb6\x11\x40\xe1\x1a\xef\x2c\x59\xe1\xae\x0d\xa0\x6d\x15\x88\xf9\x0b\x0a\x6d\x70\xdb\x17\x01\xe2\x1e\xb8\x33\x12\x0f\xd8\x6a\xaf\xea\x35\x00\x06\x73\x32\xea\x70\x84\xbd\xfc\x68\xe3\x5c\x0f\x80\xde\xdf\x82\x02\x88\xf3\xce\xb8\x6a\x7b\x7c\x06\x4b\xb2\x71\x61\xa5\xa0\xff\x88\x5f\x5e\x60\x32\x47\xd3\x12\x4f\x3e\x99\x96\x85\xc2\x03\x36\x04\xaf\xaf\x07\x0c\x59\xcc\x0d\x95\x0a\x24\xb4\x74\xa8\xae\x3e\xf0\xd8\x2f\xd9\xf5\xf8\x08\xf0\x1f\x48\xc0\xe5\x52\x17\x20\x75\x70\x6d\x55\x83\xd4\x9a\xa1\xf7\x08\x5c\xbb\xd6\x94\x90\x13\x04\xd7\x0a\x95\xdd\xc2\x74\x49\x20\x35\x0d\xd6\x4e\x70\xc3\x33\x4e\x9f\xb2\xe9\x6c\xf1\x2d\xcd\x3e\x2e\x1e\xd3\x1f\x3f\xd3\xa7\x2c\xfd\xbc\x78\x48\xb3\xe7\xd9\xe3\xd7\xc5\x7c\x9a\x3e\x9f\xdc\x02\x58\x77\xa3\xfd\xcb\xb4\x00\x4c\x61\xad\x0b\x3a\x1d\xc7\xbb\x20\xac\xce\xe0\x83\x29\x16\x94\x96\xe3\x4e\x73\xa6\xd8\x5f\x54\x90\xbc\x7b\x73\x97\x5c\xf4\x04\x43\x45\xf2\xfd\xaa\x62\xe0\x8b\xe1\xab\xdc\xfb\xfb\xb7\x37\xb8\x97\x8a\x11\x37\xde\xfd\x19\xe5\x75\xdb\xc9\xdd\x2d\xdb\x17\x8a\x31\x7e\x43\x79\xed\xdc\xea\x0f\xfc\xf7\x37\xf9\x9d\x62\xb7\xc7\x7e\x01\xd5\x59\xc4\xf6\x29\xe8\xf3\xd7\x37\x8f\xab\xd2\xf6\x17\x15\xa2\x9d\xcd\xa8\xf1\x06\x85\x14\x8c\x83\x52\x19\x97\xe3\x21\x7a\x7f\x1f\x8f\xdf\x01\x00\x00\xff\xff\x24\xf7\xfd\x61\x6d\x04\x00\x00"),
		},
		"/expose-service.yaml": &vfsgen۰CompressedFileInfo{
			name:             "expose-service.yaml",
			modTime:          time.Date(2024, 9, 4, 7, 20, 49, 155080083, time.UTC),
			uncompressedSize: 325,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x4c\x8e\x31\x6b\xc3\x30\x10\x85\x77\xff\x8a\x23\x63\xc1\x2e\x21\xe9\xa2\x2d\x53\x32\x14\x52\x6a\xa7\x6b\xb9\xca\x47\x22\x22\xe9\xc4\xdd\x35\x26\xff\xbe\x38\x36\x75\x34\x3d\xde\xfb\xf4\x71\x58\xc2\x17\x89\x06\xce\x0e\x32\xd9\xc0\x72\x0d\xf9\xdc\x04\xb5\xc0\x4d\xe0\xd7\xdb\x1a\x63\xb9\xe0\xa6\xba\x86\xdc\x3b\xd8\xa3\xd1\x80\xf7\x2a\x91\x61\x8f\x86\xae\x02\xc8\x98\xc8\x81\x17\x56\xad\x67\x45\x7d\x9e\xb9\x69\xd5\x82\x9e\x1c\x3c\xa4\xb5\xde\xd5\x28\x55\x5a\xc8\x8f\xbf\x95\x22\x79\x63\x19\x33\x4c\x88\x03\x42\xb5\x81\xd4\x16\x8d\x92\xdc\x48\x74\xa2\x6a\x28\x2c\x36\xe5\xf1\xe5\xdf\xf4\x43\xe2\x60\xfd\xb6\xdd\x6e\x96\xf6\x71\x97\x45\xfd\x6f\x8a\xb0\xb1\xe7\xe8\xa0\x7b\x6f\xe7\xd6\xa2\x2e\xa2\xc4\x3d\x39\xd8\x9d\xba\xe3\xf7\xc7\xae\x6d\xbb\xc3\xe7\xf1\xb4\x3f\xcc\xf3\x85\xd5\x9e\xd8\x1a\x56\x2f\x4d\x64\x8f\x71\xf5\x17\x00\x00\xff\xff\x54\x38\x3e\x7a\x45\x01\x00\x00"),
		},
		"/helloworld.yaml": &vfsgen۰CompressedFileInfo{
			name:             "helloworld.yaml",
			modTime:          time.Date(2024, 8, 22, 9, 9, 25, 283612495, time.UTC),
			uncompressedSize: 2045,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\xb4\x55\x4d\x6f\xd3\x40\x10\xbd\xfb\x57\x8c\xcc\x15\x37\xa9\xb8\xa0\xbd\x55\xa5\x54\x08\x5a\x45\x44\xea\x05\x71\x98\x6c\x26\xcd\xd2\xfd\x62\x77\xec\x2a\xaa\xfa\xdf\xd1\xda\xb1\xbb\x4e\x1d\x28\x45\xcc\x69\x35\x9e\x79\xf3\xde\x74\x5e\x83\x5e\xdd\x50\x88\xca\x59\x01\xcd\x69\x71\xa7\xec\x5a\xc0\x35\x1a\x8a\x1e\x25\x15\x86\x18\xd7\xc8\x28\x0a\x00\x8d\x2b\xd2\x31\xbd\x00\x54\x64\xe5\x2a\x65\x7f\x90\xe4\xb6\x97\x2c\xae\x34\xad\x0b\x00\x8b\x86\x04\x44\x34\x5e\x53\x51\x55\x55\x31\x39\xe2\xdc\xd9\x8d\xba\xbd\x42\x3f\x1a\xd1\xf5\x92\x6d\x2a\xd9\x7e\xdf\xa7\x5a\x2e\x03\x66\x5f\x6c\x28\x46\xbc\x25\x01\xe5\xc3\x03\x9c\xdc\xa0\xae\x29\x9e\x9c\xeb\x3a\x32\x85\xa4\x00\x1e\x1f\xcb\xe3\x04\x96\x14\x1a\x75\xa0\xb0\x1b\xbf\x25\xad\xdd\xbd\x0b\x7a\x3d\x39\x7e\xbc\x08\xf4\xfe\xa0\x03\x20\x76\xd0\xa3\x7c\xf4\x24\x53\x87\x77\x81\xdb\xd6\xaa\x7d\x0a\x78\x37\x9f\xcf\xdb\xae\xfd\x70\x66\x5f\x24\x08\x4d\x92\x5d\x98\x1e\x72\x28\x0a\xbd\x8f\xb3\x41\xd9\x07\xf2\xda\xed\x0c\x59\xfe\x1f\xe2\x9a\x6c\x93\xbd\xa8\x40\x5e\x2b\x89\x51\xc0\xe9\x33\xee\x06\x59\x6e\xbf\x64\xa0\xd3\xb0\x63\x60\x00\x26\xe3\x35\x32\xed\x41\x32\x1d\x29\xf4\x08\xef\x18\xe2\x21\x26\x40\x4f\x38\x85\x74\x96\x51\x59\x0a\x03\x4e\x35\xb5\xa3\x2e\x94\x69\x2f\x6d\xa5\x78\x8b\x8d\x93\xb3\xb6\xa4\x6a\x6b\x2a\xb2\xcd\x50\x17\x28\xba\x3a\x48\xca\xb8\xa5\xe4\xcf\x9a\x22\x8f\x72\x00\xd2\xd7\x02\xca\xd3\xf9\xdc\x94\xe3\x31\x8b\x5a\xeb\x85\xd3\x4a\xee\x04\x7c\xda\x5c\x3b\x5e\x04\x8a\x64\x19\xde\x9c\xe9\x7b\xdc\xc5\xa1\x7a\xb8\xa5\x9e\xfe\x20\x69\x31\x3e\xad\x14\x64\x9b\xbc\xb4\x53\x7a\x75\xb1\x5c\x9e\x5d\x5e\x64\xbc\x9a\x64\xa3\x8f\xc1\x99\x03\xb2\xbd\x5d\x3f\xd3\xee\x2b\x6d\xc6\x1f\x27\x7d\x9b\xc7\x1d\xed\x44\x6f\xd7\xd7\x38\x32\x6a\x22\xff\xd2\x7b\xed\x8b\x33\x1f\x76\xa9\xa3\x16\x7c\xff\x32\x03\x76\x28\x7f\xa2\x7f\x26\xa5\xab\x27\xad\x37\xdd\xff\x17\xde\x3d\xbe\x86\x7f\x31\x62\x8f\xfa\x0a\xc7\x65\xab\xce\x8c\xc5\x14\x8c\xb2\x98\x7e\x17\x2e\x03\x4a\x5a\x50\x50\x6e\xbd\x24\xe9\xec\x3a\x0a\xe8\x4f\x32\x8e\xf6\x75\x3d\x52\xf8\x7b\x7f\xe6\x55\x83\x35\x65\x1d\x74\xfb\x8c\xb3\xf4\x2c\x9e\x2e\xd7\x18\x4c\xab\xfd\x56\xce\x56\xca\xce\xda\xe6\xf2\x2d\x94\xca\x6e\x94\x55\xbc\x2b\xbf\xbf\xcc\x7e\x4f\xff\x54\x9c\xae\x0d\x5d\x25\xd6\x23\xf7\x99\x94\x59\x20\x6f\x05\xcc\x88\x65\x37\x69\xc6\x3a\x16\x87\x46\x89\x24\x03\x71\xd5\x01\x15\x39\xea\x73\xad\x13\xa5\xb0\xcf\xe6\x2e\xec\x32\xd9\x16\xab\x2e\x93\x95\x38\x9f\xfe\x24\xa8\x05\x70\xa8\x5b\x1f\xfe\x0a\x00\x00\xff\xff\x01\x57\x34\x04\xfd\x07\x00\x00"),
		},
		"/namespace-template.yaml": &vfsgen۰CompressedFileInfo{
			name:             "namespace-template.yaml",
			modTime:          time.Date(2024, 8, 22, 9, 9, 25, 283683288, time.UTC),
			uncompressedSize: 182,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x64\x8e\xb1\x0a\xc3\x30\x0c\x44\x77\x7f\x85\x7e\x20\x2e\x5d\xbd\x76\xef\x98\x5d\x69\x8e\x22\x62\x5b\xc1\x52\x5a\x42\xc8\xbf\x97\x84\x76\xea\xf6\xe0\xde\x1d\xc7\xb3\xf4\x68\x26\x5a\x13\xbd\xae\x61\x92\x3a\x26\xba\x73\x81\xcd\xfc\x40\x28\x70\x1e\xd9\x39\x05\xa2\xcc\x03\xb2\x1d\x44\x34\x2d\x03\x5a\x85\xc3\xa2\xe8\xe5\x67\xc5\xca\x05\x89\xc4\x5c\xb4\xb3\xd5\x1c\xe5\xb4\x5d\x67\xcd\xfa\x5c\xe3\x99\x1c\x8d\x0a\x7f\x6b\x9b\x12\x7d\xa1\xdb\x36\x8a\x3d\xe7\x05\x16\x6f\x79\x31\x47\x3b\x4e\xd0\xbe\x07\xa2\xff\xd5\x4f\x00\x00\x00\xff\xff\xb7\x76\x94\x04\xb6\x00\x00\x00"),
		},
		"/rafayremote-secret.yaml": &vfsgen۰CompressedFileInfo{
			name:             "rafayremote-secret.yaml",
			modTime:          time.Date(2024, 9, 4, 7, 21, 29, 487929970, time.UTC),
			uncompressedSize: 883,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x02\xff\x94\x92\x31\x8f\xdb\x30\x0c\x85\x77\xff\x0a\xe2\x76\xf9\xd0\x55\x5b\xe1\x6e\x05\xba\x5c\x7b\xbb\xea\x63\x52\x22\x96\x64\x90\x54\x5a\x23\xcd\x7f\x2f\x24\xd9\x8e\x9b\x34\x29\xce\x1b\x89\xc7\xc7\x8f\x4f\x76\x23\xbd\x22\x0b\xc5\x60\xe1\xf8\xa1\x39\x50\x78\xb3\xf0\x82\x3d\xa3\x36\x1e\xd5\xbd\x39\x75\xb6\x01\x70\x21\x44\x75\x4a\x31\x48\x2e\x01\x02\xea\xcf\xc8\x07\x0a\xfb\x96\x44\x29\xb6\x14\x9f\xfb\x21\x89\x22\x5b\x78\x3a\x9d\xa0\x7d\x75\x43\x42\x69\xbb\xda\xfc\xe2\x3c\xc2\xf9\xfc\xd4\x00\xf4\x8c\xc5\xe9\x2b\x79\x14\x75\x7e\xb4\x10\xd2\x30\x34\x00\x83\xfb\x8e\xc3\xbc\xa0\xb8\x3e\xfb\x34\x28\x75\xab\xaf\x72\xc2\x6c\x11\x9c\x47\x5b\x25\x86\xd1\x47\x45\x23\x05\xda\xdc\xdd\x3c\x4f\xc9\xe8\xfa\x75\x54\x26\x51\xf4\x8d\x28\x53\xd8\x7f\x9a\x4f\xbd\xeb\x60\xe1\x77\x21\xbb\xca\x2c\xb7\xe6\xcb\x67\x74\xb3\xd4\xb5\x2c\x02\x64\xa5\x1d\xf5\x4e\xd1\xb8\xa4\x3f\x22\x93\x4e\xa6\xa4\xbb\x5d\xf8\x82\x7c\x44\xee\x3e\x56\xde\xfa\x49\xe9\x59\xb8\xd5\x5d\x54\x35\x8f\x47\xb7\x03\xf4\x31\x28\xfe\xd2\x0b\x63\xad\x37\x8c\x4b\xcc\x8f\x7d\xf2\x97\xe4\x4a\xf7\x4d\xde\x8d\x93\x98\x31\xa8\x59\x30\xfe\x23\xaf\x7f\x66\x17\xc3\x8e\xf6\xa5\x31\x32\xee\x90\x31\xf4\x28\x16\x4e\x55\x94\xb1\x64\x39\xc8\xdc\x62\xfc\x4d\x39\x9f\xb1\x56\x39\x01\x2a\x48\x9b\xc7\xba\x79\xa2\xec\xd1\x21\xeb\xd6\x67\x9d\x3c\xe0\xf4\xef\x89\xcf\x38\xe5\x01\x63\x4c\xf3\x27\x00\x00\xff\xff\xa5\x77\xec\x56\x73\x03\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/controlplane.yaml"].(os.FileInfo),
		fs["/controlplane1.yaml"].(os.FileInfo),
		fs["/eastwest-gateway.yaml"].(os.FileInfo),
		fs["/expose-service.yaml"].(os.FileInfo),
		fs["/helloworld.yaml"].(os.FileInfo),
		fs["/namespace-template.yaml"].(os.FileInfo),
		fs["/rafayremote-secret.yaml"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(io.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
