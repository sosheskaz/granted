package cfaws

import (
	"fmt"
	"strconv"
	"strings"
)

const DefaultRegion = "us-east-1"

// ExpandRegion takes a string and attemps to expand it into a fully formed region e.g ue1 -> us-east-1
//
// If region is an empty string, the DefaultRegion is returned
//
// ExpandRegion does not attempt to fully validate regions and may produce regions which do not exist, for example as2 -> ap-south-2 which is not a valid region
func ExpandRegion(region string) (string, error) {
	// Region could come in one of three formats:
	// 1. No region specified
	if region == "" {
		return DefaultRegion, nil
	}
	// 2. A fully-qualified region. Assume that if there's one dash, it's valid.
	if strings.Contains(region, "-") {
		return region, nil
	}
	var major, minor, num string
	idx := 1 // Number of characters consumed from region
	// 3. Otherwise, we have a shortened region, like ue1
	if len(region) < 2 {
		return "", fmt.Errorf("region too short, needs at least two characters (eg ue)")
	}
	// Region might be one or two letters
	switch region[0] {
	case 'u':
		major = "us"
		switch region[1] {
		case 'g':
			major = "us-gov"
			idx += 1
		case 's':
			// This will break if us-southeast-1 is ever created
			idx += 1
		}
	case 'e':
		major = "eu"
		if region[1] == 'u' {
			idx += 1
		}
	case 'a':
		major = "ap"
		switch region[1] {
		case 'f':
			major = "af"
			idx += 1
		case 'p':
			idx += 1
		}
	case 'c':
		major = "ca"
		switch region[1] {
		case 'n':
			major = "cn"
			idx += 1
		case 'a':
			idx += 1
		}
	case 'm':
		major = "me"
		// This will break if me-east-1 is ever created
		if region[1] == 'e' {
			idx += 1
		}
	case 's':
		major = "sa"
		if region[1] == 'a' {
			idx += 1
		}
	default:
		return "", fmt.Errorf("unknown region major (hint: try using the first letter of the region)")
	}
	region = region[idx:]
	idx = 1
	// Location might be one or two letters (n, nw)
	switch region[0] {
	case 'n', 's':
		if region[0] == 'n' {
			minor = "north"
		} else {
			minor = "south"
		}
		if len(region) > 1 {
			switch region[1] {
			case 'w':
				minor += "west"
				idx += 1
			case 'e':
				minor += "east"
				idx += 1
			}
		}
	case 'e':
		minor = "east"
	case 'w':
		minor = "west"
	case 'c':
		minor = "central"
	default:
		return "", fmt.Errorf("unknown region minor in %s (found major: %s)", region, major)

	}
	region = region[idx:]
	if len(region) > 0 {
		_, err := strconv.Atoi(region)
		if err != nil {
			return "", fmt.Errorf("unknown region number in %s (found major: %s, minor: %s)", region, major, minor)
		}
		num = region
	} else {
		num = "1"
	}

	return fmt.Sprintf("%s-%s-%s", major, minor, num), nil
}
