package walk

import (
	"testing"

	"github.com/uconn-ling/openHouseMap/src/walk"
)

func TestGetCountries(t *testing.T) {
	c1 := walk.GetData("./data", "country_")
	if len(c1) != 1 {
		t.Errorf("getCountries(Kronos) = %d; want 1", c1)
	}
}
