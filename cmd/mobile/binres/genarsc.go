// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// Genarsc generates stripped down version of android.jar resources used
// for validation of manifest entries.
//
// Requires the selected Android SDK to support the MinSDK platform version.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"cogentcore.org/core/cmd/mobile/binres"
)

const tmpl = `// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by genarsc.go. DO NOT EDIT.

package binres

var arsc = []byte(%s)`

func main() {
	arsc, err := binres.PackResources()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("arsc.go", []byte(fmt.Sprintf(tmpl, strconv.Quote(string(arsc)))), 0644); err != nil {
		log.Fatal(err)
	}
}
