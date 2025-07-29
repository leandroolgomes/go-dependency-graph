package examples

import (
	"fmt"
	"net/http"

	"github.com/leandroolgomes/golang-dependency-graph/component"
)

// Config provides configuration values
type Config struct{
	Port int
}

func (c *Config) Start(ctx component.Context) (component.Lifecycle, error) {

	c.Port = 3000
	return c, nil
}

func (c *Config) Stop(ctx component.Context) error {

	return nil
}

// ConfigMock mock implementation for testing
type ConfigMock struct{
	Port int
}

func (c *ConfigMock) Start(ctx component.Context) (component.Lifecycle, error) {

	c.Port = 4000
	return c, nil
}

func (c *ConfigMock) Stop(ctx component.Context) error {

	return nil
}

// AppRoutes defines HTTP routes
type AppRoutes struct{
	SetupRoutes func(mux *http.ServeMux)
}

func (a *AppRoutes) Start(ctx component.Context) (component.Lifecycle, error) {

	a.SetupRoutes = func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
		fmt.Println("App routes configured!")
	}
	
	return a, nil
}

func (a *AppRoutes) Stop(ctx component.Context) error {

	return nil
}

// HttpServer sets up and runs an HTTP server
type HttpServer struct{
	Server *http.Server
}

func (h *HttpServer) Start(ctx component.Context) (component.Lifecycle, error) {

	configObj, ok := ctx["config"]
	if !ok {
		return nil, fmt.Errorf("config dependency not found")
	}
	
	appRoutesObj, ok := ctx["app_routes"]
	if !ok {
		return nil, fmt.Errorf("app_routes dependency not found")
	}
	

	config, ok := configObj.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type")
	}
	
	appRoutes, ok := appRoutesObj.(*AppRoutes)
	if !ok {
		return nil, fmt.Errorf("invalid app_routes type")
	}
	

	port := config.Port
	

	mux := http.NewServeMux()
	appRoutes.SetupRoutes(mux)
	
	addr := fmt.Sprintf(":%d", port)
	h.Server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	

	go func() {
		fmt.Printf("Example app listening on port %d\n", port)
		if err := h.Server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
	
	return h, nil
}

func (h *HttpServer) Stop(ctx component.Context) error {

	if h.Server != nil {
		fmt.Println("HTTP server closing")
		return h.Server.Close()
	}
	return nil
}
