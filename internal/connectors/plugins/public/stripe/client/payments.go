package client

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/metrics"
	"github.com/stripe/stripe-go/v79"
)

const (
	expandSource                    = "data.source"
	expandSourceCharge              = "data.source.charge"
	expandSourceDispute             = "data.source.dispute"
	expandSourcePayout              = "data.source.payout"
	expandSourceRefund              = "data.source.refund"
	expandSourceTransfer            = "data.source.transfer"
	expandSourcePaymentIntent       = "data.source.payment_intent"
	expandSourceRefundPaymentIntent = "data.source.refund.payment_intent"
)

func (c *client) GetPayments(
	ctx context.Context,
	accountID string,
	timeline Timeline,
	pageSize int64,
) (results []*stripe.BalanceTransaction, _ Timeline, hasMore bool, err error) {
	results = make([]*stripe.BalanceTransaction, 0, int(pageSize))

	if !timeline.IsCaughtUp() {
		var oldest interface{}
		oldest, timeline, hasMore, err = scanForOldest(timeline, pageSize, func(params stripe.ListParams) (stripe.ListContainer, error) {
			if accountID != "" {
				params.StripeAccount = &accountID
			}
			params.Context = metrics.OperationContext(ctx, "list_transactions_scan")
			transactionParams := &stripe.BalanceTransactionListParams{ListParams: params}
			expandBalanceTransactionParams(transactionParams)
			itr := c.balanceTransactionClient.List(transactionParams)
			return itr.BalanceTransactionList(), wrapSDKErr(itr.Err())
		})
		if err != nil {
			return results, timeline, false, err
		}
		// either there are no records or we haven't found the start yet
		if !timeline.IsCaughtUp() {
			return results, timeline, hasMore, nil
		}
		results = append(results, oldest.(*stripe.BalanceTransaction))
	}

	filters := stripe.ListParams{
		Context:      metrics.OperationContext(ctx, "list_transactions"),
		Limit:        limit(pageSize, len(results)),
		EndingBefore: &timeline.LatestID,
		Single:       true, // turn off autopagination
	}

	if accountID != "" {
		filters.StripeAccount = &accountID
	}

	params := &stripe.BalanceTransactionListParams{
		ListParams: filters,
	}
	expandBalanceTransactionParams(params)

	itr := c.balanceTransactionClient.List(params)
	results = append(results, itr.BalanceTransactionList().Data...)
	if len(results) == 0 {
		return results, timeline, itr.BalanceTransactionList().ListMeta.HasMore, wrapSDKErr(itr.Err())
	}

	timeline.LatestID = results[len(results)-1].ID
	return results, timeline, itr.BalanceTransactionList().ListMeta.HasMore, wrapSDKErr(itr.Err())
}

func expandBalanceTransactionParams(params *stripe.BalanceTransactionListParams) {
	params.AddExpand(expandSource)
	params.AddExpand(expandSourceCharge)
	params.AddExpand(expandSourceDispute)
	params.AddExpand(expandSourcePayout)
	params.AddExpand(expandSourceRefund)
	params.AddExpand(expandSourceTransfer)
	params.AddExpand(expandSourcePaymentIntent)
	params.AddExpand(expandSourceRefundPaymentIntent)
}
