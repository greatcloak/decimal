package decimal

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

func init() {
	// Auto register our Decimal encoder and decoder for bson when this package is imported.
	RegisterBSONDecimalCodec(bson.DefaultRegistry)
}

// RegisterBSONDecimalCodec registers the encoder and decoder for [Decimal].
// The encoder/decoder are automatically registered on [bson.DefaultRegistry] on package import.
func RegisterBSONDecimalCodec(registry *bsoncodec.Registry) {
	// Register custom encoder and decoder for decimal type
	registry.RegisterTypeEncoder(decimalType, bsoncodec.ValueEncoderFunc(decimalBSONEncodeValue))
	registry.RegisterTypeDecoder(decimalType, bsoncodec.ValueDecoderFunc(decimalBSONDecodeValue))
}

// decimalType is the reflected type of a [Decimal].
var decimalType = reflect.TypeOf(Decimal{})

// decimalBSONEncodeValue encodes a shopspring decimal to string.
// This is the safest format as it contains the entire value and can be decoded from.
func decimalBSONEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Type() != decimalType {
		// ShopSpring decimal
		return bsoncodec.ValueEncoderError{
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
func decimalBSONDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bsoncodec.ValueDecoderError{
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
