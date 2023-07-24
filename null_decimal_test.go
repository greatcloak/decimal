package decimal

import (
	"database/sql/driver"
	"encoding/xml"
	"testing"
)

func TestNullDecimalXML(t *testing.T) {
	// test valid values
	for _, x := range testTable {
		s := x.short
		var doc struct {
			XMLName xml.Name    `xml:"account"`
			Amount  NullDecimal `xml:"amount"`
		}
		docStr := `<account><amount>` + s + `</amount></account>`
		err := xml.Unmarshal([]byte(docStr), &doc)
		if err != nil {
			t.Errorf("error unmarshaling %s: %v", docStr, err)
		} else if doc.Amount.Decimal.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, doc.Amount.Decimal.String(),
				doc.Amount.Decimal.value.String(), doc.Amount.Decimal.exp)
		}

		out, err := xml.Marshal(&doc)
		if err != nil {
			t.Errorf("error marshaling %+v: %v", doc, err)
		} else if string(out) != docStr {
			t.Errorf("expected %s, got %s", docStr, string(out))
		}
	}

	var doc struct {
		XMLName xml.Name    `xml:"account"`
		Amount  NullDecimal `xml:"amount"`
	}

	// test for XML with empty body
	docStr := `<account><amount></amount></account>`
	err := xml.Unmarshal([]byte(docStr), &doc)
	if err != nil {
		t.Errorf("error unmarshaling: %s: %v", docStr, err)
	} else if doc.Amount.Valid {
		t.Errorf("expected null value to have Valid = false, got Valid = true and Decimal = %s (%s, %d)",
			doc.Amount.Decimal.String(),
			doc.Amount.Decimal.value.String(), doc.Amount.Decimal.exp)
	}

	expected := `<account><amount></amount></account>`
	out, err := xml.Marshal(&doc)
	if err != nil {
		t.Errorf("error marshaling %+v: %v", doc, err)
	} else if string(out) != expected {
		t.Errorf("expected %s, got %s", expected, string(out))
	}

	// test for empty XML
	docStr = `<account></account>`
	err = xml.Unmarshal([]byte(docStr), &doc)
	if err != nil {
		t.Errorf("error unmarshaling: %s: %v", docStr, err)
	} else if doc.Amount.Valid {
		t.Errorf("expected null value to have Valid = false, got Valid = true and Decimal = %s (%s, %d)",
			doc.Amount.Decimal.String(),
			doc.Amount.Decimal.value.String(), doc.Amount.Decimal.exp)
	}

	expected = `<account><amount></amount></account>`
	out, err = xml.Marshal(&doc)
	if err != nil {
		t.Errorf("error marshaling %+v: %v", doc, err)
	} else if string(out) != expected {
		t.Errorf("expected %s, got %s", expected, string(out))
	}
}

func TestNullDecimalBadXML(t *testing.T) {
	for _, testCase := range []string{
		"o_o",
		"<abc",
		"<account><amount>7",
		`<html><body></body></html>`,
		`<account><amount>nope</amount></account>`,
		`0.333`,
	} {
		var doc struct {
			XMLName xml.Name    `xml:"account"`
			Amount  NullDecimal `xml:"amount"`
		}
		err := xml.Unmarshal([]byte(testCase), &doc)
		if err == nil {
			t.Errorf("expected error, got %+v", doc)
		}
	}
}

func TestNullDecimal_Scan(t *testing.T) {
	// test the Scan method that implements the
	// sql.Scanner interface
	// check for the for different type of values
	// that are possible to be received from the database
	// drivers

	// in normal operations the db driver (sqlite at least)
	// will return an int64 if you specified a numeric format

	// Make sure handles nil values
	a := NullDecimal{}
	var dbvaluePtr interface{}
	err := a.Scan(dbvaluePtr)
	if err != nil {
		// Scan failed... no need to test result value
		t.Errorf("a.Scan(nil) failed with message: %s", err)
	} else {
		if a.Valid {
			t.Errorf("%s is not null", a.Decimal)
		}
	}

	dbvalue := 54.33
	expected := NewFromFloat(dbvalue)

	err = a.Scan(dbvalue)
	if err != nil {
		// Scan failed... no need to test result value
		t.Errorf("a.Scan(54.33) failed with message: %s", err)

	} else {
		// Scan succeeded... test resulting values
		if !a.Valid {
			t.Errorf("%s is null", a.Decimal)
		} else if !a.Decimal.Equals(expected) {
			t.Errorf("%s does not equal to %s", a.Decimal, expected)
		}
	}

	// at least SQLite returns an int64 when 0 is stored in the db
	// and you specified a numeric format on the schema
	dbvalueInt := int64(0)
	expected = New(dbvalueInt, 0)

	err = a.Scan(dbvalueInt)
	if err != nil {
		// Scan failed... no need to test result value
		t.Errorf("a.Scan(0) failed with message: %s", err)

	} else {
		// Scan succeeded... test resulting values
		if !a.Valid {
			t.Errorf("%s is null", a.Decimal)
		} else if !a.Decimal.Equals(expected) {
			t.Errorf("%v does not equal %v", a, expected)
		}
	}

	// in case you specified a varchar in your SQL schema,
	// the database driver will return byte slice []byte
	valueStr := "535.666"
	dbvalueStr := []byte(valueStr)
	expected, err = NewFromString(valueStr)
	if err != nil {
		t.Fatal(err)
	}

	err = a.Scan(dbvalueStr)
	if err != nil {
		// Scan failed... no need to test result value
		t.Errorf("a.Scan('535.666') failed with message: %s", err)

	} else {
		// Scan succeeded... test resulting values
		if !a.Valid {
			t.Errorf("%s is null", a.Decimal)
		} else if !a.Decimal.Equals(expected) {
			t.Errorf("%v does not equal %v", a, expected)
		}
	}

	// lib/pq can also return strings
	expected, err = NewFromString(valueStr)
	if err != nil {
		t.Fatal(err)
	}

	err = a.Scan(valueStr)
	if err != nil {
		// Scan failed... no need to test result value
		t.Errorf("a.Scan('535.666') failed with message: %s", err)
	} else {
		// Scan succeeded... test resulting values
		if !a.Valid {
			t.Errorf("%s is null", a.Decimal)
		} else if !a.Decimal.Equals(expected) {
			t.Errorf("%v does not equal %v", a, expected)
		}
	}
}

func TestNullDecimal_Value(t *testing.T) {
	// Make sure this does implement the database/sql's driver.Valuer interface
	var nullDecimal NullDecimal
	if _, ok := interface{}(nullDecimal).(driver.Valuer); !ok {
		t.Error("NullDecimal does not implement driver.Valuer")
	}

	// check that null is handled appropriately
	value, err := nullDecimal.Value()
	if err != nil {
		t.Errorf("NullDecimal{}.Valid() failed with message: %s", err)
	} else if value != nil {
		t.Errorf("%v is not nil", value)
	}

	// check that normal case is handled appropriately
	a := NullDecimal{Decimal: New(1234, -2), Valid: true}
	expected := "12.34"
	value, err = a.Value()
	if err != nil {
		t.Errorf("NullDecimal(12.34).Value() failed with message: %s", err)
	} else if value.(string) != expected {
		t.Errorf("%v does not equal %v", a, expected)
	}
}
