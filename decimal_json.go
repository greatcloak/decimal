package decimal

import "fmt"

// MarshalJSONWithoutQuotes should be set to true if you want the decimal to
// be JSON marshaled as a number, instead of as a string.
// WARNING: this is dangerous for decimals with many digits, since many JSON
// unmarshallers (ex: Javascript's) will unmarshal JSON numbers to IEEE 754
// double-precision floating point numbers, which means you can potentially
// silently lose precision.
//
// Set it once at program start. Changing it while other goroutines use the package is a data race.
var MarshalJSONWithoutQuotes = false

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	switch string(decimalBytes) {
	case "null":
		// TODO should this set to zero value or nil?
		return nil
	case `""`:
		// Set to zero value
		// Allow enter an empty string which is common if you work with web apps or json apis.
		// HTML typically uses a string based input which is the best representation for our Decimal.
		// The empty value for that is an empty string.
		//
		// Make sure we do not return an error when reading empty strings. We interpret them as a zero.
		*d = Decimal{}
		return nil
	}

	decimal, err := NewFromString(unquoteIfQuoted(string(decimalBytes)))
	*d = decimal
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", string(decimalBytes), err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	str := d.String()
	if MarshalJSONWithoutQuotes {
		return []byte(str), nil
	}
	b := make([]byte, 0, len(str)+2)
	b = append(b, '"')
	b = append(b, str...)
	return append(b, '"'), nil
}
