package likers

import (
	"net/http"

	"github.com/Czcan/TimeLine/internal/api/helpers"
	"github.com/Czcan/TimeLine/internal/models"
	"github.com/Czcan/TimeLine/pkg/errcode"
	"github.com/Czcan/TimeLine/pkg/validate"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h Handler) Liker(w http.ResponseWriter, r *http.Request) {
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
	liker := helpers.GetParamsBool(r, "liker")
	likers, err := models.UpdateLiker(h.DB, user.ID, id, liker)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, likers)
}

func (h Handler) Follwer(w http.ResponseWriter, r *http.Request) {
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
	follwer := helpers.GetParamsBool(r, "follwer")
	follwerCount, err := models.UpdateFollwerAndSyncCollection(h.DB, id, user.ID, follwer)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, follwerCount)
}
