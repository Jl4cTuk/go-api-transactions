package getlast

import (
	"infotex/internal/domain/model"
	resp "infotex/internal/api/response"
	"infotex/internal/logger/sl"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Transactions []model.Transaction `json:"transactions"`
}

type LastTransactionsGetter interface {
	GetLastTransactions(count int) ([]model.Transaction, error)
}

// getlast returns last N transactions between wallets
func New(log *slog.Logger, lastTransactionsGetter LastTransactionsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getlast.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		res := r.URL.Query().Get("count")
		count, err := strconv.Atoi(res)
		if err != nil {
			log.Info("missing count param")
			count = 1
		}

		log.Debug("count", slog.Int("count", count))

		transactions, err := lastTransactionsGetter.GetLastTransactions(count)

		if err != nil {
			log.Error("failed to get last transactions", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		render.JSON(w, r, Response{
			Transactions: transactions,
		})
	}
}
