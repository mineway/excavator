package routes

import (
	"context"
	"github.com/mineway/worker/internal/pkg/manager"
	"github.com/mineway/worker/internal/pkg/response"
	"net/http"
)

/*
	@Route("GET", "/ping")
	@Description("allows to check if you can ping the API")
*/
func (h Handler) Ping(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	return response.SuccessText(w, http.StatusOK, "pong")
}

/*
	@Route("GET", "/error")
	@Description("allows to check if error work")
*/
func (h Handler) Error(ctx context.Context, m *manager.Manager, w http.ResponseWriter) bool {
	return response.Error(w, http.StatusInternalServerError, "Error")
}