package router

import (
	"net/http"

	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/pkg/proxy/gateway/vald/handler/rest"
)

// NewRoutes returns REST route&method information from handler interface
func NewRoutes(h rest.Handler) []routing.Route {
	return []routing.Route{
		{
			"Index",
			[]string{
				http.MethodGet,
			},
			"/",
			h.Index,
		},
		{
			"Search",
			[]string{
				http.MethodPost,
			},
			"/search",
			h.Search,
		},
		{
			"Search By ID",
			[]string{
				http.MethodGet,
			},
			"/search/{id}",
			h.SearchByID,
		},
		{
			"Insert",
			[]string{
				http.MethodPost,
			},
			"/insert",
			h.Insert,
		},
		{
			"Multiple Insert",
			[]string{
				http.MethodPost,
			},
			"/insert/multi",
			h.MultiInsert,
		},
		{
			"Update",
			[]string{
				http.MethodPost,
				http.MethodPatch,
				http.MethodPut,
			},
			"/update",
			h.Update,
		},
		{
			"Multiple Update",
			[]string{
				http.MethodPost,
				http.MethodPatch,
				http.MethodPut,
			},
			"/update/multi",
			h.MultiUpdate,
		},
		{
			"Remove",
			[]string{
				http.MethodDelete,
			},
			"/delete/{id}",
			h.Remove,
		},
		{
			"Multiple Remove",
			[]string{
				http.MethodDelete,
				http.MethodPost,
			},
			"/delete/multi",
			h.MultiRemove,
		},
		{
			"Create Index",
			[]string{
				http.MethodGet,
			},
			"/index/create/{pool}",
			h.CreateIndex,
		},
		{
			"Save Index",
			[]string{
				http.MethodGet,
			},
			"/index/save",
			h.SaveIndex,
		},
		{
			"GetObject",
			[]string{
				http.MethodGet,
			},
			"/object/{id}",
			h.GetObject,
		},
	}
}
