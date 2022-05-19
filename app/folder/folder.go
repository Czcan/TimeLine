package folder

import (
	"net/http"

	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/errcode"
	"github.com/Czcan/TimeLine/utils/validate"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	name := r.FormValue("name")
	if !validate.ValidateStringEmpty(name) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	folder := &models.Folder{Name: name, UserID: user.ID}
	if err := h.DB.Save(&folder).Error; err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	folders := models.GetFolderList(h.DB, user.ID)
	helpers.RenderSuccessJSON(w, 200, folders)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	folders := models.GetFolderList(h.DB, user.ID)
	helpers.RenderSuccessJSON(w, 200, folders)
}

func (h Handler) Deleted(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	folderID := helpers.GetParamsInt(r, "folder_id")
	if !validate.ValidateGtInt(0, folderID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	err = models.DeletedFolderAndNote(h.DB, folderID)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "deleted failed")
		return
	}
	folders := models.GetFolderList(h.DB, user.ID)
	helpers.RenderSuccessJSON(w, 200, folders)
}
