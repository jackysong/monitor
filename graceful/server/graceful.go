package server

import (
	"net/http"

	"github.com/kisekivul/utils"
)

func Initialize_Graceful(host, directory string, port int, archivable, listable, slashing, pushState bool) {
	service := Service{
		Host: host,
		Port: port,
		Config: Config{
			Directory:  directory,
			Archivable: archivable,
			Listable:   listable,
			Slashing:   slashing,
			PushState:  pushState,
		},
	}

	handler, err := NewHandler(service.Config)
	utils.ErrorStop(err)

	http.ListenAndServe(host+":"+utils.Int2Str(port), handler)
}
