package getbalance

import (
	"errors"
	resp "infotex/internal/lib/api/response"
	"infotex/internal/lib/logger/sl"
	"infotex/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Balance float64 `json:"balance,omitempty"`
}

type BalanceGetter interface {
	GetWalletBalance(address string) (float64, error)
}

// getbalance returns balanse of the given wallet address
func New(log *slog.Logger, balanceGetter BalanceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getbalance.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		address := chi.URLParam(r, "address") // *bind to chi
		if address == "" {
			log.Info("missing address")

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Debug("got address", slog.String("address", address))

		balance, err := balanceGetter.GetWalletBalance(address)
		if errors.Is(err, storage.ErrWalletNotFound) {
			log.Info("wallet not found", slog.String("address", address))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		if err != nil {
			log.Error("failed to process address", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Debug("address found", slog.Float64("balance", balance))

		render.JSON(w, r, Response{
			Balance: balance,
		})
	}
}
