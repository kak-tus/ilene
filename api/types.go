package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/json-iterator/go"
	"github.com/kak-tus/ilene/model"
	"go.uber.org/zap"
)

// Type type
type Type struct {
	cnf  apiConfig
	enc  jsoniter.API
	lock *sync.Mutex
	log  *zap.SugaredLogger
	mdl  *model.Type
	srv  *http.Server
	tick *time.Ticker
}

type apiConfig struct {
	Addr    string
	DataDir string
	HTTPDir string
	Schema  string
}

type postResponseWriter struct {
	code int
	data []byte
	http.ResponseWriter
}
