package distance

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spilliams/colors/pkg/color"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "distance A B",
		Aliases: []string{"d"},
		Short:   "Compute the color distance between given colors",
		Long: `Compute the color distance between given colors.

This interpretation is very naive, and assumes color is represented in a
3-dimensional space with axes red, green and blue. This command will compute
the distance between two points in this space using the formula
 sqrt(R^2 + G^2 + B^2)
 where R is the difference between the two colors' red values, etc.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := color.NewFromHex("a", args[0])
			if err != nil {
				return err
			}
			b, err := color.NewFromHex("b", args[1])
			if err != nil {
				return err
			}

			fmt.Printf("A is %v\n", a.Hex())
			fmt.Printf("B is %v\n", b.Hex())
			fmt.Printf("Distance between A and B: %0.02f\n", a.DistanceTo(b))
			cr := a.ContrastRatio(b)
			fmt.Printf("Contrast ratio between A and B: %0.02f (%s)\n", cr, color.ContrastRatioDescription(cr))
			return nil
		},
	}
}
