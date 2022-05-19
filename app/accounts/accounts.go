package accounts

import (
	"net/http"
	"strconv"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/errcode"
	"github.com/Czcan/TimeLine/utils/validate"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h Handler) AccountList(w http.ResponseWriter, r *http.Request) {
	accounts := []models.Account{}
	h.DB.Order("likers desc, follwers desc").Find(&accounts)
	for i := 0; i < len(accounts); i++ {
		accounts[i].ConCatImages()
	}
	helpers.RenderSuccessJSON(w, 200, accounts)
}

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) AccoutDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	accountID, _ := strconv.Atoi(id)
	if !validate.ValidateGtInt(0, accountID) || !validate.ValidateStringEmpty(id) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	account, comments, err := models.FindAccountDetail(h.DB, accountID)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, entries.AccountDetail{
		Account:  account,
		Comments: comments,
	})
}
