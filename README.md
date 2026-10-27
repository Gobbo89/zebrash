# Zebrash

> [!IMPORTANT]
> [Fedir Kryvytskyi](https://github.com/fkryvyts-ingrid) left [Ingrid](https://github.com/ingridhq), and the [original repo](https://github.com/ingridhq/zebrash) is not actively maintained.  
> This fork was created to add EAN 8 support.

Library for rendering ZPL (Zebra Programming Language) files as raster images

- Partially based on https://github.com/BinaryKits/BinaryKits.Zpl
- Uses slightly modified implementations of:
	- PDF417, Aztec, Code 2 of 5, EAN13 and EAN8 barcodes from https://github.com/boombuler/barcode/
	- Code 128, Code 39, QR Code and DataMatrix from https://github.com/makiuchi-d/gozxing

## Description

This library emulates subset of ZPL engine and allows you to view most of the ZPL labels that are used by carriers such as Fedex, UPS or DHL as PNGs without the need to possess physical Zebra-compatible printer.
Think of https://labelary.com/viewer.html except it is completely free for commercial use, has no API limits and can easily be self-hosted or plugged into existing Go application so you don't need to send labels with real customers information to some 3rd-party servers

Example of the output (more examples can be found inside `testdata` folder):

![UPS label](testdata/ups_grayscale.png)


## Rendering compatibility

PDF417 (`^B7`) uses the module width supplied by `^BY`. A fresh parser defaults
to width 2 when `^BY` is omitted. Explicit zero or negative widths use a minimum
of 1, consistent with the Code 128 and EAN encoders. Positive widths are preserved,
including values above 10. This deliberately renders a barcode for `^BY0`, unlike
the tested Zebra ZD230t (firmware V89.21.46Z), which produces no visible barcode;
Labelary also substitutes width 1 for that input.

Rotated `^FO`/`^FB` text snapshots cover full, underfilled, and single-line blocks
at line spacings -10, 0, and +6. Their placement was checked against that printer.
The fixtures use the renderer's fonts, so they do not require pixel-identical
printer glyphs. The declared maximum line count determines the initial line
position even when the block contains fewer lines.

## Usage:

```go

	file, err := os.ReadFile("./testdata/label.zpl")
	if err != nil {
		t.Fatal(err)
	}

	parser := zebrash.NewParser()

	res, err := parser.Parse(file)
	if err != nil {
		t.Fatal(err)
	}

	var buff bytes.Buffer

	drawer := zebrash.NewDrawer()

	err = drawer.DrawLabelAsPng(res[0], &buff, drawers.DrawerOptions{
		LabelWidthMm:         101.6,
		LabelHeightMm:        203.2,
		Dpmm:                 8,
		EnableInvertedLabels: true,
		GrayscaleOutput:      true
	})
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("./testdata/label.png", buff.Bytes(), 0744)
	if err != nil {
		t.Fatal(err)
	}

```
