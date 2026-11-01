package decimal

import "database/sql/driver"

// NullDecimal represents a nullable decimal with compatibility for
// scanning null values from the database.
type NullDecimal struct {
	Decimal Decimal

	// Valid is true if the Decimal field is non null.
	Valid bool
}

// NewNullDecimal returns a valid NullDecimal holding d.
func NewNullDecimal(d Decimal) NullDecimal {
	return NullDecimal{
		Decimal: d,
		Valid:   true,
	}
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *NullDecimal) Scan(value interface{}) error {
	if value == nil {
		d.Valid = false
		return nil
	}
	err := d.Decimal.Scan(value)
	d.Valid = err == nil
	return err
}

// Value implements the driver.Valuer interface for database serialization.
func (d NullDecimal) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Decimal.Value()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization
func (d *NullDecimal) UnmarshalText(text []byte) error {
	str := string(text)

	// check for empty XML or XML without body e.g., <tag></tag>
	if str == "" {
		d.Valid = false
		return nil
	}
	if err := d.Decimal.UnmarshalText(text); err != nil {
		d.Valid = false
		return err
	}
	d.Valid = true
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d NullDecimal) MarshalText() (text []byte, err error) {
	if !d.Valid {
		return []byte{}, nil
	}
	return d.Decimal.MarshalText()
}

// DecodeSpanner decodes a Spanner value into a Decimal
func (d *NullDecimal) DecodeSpanner(value interface{}) error {
	switch t := value.(type) {
	case nil:
		d.Valid = false
		return nil
	case *string:
		if t == nil {
			d.Valid = false
			return nil
		}
		value = *t
	}

	err := d.Decimal.Scan(value)
	d.Valid = err == nil
	return err
}

// EncodeSpanner encodes a Decimal into a Spanner value
func (d NullDecimal) EncodeSpanner() (interface{}, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Decimal.String(), nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface of the MongoDB Go driver v2.
// An invalid NullDecimal is stored as BSON null, see Decimal.MarshalBSONValue for the rest.
func (d NullDecimal) MarshalBSONValue() (byte, []byte, error) {
	if !d.Valid {
		return bsonTypeNull, nil, nil
	}
	return d.Decimal.MarshalBSONValue()
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface of the MongoDB Go driver v2.
// BSON null and undefined make the NullDecimal invalid, see Decimal.UnmarshalBSONValue for the rest.
func (d *NullDecimal) UnmarshalBSONValue(typ byte, data []byte) error {
	err := d.Decimal.UnmarshalBSONValue(typ, data)
	d.Valid = err == nil && typ != bsonTypeNull && typ != bsonTypeUndefined
	return err
}
