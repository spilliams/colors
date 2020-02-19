package contrastratio

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/colors/pkg/color"
)

const (
	unreportedContrastRatioThreshold = 3
	inFileColumnSeparator            = "\t"
	inFileRowSeparator               = "\n"
	outFileColumnSeparator           = ","
	outFileRowSeparator              = "\n"
)

func NewCmd() *cobra.Command {
	var flags struct {
		inFile  string
		outFile string
	}
	exampleInput := fmt.Sprintf("red%sff0000%sgreen%s00ff00%sblue%s0000ff", inFileColumnSeparator, inFileRowSeparator, inFileColumnSeparator, inFileRowSeparator, inFileColumnSeparator)

	cmd := &cobra.Command{
		Use:     "contrast-ratio --in INFILE",
		Aliases: []string{"cr"},
		Short:   "Compute the contrast ratios between given colors",
		Long: fmt.Sprintf(`Compute the contrast ratios between given colors.

This command takes an input file, and produces output either on stdout or to an
output file.

The input file must be like

%s

with %q between each line, and %q between the name and hex value. Note that the
hex value must not have a "#" character.

The command will provide contrast ratios between all input colors, as well as
white (ffffff) and black (000000). It uses the formulas found here to compute:
https://medium.muz.li/the-science-of-color-contrast-an-expert-designers-guide-33e84c41d156
Namely:

	[1] Contrast Ratio is (L_1 + 0.05) / (L_2 + 0.05), where L_1 is the relative
	    luminance of the lighter of the two colors
	[2] Relative Luminance is 0.2126 * R + 0.7152 * G + 0.0722 * B, where R, G
	    and B are the sRGB values of the color
	[3] sRGB of a color value is V / 12.92 if V <= 0.03928. Otherwise, sRGB is
	    ((V + 0.055) / 1.055) ^ 2.4

Note that this means your input colors should be in native RGB values, not sRGB.

The output is CSV-formatted. By default it's printed to stdout, but if you
provide an --out flag it will attempt to write the output to file.`, exampleInput, inFileRowSeparator, inFileColumnSeparator),
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

			if err := contrastSet(w, colors); err != nil {
				return err
			}

			w.Flush()
			if len(flags.outFile) > 0 {
				log.Infof("Output is in file %s", flags.outFile)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&flags.inFile, "in", "", "The name of the file with all the colors in it (required)")
	_ = cmd.MarkFlagRequired("in")
	cmd.Flags().StringVar(&flags.outFile, "out", "", "The name of the file to use for output. If blank, this command will use stdout")

	return cmd
}

func readInFile(in string) ([]*color.Color, error) {
	inBytes, err := ioutil.ReadFile(in)
	if err != nil {
		return nil, err
	}
	inLines := strings.Split(string(inBytes), inFileRowSeparator)
	colors := make([]*color.Color, 0)
	for i, line := range inLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, inFileColumnSeparator)
		if len(parts) != 2 {
			return nil, fmt.Errorf("syntax error on line %d: line \"%s\" must be in the format %q", i, line, fmt.Sprintf("name%shex", inFileColumnSeparator))
		}
		c, err := color.NewFromHex(parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		colors = append(colors, c)
	}
	return colors, nil
}

func contrastSet(w io.Writer, colors []*color.Color) error {
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
	rows := make([][]string, len(colors))
	for i, c := range colors {
		name := c.String()
		headers = append(headers, name)

		data := []string{name}
		for _, vs := range colors {
			datum := contrast(c, vs)
			data = append(data, datum)
		}
		rows[i] = data
	}

	if err := writeLine(w, headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writeLine(w, row); err != nil {
			return err
		}
	}

	return nil
}

func writeLine(w io.Writer, parts []string) error {
	line := fmt.Sprintf("%s%s", strings.Join(parts, outFileColumnSeparator), outFileRowSeparator)
	_, err := w.Write([]byte(line))
	return err
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
