package matching

import (
	"errors"
	"fmt"
	"github.com/danvolchek/ridepanda-notifier/internal"
	"github.com/danvolchek/ridepanda-notifier/internal/ridePandaAPI"
	"log"
	"strings"
)

// A Match is the result of a Target matching desired vehicles.
type Match struct {
	// The Target name.
	Name string

	// The vehicles that matched the criteria.
	Vehicles []ridePandaAPI.Vehicle
}

// A Matcher matches vehicles against criteria.
type Matcher struct {
	config Config

	log *log.Logger
}

func NewMatcher(config Config) *Matcher {
	return &Matcher{
		config: config,
		log:    internal.NewLogger("matcher"),
	}
}

// Match matches vehicles against the matcher's config.
func (m *Matcher) Match(vehicles []ridePandaAPI.Vehicle) ([]Match, error) {
	var results []Match

	for _, target := range m.config.Targets {
		matches, found, err := m.matches(target.Criteria, vehicles)
		if err != nil {
			return nil, fmt.Errorf("failed match for criteria '%s': %s", target.Name, err)
		}

		if found {
			results = append(results, Match{
				Name:     target.Name,
				Vehicles: matches,
			})

			m.log.Printf("Target '%s' matched vehicles: %s\n", target, strings.Join(internal.Map(matches, ridePandaAPI.Vehicle.NameWithVariants), ", "))
		} else {
			m.log.Printf("Target '%s' did not match any vehicles", target.Name)
		}
	}

	m.log.Printf("%d targets matched\n", len(results))

	return results, nil
}

func (m *Matcher) matches(c Criterion, vehicles []ridePandaAPI.Vehicle) ([]ridePandaAPI.Vehicle, bool, error) {
	var result []ridePandaAPI.Vehicle

	anyMatchedNotConsideringStock := false
	for _, vehicle := range vehicles {
		match, matchedNotConsideringStock, found, err := m.matchesSingle(c, vehicle)

		if err != nil {
			return nil, false, err
		}

		if found {
			result = append(result, match)
		}

		anyMatchedNotConsideringStock = anyMatchedNotConsideringStock || matchedNotConsideringStock
	}

	if !anyMatchedNotConsideringStock {
		return nil, false, errors.New("criteria matched no bikes regardless of stock - check if it was removed")
	}

	return result, len(result) != 0, nil
}

func (m *Matcher) matchesSingle(c Criterion, vehicle ridePandaAPI.Vehicle) (ridePandaAPI.Vehicle, bool, bool, error) {
	if !matches(vehicle.Name(), c.Name) {
		return ridePandaAPI.Vehicle{}, false, false, nil
	}

	result := ridePandaAPI.Vehicle{
		InStock:  vehicle.InStock,
		Brand:    vehicle.Brand,
		Model:    vehicle.Model,
		Variants: nil,
	}

	matchedNotConsideringStock := false

	for _, variant := range vehicle.Variants {
		sizeCode, err := variant.SizeCode()
		if err != nil {
			return ridePandaAPI.Vehicle{}, false, false, err
		}

		if !matchesExact(sizeCode, c.Size) {
			continue
		}

		if !matches(variant.Color, c.Color) {
			continue
		}

		// m.log.Printf("criteria '%s' matched vehicle %s %s regardless of stock\n", c.Name, vehicle.Name(), variant.Name())

		matchedNotConsideringStock = true

		if !variant.InStock {
			continue
		}

		result.Variants = append(result.Variants, ridePandaAPI.Variant{
			Color:   variant.Color,
			Size:    variant.Size,
			InStock: variant.InStock,
		})
	}

	return result, matchedNotConsideringStock, len(result.Variants) != 0, nil
}

func matches(s, substr string) bool {
	s, substr = strings.TrimSpace(s), strings.TrimSpace(substr)

	if substr == "*" {
		return true
	}

	for _, part := range strings.Split(substr, "|") {
		if strings.Contains(s, strings.TrimSpace(part)) {
			return true
		}
	}

	return false
}

func matchesExact(s, substr string) bool {
	s, substr = strings.TrimSpace(s), strings.TrimSpace(substr)

	if substr == "*" {
		return true
	}

	for _, part := range strings.Split(substr, "|") {
		if s == strings.TrimSpace(part) {
			return true
		}
	}

	return false
}
