// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xyz

import (
	"testing"
)

func TestScene(t *testing.T) {
	t.Skip("todo: fixme")
	sc := NewOffscreenScene()
	sc.AssertImage(t, "scene")
}
