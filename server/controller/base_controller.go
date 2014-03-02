package controller

import (
	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type Controller interface {
	RegisterHandlers(r *mux.Router)
}
