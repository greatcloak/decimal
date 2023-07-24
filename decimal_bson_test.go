package decimal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// TestBSONDecimal tests marshalling and unmarshalling decimals as values and pointers.
func TestBSONDecimal(t *testing.T) {

	type TestData struct {
		// Test both as a value and as a pointer
		DVal     Decimal
		DPointer *Decimal

		// Ensure regular fields are marshalling/unmarshalling correctly
		Name   string
		Number int
	}

	td := TestData{Name: "tester", Number: 23}
	x := NewFromFloat(32.91)
	td.DPointer = &x

	td.DVal = NewFromInt(99)

	out, err := bson.Marshal(td)
	require.NoError(t, err, "should marshal bson")

	var td2 TestData
	err = bson.Unmarshal(out, &td2)
	require.NoError(t, err, "should unmarshal bson")

	require.Equal(t, td.DPointer.String(), td2.DPointer.String())
	require.Equal(t, td.DVal.String(), td2.DVal.String())

	// Make sure other non decimal and non string fields continue to marshal and unmarshal correctly
	// ie we didn't break the registry
	require.Equal(t, td.Name, td2.Name)
	require.Equal(t, td.Number, td2.Number)
}
