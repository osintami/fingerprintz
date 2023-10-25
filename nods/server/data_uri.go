// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"strings"
)

type DataURI struct {
	CategoryName string
	SourceName   string
	ItemName     string
	URI          string
}

func NewItemSplitter(itemName string) *DataURI {
	parts := strings.SplitN(itemName, "/", 3)
	if len(parts) < 3 || parts[0] == "" {
		return NewDataURI("unknown", "unknown", "unknown")
	}
	return NewDataURI(parts[0], parts[1], parts[2])
}

func NewDataURI(categoryName, sourceName, itemName string) *DataURI {
	return &DataURI{
		CategoryName: categoryName,
		SourceName:   sourceName,
		ItemName:     itemName,
		URI:          categoryName + "/" + sourceName + "/" + itemName}
}

func (x DataURI) Key() string {
	return x.URI
}
