package account

import (
	"aws-account-service/internal/model"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

type ListAccountsAPIClient interface {
	ListAccounts(
		context.Context,
		*organizations.ListAccountsInput,
		...func(*organizations.Options),
	) (*organizations.ListAccountsOutput, error)
}

func ListAccounts(ctx context.Context, client ListAccountsAPIClient) ([]model.Account, error) {
	out, err := client.ListAccounts(
		ctx,
		&organizations.ListAccountsInput{},
	)
	if err != nil {
		return nil, err
	}

	accounts := make([]model.Account, 0, len(out.Accounts))

	for _, a := range out.Accounts {
		accounts = append(accounts, model.Account{
			ID:              aws.ToString(a.Id),
			Name:            aws.ToString(a.Name),
			Email:           aws.ToString(a.Email),
			JoinedMethod:    string(a.JoinedMethod),
			JoinedTimestamp: *a.JoinedTimestamp,
			State:           string(a.State),
		})
	}

	return accounts, nil
}
