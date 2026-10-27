package handlers

import (
	"encoding/csv"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CollegeBTK-Projects/my-node-go-app/models"
	"github.com/CollegeBTK-Projects/my-node-go-app/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (u *UserHandler) InsertUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	id, err := u.repo.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: ": err.Error()})
		return
	}
	user.ID = id
	c.Header("X-Group", "ББМО-01-23")
	c.JSON(http.StatusCreated, user)

	go SendEmail(user.Email)
}

func (u *UserHandler) FindById(c *gin.Context) {
	id := c.Param("id")

	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	um, err := u.repo.FindById(c.Request.Context(), idHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-Group", "ББМО-01-23")
	c.JSON(http.StatusOK, um)
}

func (u *UserHandler) FindByAll(c *gin.Context) {
	um, err := u.repo.FindByAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("X-Group", "ББМО-01-23")
	c.JSON(http.StatusOK, um)
}

func (u *UserHandler) UpdateById(c *gin.Context) {
	id := c.Param("id")

	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}
	upd, err := u.repo.UpdateById(c.Request.Context(), idHex, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("X-Group", "ББМО-01-23")
	c.JSON(http.StatusOK, gin.H{"count update:": upd})
}

func (u *UserHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")

	idHex, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.repo.DeleteById(c.Request.Context(), idHex)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": err.Error()})
		return
	}

	c.Header("X-Group", "ББМО-01-23")
	c.JSON(http.StatusAccepted, gin.H{"count deleted": 1})
}

func SendEmail(u string) {
	time.Sleep(2 * time.Second)
	log.Printf("Письмо отправлено на %s", u)
}

func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email обязателен"})
		return
	}

	u, err := h.repo.FindByEmail(c.Request.Context(), email)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func bindQuery(c *gin.Context) (models.UserQuery, bool) {
	var q models.UserQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := map[string]string{}
			for _, fe := range ve {
				switch fe.Tag() {
				case "oneof":
					errs[fe.Field()] = "допустимо: " + fe.Param()
				case "group":
					errs[fe.Field()] = "формат: ББМО-XX-XX"
				default:
					errs[fe.Field()] = "неверное значение"
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			return q, false
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный тип одного из query-параметров"})
		return q, false
	}
	if err := q.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return q, false
	}
	return q, true
}

func (h *UserHandler) List(c *gin.Context) {
	q, ok := bindQuery(c)
	if !ok {
		return
	}
	h.respondList(c, q)
}

func (h *UserHandler) respondList(c *gin.Context, q models.UserQuery) {
	users, total, err := h.repo.List(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	page, limit := q.PageLimit()
	c.JSON(http.StatusOK, models.UserList{
		Data: users,
		Meta: models.Meta{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *UserHandler) Search(c *gin.Context) {
	text := strings.TrimSpace(c.Query("q"))
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр q обязателен"})
		return
	}
	q, ok := bindQuery(c)
	if !ok {
		return
	}
	q.Search = text
	h.respondList(c, q)
}

func (h *UserHandler) Stats(c *gin.Context) {
	st, err := h.repo.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *UserHandler) Export(c *gin.Context) {
	q, ok := bindQuery(c)
	if !ok {
		return
	}
	users, err := h.repo.ListAll(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="users.csv"`)
	c.Writer.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"id", "name", "email", "group_name", "age", "course", "created_at"})
	for _, u := range users {
		_ = w.Write([]string{
			u.ID.Hex(),
			u.Name,
			u.Email,
			u.GroupName,
			strconv.Itoa(u.Age),
			strconv.Itoa(u.Course),
			u.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Flush()
}
