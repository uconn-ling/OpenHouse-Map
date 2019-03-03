package walk

import (
  "fmt"
  "testing"
  "../src/walk"
)

func TestGetCountries (t *testing.T) {
  c1 := walk.GetCountries("./data", "country_")
  if len(c1) != 1 {
    t.Errorf("getCountries(Kronos) = %d; want 1", c1)
}
