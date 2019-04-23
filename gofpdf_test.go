package gofpdf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFont(t *testing.T) {
	ttfr, err := os.Open("test/res/times.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	font, err := SubsetFontByReader(ttfr)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := font.Serialize()
	if err != nil {
		t.Error(err)
		return
	}

	font2, err := DeserializeSubsetFont(b)
	if err != nil {
		t.Error(err)
		return
	}

	b, err = font.SerializeJSON()
	if err != nil {
		t.Error(err)
		return
	}

	font2, err = DeserializeJSONSubsetFont(b)
	if err != nil {
		t.Error(err)
		return
	}

	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)
	if err != nil {
		t.Error(err)
		return
	}

	if err := pdf.AddTTFFontBySubsetFont("test", font2); err != nil {
		t.Error(err)
		return
	}
}

func TestImageObj(t *testing.T) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)
	pdf.AddPage()

	holder, err := ImageHolderByPath("test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
	}

	img, err := NewImageObj(holder)
	if err != nil {
		t.Error(err)
	}

	b, err := img.Serialize()
	if err != nil {
		t.Error(err)
	}

	img2, err := DeserializeImage(b)
	if err != nil {
		t.Error(err)
	}

	b, err = img.SerializeJSON()
	if err != nil {
		t.Error(err)
	}

	img2, err = DeserializeJSONImage(b)
	if err != nil {
		t.Error(err)
	}

	err = pdf.ImageByObj(img2, 0, 0, Rect{W: 0, H: 0})
	if err != nil {
		t.Error(err)
	}
}

func TestTemplateAutoPage(t *testing.T) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)

	if err != nil {
		t.Error(err)
	}

	pdf.AddPage()
	pdf.AddTTFFont("a", "test/res/times.ttf")
	pdf.SetFont("a", "", 12)
	var b string

	for x := 0; x < 10000; x++ {
		b = fmt.Sprintf("something %s", b)
	}

	pdf.WriteText(12, b)

	tmpl, err := pdf.Template(Point{})
	if err != nil {
		t.Error(err)
	}

	bb, err := tmpl.Serialize()
	if err != nil {
		t.Error(err)
	}

	tmpl2, err := DeserializeTemplate(bb)
	if err != nil {
		t.Error(err)
	}

	bb, err = tmpl.SerializeJSON()
	if err != nil {
		t.Error(err)
	}

	tmpl2, err = DeserializeJSONTemplate(bb)
	if err != nil {
		t.Error(err)
	}

	if tmpl2.NumPages() != 668 {
		t.Error(fmt.Errorf("number of pages %d", tmpl2.NumPages()))
	}
}

func TestTemplatePages(t *testing.T) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)

	if err != nil {
		t.Error(err)
	}

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	tmpl, err := pdf.Template(Point{})
	if err != nil {
		t.Error(err)
	}

	if tmpl.NumPages() != 2 {
		t.Error(errors.New("there should be more pages"))
	}
}

func TestAutoWidth(t *testing.T) {
	pdf, err := New(PdfOptionPageSize(250, 250))
	if err != nil {
		t.Error(err)
	}

	pdf.AddPage()

	if err := pdf.AddTTFFont("times", "test/res/times.ttf"); err != nil {
		t.Error(err)
	}

	if err := pdf.SetFont("times", "", 12); err != nil {
		t.Error(err)
	}

	err = pdf.MultiCellOpts(0, 10, "something here", CellOption{
		Align:  Top | Left,
		Border: 0,
		Float:  Left,
	})

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkTemplateSerializtion(b *testing.B) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)

	pdf.AddPage()
	pdf.AddTTFFont("a", "test/res/times.ttf")
	pdf.SetFont("a", "", 12)

	if err != nil {
		b.Error(err)
	}

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	for x := 0; x < 5; x++ {
		holder, err := newImageBuffByPath("test/res/chilli.jpg")
		if err != nil {
			b.Error(err)
		}

		err = pdf.ImageByHolder(holder, 0, 0, Rect{W: 20, H: 20})
		if err != nil {
			b.Error(err)
		}
	}

	tmpl, err := pdf.Template(Point{})
	if err != nil {
		b.Error(err)
	}

	bb, err := tmpl.Serialize()
	if err != nil {
		b.Error(err)
	}

	b.Run("Serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := tmpl.Serialize()
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("Deserialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := DeserializeTemplate(bb)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkFontSerialization(b *testing.B) {
	ttfr, err := os.Open("test/res/times.ttf")
	if err != nil {
		b.Error(err)
	}

	font, err := SubsetFontByReader(ttfr)
	if err != nil {
		b.Error(err)
	}

	bb, err := font.Serialize()
	if err != nil {
		b.Error(err)
	}

	b.Run("Serialization", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := font.Serialize()
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("Deserialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := DeserializeSubsetFont(bb)
			if err != nil {
				b.Error(err)
			}
		}
	})

}

func BenchmarkPdfWithImageObj(b *testing.B) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)
	if err != nil {
		b.Error(err)
	}
	pdf.AddPage()

	holder, err := ImageHolderByPath("test/res/gopher01.jpg")
	if err != nil {
		b.Error(err)
	}

	img, err := NewImageObj(holder)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		err = pdf.ImageByObj(img, 0, 0, Rect{W: 0, H: 0})
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkPdfWithImageHolder(b *testing.B) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)
	if err != nil {
		b.Error(err)
	}
	pdf.AddPage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		holder, err := newImageBuffByPath("test/res/chilli.jpg")
		if err != nil {
			b.Error(err)
		}

		err = pdf.ImageByHolder(holder, 0, 0, Rect{W: 20, H: 20})
		if err != nil {
			b.Error(err)
		}
	}
}

func initTesting() error {
	err := os.MkdirAll("./test/out", 0777)
	if err != nil {
		return err
	}
	return nil
}

func TestPdfWithImageHolder(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf, err := New(PdfOptionPageSize(595.28, 841.89)) //595.28, 841.89 = A4
	if err != nil {
		t.Error(err)
	}
	pdf.AddPage()
	err = pdf.AddTTFFont("loma", "./test/res/times.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("loma", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	bytesOfImg, err := ioutil.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.ImageByHolder(imgH, 20.0, 20, Rect{W: 20, H: 20})
	if err != nil {
		t.Error(err)
		return
	}

	// because this uses a reader it's pointer is in the wrong place when used
	// for a second time. we might need to add some extra stuff to reset the
	// pointer after an image holder is consumed
	imgH, err = ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.ImageByHolder(imgH, 20.0, 200, Rect{W: 20, H: 20})
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(20, 20, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test.pdf")
}

func BenchmarkSerializationFont(b *testing.B) {
	ttfr, err := os.Open("test/res/times.ttf")
	if err != nil {
		b.Error(err)
	}

	font, err := SubsetFontByReader(ttfr)
	if err != nil {
		b.Error(err)
	}

	b.Run("gob", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := font.Serialize()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeSubsetFont(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("json", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := font.SerializeJSON()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeJSONSubsetFont(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkSerializationImage(b *testing.B) {
	holder, err := ImageHolderByPath("test/res/gopher01.jpg")
	if err != nil {
		b.Error(err)
	}

	img, err := NewImageObj(holder)
	if err != nil {
		b.Error(err)
	}

	b.Run("gob", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := img.Serialize()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeImage(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("json", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := img.SerializeJSON()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeJSONImage(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkSerializationTemplate(b *testing.B) {
	pdf, err := New(
		PdfOptionUnit(Unit_IN),
		PdfOptionPageSize(12, 12),
	)

	pdf.AddPage()
	pdf.AddTTFFont("a", "test/res/times.ttf")
	pdf.SetFont("a", "", 12)

	if err != nil {
		b.Error(err)
	}

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	pdf.AddPage()
	pdf.Line(0, 0, 12, 12)

	for x := 0; x < 5; x++ {
		holder, err := newImageBuffByPath("test/res/chilli.jpg")
		if err != nil {
			b.Error(err)
		}

		err = pdf.ImageByHolder(holder, 0, 0, Rect{W: 20, H: 20})
		if err != nil {
			b.Error(err)
		}
	}

	tmpl, err := pdf.Template(Point{})
	if err != nil {
		b.Error(err)
	}

	b.Run("gob", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := tmpl.Serialize()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeTemplate(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("json", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			bs, err := tmpl.SerializeJSON()
			if err != nil {
				b.Error(err)
			}

			_, err = DeserializeJSONTemplate(bs)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
