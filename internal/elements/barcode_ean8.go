package elements

type BarcodeEan8 struct {
	// The bar code orientation to use.
	Orientation FieldOrientation
	// The bar code height, in dots.
	Height int
	// Whether or not to include human-readable text with the bar code.
	Line bool
	// Whether or not to place the human-readable text above the bar code.
	LineAbove bool
}

type BarcodeEan8WithData struct {
	ReversePrint
	BarcodeEan8
	Width    int
	Position LabelPosition
	Data     string
}
