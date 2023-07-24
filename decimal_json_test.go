package decimal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	for _, x := range testTable {
		s := x.short
		var doc struct {
			Amount Decimal `json:"amount"`
		}
		docStr := `{"amount":"` + s + `"}`
		docStrNumber := `{"amount":` + s + `}`
		err := json.Unmarshal([]byte(docStr), &doc)
		if err != nil {
			t.Errorf("error unmarshaling %s: %v", docStr, err)
		} else if doc.Amount.String() != s {
			t.Errorf("expected %s, got %s (%s, %d)",
				s, doc.Amount.String(),
				doc.Amount.value.String(), doc.Amount.exp)
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
}

func TestUnmarshalJSONNull(t *testing.T) {
	var doc struct {
		Amount Decimal `json:"amount"`
	}
	docStr := `{"amount": null}`
	err := json.Unmarshal([]byte(docStr), &doc)
	if err != nil {
		t.Errorf("error unmarshaling %s: %v", docStr, err)
	} else if !doc.Amount.Equal(Zero) {
		t.Errorf("expected Zero, got %s (%s, %d)",
			doc.Amount.String(),
			doc.Amount.value.String(), doc.Amount.exp)
	}
}

func TestBadJSON(t *testing.T) {
	for _, testCase := range []string{
		"]o_o[",
		"{",
		`{"amount":""`,
		`{"amount":"nope"}`,
		`0.333`,
	} {
		var doc struct {
			Amount Decimal `json:"amount"`
		}
		err := json.Unmarshal([]byte(testCase), &doc)
		if err == nil {
			t.Errorf("expected error, got %+v", doc)
		}
	}
}

func TestReadEmptyStringFromJSON(t *testing.T) {
	testCase := `{"amount":""}`
	var doc struct {
		Amount Decimal `json:"amount"`
	}

	err := json.Unmarshal([]byte(testCase), &doc)
	require.NoError(t, err)
	require.True(t, Zero.Equal(doc.Amount))
}
