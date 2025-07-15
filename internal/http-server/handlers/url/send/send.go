package send

import (
	resp "infotex/internal/lib/api/response"
	"infotex/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Sender  string `json:"from" validate:"required"`
	Reciver string `json:"to" validate:"required"`
	Amount  int    `json:"amount" validate:"required"`
}

type Response struct {
	resp.Response
}

type TransactionProcesser interface {
	ProcessTransactions(senderAdress, receiverAdress string, amount int) error
}

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

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		// check required params
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		if req.Amount <= 0 || req.Sender == req.Reciver { // TODO: implement in validator
			log.Debug("stupid user oh god")
			render.JSON(w, r, resp.Error("check ur request"))
			return
		}
		err = transactionProcesser.ProcessTransactions(req.Sender, req.Reciver, req.Amount)

		if err != nil { // TODO: validate errors from db
			log.Error("failed to process transaction", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to process transaction"))
			return
		}

		log.Debug("successful transaction")

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
