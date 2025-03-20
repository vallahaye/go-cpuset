package cpuset

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
	"strings"
)

const (
	wordBase    = 16
	wordBitSize = 32
)

// ParseMask decodes s into a [CPUSet]. It returns an error if s is not a valid
// mask string, as specified in the [Linux cpuset(7) man page] (see "Mask
// Format").
//
// [Linux cpuset(7) man page]: https://man7.org/linux/man-pages/man7/cpuset.7.html
func ParseMask(s string) (cset CPUSet, _ error) {
	if s == "" {
		return CPUSet{}, nil
	}

	words := strings.Split(s, ",")
	overflow, offset := bits.Mul(uint(len(words)), wordBitSize)
	if overflow > 0 {
		return CPUSet{}, formatParseError(s, "offset value out of range")
	}

	for _, word := range words {
		offset -= wordBitSize

		ui64, err := strconv.ParseUint(word, wordBase, wordBitSize)
		if err != nil {
			return CPUSet{}, formatParseError(s, fmt.Sprintf("invalid 32-bit word %q", word))
		}

		cpuMask := uint32(ui64)

		for pos := range uint(wordBitSize) {
			if cpuMask&(1<<pos) != 0 {
				cset.Add(pos + offset)
			}
		}
	}

	return cset, nil
}

// MaskString encodes s into a mask string.
func (s *CPUSet) MaskString() string {
	if s.Len() == 0 {
		return ""
	}

	cpus := s.UnsortedList()
	cpuMasks := make([]uint32, slices.Max(cpus)/wordBitSize+1)
	for _, cpu := range cpus {
		i, pos := int(cpu/wordBitSize), cpu%wordBitSize
		cpuMasks[len(cpuMasks)-i-1] |= 1 << pos
	}

	words := make([]string, len(cpuMasks))
	for i, cpuMask := range cpuMasks {
		words[i] = fmt.Sprintf("%08x", cpuMask)
	}

	return strings.Join(words, ",")
}
