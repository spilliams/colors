package color

import "testing"

func TestContrastRatio(t *testing.T) {
	cases := []struct {
		name string
		one  string
		two  string
		cr   float64
	}{
		{
			name: "yellows",
			one:  "fff6dd",
			two:  "a95f09",
			cr:   4.5007294383,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c1, err := NewFromHex("one", c.one)
			if err != nil {
				t.Fatal(err)
			}
			c2, err := NewFromHex("two", c.two)
			if err != nil {
				t.Fatal(err)
			}
			cr := c1.ContrastRatio(c2)
			// truncate because otherwise we'll be here all day
			cr = float64(int(cr*1e10)) / 1e10
			if cr != c.cr {
				t.Errorf("Expected cr %0.10f, got cr %0.10f", c.cr, cr)
			}
		})
	}
}
