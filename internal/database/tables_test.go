package database

import (
	"golang-test-task/internal/entities"
	"reflect"
	"testing"
)

func TestAdAPIListItem_CreateMap(t *testing.T) {
	item := AdListItem{}
	expected := entities.APIAdListItem{}
	actual := item.CreateMap()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapEmptyItem(t *testing.T) {
	item := AdItem{}
	expected := entities.APIAdItem{}
	actual := item.CreateMap([]string{})
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapWithMainImage(t *testing.T) {
	s := ""
	var url = &s
	mainURL := ImageURL{URL: *url}
	item := AdItem{MainImageURL: &mainURL}
	expected := entities.APIAdItem{MainImageURL: url}
	actual := item.CreateMap([]string{})
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}

func TestAdItem_CreateMapEmptyItemAdditionalFieldsEmptyImageList(t *testing.T) {
	description := "description"
	item := AdItem{Description: description}
	expected := entities.APIAdItem{Description: description}
	actual := item.CreateMap([]string{"description"})
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected != actual ; expected = %v ; actual = %v", expected, actual)
	}
}
