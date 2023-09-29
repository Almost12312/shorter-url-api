package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"short-url-api/internal/lib/api/response"
	"short-url-api/internal/lib/logger/sl"
	"short-url-api/internal/storage"
)

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//get url query
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias empty")

			render.JSON(w, r,
				response.Error("invalid alias"),
			)

			return
		}

		//db request
		url, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			w.WriteHeader(404)
			render.Status(r, 404)
			render.JSON(w, r,
				response.Error("not found"),
			)

			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r,
				response.Error("internal error"),
			)
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
