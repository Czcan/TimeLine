package folder

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

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	kind := helpers.GetParamsInt(r, "kind")
	if kind < 0 || kind > 1 {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	folders := models.GetFolderList(h.DB, user.ID, kind)
	helpers.RenderSuccessJSON(w, 200, folders)
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
	kind := helpers.GetParamsInt(r, "kind")
	if kind < 0 || kind > 1 {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	folder, err := models.CreateFolder(h.DB, user.ID, name, kind)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, folder)
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
	kind := helpers.GetParamsInt(r, "kind")
	if kind < 0 || kind > 1 {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	folders, err := models.DeletedFolder(h.DB, folderID, user.ID, kind)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "deleted failed")
		return
	}
	helpers.RenderSuccessJSON(w, 200, folders)
}