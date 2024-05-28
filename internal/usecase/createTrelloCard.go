package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/adlio/trello"
)

var (
	boardIDSByResponsible = map[string]string{
		"ANA":      "PE2najNK",
		"FERNANDA": "ZJFk3iEG",
		"EVERSON":  "oil3x0mt",
		"MARCELA":  "sLPtCxtG",
		"SARA":     "LyDxxGqi",
		// "FELIPE":   "VwdF3dKW",
	}
)

type createTrelloCardUseCase struct {
	trelloClient *trello.Client
}

type CreateTrelloCardInput struct {
	Cards []Card `json:""`
}

type Card struct {
	Comment     string
	Description string
	Title       string
	Board       string
}

func NewCreateTrelloCardUseCase(trelloClient *trello.Client) UseCase {
	return &createTrelloCardUseCase{
		trelloClient: trelloClient,
	}
}

func (u *createTrelloCardUseCase) Execute(inp any) error {
	input, ok := inp.(CreateTrelloCardInput)
	if !ok {
		return errors.New("wrong input format")
	}

	for _, card := range input.Cards {
		board, err := u.trelloClient.GetBoard(card.Board, trello.Defaults())
		if err != nil {
			return err
		}

		lists, err := board.GetLists(trello.Defaults())
		if err != nil {
			return err
		}

		var cardErrorsList []error
		for _, list := range lists {
			if strings.ToUpper(list.Name) == "REMANEJAR" {
				cardToCreate := mapToTrelloCard(card)
				err := list.AddCard(cardToCreate, trello.Defaults())
				if err != nil {
					cardErrorsList = append(cardErrorsList, err)
				}
				cardToCreate.AddComment(card.Comment)
				break
			}
		}
		if len(cardErrorsList) > 0 {
			return errors.Join(cardErrorsList...)
		}
	}

	return nil
}

func mapToTrelloCard(inputCard Card) *trello.Card {
	now := time.Now()
	y, m, d := now.Date()
	due := time.Date(y, m, d+1, 0, 0, 0, -1, now.Location())
	return &trello.Card{
		Name: inputCard.Title,
		Desc: inputCard.Description,
		Due:  &due,
	}
}

func GetBoardByResponsible(owner string) (string, bool) {
	board, ok := boardIDSByResponsible[owner]
	return board, ok
}
