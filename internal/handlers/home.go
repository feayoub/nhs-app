package handlers

import (
	"net/http"

	"github.com/feayoub/nhs-app/internal/templates"
)

type homeHandler struct {
}

func NewHomeHandler() http.Handler {
	return &homeHandler{}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Home()
	templates.Layout(c, "App").Render(r.Context(), w)
}
