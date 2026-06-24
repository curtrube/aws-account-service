package account

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/curtrube/aws-account-service/internal/account"
)

type ListAccountsAPIClient interface {
	ListAccounts(
		context.Context,
		*organizations.ListAccountsInput,
		...func(*organizations.Options),
	) (*organizations.ListAccountsOutput, error)
}

// ListAccounts returns a list of accounts from the AWS organization.
func ListAccounts(ctx context.Context, client ListAccountsAPIClient) ([]account.Account, error) {
	out, err := client.ListAccounts(
		ctx,
		&organizations.ListAccountsInput{},
	)
	if err != nil {
		return nil, err
	}

	accounts := make([]account.Account, 0, len(out.Accounts))

	for _, a := range out.Accounts {
		accounts = append(accounts, account.Account{
			ID:              aws.ToString(a.Id),
			Name:            aws.ToString(a.Name),
			Email:           aws.ToString(a.Email),
			JoinedMethod:    string(a.JoinedMethod),
			JoinedTimestamp: aws.ToTime(a.JoinedTimestamp),
			State:           string(a.State),
		})
	}
	return accounts, nil
}
