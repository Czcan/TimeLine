package likers

import (
	"net/http"

	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h Handler) Liker(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamsInt(r, "id")
	if id <= 0 {
		helpers.RenderFailureJSON(w, 400, "invalid params")
		return
	}
	liker := helpers.GetParamsBool(r, "liker")
	account := &models.Account{}
	h.DB.Where("id = ?", id).First(account)
	if liker {
		account.Likers += 1
	} else {
		account.Likers -= 1
	}
	h.DB.Save(&account)
	helpers.RenderSuccessJSON(w, 200, account.Likers)
}

func (h Handler) Follwer(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	id := helpers.GetParamsInt(r, "id")
	if id <= 0 {
		helpers.RenderFailureJSON(w, 400, "invalid params")
		return
	}
	follwer := helpers.GetParamsBool(r, "follwer")
	follwerCount, err := models.UpdateFollwerAndSyncCollection(h.DB, id, user.ID, follwer)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, follwerCount)
}
