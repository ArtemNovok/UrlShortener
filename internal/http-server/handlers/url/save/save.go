package save

import (
	"errors"
	"log/slog"
	"net/http"
	myvalidator "url-shortener/internal/lib/api/validator"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty`
	Alias  string `json:"alias,omitempty"`
}

const (
	aliasLength = 5
)

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With("op", op, slog.String("request_id", middleware.GetReqID(r.Context())))
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "Failed to decode request",
			})
			return
		}
		log.Info("request body decoded", slog.Any("req", req))
		if err := validator.New().Struct(req); err != nil {
			validErrors := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.String("error", err.Error()))
			render.JSON(w, r, myvalidator.ValidationError(validErrors))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomSTR(aliasLength)
		}
		err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, Response{
					Status: "Error",
					Error:  "Url already exists",
				})
				return
			}
			log.Error("failed to save url", slog.String("error", err.Error()))
			render.JSON(w, r, Response{
				Status: "Error",
				Error:  "Failed to save url",
			})
			return
		}
		log.Info("url saved", slog.String("url", req.URL), slog.String("alias", alias))
		render.JSON(w, r, Response{
			Status: "OK",
			Alias:  alias,
		})
	}
}
