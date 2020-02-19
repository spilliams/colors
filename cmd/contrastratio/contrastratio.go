package contrastratio

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spilliams/colors/pkg/color"
)

const (
	unreportedContrastRatioThreshold = 3
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contrast-ratio",
		Aliases: []string{"cr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// each row has 9 colors

			// to get these hex values from Jasmine:
			// Digital Color Meter: display native values (in hex)
			// open the svg in firefox
			rows := [][]color.Color{
				[]color.Color{
					color.NewFromHex("neutral-10", "FDFEFF"),
					color.NewFromHex("neutral-20", "DFE8EE"),
					color.NewFromHex("neutral-30", "C6D2DB"),
					color.NewFromHex("neutral-40", "a9b5c1"),
					color.NewFromHex("neutral-50", "858F9D"),
					color.NewFromHex("neutral-60", "68707C"),
					color.NewFromHex("neutral-70", "4A505A"),
					color.NewFromHex("neutral-80", "33373E"),
					color.NewFromHex("neutral-90", "23252B"),
				},
				[]color.Color{
					color.NewFromHex("blue-10", "D8EBF2"),
					color.NewFromHex("blue-20", "A7D7EE"),
					color.NewFromHex("blue-30", "7ABEE7"),
					color.NewFromHex("blue-40", "53A1DD"),
					color.NewFromHex("blue-50", "3584D0"),
					color.NewFromHex("blue-60", "2067be"),
					color.NewFromHex("blue-70", "104CA6"),
					color.NewFromHex("blue-80", "073487"),
					color.NewFromHex("blue-90", "02216A"),
				},
				[]color.Color{
					color.NewFromHex("green-10", "e1fae6"),
					color.NewFromHex("green-20", "bcebcd"),
					color.NewFromHex("green-30", "8ed7a9"),
					color.NewFromHex("green-40", "6ec28e"),
					color.NewFromHex("green-50", "409e6f"),
					color.NewFromHex("green-60", "237949"),
					color.NewFromHex("green-70", "0a633c"),
					color.NewFromHex("green-80", "044c2e"),
					color.NewFromHex("green-90", "013d27"),
				},
				[]color.Color{
					color.NewFromHex("yellow-10", "FFF6DD"),
					color.NewFromHex("yellow-20", "FEDA8D"),
					color.NewFromHex("yellow-30", "F9C650"),
					color.NewFromHex("yellow-40", "EDA115"),
					color.NewFromHex("yellow-50", "C87F00"),
					color.NewFromHex("yellow-60", "A66203"),
					color.NewFromHex("yellow-70", "854A02"),
					color.NewFromHex("yellow-80", "6E3808"),
					color.NewFromHex("yellow-90", "602E01"),
				},
				[]color.Color{
					color.NewFromHex("orange-10", "FFE5CF"),
					color.NewFromHex("orange-20", "fec391"),
					color.NewFromHex("orange-30", "FB9853"),
					color.NewFromHex("orange-40", "EB782A"),
					color.NewFromHex("orange-50", "da6018"),
					color.NewFromHex("orange-60", "b14a12"),
					color.NewFromHex("orange-70", "993b09"),
					color.NewFromHex("orange-80", "6c2400"),
					color.NewFromHex("orange-90", "431a00"),
				},
				[]color.Color{
					color.NewFromHex("red-10", "FFDAE5"),
					color.NewFromHex("red-20", "FDBFCC"),
					color.NewFromHex("red-30", "FF99AE"),
					color.NewFromHex("red-40", "FC5D7D"),
					color.NewFromHex("red-50", "DF3655"),
					color.NewFromHex("red-60", "C02A40"),
					color.NewFromHex("red-70", "AC1F34"),
					color.NewFromHex("red-80", "80010E"),
					color.NewFromHex("red-90", "5A0B0D"),
				},
				[]color.Color{
					color.NewFromHex("purple-10", "F1E7FF"),
					color.NewFromHex("purple-20", "DCC3F6"),
					color.NewFromHex("purple-30", "CC95FC"),
					color.NewFromHex("purple-40", "C076F9"),
					color.NewFromHex("purple-50", "AA41EC"),
					color.NewFromHex("purple-60", "921FD1"),
					color.NewFromHex("purple-70", "7615AB"),
					color.NewFromHex("purple-80", "5F0D85"),
					color.NewFromHex("purple-90", "3D054C"),
				},
			}

			// each color in a row must be contrasted with the other 8 colors in its
			// row, plus white and black
			for _, row := range rows {
				table := tablewriter.NewWriter(os.Stderr)
				table.SetRowLine(true)
				table.SetAutoFormatHeaders(false)
				table.SetAlignment(tablewriter.ALIGN_CENTER)
				contrastSet(table, row, true)
			}

			// or perhaps each color is contrasted with EVERY OTHER COLOR
			colors := []color.Color{}
			for _, row := range rows {
				for _, c := range row {
					colors = append(colors, c)
				}
			}
			csvTable := tablewriter.NewWriter(os.Stdout)
			csvTable.SetAutoWrapText(false)
			csvTable.SetAutoFormatHeaders(false)
			csvTable.SetBorder(false)
			csvTable.SetHeaderLine(false)
			csvTable.SetRowSeparator("+")
			csvTable.SetColumnSeparator(",")
			csvTable.SetCenterSeparator("^")
			csvTable.SetTablePadding(" ")
			csvTable.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			csvTable.SetAlignment(tablewriter.ALIGN_LEFT)

			contrastSet(csvTable, colors, false)

			return nil
		},
	}

	return cmd
}

func contrastSet(table *tablewriter.Table, colors []color.Color, multilineRows bool) {
	white := color.NewFromHex("white", "ffffff")
	black := color.NewFromHex("black", "000000")
	c := []color.Color{white}
	c = append(c, colors...)
	c = append(c, black)
	colors = c

	headers := []string{""}
	for _, c := range colors {
		name := c.String()
		if multilineRows {
			name = breakline(name)
		}
		headers = append(headers, name)

		data := []string{name}
		for _, vs := range colors {
			datum := contrast(c, vs)
			if multilineRows {
				datum = breakline(datum)
			}
			data = append(data, datum)
		}
		table.Append(data)
	}
	table.SetHeader(headers)
	table.Render() // Send output
}

func breakline(s string) string {
	return strings.Join(strings.Split(s, " "), "\n")
}

func contrast(fg, bg color.Color) string {
	cr := contrastRatio(fg.Luminance(), bg.Luminance())
	if cr < unreportedContrastRatioThreshold {
		return "--"
	}
	name := "AAA"
	if cr < 4.5 {
		name = "AA+"
	} else if cr < 7 {
		name = "AA"
	}
	return fmt.Sprintf("%0.02f %s", cr, name)
}

// https://medium.muz.li/the-science-of-color-contrast-an-expert-designers-guide-33e84c41d156
func contrastRatio(l1, l2 float64) float64 {
	lighter := l1
	darker := l2
	if l2 > l1 {
		lighter = l2
		darker = l1
	}
	return (lighter + 0.05) / (darker + 0.05)
}
