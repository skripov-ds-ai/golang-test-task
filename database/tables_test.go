package database

import (
	"testing"
)

func TestCreateMapFromAdItemEmptyItem(t *testing.T) {
	item := AdItem{}
	var url *string
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": ""}
	actual := item.CreateMapFromAdItem([]string{})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	for k, v := range expected {
		val, ok := actual[k]
		if !ok {
			t.Fatalf("k = %s not found in actual ; expected = %v ; actual = %v", k, expected, actual)
		}
		if v == nil && val != nil || v != nil && v != val {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}
