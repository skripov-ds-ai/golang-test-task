package database

import "testing"

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

func TestCreateMapFromAdItemNotNilMainImage(t *testing.T) {
	s := ""
	var url = &s
	mainURL := ImageURL{URL: *url}
	item := AdItem{MainImageURL: &mainURL}
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
		if r2, ok := val.(*string); ok {
			r1, ok := v.(*string)
			if ok {
				if *r1 == *r2 {
					continue
				}
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, *r1, *r2, expected, actual)
			}
		}
		if v == nil && val != nil || v != nil && v != val {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}

func TestCreateMapFromAdItemEmptyItemAdditionalFieldsEmptyImageList(t *testing.T) {
	var url *string
	item := AdItem{}
	expected := map[string]interface{}{"id": 0, "main_image_url": url, "title": "", "description": "", "image_urls": []string{}}
	actual := item.CreateMapFromAdItem([]string{"description", "image_urls"})

	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual) ; expected = %v ; actual = %v", expected, actual)
	}
	for k, v := range expected {
		val, ok := actual[k]
		if !ok {
			t.Fatalf("k = %s not found in actual ; expected = %v ; actual = %v", k, expected, actual)
		}
		if k == "image_urls" {
			r1, ok1 := v.([]string)
			r2, ok2 := val.([]string)
			if !ok1 || !ok2 || len(r1) != len(r2) {
				t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
			}
		} else if v == nil && val != nil || v != nil && v != val {
			t.Fatalf("expected[%s] != actual[%s] ; %v != %v ; expected = %v ; actual = %v", k, k, v, val, expected, actual)
		}
	}
}
