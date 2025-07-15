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
	resp.Response
	Balance int64 `json:"balance,omitempty"`
}

type BalanceGetter interface {
	GetWalletBalance(address string) (int64, error)
}

// GET /api/wallet/{address}/balance
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

			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Info("got address", slog.String("address", address))

		balance, err := balanceGetter.GetWalletBalance(address)
		if errors.Is(err, storage.ErrWalletNotFound) { // TODO: properly validate this
			log.Info("wallet not found", slog.String("address", address))

			render.JSON(w, r, resp.Error("not found"))
			return
		}

		if err != nil {
			log.Error("failed to process address", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("address found", slog.Int64("balance", balance))

		render.JSON(w, r, Response{
			Balance: balance,
		})
	}
}
