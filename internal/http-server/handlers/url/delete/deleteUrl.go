package deleteUrl

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"short-url-api/internal/lib/api/response"
	"short-url-api/internal/lib/logger/sl"
)

type UrlDelete interface {
	DeleteUrl(alias string) (bool, error)
}

func New(log *slog.Logger, urlDelete UrlDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// get alias
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias empty")

			render.JSON(w, r,
				response.Error("invalid alias"),
			)

			return
		}

		// db request
		ok, err := urlDelete.DeleteUrl(alias)
		if err != nil {
			log.Error("cant delete url", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r,
				response.Error("internal error"),
			)

			return
		}

		if !ok {
			log.Error("url not found")

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r,
				response.Error("url not found or already delete"),
			)

			return
		}

		// response
		responseOk(w, r)
		return
	}
}

func responseOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r,
		response.OK(),
	)
}
