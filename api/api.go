package api

import (
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/applog"
	"git.aqq.me/go/app/event"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	oas "github.com/hypnoglow/oas2"
	"github.com/iph0/conf"
	jsoniter "github.com/json-iterator/go"
	"github.com/kak-tus/ilene/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var obj *Type

func init() {
	event.Init.AddHandler(
		func() error {
			cnfMap := appconf.GetConfig()["api"]

			var cnf apiConfig
			err := conf.Decode(cnfMap, &cnf)
			if err != nil {
				return err
			}

			srv := &http.Server{
				Addr: cnf.Addr,
			}

			obj = &Type{
				cnf:  cnf,
				enc:  jsoniter.Config{UseNumber: true}.Froze(),
				lock: &sync.Mutex{},
				log:  applog.GetLogger().Sugar(),
				mdl:  model.Get(),
				srv:  srv,
				tick: time.NewTicker(time.Second * 30),
			}

			doc, err := oas.LoadFile(filepath.Join(cnf.DataDir, cnf.Schema))
			if err != nil {
				return err
			}

			handlers := oas.OperationHandlers{
				"streams": http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					obj.streams(w, r)
				}),
			}

			crs := cors.New(cors.Options{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Content-Type"},
				AllowCredentials: true,
				MaxAge:           300,
			})

			baseRouter := chi.NewRouter()
			baseRouter.Use(crs.Handler)
			baseRouter.Use(middleware.SetHeader("Content-Type", "application/json"))

			dir := http.Dir(filepath.Join(cnf.DataDir, cnf.HTTPDir))
			baseRouter.Mount("/*", http.StripPrefix("/", http.FileServer(dir)))

			baseRouter.Mount("/metrics", promhttp.Handler())

			errReqHandler := obj.middlewareRequestErrorHandler()
			queryValidator := oas.QueryValidator(errReqHandler)
			bodyValidator := oas.BodyValidator(errReqHandler)

			errRespHandler := obj.middlewareResponseErrorHandler()
			respValidator := oas.ResponseBodyValidator(errRespHandler)

			router, err := oas.NewRouter(
				doc,
				handlers,
				oas.Base(oas.ChiAdapter(baseRouter)),
				oas.Use(queryValidator),
				oas.Use(bodyValidator),
				oas.Use(obj.middlewarePostResponse()),
				oas.Use(respValidator),
			)

			if err != nil {
				return err
			}

			obj.srv.Handler = router

			obj.log.Info("Started API")

			return nil
		},
	)
	event.Stop.AddHandler(
		func() error {
			obj.log.Info("Stop API")

			err := obj.srv.Close()
			if err != nil {
				return err
			}

			obj.tick.Stop()
			obj.lock.Lock()

			obj.log.Info("Stopped API")
			return nil
		},
	)
}

// GetAPI return instance
func GetAPI() *Type {
	return obj
}

// Start API
func (a *Type) Start() {
	a.metrics()

	err := a.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.log.Error(err)
	}
}

func (a *Type) streams(w http.ResponseWriter, r *http.Request) {
	info := a.mdl.ListStreams()

	err := jsoniter.NewEncoder(w).Encode(info)
	if err != nil {
		a.log.Error(err)
	}
}
