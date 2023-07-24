package decimal

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var _ bson.ValueUnmarshaler = (*Decimal)(nil)

func (dec *Decimal) UnmarshalBSONValue(bt bsontype.Type, raw []byte) error {

	// Only read from string
	// TODO should we support mongodb's Decimal128
	if bt != bson.TypeString {
		// TODO we can also add support for reading in Decimal128 or other number types from mongodb.
		return fmt.Errorf("expected bson string for Decimal type but received: %s", bt)
	}

	var valAsString string

	err := bson.UnmarshalValue(bt, raw, &valAsString)
	if err != nil {
		return err
	}

	result, err := NewFromString(valAsString)
	if err != nil {
		return fmt.Errorf("error creating new decimal from string: %w", err)
	}

	// val.Set(reflect.ValueOf(result))
	*dec = result

	return nil
}

var _ bson.ValueMarshaler = Decimal{}

// MarshalBSONValue implments marshalling a GridStorageFile to bson as a string.
//
// Mongodb also supports Decimal128 which has 31 digits of precision.
// You can register a custom codec w/ bson if this is enough precision.
// For finacial systems it usually is not.
func (dec Decimal) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(dec.String())
	// return bson.Marshal(dec.String())
}
