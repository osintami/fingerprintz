// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/osintami/fingerprintz/common"
)

type ItemInfo struct {
	Category    string `json:"-"`
	Source      string `json:"-"`
	Item        string
	Type        string
	Enabled     bool
	Description string
}

func (x *NormalizedDataServer) DictionaryHandler(w http.ResponseWriter, r *http.Request) {
	wantCSV, _ := strconv.ParseBool(r.URL.Query().Get("csv"))
	bAdmin := r.URL.Query().Get("role") == "admin"
	bUser := r.URL.Query().Get("role") == "user"

	if !bAdmin && !bUser {
		common.SendError(w, ErrInvalidUserRole, http.StatusForbidden)
		return
	}

	adminList := x.schema.ListItems()
	sort.Slice(adminList, func(i, j int) bool {
		return adminList[i].Path < adminList[j].Path
	})

	if bAdmin {
		if !wantCSV {
			common.SendJSON(w, adminList)
		} else {
			var userList []ItemInfo
			for _, item := range adminList {
				uri := NewItemSplitter(item.Path)
				userList = append(userList, ItemInfo{
					Category:    uri.CategoryName,
					Source:      uri.SourceName,
					Item:        uri.ItemName,
					Enabled:     item.Enabled,
					Type:        item.Type.String(),
					Description: item.Description})
			}
			x.ItemToCSV(w, userList)
		}
		return
	}

	if bUser {
		var userList []ItemInfo
		for _, item := range adminList {
			if item.Enabled && item.SourceName == SOURCE_OSINTAMI_NAME {
				uri := NewItemSplitter(item.Path)
				userList = append(userList, ItemInfo{
					Category:    uri.CategoryName,
					Source:      item.SourceName,
					Item:        item.Path,
					Enabled:     item.Enabled,
					Type:        item.Type.String(),
					Description: item.Description})
			}
		}
		if !wantCSV {
			common.SendJSON(w, userList)
		} else {
			x.ItemToCSV(w, userList)
		}
		return
	}
}

func (x *NormalizedDataServer) ItemToCSV(w http.ResponseWriter, items []ItemInfo) {
	writer := tabwriter.NewWriter(w, 0, 8, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(writer, "CATEGORY,\tITEM,\tSOURCE,\tENABLED,\tTYPE,\tDESCRIPTION\n")
	for _, item := range items {
		_, _ = fmt.Fprintf(writer, "%s,\t%s,\t%s,\t%t,\t%s,\t%s\n",
			item.Category,
			item.Source,
			item.Item,
			item.Enabled,
			item.Type,
			item.Description)
	}
	_ = writer.Flush()
}
