package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty`
	Alias  string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.New"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "invalid request",
			})
			return
		}
		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				render.JSON(w, r, Response{
					Status: "Error",
					Error:  "not found",
				})
				return
			}
			log.Error("failed to get url", slog.String("error", err.Error()))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}
		log.Info("got url", slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
