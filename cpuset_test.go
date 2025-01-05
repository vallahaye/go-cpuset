package cpuset

import (
	"errors"
	"maps"
	"slices"
	"testing"
)

func TestOf(t *testing.T) {
	for _, params := range []struct {
		name string
		cpus []uint
	}{
		{
			name: "empty",
			cpus: nil,
		},
		{
			name: "not empty",
			cpus: []uint{0, 1, 2, 3},
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			m := make(map[uint]struct{}, len(params.cpus))
			for _, cpu := range params.cpus {
				m[cpu] = struct{}{}
			}

			s := Of(params.cpus...)
			if !maps.Equal(m, s.m) {
				t.Errorf("unexpected map: got %v, want %v", m, s.m)
			}
		})
	}
}

func TestCPUSetAdd(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		cpu  uint
		want bool
	}{
		{
			name: "empty not present",
			s:    CPUSet{},
			cpu:  0,
			want: true,
		},
		{
			name: "not empty present",
			s:    Of(0, 1, 2, 3),
			cpu:  0,
			want: false,
		},
		{
			name: "not empty not present",
			s:    Of(0, 1, 2, 3),
			cpu:  4,
			want: true,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			switch got := params.s.Add(params.cpu); {
			case !params.s.Contains(params.cpu):
				t.Error("not added")
			case got != params.want:
				t.Errorf("unexpected presence report: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetDelete(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		cpu  uint
		want bool
	}{
		{
			name: "empty not present",
			s:    CPUSet{},
			cpu:  0,
			want: false,
		},
		{
			name: "not empty present",
			s:    Of(0, 1, 2, 3),
			cpu:  0,
			want: true,
		},
		{
			name: "not empty not present",
			s:    Of(0, 1, 2, 3),
			cpu:  4,
			want: false,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			switch got := params.s.Delete(params.cpu); {
			case params.s.Contains(params.cpu):
				t.Error("not deleted")
			case got != params.want:
				t.Errorf("unexpected presence report: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetContains(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		cpu  uint
		want bool
	}{
		{
			name: "empty not present",
			s:    CPUSet{},
			cpu:  0,
			want: false,
		},
		{
			name: "not empty present",
			s:    Of(0, 1, 2, 3),
			cpu:  0,
			want: true,
		},
		{
			name: "not empty not present",
			s:    Of(0, 1, 2, 3),
			cpu:  4,
			want: false,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.Contains(params.cpu); got != params.want {
				t.Errorf("unexpected presence report: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetUnsortedList(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		want []uint
	}{
		{
			name: "empty",
			s:    CPUSet{},
			want: []uint{},
		},
		{
			name: "not empty",
			s:    Of(0, 1, 2, 3),
			want: []uint{0, 1, 2, 3},
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			got := params.s.UnsortedList()
			slices.Sort(got)
			if !slices.Equal(got, params.want) {
				t.Errorf("unexpected list: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetEqual(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		s2   CPUSet
		want bool
	}{
		{
			name: "empty equal empty",
			s:    CPUSet{},
			s2:   CPUSet{},
			want: true,
		},
		{
			name: "empty not equal not empty",
			s:    CPUSet{},
			s2:   Of(0, 1, 2, 3),
			want: false,
		},
		{
			name: "not empty not equal empty",
			s:    Of(0, 1, 2, 3),
			s2:   CPUSet{},
			want: false,
		},
		{
			name: "not empty not equal not empty",
			s:    Of(0, 1, 2, 3),
			s2:   Of(0),
			want: false,
		},
		{
			name: "not empty equal not empty",
			s:    Of(0, 1, 2, 3),
			s2:   Of(0, 1, 2, 3),
			want: true,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.Equal(params.s2); got != params.want {
				t.Errorf("unexpected equal report: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetClear(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
	}{
		{
			name: "empty",
			s:    CPUSet{},
		},
		{
			name: "not empty",
			s:    Of(0, 1, 2, 3),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			params.s.Clear()
			if params.s.Len() != 0 {
				t.Error("not cleared")
			}
		})
	}
}

func TestCPUSetClone(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		want CPUSet
	}{
		{
			name: "empty",
			s:    CPUSet{},
			want: CPUSet{},
		},
		{
			name: "not empty",
			s:    Of(0, 1, 2, 3),
			want: Of(0, 1, 2, 3),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.Clone(); !got.Equal(params.want) {
				t.Errorf("unexpected clone: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestCPUSetLen(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		want int
	}{
		{
			name: "empty",
			s:    CPUSet{},
			want: 0,
		},
		{
			name: "not empty",
			s:    Of(0, 1, 2, 3),
			want: 4,
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.Len(); got != params.want {
				t.Errorf("unexpected len: got %d, want %d", got, params.want)
			}
		})
	}
}

func TestCPUSetString(t *testing.T) {
	for _, params := range []struct {
		name string
		s    CPUSet
		want string
	}{
		{
			name: "empty",
			s:    CPUSet{},
			want: "",
		},
		{
			name: "not empty",
			s:    Of(0, 1, 2, 3),
			want: "0-3",
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := params.s.String(); got != params.want {
				t.Errorf("unexpected string: got %q, want %q", got, params.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	for _, params := range []struct {
		name string
		s1   CPUSet
		s2   CPUSet
		want CPUSet
	}{
		{
			name: "empty difference empty",
			s1:   CPUSet{},
			s2:   CPUSet{},
			want: CPUSet{},
		},
		{
			name: "empty difference not empty",
			s1:   CPUSet{},
			s2:   Of(0, 1, 2, 3),
			want: CPUSet{},
		},
		{
			name: "not empty difference empty",
			s1:   Of(0, 1, 2, 3),
			s2:   CPUSet{},
			want: Of(0, 1, 2, 3),
		},
		{
			name: "not empty difference not empty",
			s1:   Of(0, 1, 2, 3),
			s2:   Of(0),
			want: Of(1, 2, 3),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := Difference(params.s1, params.s2); !got.Equal(params.want) {
				t.Errorf("unexpected difference: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	for _, params := range []struct {
		name string
		s1   CPUSet
		s2   CPUSet
		want CPUSet
	}{
		{
			name: "empty intersection empty",
			s1:   CPUSet{},
			s2:   CPUSet{},
			want: CPUSet{},
		},
		{
			name: "empty intersection not empty",
			s1:   CPUSet{},
			s2:   Of(0, 1, 2, 3),
			want: CPUSet{},
		},
		{
			name: "not empty intersection empty",
			s1:   Of(0, 1, 2, 3),
			s2:   CPUSet{},
			want: CPUSet{},
		},
		{
			name: "not empty intersection not empty",
			s1:   Of(0, 1, 2, 3),
			s2:   Of(0),
			want: Of(0),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := Intersection(params.s1, params.s2); !got.Equal(params.want) {
				t.Errorf("unexpected intersection: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	for _, params := range []struct {
		name string
		s1   CPUSet
		s2   CPUSet
		want CPUSet
	}{
		{
			name: "empty union empty",
			s1:   CPUSet{},
			s2:   CPUSet{},
			want: CPUSet{},
		},
		{
			name: "empty union not empty",
			s1:   CPUSet{},
			s2:   Of(0, 1, 2, 3),
			want: Of(0, 1, 2, 3),
		},
		{
			name: "not empty union empty",
			s1:   Of(0, 1, 2, 3),
			s2:   CPUSet{},
			want: Of(0, 1, 2, 3),
		},
		{
			name: "not empty union not empty",
			s1:   Of(0, 1, 2, 3),
			s2:   Of(0),
			want: Of(0, 1, 2, 3),
		},
	} {
		t.Run(params.name, func(t *testing.T) {
			if got := Union(params.s1, params.s2); !got.Equal(params.want) {
				t.Errorf("unexpected union: got %v, want %v", got, params.want)
			}
		})
	}
}

func TestFormatParseError(t *testing.T) {
	want := errors.New(`cpuset: parsing "s": error`)
	if got := formatParseError("s", "error"); got.Error() != want.Error() {
		t.Errorf("unexpected error: got %v, want %v", got, want)
	}
}
