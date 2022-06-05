package comments

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

func (h Handler) Comment(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	content := r.FormValue("content")
	accountID := helpers.GetParamsInt(r, "id")
	if !validate.ValidateStringEmpty(content) || !validate.ValidateGtInt(0, accountID) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	err = models.SaveComment(h.DB, user.ID, accountID, content)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "comment failed")
		return
	}
	comments := models.FindComments(h.DB, accountID)
	helpers.RenderSuccessJSON(w, 200, comments)
}
