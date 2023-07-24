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
		// Output nothing because we don't set value and it's omitempty
		DPointerEmpty *Decimal `bson:",omitempty"`
		// Output a null because we don't set value and there is no omitempty
		DPointerEmptyButNull *Decimal

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
	// require.Nil(t, td.DPointerEmpty, "omitempty bson should work")
	require.Equal(t, td.DVal.String(), td2.DVal.String())

	// Make sure other non decimal and non string fields continue to marshal and unmarshal correctly
	// ie we didn't break the registry
	require.Equal(t, td.Name, td2.Name)
	require.Equal(t, td.Number, td2.Number)

	var td3 bson.M
	err = bson.Unmarshal(out, &td3)
	require.NoError(t, err)

	_, isOmitemptyFieldInMap := td3["dpointerempty"]
	require.False(t, isOmitemptyFieldInMap, "bson omitempty should not be marshalled out as null in map, it should be omitted")

	_, isInMap := td3["dpointer"]
	require.True(t, isInMap, "bson pointer decimal field without omitempty should be marshalled out as null")

	dbpointerEmptyNoOmit, isInMap := td3["dpointeremptybutnull"]
	require.True(t, isInMap, "bson pointer decimal field without omitempty should be marshalled out as null")
	require.Nil(t, dbpointerEmptyNoOmit)

}
