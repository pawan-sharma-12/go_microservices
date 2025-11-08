package account
import
 ( 
		"context"
 		"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string, email string) (*Account, error)
	GetAccountByID  (ctx context.Context, id string) (*Account, error)
	GetAccounts (ctx context.Context, skip uint64, take uint64) ([]Account, error)

}


type Account struct {
	ID	string `json:"id"`
	Name	string `json:"name"`
	Email	string `json:"email"`
}
type accountService struct {
	repo Repository
}

func NewAccountService(repo Repository) Service {
	return &accountService{
		repo: repo,
	}
}

func (s *accountService) PostAccount(ctx context.Context, name string, email string) (*Account, error) {
	account := Account{
		ID:    ksuid.New().String(),
		Name:  name,
		Email: email,
	}
	if err := s.repo.PutAccount(ctx, account); err != nil {
		return nil, err
	}
	return &account, nil
}
func (s *accountService) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	return s.repo.GetAccountByID(ctx, id)
}
func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repo.ListAccounts(ctx, skip, take)
}