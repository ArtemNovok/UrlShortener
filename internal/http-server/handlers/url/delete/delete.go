package delete

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty`
	Alias  string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log = log.With(slog.String("op", op))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("request_id", middleware.GetReqID(r.Context())))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}
		if req.Alias == "" {
			log.Info("empty alias", slog.String("request_id", middleware.GetReqID(r.Context())))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "bad request",
			})
			return
		}
		err = urlDeleter.DeleteURL(req.Alias)
		if err != nil {
			log.Error("failed to delete url", slog.String("error", err.Error()))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}
		render.JSON(w, r, Response{
			Status: "OK",
			Alias:  req.Alias,
		})
	}
}
