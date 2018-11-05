package api

import (
	"net/http"

	"github.com/json-iterator/go"
	"github.com/kak-tus/ilene/model"
	"go.uber.org/zap"
)

// Type type
type Type struct {
	cnf apiConfig
	log *zap.SugaredLogger
	enc jsoniter.API
	srv *http.Server
	mdl *model.Type
}

type apiConfig struct {
	Schema string
	Addr   string
}
