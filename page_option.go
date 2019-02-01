package gofpdf

import (
	"fmt"
	"io"
)

// Page boundary types
const (
	PageBoundaryMedia = iota
	PageBoundaryCrop
	PageBoundaryBleed
	PageBoundaryTrim
	PageBoundaryArt
)

//PageOption option of page
type PageOption struct {
	PageBoundaries []*PageBoundary
}

func (po *PageOption) AddPageBoundary(pb *PageBoundary) {
	for x := 0; x < len(po.PageBoundaries); x++ {
		if po.PageBoundaries[x].Type == pb.Type {
			po.PageBoundaries[x] = pb
			return
		}
	}

	po.PageBoundaries = append(po.PageBoundaries, pb)
}

func (po *PageOption) writePageBoundaries(w io.Writer) error {
	for x := 0; x < len(po.PageBoundaries); x++ {
		if err := po.PageBoundaries[x].write(w); err != nil {
			return err
		}
	}

	return nil
}

func (po *PageOption) GetBoundary(t int) (pb *PageBoundary) {
	for x := 0; x < len(po.PageBoundaries); x++ {
		if po.PageBoundaries[x].Type == t {
			pb = po.PageBoundaries[x]
			break
		}
	}
	return
}

func (po PageOption) merge(po2 PageOption) PageOption {
	var pageOpt PageOption
	copy(pageOpt.PageBoundaries, po2.PageBoundaries)
	for x := 0; x < len(po.PageBoundaries); x++ {
		pageOpt.AddPageBoundary(po.PageBoundaries[x])
	}
	return pageOpt
}

type PageBoundary struct {
	Type     int
	Position Point
	Size     Rect
}

func (pb *PageBoundary) write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "/%s [%.2f %.2f %.2f %.2f]\n", pb.TypeString(), pb.Position.X, pb.Position.Y, pb.Size.W, pb.Size.H)
	return err
}

func (pb *PageBoundary) TypeString() string {
	switch pb.Type {
	case PageBoundaryMedia:
		return "MediaBox"
	case PageBoundaryCrop:
		return "CropBox"
	case PageBoundaryBleed:
		return "BleedBox"
	case PageBoundaryTrim:
		return "TrimBox"
	case PageBoundaryArt:
		return "ArtBox"
	}

	return ""
}

func NewPageBoundary(u int, t int, x, y, w, h float64) *PageBoundary {
	UnitsToPointsVar(u, &x, &y, &w, &h)
	return &PageBoundary{
		Type:     t,
		Position: Point{X: x, Y: y},
		Size:     Rect{W: w, H: h},
	}
}

func (gp *Fpdf) NewPageBoundary(t int, x, y, w, h float64) *PageBoundary {
	return NewPageBoundary(gp.config.Unit, t, x, y, w, h)
}

func NewPageSizeBoundary(u int, w, h float64) *PageBoundary {
	return NewPageBoundary(u, PageBoundaryMedia, 0, 0, w, h)
}

func (gp *Fpdf) NewPageSizeBoundary(w, h float64) *PageBoundary {
	return gp.NewPageBoundary(PageBoundaryMedia, 0, 0, w, h)
}

func NewCropPageBoundary(u int, x, y, w, h float64) *PageBoundary {
	return NewPageBoundary(u, PageBoundaryCrop, x, y, w, h)
}

func (gp *Fpdf) NewCropPageBoundary(x, y, w, h float64) *PageBoundary {
	return gp.NewPageBoundary(PageBoundaryCrop, x, y, w, h)
}

func NewBleedPageBoundary(u int, x, y, w, h float64) *PageBoundary {
	return NewPageBoundary(u, PageBoundaryBleed, x, y, w, h)
}

func (gp *Fpdf) NewBleedPageBoundary(x, y, w, h float64) *PageBoundary {
	return gp.NewPageBoundary(PageBoundaryBleed, x, y, w, h)
}

func NewTrimPageBoundary(u int, x, y, w, h float64) *PageBoundary {
	return NewPageBoundary(u, PageBoundaryTrim, x, y, w, h)
}

func (gp *Fpdf) NewTrimPageBoundary(x, y, w, h float64) *PageBoundary {
	return gp.NewPageBoundary(PageBoundaryTrim, x, y, w, h)
}

func NewArtPageBoundary(u int, x, y, w, h float64) *PageBoundary {
	return NewPageBoundary(u, PageBoundaryArt, x, y, w, h)
}

func (gp *Fpdf) NewArtPageBoundary(x, y, w, h float64) *PageBoundary {
	return gp.NewPageBoundary(PageBoundaryArt, x, y, w, h)
}
