package contrastratio

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/colors/pkg/color"
)

const (
	unreportedContrastRatioThreshold = 3
)

func NewCmd() *cobra.Command {
	var flags struct {
		inFile  string
		outFile string
	}

	cmd := &cobra.Command{
		Use:     "contrast-ratio --in INFILE",
		Aliases: []string{"cr"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// read the infile
			colors, err := readInFile(flags.inFile)
			if err != nil {
				return err
			}

			var out io.Writer
			out = os.Stdout
			if len(flags.outFile) > 0 {
				f, err := os.Create(flags.outFile)
				defer f.Close()
				if err != nil {
					return err
				}
				out = f
			}
			w := bufio.NewWriter(out)

			csvTable := tablewriter.NewWriter(w)
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

			if err := contrastSet(csvTable, colors); err != nil {
				return err
			}

			w.Flush()
			if len(flags.outFile) > 0 {
				log.Infof("Output is in file %s", flags.outFile)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&flags.inFile, "in", "", "The name of the file with all the colors in it")
	_ = cmd.MarkFlagRequired("in")
	cmd.Flags().StringVar(&flags.outFile, "out", "", "The name of the file to use for output. If blank, this command will use stdout")

	return cmd
}

func readInFile(in string) ([]*color.Color, error) {
	inBytes, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, err
	}
	inLines := strings.Split(string(inBytes), "\n")
	colors := make([]*color.Color, 0)
	for i, line := range inLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("syntax error on line %d: line '%s' must be in the format 'name hex'", i, line)
		}
		c, err := color.NewFromHex(parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		colors = append(colors, c)
	}
	return colors, nil
}

func contrastSet(table *tablewriter.Table, colors []*color.Color) error {
	white, err := color.NewFromHex("white", "ffffff")
	if err != nil {
		return err
	}
	black, err := color.NewFromHex("black", "000000")
	if err != nil {
		return err
	}

	c := []*color.Color{white}
	c = append(c, colors...)
	c = append(c, black)
	colors = c

	headers := []string{""}
	for _, c := range colors {
		name := c.String()
		headers = append(headers, name)

		data := []string{name}
		for _, vs := range colors {
			datum := contrast(c, vs)
			data = append(data, datum)
		}
		table.Append(data)
	}
	table.SetHeader(headers)
	table.Render() // Send output

	return nil
}

func contrast(fg, bg *color.Color) string {
	cr := fg.ContrastRatio(bg)
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
