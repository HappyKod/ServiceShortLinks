package server

import (
	"net/http"
)

//Server запуск сервера
func Server(server *http.Server) error {
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
