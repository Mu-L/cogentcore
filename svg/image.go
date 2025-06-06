// Copyright (c) 2021, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

import (
	"errors"
	"image"
	"log"

	"cogentcore.org/core/base/iox/imagex"
	"cogentcore.org/core/base/slicesx"
	"cogentcore.org/core/math32"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

// Image is an SVG image (bitmap)
type Image struct {
	NodeBase

	// position of the top-left of the image
	Pos math32.Vector2 `xml:"{x,y}"`

	// rendered size of the image (imposes a scaling on image when it is rendered)
	Size math32.Vector2 `xml:"{width,height}"`

	// file name of image loaded -- set by OpenImage
	Filename string

	// how to scale and align the image
	ViewBox ViewBox `xml:"viewbox"`

	// Pixels are the image pixels, which has imagex.WrapJS already applied.
	Pixels image.Image `xml:"-" json:"-" display:"-"`
}

func (g *Image) SVGName() string { return "image" }

func (g *Image) SetNodePos(pos math32.Vector2) {
	g.Pos = pos
}

func (g *Image) SetNodeSize(sz math32.Vector2) {
	g.Size = sz
}

// pixelsOfSize returns the Pixels as an imagex.Image of given size.
// makes a new one if not already the correct size.
func (g *Image) pixelsOfSize(nwsz image.Point) image.Image {
	if nwsz.X == 0 || nwsz.Y == 0 {
		return nil
	}
	if g.Pixels != nil && g.Pixels.Bounds().Size() == nwsz {
		return g.Pixels
	}
	g.Pixels = imagex.WrapJS(image.NewRGBA(image.Rectangle{Max: nwsz}))
	return g.Pixels
}

// SetImage sets an image for the bitmap, and resizes to the size of the image
// or the specified size. Pass 0 for width and/or height to use the actual image size
// for that dimension. Copies from given image into internal image for this bitmap.
func (g *Image) SetImage(img image.Image, width, height float32) {
	if img == nil {
		return
	}
	img = imagex.Unwrap(img)
	sz := img.Bounds().Size()
	if width <= 0 && height <= 0 {
		cp := imagex.CloneAsRGBA(img)
		g.Pixels = imagex.WrapJS(cp)
		if g.Size.X == 0 && g.Size.Y == 0 {
			g.Size = math32.FromPoint(sz)
		}
	} else {
		tsz := sz
		transformer := draw.BiLinear
		scx := float32(1)
		scy := float32(1)
		if width > 0 {
			scx = width / float32(sz.X)
			tsz.X = int(width)
		}
		if height > 0 {
			scy = height / float32(sz.Y)
			tsz.Y = int(height)
		}
		pxi := g.pixelsOfSize(tsz)
		px := imagex.Unwrap(pxi).(*image.RGBA)
		m := math32.Scale2D(scx, scy)
		s2d := f64.Aff3{float64(m.XX), float64(m.XY), float64(m.X0), float64(m.YX), float64(m.YY), float64(m.Y0)}
		transformer.Transform(px, s2d, img, img.Bounds(), draw.Over, nil)
		if g.Size.X == 0 && g.Size.Y == 0 {
			g.Size = math32.FromPoint(tsz)
		}
	}
}

func (g *Image) DrawImage(sv *SVG) {
	if g.Pixels == nil {
		return
	}
	pc := g.Painter(sv)
	pc.DrawImageScaled(g.Pixels, g.Pos.X, g.Pos.Y, g.Size.X, g.Size.Y)
}

func (g *Image) LocalBBox(sv *SVG) math32.Box2 {
	bb := math32.Box2{}
	bb.Min = g.Pos
	bb.Max = g.Pos.Add(g.Size)
	return bb.Canon()
}

func (g *Image) Render(sv *SVG) {
	vis := g.IsVisible(sv)
	if !vis {
		return
	}
	g.DrawImage(sv)
	g.RenderChildren(sv)
}

// ApplyTransform applies the given 2D transform to the geometry of this node
// each node must define this for itself
func (g *Image) ApplyTransform(sv *SVG, xf math32.Matrix2) {
	rot := xf.ExtractRot()
	if rot != 0 || !g.Paint.Transform.IsIdentity() {
		g.Paint.Transform.SetMul(xf)
		g.SetProperty("transform", g.Paint.Transform.String())
	} else {
		g.Pos = xf.MulVector2AsPoint(g.Pos)
		g.Size = xf.MulVector2AsVector(g.Size)
	}
}

// ApplyDeltaTransform applies the given 2D delta transforms to the geometry of this node
// relative to given point.  Trans translation and point are in top-level coordinates,
// so must be transformed into local coords first.
// Point is upper left corner of selection box that anchors the translation and scaling,
// and for rotation it is the center point around which to rotate
func (g *Image) ApplyDeltaTransform(sv *SVG, trans math32.Vector2, scale math32.Vector2, rot float32, pt math32.Vector2) {
	crot := g.Paint.Transform.ExtractRot()
	if rot != 0 || crot != 0 {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, false) // exclude self
		g.Paint.Transform.SetMulCenter(xf, lpt)
		g.SetProperty("transform", g.Paint.Transform.String())
	} else {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, true) // include self
		g.Pos = xf.MulVector2AsPointCenter(g.Pos, lpt)
		g.Size = xf.MulVector2AsVector(g.Size)
	}
}

// WriteGeom writes the geometry of the node to a slice of floating point numbers
// the length and ordering of which is specific to each node type.
// Slice must be passed and will be resized if not the correct length.
func (g *Image) WriteGeom(sv *SVG, dat *[]float32) {
	*dat = slicesx.SetLength(*dat, 4+6)
	(*dat)[0] = g.Pos.X
	(*dat)[1] = g.Pos.Y
	(*dat)[2] = g.Size.X
	(*dat)[3] = g.Size.Y
	g.WriteTransform(*dat, 4)
}

// ReadGeom reads the geometry of the node from a slice of floating point numbers
// the length and ordering of which is specific to each node type.
func (g *Image) ReadGeom(sv *SVG, dat []float32) {
	g.Pos.X = dat[0]
	g.Pos.Y = dat[1]
	g.Size.X = dat[2]
	g.Size.Y = dat[3]
	g.ReadTransform(dat, 4)
}

// OpenImage opens an image for the bitmap, and resizes to the size of the image
// or the specified size -- pass 0 for width and/or height to use the actual image size
// for that dimension
func (g *Image) OpenImage(filename string, width, height float32) error {
	img, _, err := imagex.Open(filename)
	if err != nil {
		log.Printf("svg.OpenImage -- could not open file: %v, err: %v\n", filename, err)
		return err
	}
	g.Filename = filename
	g.SetImage(img, width, height)
	return nil
}

// SaveImage saves current image to a file
func (g *Image) SaveImage(filename string) error {
	if g.Pixels == nil {
		return errors.New("svg.SaveImage Pixels is nil")
	}
	return imagex.Save(g.Pixels, filename)
}
