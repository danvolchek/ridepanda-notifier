package ridePandaAPI

import (
	"fmt"
	"strings"
)

// Vehicle represents a RidePanda vehicle.
type Vehicle struct {
	// Whether the vehicle is in stock. Presumably whether at least one variant is in stock.
	InStock bool `graphql:"inStock"`

	// The brand name of the vehicle, arbitrary string.
	Brand string `graphql:"brand"`

	// The model name of the vehicle, arbitrary string.
	Model string `graphql:"model"`

	// The available vehicle variants.
	Variants []Variant `graphql:"variants"`
}

// Name is the human-readable name of the vehicle: brand plus model.
func (v Vehicle) Name() string {
	return v.Brand + " " + v.Model
}

// NameWithVariants is the Name of the vehicle, plus the name of all the variants.
func (v Vehicle) NameWithVariants() string {
	variantNames := make([]string, len(v.Variants))
	for i, variant := range v.Variants {
		variantNames[i] = variant.Name()
	}
	return v.Name() + " [" + strings.Join(variantNames, ", ") + "]"
}

// Variant represents a RidePands main.Vehicle variant.
type Variant struct {
	// The vehicle color, arbitrary string. E.g. "Silver".
	Color string `graphql:"color"`

	// The vehicle size, generally in format `Name (MinHeight - MaxHeight)` where
	//   - Name is a string like S/M/L/One Size
	//   - MinHeight/MaxHeight are strings like 5'11" or 6'0"
	// But spacing can be variable - so parse liberally.
	Size string `graphql:"size"`

	// Whether the variant is in stock.
	InStock bool `graphql:"inStock"`
}

// Name is the human-readable name of the variant: color plus size.
func (v Variant) Name() string {
	return v.Color + " " + v.Size
}

// SizeCode is the size classification of the variant, without the height range. E.g. "L".
func (v Variant) SizeCode() (string, error) {
	sizeCode, _, found := strings.Cut(v.Size, "(")
	if found {
		return strings.TrimSpace(sizeCode), nil
	}

	return "", fmt.Errorf("unexpected size format '%s'", v.Size)
}
