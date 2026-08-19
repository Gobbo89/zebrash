package parsers

import (
	"github.com/Gobbo89/zebrash/internal/printers"
)

func NewFieldReversePrintParser() *CommandParser {
	const code = "^FR"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			printer.NextElementFieldReverse = true

			return nil, nil
		},
	}
}
