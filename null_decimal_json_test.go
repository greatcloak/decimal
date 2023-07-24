package decimal

import (
	"encoding/json"
	"testing"
)

func TestNullDecimalJSON(t *testing.T) {
	for _, x := range testTable {
		s := x.short
		var doc struct {
			Amount NullDecimal `json:"amount"`
		}
		docStr := `{"amount":"` + s + `"}`
		docStrNumber := `{"amount":` + s + `}`
		err := json.Unmarshal([]byte(docStr), &doc)
		if err != nil {
			t.Errorf("error unmarshaling %s: %v", docStr, err)
		} else {
			if !doc.Amount.Valid {
				t.Errorf("expected %s to be valid (not NULL), got Valid = false", s)
			}
			if doc.Amount.Decimal.String() != s {
				t.Errorf("expected %s, got %s (%s, %d)",
					s, doc.Amount.Decimal.String(),
					doc.Amount.Decimal.value.String(), doc.Amount.Decimal.exp)
			}
		}

		out, err := json.Marshal(&doc)
		if err != nil {
			t.Errorf("error marshaling %+v: %v", doc, err)
		} else if string(out) != docStr {
			t.Errorf("expected %s, got %s", docStr, string(out))
		}

		// make sure unquoted marshalling works too
		MarshalJSONWithoutQuotes = true
		out, err = json.Marshal(&doc)
		if err != nil {
			t.Errorf("error marshaling %+v: %v", doc, err)
		} else if string(out) != docStrNumber {
			t.Errorf("expected %s, got %s", docStrNumber, string(out))
		}
		MarshalJSONWithoutQuotes = false
	}

	var doc struct {
		Amount NullDecimal `json:"amount"`
	}
	docStr := `{"amount": null}`
	err := json.Unmarshal([]byte(docStr), &doc)
	if err != nil {
		t.Errorf("error unmarshaling %s: %v", docStr, err)
	} else if doc.Amount.Valid {
		t.Errorf("expected null value to have Valid = false, got Valid = true and Decimal = %s (%s, %d)",
			doc.Amount.Decimal.String(),
			doc.Amount.Decimal.value.String(), doc.Amount.Decimal.exp)
	}

	expected := `{"amount":null}`
	out, err := json.Marshal(&doc)
	if err != nil {
		t.Errorf("error marshaling %+v: %v", doc, err)
	} else if string(out) != expected {
		t.Errorf("expected %s, got %s", expected, string(out))
	}

	// make sure unquoted marshalling works too
	MarshalJSONWithoutQuotes = true
	expectedUnquoted := `{"amount":null}`
	out, err = json.Marshal(&doc)
	if err != nil {
		t.Errorf("error marshaling %+v: %v", doc, err)
	} else if string(out) != expectedUnquoted {
		t.Errorf("expected %s, got %s", expectedUnquoted, string(out))
	}
	MarshalJSONWithoutQuotes = false
}

func TestNullDecimalBadJSON(t *testing.T) {
	for _, testCase := range []string{
		"]o_o[",
		"{",
		`{"amount":""`,
		`{"amount":"nope"}`,
		`{"amount":nope}`,
		`0.333`,
	} {
		var doc struct {
			Amount NullDecimal `json:"amount"`
		}
		err := json.Unmarshal([]byte(testCase), &doc)
		if err == nil {
			t.Errorf("expected error, got %+v", doc)
		}
	}
}
