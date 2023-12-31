// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!cgo && darwin) || (cgo && !arm64 && darwin)
// +build !cgo,darwin cgo,!arm64,darwin

// Use non-cgo on Darwin to prevent duplicate symbols on cgo

package x509

func loadSystemRoots() (*CertPool, error) {
	return execSecurityRoots()
}
