package cpuset

import (
	"fmt"
	"slices"
	"strings"
)

// ParseList decodes s into a [CPUSet]. It returns an error if s is not a valid
// list string, as specified in the [Linux cpuset(7) man page] (see "List
// Format").
//
// [Linux cpuset(7) man page]: https://man7.org/linux/man-pages/man7/cpuset.7.html
func ParseList(s string) (CPUSet, error) {
	if s == "" {
		return CPUSet{}, nil
	}

	var cpus []uint
	for _, elem := range strings.Split(s, ",") {
		parts := strings.Split(elem, "-")
		switch len(parts) {
		case 1:
			var cpu uint
			if _, err := fmt.Sscan(parts[0], &cpu); err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid element %q", elem))
			}

			cpus = append(cpus, cpu)

		case 2:
			var lowerBound uint
			if _, err := fmt.Sscan(parts[0], &lowerBound); err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid lower bound %q in range %q", parts[0], elem))
			}

			var upperBound uint
			if _, err := fmt.Sscan(parts[1], &upperBound); err != nil {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid upper bound %q in range %q", parts[1], elem))
			}

			if upperBound < lowerBound {
				return CPUSet{}, formatParseError(s, fmt.Sprintf("negative range %q", elem))
			}

			for i := range upperBound - lowerBound + 1 {
				cpus = append(cpus, lowerBound+i)
			}

		default:
			return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid element %q", elem))
		}
	}

	return Of(cpus...), nil
}

// ListString encodes s into a list string.
func (s *CPUSet) ListString() string {
	if len(s.m) == 0 {
		return ""
	}

	cpus := s.UnsortedList()
	slices.Sort(cpus)

	var (
		elems      []string
		lowerBound = cpus[0]
		upperBound = cpus[0]
	)

	for _, cpu := range cpus[1:] {
		if cpu != upperBound+1 {
			elems = append(elems, formatListElem(lowerBound, upperBound))
			lowerBound = cpu
		}

		upperBound = cpu
	}

	elems = append(elems, formatListElem(lowerBound, upperBound))

	return strings.Join(elems, ",")
}

func formatListElem(lowerBound, upperBound uint) string {
	if lowerBound == upperBound {
		return fmt.Sprint(lowerBound)
	}

	return fmt.Sprint(lowerBound, "-", upperBound)
}
