package column

import (
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/public/column/client"
	"github.com/formancehq/payments/internal/models"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Column Plugin Payments", func() {
	var (
		plg *Plugin
	)

	BeforeEach(func() {
		plg = &Plugin{}
	})

	Context("fetching next payments", func() {
		var (
			mockHTTPClient *client.MockHTTPClient
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			mockHTTPClient = client.NewMockHTTPClient(ctrl)
			plg.client = client.New("test", "aseplye", "https://test.com")
			plg.client.SetHttpClient(mockHTTPClient)
		})

		Context("validateReversePayout", func() {
			It("should validate a valid reverse payout request", func() {
				pr := models.PSPPaymentInitiationReversal{
					Metadata: map[string]string{
						client.ColumnReasonMetadataKey: "incorrect_amount",
					},
					RelatedPaymentInitiation: models.PSPPaymentInitiation{
						Reference: "test-reference",
					},
				}

				err := plg.validateReversePayout(pr)
				Expect(err).To(BeNil())
			})

			It("should return error when metadata is nil", func() {
				pr := models.PSPPaymentInitiationReversal{
					Metadata: nil,
				}

				err := plg.validateReversePayout(pr)
				Expect(err).To(MatchError("validation error occurred for field metadata: required field metadata must be provided"))
			})

			It("should return error when relatedPaymentInitiation.reference is missing", func() {
				pr := models.PSPPaymentInitiationReversal{
					Metadata: map[string]string{
						client.ColumnReasonMetadataKey: "incorrect_amount",
					},
					RelatedPaymentInitiation: models.PSPPaymentInitiation{},
				}

				err := plg.validateReversePayout(pr)
				Expect(err).To(MatchError("validation error occurred for field relatedPaymentInitiation.reference: required field relatedPaymentInitiation.reference must be provided"))
			})

			It("should return error when reason is missing", func() {
				pr := models.PSPPaymentInitiationReversal{
					Metadata: map[string]string{},
					RelatedPaymentInitiation: models.PSPPaymentInitiation{
						Reference: "test-reference",
					},
				}

				err := plg.validateReversePayout(pr)
				Expect(err).To(MatchError("validation error occurred for field com.column.spec/reason: required metadata field com.column.spec/reason must be provided"))
			})

			It("should return error when reason is invalid", func() {
				pr := models.PSPPaymentInitiationReversal{
					Metadata: map[string]string{
						client.ColumnReasonMetadataKey: "invalid-reason",
					},
				}

				err := plg.validateReversePayout(pr)
				Expect(err).To(MatchError("validation error occurred for field com.column.spec/reason: required metadata field com.column.spec/reason must be a valid reason"))
			})
		})

		Context("JSON Marshaling Errors", func() {
			It("should return an error when marshaling reverse payout request fails", func(ctx SpecContext) {
				req := models.ReversePayoutRequest{
					PaymentInitiationReversal: models.PSPPaymentInitiationReversal{
						Metadata: map[string]string{
							client.ColumnReasonMetadataKey: "incorrect_amount",
						},
						RelatedPaymentInitiation: models.PSPPaymentInitiation{
							Reference: "test-reference",
						},
					},
				}

				mockHTTPClient.EXPECT().Do(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					500,
					fmt.Errorf("mock marshal error"),
				)

				res, err := plg.ReversePayout(ctx, req)
				Expect(err).ToNot(BeNil())
				Expect(res).To(Equal(models.ReversePayoutResponse{}))
				Expect(err.Error()).To(ContainSubstring("mock marshal error"))
			})
		})

		Context("HTTP Request Creation Errors", func() {
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				mockHTTPClient = client.NewMockHTTPClient(ctrl)
				plg.client = client.New("test", "aseplye", "http://invalid:port")
				plg.client.SetHttpClient(mockHTTPClient)
			})

			It("should return an error when reverse payout request URL is invalid", func(ctx SpecContext) {
				req := models.ReversePayoutRequest{
					PaymentInitiationReversal: models.PSPPaymentInitiationReversal{
						Metadata: map[string]string{
							client.ColumnReasonMetadataKey: "incorrect_amount",
						},
						RelatedPaymentInitiation: models.PSPPaymentInitiation{
							Reference: "test-reference",
						},
					},
				}

				res, err := plg.ReversePayout(ctx, req)
				Expect(err).ToNot(BeNil())
				Expect(res).To(Equal(models.ReversePayoutResponse{}))
				Expect(err.Error()).To(ContainSubstring("failed to create reverse payout request"))
			})
		})

		Context("CreatedAt Timestamp Parsing", func() {
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				mockHTTPClient = client.NewMockHTTPClient(ctrl)
				plg.client = client.New("test", "aseplye", "https://test.com")
				plg.client.SetHttpClient(mockHTTPClient)
			})

			It("should successfully parse a valid timestamp", func(ctx SpecContext) {
				req := models.ReversePayoutRequest{
					PaymentInitiationReversal: models.PSPPaymentInitiationReversal{
						Metadata: map[string]string{
							client.ColumnReasonMetadataKey: "incorrect_amount",
						},
						RelatedPaymentInitiation: models.PSPPaymentInitiation{
							Reference: "test-reference",
						},
					},
				}

				mockHTTPClient.EXPECT().Do(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					200,
					nil,
				).SetArg(2, client.ReversePayoutResponse{
					CreatedAt:     "2024-03-04T12:00:00Z",
					ID:            "test-id",
					Amount:        100,
					CurrencyCode:  "USD",
					BankAccountID: "test-bank",
					Description:   "test description",
					Status:        "completed",
				})

				res, err := plg.ReversePayout(ctx, req)
				Expect(err).To(BeNil())
				Expect(res.Payment.CreatedAt).ToNot(BeZero())
				Expect(res.Payment.CreatedAt.Year()).To(Equal(2024))
				Expect(res.Payment.CreatedAt.Month()).To(Equal(time.March))
				Expect(res.Payment.CreatedAt.Day()).To(Equal(4))
			})

			It("should return an error when timestamp is invalid", func(ctx SpecContext) {
				req := models.ReversePayoutRequest{
					PaymentInitiationReversal: models.PSPPaymentInitiationReversal{
						Metadata: map[string]string{
							client.ColumnReasonMetadataKey: "incorrect_amount",
						},
						RelatedPaymentInitiation: models.PSPPaymentInitiation{
							Reference: "test-reference",
						},
					},
				}

				mockHTTPClient.EXPECT().Do(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(
					200,
					nil,
				).SetArg(2, client.ReversePayoutResponse{
					CreatedAt:     "invalid-timestamp",
					ID:            "test-id",
					Amount:        100,
					CurrencyCode:  "USD",
					BankAccountID: "test-bank",
					Description:   "test description",
					Status:        "completed",
				})

				res, err := plg.ReversePayout(ctx, req)
				Expect(err).ToNot(BeNil())
				Expect(res).To(Equal(models.ReversePayoutResponse{}))
				Expect(err.Error()).To(ContainSubstring("parsing time"))
			})
		})
	})
})
