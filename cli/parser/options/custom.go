package options

import (
	"fmt"
	"strings"
)

const CategoryCustom = "custom"

// ParseCustom parses a custom option.
func ParseCustom(option string) (*Option, error) {
	splitoption := strings.Fields(option)

	if len(splitoption) == 0 {
		return nil, fmt.Errorf("there is an unspecified %s option at an unknown line", CategoryCustom)
	} else if len(splitoption) == 1 {
		return nil, fmt.Errorf("there is a misconfigured %s option: %q.\nIs it in format %s?", CategoryCustom, option, "<option><whitespaces><value>")
	}

	return &Option{
		Category: splitoption[0],
		Value:    splitoption[:1],
	}, nil
}

// MapCustomOption maps a custom option in an optionmap[category][]values.
func MapCustomOption(optionmap map[string][]string, option *Option) error {
	if optionmap == nil {
		optionmap = make(map[string][]string)
	}

	switch option.Category {
	case CategoryConvert, CategoryDeepcopy, CategoryDepth, CategoryMap:
	default:
		if v, ok := option.Value.(string); ok {
			optionmap[option.Category] = append(optionmap[option.Category], v)
		} else {
			return fmt.Errorf("failed to map custom option: %v", option.Value)
		}
	}

	return nil
}

// MapCustomOptions maps options with custom categories in a list of options to a customoptionmap[category][]value.
func MapCustomOptions(options []*Option) (map[string][]string, error) {
	var optionmap map[string][]string
	for _, option := range options {
		err := MapCustomOption(optionmap, option)
		if err != nil {
			return nil, err
		}
	}

	return optionmap, nil
}
