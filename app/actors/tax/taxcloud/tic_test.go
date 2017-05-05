package taxcloud_test

import (
	"testing"
	"github.com/ottemo/foundation/app/actors/tax/taxcloud"
)

func TestGetTICs(t *testing.T) {
	result, err := testGateway.GetTICs()
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if result == nil {
		t.Fatal("result is empty")
	} else if result.TICs == nil {
		t.Fatal("no TICs found")
	} else if len(result.TICs) <= 0 {
		t.Fatal("expected array of TICs")
	}
}

func TestGetTICGroups(t *testing.T) {
	result, err := testGateway.GetTICGroups()
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if result == nil {
		t.Fatal("result is empty")
	} else if result.TICGroups == nil {
		t.Fatal("no TICGroups found")
	} else if len(result.TICGroups) <= 0 {
		t.Fatal("expected array of TICGroups")
	}
}

func TestGetTICsByUnknownGroup(t *testing.T) {
	result, err := testGateway.GetTICsByGroup(taxcloud.GetTICsByGroupParams{
		GroupID: 1, // non real value
	})
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if result == nil {
		t.Fatal("result is empty")
	} else if result.TICs == nil {
		t.Fatal("no TICs found")
	} else if len(result.TICs) != 0 {
		t.Fatal("expected empty array of TICs")
	}
}

func TestGetTICsByGroup(t *testing.T) {
	groups, err := testGateway.GetTICGroups()
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	tics, err := testGateway.GetTICsByGroup(taxcloud.GetTICsByGroupParams{
		GroupID: groups.TICGroups[0].GroupID,
	})
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	if tics == nil {
		t.Fatal("result is empty")
	} else if tics.TICs == nil {
		t.Fatal("no TICs found")
	} else if len(tics.TICs) <= 0 {
		t.Fatal("expected array of TICs")
	}
}
