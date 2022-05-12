package notes

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

func (h Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	folderID := helpers.GetParamsInt(r, "folder_id")
	content := r.FormValue("content")
	if folderID == 0 || content == "" {
		helpers.RenderFailureJSON(w, 400, "invalid params")
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

func (h Handler) NoteList(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	folderID := helpers.GetParamsInt(r, "folder_id")
	if folderID == 0 {
		helpers.RenderFailureJSON(w, 400, "invalid params")
		return
	}
	notes := models.GetNoteList(h.DB, user.ID, folderID)
	helpers.RenderSuccessJSON(w, 200, notes)
}

func (h Handler) FinishNote(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	noteID := helpers.GetParamsInt(r, "note_id")
	if noteID == 0 {
		helpers.RenderFailureJSON(w, 400, "invalid params")
		return
	}
	status := helpers.GetParamsBool(r, "status")
	h.DB.Model(&models.Note{}).Where("id = ? AND user_id = ?", noteID, user.ID).Update("status", status)
	helpers.RenderSuccessJSON(w, 200, "updated successed")
}

func (h Handler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	name := r.FormValue("name")
	if name == "" {
		helpers.RenderFailureJSON(w, 400, "invalid params")
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

func (h Handler) FolderList(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	folders := models.GetFolderList(h.DB, user.ID)
	helpers.RenderSuccessJSON(w, 200, folders)
}
