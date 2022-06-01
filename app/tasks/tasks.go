package tasks

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

func (h Handler) TaskList(w http.ResponseWriter, r *http.Request) {
	_, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	folderID := helpers.GetParamsInt(r, "folder_id")
	if !validate.ValidateGtInt(0, folderID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	tasks := models.TaskList(h.DB, folderID)
	helpers.RenderSuccessJSON(w, 200, tasks)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	folderID := helpers.GetParamsInt(r, "folder_id")
	startAt := helpers.GetParamsInt(r, "start_at")
	if !validate.ValidateGtInt(0, folderID, startAt) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	content := r.FormValue("content")
	description := r.FormValue("desc")
	if !validate.ValidateStringEmpty(content, description) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	tasks, err := models.CreateTask(h.DB, user.ID, folderID, content, description, startAt)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err)
		return
	}
	helpers.RenderSuccessJSON(w, 200, tasks)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	id := helpers.GetParamsInt(r, "task_id")
	folderID := helpers.GetParamsInt(r, "folder_id")
	if !validate.ValidateGtInt(0, id, folderID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	tasks, err := models.DeleteTask(h.DB, id, folderID)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err)
		return
	}
	helpers.RenderSuccessJSON(w, 200, tasks)
}
