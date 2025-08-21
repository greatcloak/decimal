package decimal

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// RegisterBSONDecimalCodec registers the encoder and decoder for [Decimal].
// The encoder/decoder are automatically registered on [bson.DefaultRegistry] on package import.
func RegisterBSONDecimalCodec(registry *bson.Registry) {
	// Register custom encoder and decoder for decimal type
	registry.RegisterTypeEncoder(decimalType, bson.ValueEncoderFunc(decimalBSONEncodeValue))
	registry.RegisterTypeDecoder(decimalType, bson.ValueDecoderFunc(decimalBSONDecodeValue))
}

// decimalType is the reflected type of a [Decimal].
var decimalType = reflect.TypeOf(Decimal{})

// decimalBSONEncodeValue encodes a shopspring decimal to string.
// This is the safest format as it contains the entire value and can be decoded from.
func decimalBSONEncodeValue(ec bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	if val.Type() != decimalType {
		// ShopSpring decimal
		return bson.ValueEncoderError{
			// prefix with GC to avoid any mongodb name collisions
			// They have a Decimal128 type which this is not.
			Name:     "GCDecimalEncodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	// Convert to a decimal
	dec := val.Interface().(Decimal)

	return vw.WriteString(dec.String())
}

// decimalBSONDecodeValue decodes a string into a decimal.
func decimalBSONDecodeValue(dc bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {

	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bson.ValueDecoderError{
			// prefix with GC to avoid any mongodb name collisions
			// They have a Decimal128 type which this is not.
			Name:     "GCDecimalDecodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}

	// Only read from string
	if vr.Type() != bson.TypeString {
		return fmt.Errorf("received invalid BSON type to decode into Decimal: %s", vr.Type())
	}
	b, err := vr.ReadString()
	if err != nil {
		return err
	}

	result, err := NewFromString(b)
	if err != nil {
		return fmt.Errorf("error creating new decimal from string: %w", err)
	}

	val.Set(reflect.ValueOf(result))

	return nil
}
