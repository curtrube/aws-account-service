package aws

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/curtrube/aws-account-service/internal/account"
	"github.com/stretchr/testify/require"
)

type listAccountsStub struct {
	responses []*organizations.ListAccountsOutput
	err       error
	calls     int
	inputs    []*organizations.ListAccountsInput
}

func (s *listAccountsStub) ListAccounts(
	_ context.Context,
	input *organizations.ListAccountsInput,
	_ ...func(*organizations.Options),
) (*organizations.ListAccountsOutput, error) {
	s.inputs = append(s.inputs, input)

	if s.err != nil {
		return nil, s.err
	}
	if s.calls >= len(s.responses) {
		return nil, errors.New("unexpected ListAccounts call")
	}
	response := s.responses[s.calls]
	s.calls++
	return response, nil
}

func TestListAccounts(t *testing.T) {
	t.Parallel()

	joinedAt := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	nextToken := "next-token"

	client := &listAccountsStub{
		responses: []*organizations.ListAccountsOutput{
			{
				Accounts: []types.Account{
					{
						Id:              awssdk.String("123456789012"),
						Name:            awssdk.String("Prod"),
						Arn:             awssdk.String("arn:aws:organizations::123456789012:account/o-example/123456789012"),
						Email:           awssdk.String("prod@example.com"),
						State:           types.AccountStateActive,
						JoinedMethod:    types.AccountJoinedMethodCreated,
						JoinedTimestamp: &joinedAt,
					},
				},
				NextToken: &nextToken,
			},
			{
				Accounts: []types.Account{
					{
						Id:    aws.String("210987654321"),
						Name:  awssdk.String("Dev"),
						State: types.AccountStateSuspended,
					},
				},
			},
		},
	}

	got, err := ListAccounts(context.Background(), client)

	require.NoError(t, err)
	require.Equal(t, []account.Account{
		{
			AccountID:    "123456789012",
			AccountName:  "Prod",
			AccountArn:   "arn:aws:organizations::123456789012:account/o-example/123456789012",
			Email:        "prod@example.com",
			AccountState: "ACTIVE",
			JoinedMethod: "CREATED",
			JoinedAt:     joinedAt,
		},
		{
			AccountID:    "210987654321",
			AccountName:  "Dev",
			AccountState: "SUSPENDED",
		},
	}, got)

	require.Len(t, client.inputs, 2)
	require.Empty(t, awssdk.ToString(client.inputs[0].NextToken))
	require.Equal(t, nextToken, awssdk.ToString(client.inputs[1].NextToken))
}

func TestListAccounts_ReturnsError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("AWS organizations list accounts failed")
	client := &listAccountsStub{err: wantErr}

	_, err := ListAccounts(context.Background(), client)

	require.ErrorIs(t, err, wantErr)
}
