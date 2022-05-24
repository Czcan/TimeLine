package upload

import (
	"net/http"
	"strconv"

	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/config"
	"github.com/Czcan/TimeLine/utils/file"
	"gorm.io/gorm"
)

type Handler struct {
	DB         *gorm.DB
	UploadPath string
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db, UploadPath: config.MustGetAppConfig().AvatarPath}
}

func (h Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	r.ParseMultipartForm(32 << 20)
	_, f, err := r.FormFile("image")
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	imageUrl, err := file.SaveUploadFile(f, h.UploadPath, strconv.Itoa(user.ID))
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, imageUrl)
}
