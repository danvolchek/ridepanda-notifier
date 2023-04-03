package matching

type Config struct {
	// A list of notifications to notify on.
	Targets []Target `json:"targets"`
}

// Target represents a desired bike criteria.
type Target struct {
	// Human-readable name for this target, explains what the criteria are.
	Name string `json:"name"`

	// The criteria to match vehicles against. For each criterion, matching vehicles are notified.
	Criteria Criterion `json:"criteria"`
}

// Criterion represents a selector for matching an ridePandaAPI.Vehicle.
//
// All fields must match for a vehicle to be matched.
//
// For Name and Color, matching is done by substring check after removing leading/trailing whitespace.
// For Size, matching is exact.
//
// Empty fields are invalid. Use "*" to indicate any value. Use "|" to indicate or.
type Criterion struct {
	// The name of the vehicle to match, which is `Brand Model`. E.g. "Ride1Up Roadster".
	Name string `json:"name"`

	// The size of the vehicle to match (excluding height measurements). E.g. "L".
	Size string `json:"size"`

	// The color of the vehicle to match. E.g. "Silver".
	Color string `json:"color"`
}
