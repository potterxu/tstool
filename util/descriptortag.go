package util

import "fmt"

type DescriptorTag int

const (
	ZERO DescriptorTag = iota
)

var (
	descriptorTagString = []string{}
)

func (d DescriptorTag) String() string {
	if int(d) >= len(descriptorTagString) {
		return fmt.Sprintf("0x%02x", int(d))
	}
	return descriptorTagString[d]
}
