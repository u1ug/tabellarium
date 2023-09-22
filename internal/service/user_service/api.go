package user_service

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"tabellarium/internal/infrastructure/storage"
	"tabellarium/pkg/logging"
)

// UserService is a service for handling user devices registration and management
type UserService struct {
	settings Settings
	router   *httprouter.Router
	store    *storage.DeviceStorage
	logger   *logging.Logger
}

// Settings for the service.
type Settings struct {
	// Service configuration.
	Address string
	Port    int
	// Redis configuration.
	RedisAddress  string
	RedisPort     uint
	RedisDB       int
	RedisPassword string
	RedisTTL      int
}

// initRoutes registers http server routes for UserService.
func (u *UserService) initRoutes() {
	u.router.POST("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		res, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			u.logger.Errorln(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload, err := ParseRegisterDevicePayload(res)
		if err != nil {
			u.logger.Errorln(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !payload.IsValid() {
			writer.WriteHeader(http.StatusBadRequest)
			_, err := fmt.Fprintf(writer, "invalid payload received")
			if err != nil {
				u.logger.Errorln(err)
			}
			return
		}
		err = u.store.RegisterDevice(payload.UserID, payload.Token)
		if err != nil {
			u.logger.Errorln(err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
		u.logger.Debugf("payload: %s, %s", payload.Token, payload.UserID)
		writer.WriteHeader(http.StatusOK)
	})
}

// NewUserService inits UserService with
func NewUserService(s Settings) *UserService {
	store := storage.NewDeviceStorage(storage.Settings{
		Address:  s.RedisAddress,
		Port:     s.RedisPort,
		DB:       s.RedisDB,
		Password: s.RedisPassword,
		TTL:      s.RedisTTL,
	})
	u := &UserService{
		settings: s,
		store:    store,
		router:   httprouter.New(),
		logger:   logging.GetLogger(),
	}
	u.initRoutes()
	return u
}

// Listen starts UserService main loop.
func (u *UserService) Listen() error {
	addr := fmt.Sprintf("%s:%d", u.settings.Address, u.settings.Port)
	return http.ListenAndServe(addr, u.router)
}
