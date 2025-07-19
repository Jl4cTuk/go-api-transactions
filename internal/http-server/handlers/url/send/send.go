package send

import (
	"errors"
	resp "infotex/internal/api/response"
	"infotex/internal/logger/sl"
	"infotex/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Sender  string  `json:"from" validate:"required"`
	Reciver string  `json:"to" validate:"required,nefield=Sender"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
}

type Response struct {
	resp.Response
}

type TransactionProcesser interface {
	ProcessTransactions(senderAdress, receiverAdress string, amount float64) error
}

// send process transaction between two adresses and given amount
func New(log *slog.Logger, transactionProcesser TransactionProcesser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.send.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		// check required params
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		err = transactionProcesser.ProcessTransactions(req.Sender, req.Reciver, req.Amount)

		if errors.Is(err, storage.ErrInvalidWallet) {
			log.Info("wrong wallet adress", slog.String("wallet", req.Sender))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("wrong wallet adress"))
			return
		}
		if errors.Is(err, storage.ErrInsufficientFunds) {
			log.Info("insufficient funds", slog.String("wallet", req.Sender))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("insufficient funds"))
			return
		}
		if err != nil {
			log.Error("failed to process transaction", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to process transaction"))
			return
		}

		render.Status(r, http.StatusNoContent)
	}
}
