// Iris - Distributed Messaging Framework
// Copyright 2013 Peter Szilagyi. All rights reserved.
//
// Iris is dual licensed: you can redistribute it and/or modify it under the
// terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// The framework is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// Alternatively, the Iris framework may be used in accordance with the terms
// and conditions contained in a signed written agreement between you and the
// author(s).
//
// Author: peterke@gmail.com (Peter Szilagyi)

package overlay

import (
	"crypto/x509"
	"testing"
	"time"
)

// Another private key to check security negotiation
var privKeyDerBad = []byte{
	0x30, 0x82, 0x01, 0x3a, 0x02, 0x01, 0x00, 0x02,
	0x41, 0x00, 0xe1, 0xb7, 0xeb, 0x69, 0x12, 0x77,
	0xed, 0xe1, 0x89, 0x55, 0x11, 0x0a, 0x28, 0xee,
	0xfe, 0x2d, 0xc7, 0xe3, 0x9c, 0x03, 0xd1, 0x41,
	0x6b, 0xd5, 0x8e, 0xab, 0x6f, 0x79, 0x99, 0x7e,
	0x0e, 0xdf, 0xe9, 0xdf, 0xa0, 0x85, 0x1f, 0x87,
	0xed, 0xf8, 0xdc, 0xbf, 0x74, 0xcc, 0xfa, 0x3d,
	0xe6, 0x33, 0xcd, 0xd3, 0xeb, 0x4a, 0x90, 0xba,
	0x76, 0x97, 0x72, 0x43, 0x5c, 0x11, 0xb8, 0x6f,
	0x1f, 0xb7, 0x02, 0x03, 0x01, 0x00, 0x01, 0x02,
	0x41, 0x00, 0x91, 0x7d, 0x9e, 0x64, 0xdc, 0xbd,
	0xa6, 0xfc, 0x80, 0x2b, 0xef, 0xce, 0xcf, 0xb1,
	0xb4, 0xee, 0xc7, 0x70, 0x43, 0xc9, 0x06, 0x21,
	0x98, 0x23, 0x50, 0x51, 0xda, 0x89, 0xd3, 0xd1,
	0x38, 0x3e, 0x10, 0x19, 0x0f, 0x18, 0x09, 0x6f,
	0x96, 0xb8, 0xad, 0xdb, 0xe2, 0xba, 0x90, 0xf4,
	0xe6, 0x61, 0x56, 0xde, 0x51, 0x02, 0xe9, 0xf1,
	0x7b, 0xc3, 0x45, 0x19, 0x64, 0xa4, 0x3a, 0x97,
	0xaa, 0xa9, 0x02, 0x21, 0x00, 0xf7, 0x80, 0x1c,
	0x65, 0x71, 0x96, 0xbc, 0x2e, 0x3b, 0xa3, 0x94,
	0x10, 0xe6, 0x8a, 0x1e, 0xdb, 0xe5, 0x79, 0xbb,
	0xc4, 0xfe, 0x93, 0x97, 0x67, 0xe3, 0x9c, 0x55,
	0x68, 0x06, 0x72, 0xa4, 0x23, 0x02, 0x21, 0x00,
	0xe9, 0x78, 0x50, 0x2c, 0xe3, 0x2e, 0xde, 0x32,
	0x22, 0x34, 0x94, 0xe6, 0xa5, 0x42, 0x40, 0x59,
	0xfd, 0xa6, 0x5b, 0x51, 0x37, 0xa2, 0x78, 0x8b,
	0xcb, 0x40, 0x0c, 0x0a, 0x60, 0xbc, 0xf5, 0x5d,
	0x02, 0x20, 0x3a, 0xd9, 0x64, 0x67, 0x1e, 0x96,
	0x27, 0xd2, 0x9a, 0x6d, 0xaf, 0xc4, 0x40, 0xfc,
	0xf0, 0x8f, 0x7a, 0xc5, 0xf2, 0x46, 0xc9, 0xfd,
	0x90, 0x0f, 0xac, 0xc8, 0x3c, 0x6a, 0x8a, 0xb5,
	0xf1, 0x9f, 0x02, 0x20, 0x5a, 0x14, 0x47, 0xaa,
	0xea, 0x04, 0xbd, 0x4a, 0x4e, 0x94, 0x47, 0x17,
	0x2e, 0x8f, 0xff, 0x50, 0x39, 0x63, 0xa4, 0x3c,
	0x13, 0xe1, 0x25, 0xed, 0x9a, 0x00, 0x1e, 0x25,
	0x2a, 0xe9, 0xfd, 0x29, 0x02, 0x20, 0x07, 0x0c,
	0x31, 0x52, 0xda, 0x87, 0x67, 0xb7, 0xfb, 0x88,
	0x3a, 0xe7, 0xeb, 0x49, 0x39, 0x6b, 0xfd, 0x99,
	0xbc, 0xf2, 0x4c, 0x69, 0x06, 0x67, 0x3f, 0x57,
	0x8e, 0x51, 0xf4, 0x74, 0xa2, 0xaa,
}

// Id for connection filtering
var appIdBad = "overlay.test.bad"

func TestHandshake(t *testing.T) {
	// Make sure cleanups terminate before returning
	defer time.Sleep(3 * time.Second)

	key, _ := x509.ParsePKCS1PrivateKey(privKeyDer)
	bad, _ := x509.ParsePKCS1PrivateKey(privKeyDerBad)

	// Start first overlay node
	alice := New(appId, key, nil)
	if err := alice.Boot(); err != nil {
		t.Errorf("failed to boot alice: %v.", err)
	}
	defer alice.Shutdown()

	// Start second overlay node
	bob := New(appId, key, nil)
	if err := bob.Boot(); err != nil {
		t.Errorf("failed to boot bob: %v.", err)
	}
	defer bob.Shutdown()

	// Wait a while for the handshakes to complete
	time.Sleep(time.Second)

	// Verify that they found each other
	if size := len(alice.pool); size != 1 {
		t.Errorf("invalid pool size for alice: have %v, want %v.", size, 1)
	} else if _, ok := alice.pool[bob.nodeId.String()]; !ok {
		t.Errorf("bob (%v) missing from the pool of alice: %v.", bob.nodeId, alice.pool)
	}
	if size := len(bob.pool); size != 1 {
		t.Errorf("invalid pool size for bob: have %v, want %v.", size, 1)
	} else if _, ok := bob.pool[alice.nodeId.String()]; !ok {
		t.Errorf("alice (%v) missing from the pool of bob: %v.", alice.nodeId, bob.pool)
	}

	// Start a second application
	eve := New(appIdBad, key, nil)
	if err := eve.Boot(); err != nil {
		t.Errorf("failed to boot eve: %v.", err)
	}

	// Wait a while for any handshakes to complete
	time.Sleep(time.Second)

	// Ensure that eve hasn't been added (app filtering)
	if len(eve.pool) != 0 {
		t.Errorf("invalid pool contents for eve: %v.", eve.pool)
	}
	if _, ok := alice.pool[eve.nodeId.String()]; ok {
		t.Errorf("eve (%v) found in the pool of alice: %v.", eve.nodeId, alice.pool)
	}
	if _, ok := bob.pool[eve.nodeId.String()]; ok {
		t.Errorf("eve (%v) found in the pool of bob: %v.", eve.nodeId, bob.pool)
	}

	// Terminate eve and wait a while to finish
	eve.Shutdown()
	time.Sleep(time.Second)

	// Start a malicious node impersonating the app but invalid key
	mallory := New(appId, bad, nil)
	if err := mallory.Boot(); err != nil {
		t.Errorf("failed to boot mallory: %v.", err)
	}
	defer mallory.Shutdown()

	// Wait a while for any handshakes to complete
	time.Sleep(time.Second)

	// Ensure that mallory hasn't been added (invalid security key)
	if len(mallory.pool) != 0 {
		t.Errorf("invalid pool contents for mallory: %v.", mallory.pool)
	}
	if _, ok := alice.pool[mallory.nodeId.String()]; ok {
		t.Errorf("mallory (%v) found in the pool of alice: %v.", mallory.nodeId, alice.pool)
	}
	if _, ok := bob.pool[mallory.nodeId.String()]; ok {
		t.Errorf("mallory (%v) found in the pool of bob: %v.", mallory.nodeId, bob.pool)
	}
}