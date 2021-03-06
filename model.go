package medtronic

import (
	"log"
	"strconv"
)

// Family represents a pump family.
// Use int8 so the compiler will warn about
// accidental uses of 523 instead of 23, etc.
type Family int8

// Model requests the model number from the pump and returns it,
// caching the pump family as a side effect.
// Use Family to avoid contacting the pump more than once.
func (pump *Pump) Model() string {
	data := pump.Execute(model)
	if pump.Error() != nil {
		return ""
	}
	if len(data) < 2 {
		pump.BadResponse(model, data)
		return ""
	}
	n := int(data[1])
	if len(data) < 2+n {
		pump.BadResponse(model, data)
		return ""
	}
	m := string(data[2 : 2+n])
	pump.cacheFamily(m)
	return m
}

func (pump *Pump) cacheFamily(model string) {
	if pump.family != 0 {
		return
	}
	log.Printf("model %s pump", model)
	family := -1
	n, err := strconv.Atoi(model)
	if err != nil {
		log.Printf("%v", err)
	} else if 500 < n && n < 600 {
		family = n - 500
	} else if 700 < n && n < 800 {
		family = n - 700
	} else {
		log.Printf("unsupported pump model %d", n)
	}
	pump.family = Family(family)
}

// Family returns 22 for 522/722 pumps, 23 for 523/723 pumps, etc.,
// and returns -1 for an unrecognized model.  It calls Model once and
// caches the result.
func (pump *Pump) Family() Family {
	if pump.family == 0 {
		pump.Model()
		if pump.Error() != nil {
			return -1
		}
	}
	return pump.family
}
