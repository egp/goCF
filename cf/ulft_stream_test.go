// ulft_stream_test.go v2
package cf

import "testing"

func collectAll(cf ContinuedFraction) []int64 {
	var out []int64
	for {
		v, ok := cf.Next()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

func TestULFTStream_Golden_RationalInputs(t *testing.T) {
	type tc struct {
		name string
		t    ULFT
		x    Rational
	}

	cases := []tc{
		{
			name: "identity_355_113",
			t:    NewULFT(1, 0, 0, 1),
			x:    mustRat(355, 113),
		},
		{
			name: "x_plus_1_on_3_2",
			t:    NewULFT(1, 1, 0, 1), // x + 1
			x:    mustRat(3, 2),
		},
		{
			name: "reciprocal_on_7_5",
			t:    NewULFT(0, 1, 1, 0), // 1/x
			x:    mustRat(7, 5),
		},
		{
			name: "general_on_1_2",
			t:    NewULFT(2, 1, 3, 4), // (2x+1)/(3x+4)
			x:    mustRat(1, 2),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Reference exact rational y = T(x), then CF-expand y completely.
			y, err := c.t.ApplyRat(c.x)
			if err != nil {
				t.Fatal(err)
			}
			want := collectAll(NewRationalCF(y))

			// Streamed output from ULFTStream(T, CF(x))
			src := NewRationalCF(c.x)
			s := NewULFTStream(c.t, src, ULFTStreamOptions{
				DetectCycles: true,
				MaxRepeats:   2,
			})

			got := collectAll(s)
			if s.Err() != nil {
				t.Fatalf("stream error: %v", s.Err())
			}

			if !equalSlice(got, want) {
				t.Fatalf("got %v, want %v (y=%v)", got, want, y)
			}
		})
	}
}

func TestULFTStream_Identity_ReproducesInputCF(t *testing.T) {
	x := mustRat(355, 113) // [3;7,16]
	src := NewRationalCF(x)
	want := collectAll(NewRationalCF(x))

	id := NewULFT(1, 0, 0, 1)
	s := NewULFTStream(id, src, ULFTStreamOptions{DetectCycles: true})

	got := collectAll(s)
	if s.Err() != nil {
		t.Fatalf("stream error: %v", s.Err())
	}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// ulft_stream_test.go v2
