package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"short-url-api/internal/lib/api/response"
	"short-url-api/internal/lib/logger/sl"
	"short-url-api/internal/lib/random"
	"short-url-api/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: to env/db
const aliasLength = 6

type URLSaver interface {
	SaveUrl(urlToSave, alias string) (id int64, err error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed decode json"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidatorErrors(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)
		if errors.Is(err, storage.ErrURLExist) {
			log.Info("cant save, url already exist!", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exist!"))
			return
		}

		log.Info("url added", slog.Int64("id", id), slog.String("request_id", middleware.GetReqID(r.Context())))

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r,
		Response{
			Response: response.OK(),
			Alias:    alias,
		})
}
