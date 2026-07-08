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

type AccountsHandler struct {
	usecase usecase.AccountsUsecaseInterface
}

func NewAccountsHandler(usecase usecase.AccountsUsecaseInterface) *AccountsHandler {
	return &AccountsHandler{
		usecase: usecase,
	}
}

func (h *AccountsHandler) Create(c *gin.Context) {
	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	account := &models.Account{
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	err := h.usecase.Create(c, account)
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, appresult.NewAppSuccess("Successfully created", response.ToAccounResponse(account)))
}

func (h *AccountsHandler) GetByID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	account, err := h.usecase.GetByID(c, int64(id))
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, appresult.NewAppSuccess("Successfully", response.ToAccounResponse(account)))
}

func (h *AccountsHandler) Delete(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		appresult.RespondAppError(c, appError.ErrNotAccountId)
		return
	}

	err = h.usecase.Delete(c, int64(id))
	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, appresult.NewAppSuccess("Successfully deleted", nil))
}

func (h *AccountsHandler) ListAccounts(c *gin.Context) {
	var req request.ListAccountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		appresult.RespondAppError(c, appError.ErrInvalidRequest)
		return
	}

	if err := req.Validate(); err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	params := &models.ListAccounts{
		Currency: req.Currency,
		Page:     req.Page,
		Limit:    req.Limit,
	}
	accounts, total, err := h.usecase.GetAll(c, params)

	if err != nil {
		appresult.RespondAppError(c, err)
		return
	}
	c.JSON(http.StatusOK, appresult.NewAppSuccess("Successfully got all accounts", response.ToListAccountsResponse(accounts, total, *params.Page, *params.Limit)))
}
