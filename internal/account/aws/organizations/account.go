package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	orgsdk "github.com/aws/aws-sdk-go-v2/service/organizations"
	awstypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/curtrube/aws-account-service/internal/account"
)

type ListAccountsAPIClient interface {
	ListAccounts(
		context.Context,
		*orgsdk.ListAccountsInput,
		...func(*orgsdk.Options),
	) (*orgsdk.ListAccountsOutput, error)
}

// ListAccounts returns a list of accounts from the AWS organization.
func ListAccounts(ctx context.Context, client ListAccountsAPIClient) ([]account.Account, error) {
	var awsAccounts []awstypes.Account

	paginator := orgsdk.NewListAccountsPaginator(client, &orgsdk.ListAccountsInput{})
	for page := 1; paginator.HasMorePages(); page++ {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("list organization accounts, page %d: %w", page, err)
		}

		awsAccounts = append(awsAccounts, output.Accounts...)
	}

	out := make([]account.Account, 0, len(awsAccounts))

	for _, a := range awsAccounts {
		out = append(out, account.Account{
			AccountID:    awssdk.ToString(a.Id),
			AccountName:  awssdk.ToString(a.Name),
			AccountArn:   awssdk.ToString(a.Arn),
			Email:        awssdk.ToString(a.Email),
			AccountState: string(a.State),
			JoinedMethod: string(a.JoinedMethod),
			JoinedAt:     awssdk.ToTime(a.JoinedTimestamp),
		})
	}

	return out, nil
}
