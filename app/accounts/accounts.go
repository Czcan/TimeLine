package accounts

import (
	"net/http"
	"strconv"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func (h Handler) AccountList(w http.ResponseWriter, r *http.Request) {
	accounts := []models.Account{}
	h.DB.Order("likers desc, follwers desc").Find(&accounts)
	for i := 0; i < len(accounts); i++ {
		accounts[i].ConCatImages()
	}
	helpers.RenderSuccessJSON(w, 200, accounts)
}

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
}

func (h Handler) AccoutDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	accountID, _ := strconv.Atoi(id)
	if accountID <= 0 || id == "" {
		helpers.RenderFailureJSON(w, 400, "invalid param")
		return
	}
	var (
		account  = models.Account{}
		comments = []entries.Comment{}
	)
	if err := h.DB.Where("id = ?", accountID).First(&account).Error; err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	account.ConCatImages()
	selectSQL := `
		SELECT comments.content, comments.created_at AS date, users.nick_name, CONCAT('/images/', users.id, '.jpg') AS avatar_url
		FROM comments
		LEFT JOIN users ON comments.user_id = users.id
		WHERE comments.account_id = ?
	`
	h.DB.Raw(selectSQL, accountID).Scan(&comments)
	helpers.RenderSuccessJSON(w, 200, entries.AccountDetail{
		Account:  account,
		Comments: comments,
	})
}
