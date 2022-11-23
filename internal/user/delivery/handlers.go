package delivery

import "github.com/sornick01/UserAPI/internal/user"

type Handlers struct {
	useCase user.UseCase
}

func NewHandlers(useCase user.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}
