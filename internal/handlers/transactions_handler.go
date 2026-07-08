package handlers

import (
	"net/http"
	"strconv"

	"github.com/The-Ogulgozel/Banking-system/internal/appresult"
	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"
	"github.com/The-Ogulgozel/Banking-system/internal/handlers/request"
	"github.com/The-Ogulgozel/Banking-system/internal/handlers/response"
	"github.com/The-Ogulgozel/Banking-system/internal/models"
	"github.com/The-Ogulgozel/Banking-system/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	uc usecase.TransactionUseCaseInterface
}

func NewTransactionHandler(tcase usecase.TransactionUseCaseInterface) *TransactionHandler {
	return &TransactionHandler{
		uc: tcase,
	}
}

func (h *TransactionHandler) Deposit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	var req request.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	t, err := h.uc.Deposit(c.Request.Context(), id, req.Amount)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusCreated, appresult.NewAppSuccess("Successfully created", response.ToTransactionResponse(t)))

}

func (h *TransactionHandler) Withdraw(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	var req request.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	t, err := h.uc.Withdraw(c.Request.Context(), id, req.Amount)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusCreated, appresult.NewAppSuccess("Successfully created", response.ToTransactionResponse(t)))

}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	var req request.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	t, err := h.uc.Transfer(c.Request.Context(), id, req.ToAccountID, req.Amount)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusCreated, appresult.NewAppSuccess("Successfully created", response.ToTransactionResponse(t)))
}

func (h *TransactionHandler) ListTransactionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	var req request.ListTransactions
	if err := c.ShouldBindQuery(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	startDate, err := mustParseDate(req.StartDate)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	endDate, err := mustParseDate(req.EndDate)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	var typePtr *models.TransactionType
	if req.Type != nil {
		t := models.TransactionType(*req.Type)
		typePtr = &t
	}
	dto := models.ListTransactions{
		AccountID: id,
		Page:      req.Page,
		Limit:     req.Limit,
		Type:      typePtr,
		StartDate: startDate,
		EndDate:   endDate,
	}

	t, total, err := h.uc.ListByAccountID(c.Request.Context(), &dto)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, appresult.NewAppSuccess("Successfully listed", response.ToListTransactionsResponse(t, total, *dto.Page, *dto.Limit)))
}
