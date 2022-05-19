package notes

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
	folderID := helpers.GetParamsInt(r, "folder_id")
	content := r.FormValue("content")
	if !validate.ValidateGtInt(0, folderID) || !validate.ValidateStringEmpty(content) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	note := &models.Note{FolderID: folderID, Content: content, UserID: user.ID}
	if err := h.DB.Save(&note).Error; err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	notes := models.GetNoteList(h.DB, folderID, user.ID)
	helpers.RenderSuccessJSON(w, 200, notes)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
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
	notes := models.GetNoteList(h.DB, user.ID, folderID)
	helpers.RenderSuccessJSON(w, 200, notes)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	noteID := helpers.GetParamsInt(r, "note_id")
	if !validate.ValidateGtInt(0, noteID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	status := helpers.GetParamsBool(r, "status")
	h.DB.Model(&models.Note{}).Where("id = ? AND user_id = ?", noteID, user.ID).Update("status", status)
	helpers.RenderSuccessJSON(w, 200, "updated successed")
}

func (h Handler) Deleted(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	noteID := helpers.GetParamsInt(r, "note_id")
	folderID := helpers.GetParamsInt(r, "folder_id")
	if !validate.ValidateGtInt(0, noteID, folderID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	h.DB.Where("id = ?", noteID).Delete(&models.Note{})
	notes := models.GetNoteList(h.DB, user.ID, folderID)
	helpers.RenderSuccessJSON(w, 200, notes)
}
