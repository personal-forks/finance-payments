package bankingcircle

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/formancehq/go-libs/v2/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/public/bankingcircle/client"
	"github.com/formancehq/payments/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("BankingCircle Plugin Payouts Creation", func() {
	var (
		plg *Plugin
	)

	BeforeEach(func() {
		plg = &Plugin{}
	})

	Context("create payout", func() {
		var (
			m                          *client.MockClient
			samplePSPPaymentInitiation models.PSPPaymentInitiation
			now                        time.Time
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			m = client.NewMockClient(ctrl)
			plg.client = m
			now = time.Now().UTC()

			samplePSPPaymentInitiation = models.PSPPaymentInitiation{
				Reference:   "test1",
				CreatedAt:   now.UTC(),
				Description: "test1",
				SourceAccount: &models.PSPAccount{
					Reference:    "acc1",
					CreatedAt:    now.Add(-time.Duration(50) * time.Minute).UTC(),
					Name:         pointer.For("acc1"),
					DefaultAsset: pointer.For("EUR/2"),
				},
				DestinationAccount: &models.PSPAccount{
					Reference:    "acc2",
					CreatedAt:    now.Add(-time.Duration(49) * time.Minute).UTC(),
					Name:         pointer.For("acc2"),
					DefaultAsset: pointer.For("EUR/2"),
				},
				Amount: big.NewInt(100),
				Asset:  "EUR/2",
				Metadata: map[string]string{
					"foo": "bar",
				},
			}
		})

		It("should return an error - validation error - source account", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			req.PaymentInitiation.SourceAccount = nil

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("source account is required: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - validation error - destination account", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			req.PaymentInitiation.DestinationAccount = nil

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("destination account is required: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - validation error - asset not supported", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			req.PaymentInitiation.Asset = "HUF/2"

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get currency and precision from asset: missing currencies: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - get account error", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(nil, errors.New("test error"))

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get source account: test error: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - missing source account identifiers error", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{},
				}, nil)

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("no account identifiers provided for source account: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - get account 2 error", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{}},
				}, nil)

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.DestinationAccount.Reference).
				Return(nil, errors.New("test error"))

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("failed to get destination account: test error: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - missing destination account identifiers error", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{}},
				}, nil)

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.DestinationAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{},
				}, nil)

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("no account identifiers provided for destination account: invalid request"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should return an error - initiate payout error", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{
						Account:              "123456789",
						FinancialInstitution: "test",
						Country:              "US",
					}},
				}, nil)

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.DestinationAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{
						Account:              "987654321",
						FinancialInstitution: "test 2",
						Country:              "US",
					}},
				}, nil)

			m.EXPECT().InitiateTransferOrPayouts(ctx, &client.PaymentRequest{
				IdempotencyKey:         samplePSPPaymentInitiation.Reference,
				RequestedExecutionDate: samplePSPPaymentInitiation.CreatedAt,
				DebtorAccount: client.PaymentAccount{
					Account:              "123456789",
					FinancialInstitution: "test",
					Country:              "US",
				},
				DebtorReference:    samplePSPPaymentInitiation.Description,
				CurrencyOfTransfer: "EUR",
				Amount: client.Amount{
					Currency: "EUR",
					Amount:   "1.00",
				},
				ChargeBearer: "SHA",
				CreditorAccount: &client.PaymentAccount{
					Account:              "987654321",
					FinancialInstitution: "test 2",
					Country:              "US",
				},
			}).Return(nil, errors.New("test error"))

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).ToNot(BeNil())
			Expect(err).To(MatchError("test error"))
			Expect(resp).To(Equal(models.CreatePayoutResponse{}))
		})

		It("should be ok", func(ctx SpecContext) {
			req := models.CreatePayoutRequest{
				PaymentInitiation: samplePSPPaymentInitiation,
			}

			pr := client.PaymentResponse{
				PaymentID: "p1",
			}
			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.SourceAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{
						Account:              "123456789",
						FinancialInstitution: "test",
						Country:              "US",
					}},
				}, nil)

			m.EXPECT().GetAccount(ctx, samplePSPPaymentInitiation.DestinationAccount.Reference).
				Return(&client.Account{
					AccountIdentifiers: []client.AccountIdentifier{{
						Account:              "987654321",
						FinancialInstitution: "test 2",
						Country:              "US",
					}},
				}, nil)

			m.EXPECT().InitiateTransferOrPayouts(ctx, &client.PaymentRequest{
				IdempotencyKey:         samplePSPPaymentInitiation.Reference,
				RequestedExecutionDate: samplePSPPaymentInitiation.CreatedAt,
				DebtorAccount: client.PaymentAccount{
					Account:              "123456789",
					FinancialInstitution: "test",
					Country:              "US",
				},
				DebtorReference:    samplePSPPaymentInitiation.Description,
				CurrencyOfTransfer: "EUR",
				Amount: client.Amount{
					Currency: "EUR",
					Amount:   "1.00",
				},
				ChargeBearer: "SHA",
				CreditorAccount: &client.PaymentAccount{
					Account:              "987654321",
					FinancialInstitution: "test 2",
					Country:              "US",
				},
			}).Return(&pr, nil)

			paymentResponse := client.Payment{
				PaymentID:                    "p1",
				TransactionReference:         "transaction-p1",
				Status:                       "Processed",
				Classification:               "Outgoing",
				ProcessedTimestamp:           now.UTC(),
				LatestStatusChangedTimestamp: now.UTC(),
				DebtorInformation: client.DebtorInformation{
					AccountID: "123",
				},
				Transfer: client.Transfer{
					Amount: client.Amount{
						Currency: "EUR",
						Amount:   "1.00",
					},
				},
				CreditorInformation: client.CreditorInformation{
					AccountID: "321",
				},
			}

			m.EXPECT().GetPayment(ctx, "p1").Return(&paymentResponse, nil)

			raw, err := json.Marshal(&paymentResponse)
			Expect(err).To(BeNil())

			resp, err := plg.CreatePayout(ctx, req)
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(models.CreatePayoutResponse{
				Payment: models.PSPPayment{
					Reference:                   "p1",
					CreatedAt:                   now.UTC(),
					Type:                        models.PAYMENT_TYPE_PAYOUT,
					Amount:                      big.NewInt(100),
					Asset:                       "EUR/2",
					Scheme:                      models.PAYMENT_SCHEME_OTHER,
					Status:                      models.PAYMENT_STATUS_SUCCEEDED,
					SourceAccountReference:      pointer.For("123"),
					DestinationAccountReference: pointer.For("321"),
					Raw:                         raw,
				},
			}))
		})
	})
})
