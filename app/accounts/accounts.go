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

func (h Handler) AcccountPersonal(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	accounts := []models.Account{}
	h.DB.Where("user_id = ?", user.ID).Find(&accounts)
	for i := 0; i < len(accounts); i++ {
		accounts[i].ConCatImages()
	}
	helpers.RenderSuccessJSON(w, 200, accounts)
}

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	content := r.FormValue("content")
	title := r.FormValue("title")
	if !validate.ValidateStringEmpty(content) || !validate.ValidateStringEmpty(title) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	r.ParseMultipartForm(20 << 32)
	err = models.CreatedAccount(h.DB, user.ID, content, title, r.MultipartForm.File["image"])
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, "success created")
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

func (h Handler) AccoutDelted(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	id := helpers.GetParamsInt(r, "id")
	if !validate.ValidateGtInt(0, id) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	h.DB.Where("id = ?", id).Delete(&models.Account{})
	accounts := []models.Account{}
	h.DB.Where("user_id = ?", user.ID).Find(&accounts)
	helpers.RenderSuccessJSON(w, 200, accounts)
}
