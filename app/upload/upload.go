package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/utils"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB         *gorm.DB
	UploadPath string
}

func New(db *gorm.DB, path string) Handler {
	return Handler{DB: db, UploadPath: path}
}

func (h Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	defer file.Close()
	if ok := utils.Exists(h.UploadPath); !ok {
		os.MkdirAll(h.UploadPath, os.ModePerm)
	}
	savePath := fmt.Sprintf("%s/%d.jpg", h.UploadPath, 1)
	fmt.Println(filepath.Abs(savePath))
	f, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, fmt.Sprintf("/images/%d.jpg", user.ID))
}
