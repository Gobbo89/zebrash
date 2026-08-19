package parsers

import (
	"github.com/Gobbo89/zebrash/internal/printers"
)

func NewPrintOrientationParser() *CommandParser {
	const code = "^PO"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			text := commandText(command, code)

			printer.LabelInverted = (text == "I")

			return nil, nil
		},
	}
}
