package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/adlio/trello"
)

var (
	boardIDSByResponsible = map[string]string{
		"ANA":      "1",
		"FERNANDA": "2",
		"EVERSON":  "3",
		"MARCELA":  "4",
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
	Owner       string
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
		boardID := boardIDSByResponsible[card.Owner]
		board, err := u.trelloClient.GetBoard(boardID, trello.Defaults())
		if err != nil {
			return err
		}

		lists, err := board.GetLists(trello.Defaults())
		if err != nil {
			return err
		}

		for _, list := range lists {
			if strings.ToUpper(list.Name) == "REMANEJAR" {
				cardToCreate := mapToTrelloCard(card)
				list.AddCard(cardToCreate, trello.Defaults())
				cardToCreate.AddComment(card.Comment)
				break
			}
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
