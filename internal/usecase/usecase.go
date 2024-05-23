package usecase

type UseCase interface {
	Execute(input any) error
}
