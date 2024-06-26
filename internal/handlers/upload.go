package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"code.sajari.com/docconv"
	"github.com/feayoub/nhs-app/internal/templates"
	"github.com/feayoub/nhs-app/internal/usecase"
)

type uploadHandler struct {
	createTrelloCardUseCase usecase.UseCase
}

func NewUploadHandler(createTrelloCardUseCase usecase.UseCase) http.Handler {
	return &uploadHandler{
		createTrelloCardUseCase: createTrelloCardUseCase,
	}
}

func (h *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read form file
	_, header, err := r.FormFile("file")
	if err != nil {
		c := templates.UploadError(err.Error())
		c.Render(r.Context(), w)
		return
	}

	// Source
	src, err := header.Open()
	if err != nil {
		c := templates.UploadError(err.Error())
		c.Render(r.Context(), w)
		return
	}
	defer src.Close()

	if !strings.HasSuffix(header.Filename, ".doc") && !strings.HasSuffix(header.Filename, ".docx") {
		c := templates.UploadError("Formato de arquivo inválido")
		c.Render(r.Context(), w)
		return
	}

	// Destination
	dst, err := os.Create(header.Filename)
	if err != nil {
		c := templates.UploadError(err.Error())
		c.Render(r.Context(), w)
		return
	}
	defer dst.Close()
	defer os.Remove(header.Filename)

	// Copy
	io.Copy(dst, src)

	pubsString := extractPublicationsString(header.Filename)

	useCaseInput := convertToUseCaseInput(pubsString)
	_ = useCaseInput

	if err = h.createTrelloCardUseCase.Execute(useCaseInput); err != nil {
		c := templates.UploadError(err.Error())
		c.Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	c := templates.UploadSuccess()
	c.Render(r.Context(), w)
}

func extractPublicationsString(fileName string) []string {
	res, err := docconv.ConvertPath(fileName)
	if err != nil {
		fmt.Println(err)
	}
	splits := strings.Split(res.Body, "\n")
	chunk := ""
	var startChunk bool
	var publications []string
	r := regexp.MustCompile(`^\d+\.`)
	for i, s := range splits {
		if strings.TrimSpace(s) == "" || strings.TrimSpace(s) == "COMENTÁRIOS" {
			if i == len(splits)-1 && chunk != "" {
				publications = append(publications, chunk)
				chunk = ""
			}
			continue
		}
		if r.Match([]byte(s)) {
			if len(publications) == 0 {
				startChunk = true
			}
			if chunk != "" {
				publications = append(publications, chunk)
				chunk = ""
			}
		}
		if startChunk {
			chunk += s + "\n"
		}
	}
	return publications
}

func convertToUseCaseInput(pubsString []string) usecase.CreateTrelloCardInput {
	input := usecase.CreateTrelloCardInput{}
	for _, pub := range pubsString {
		pubLine := strings.Split(pub, "\n")
		var comment, board, title, description string
		for i, line := range pubLine {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if i <= 5 {
				comment += line + "\n"
				continue
			} else if i == 6 {
				title = line
				if strings.ToUpper(title) == "NADA" {
					comment, board, title, description = "", "", "", ""
					break
				}
			}
			if b, ok := usecase.GetBoardByResponsible(strings.ToUpper(line)); ok {
				if board != "" {
					card := createCard(comment, board, title, description)
					input.Cards = append(input.Cards, card)
					description = ""
				}
				board = b
				continue
			}
			description += line + "\n"
		}
		cardInput := createCard(comment, board, title, description)
		input.Cards = append(input.Cards, cardInput)
	}
	return input
}

func createCard(comment, board, title, description string) usecase.Card {
	return usecase.Card{
		Comment:     comment,
		Board:       board,
		Title:       title,
		Description: strings.TrimSuffix(description, "\n"),
	}
}
