package account

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/stretchr/testify/require"
)

type MockOrganizations struct {
    Output *organizations.ListAccountsOutput
    Err    error
}

func (m *MockOrganizations) ListAccounts(
    ctx context.Context,
    input *organizations.ListAccountsInput,
    optFns ...func(*organizations.Options),
) (*organizations.ListAccountsOutput, error) {
    return m.Output, m.Err
}

func TestListAccounts(t *testing.T) {
    mock := &MockOrganizations{
        Output: &organizations.ListAccountsOutput{
            Accounts: []types.Account{
                {
                    Id:   aws.String("123456789012"),
                    Name: aws.String("Prod"),
                },
                {
                    Id:   aws.String("210987654321"),
                    Name: aws.String("Dev"),
                },
            },
        },
    }

    accounts, err := ListAccounts(context.Background(), mock)

    require.NoError(t, err)
    require.Len(t, accounts, 2)
}
