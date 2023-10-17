package common

import "github.com/polera/gorbl"

type IRealtimeBlackholeList interface {
	Lookup(host, addr string) bool
}

type RealtimeBlackholeList struct {
}

func NewRealtimeBlackholeList() IRealtimeBlackholeList {
	return &RealtimeBlackholeList{}
}

func (x *RealtimeBlackholeList) Lookup(rblServer, blacklistedHost string) bool {
	info := gorbl.Lookup(rblServer, blacklistedHost)
	return len(info.Results) > 0 && info.Results[0].Listed
}
