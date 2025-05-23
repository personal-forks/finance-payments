package column

import (
	"errors"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/public/column/client"
	"github.com/formancehq/payments/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Column Plugin Uninstall", func() {
	var (
		plg            *Plugin
		mockHTTPClient *client.MockHTTPClient
		now            time.Time
	)

	BeforeEach(func() {
		now = time.Now().UTC()
		ctrl := gomock.NewController(GinkgoT())
		mockHTTPClient = client.NewMockHTTPClient(ctrl)
		plg = &Plugin{
			client: client.New("test", "aseplye", "https://test.com"),
		}
		plg.client.SetHttpClient(mockHTTPClient)
	})

	Context("uninstalling connector", func() {
		It("should handle empty webhooks list", func(ctx SpecContext) {
			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				200,
				nil,
			).SetArg(
				2,
				client.ListWebhookResponseWrapper[[]*client.EventSubscription]{
					WebhookEndpoints: []*client.EventSubscription{},
				},
			)

			resp, err := plg.uninstall(ctx, models.UninstallRequest{
				ConnectorID: "test-connector",
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(models.UninstallResponse{}))
		})

		It("should handle webhook list event error", func(ctx SpecContext) {
			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				500,
				errors.New("list failed"),
			)
			resp, err := plg.uninstall(ctx, models.UninstallRequest{
				ConnectorID: "test-connector",
			})
			Expect(err).To(MatchError("failed to list web hooks: list failed : "))
			Expect(resp).To(Equal(models.UninstallResponse{}))
		})

		It("should handle webhook deletion error", func(ctx SpecContext) {
			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				200,
				nil,
			).SetArg(
				2,
				client.ListWebhookResponseWrapper[[]*client.EventSubscription]{
					WebhookEndpoints: []*client.EventSubscription{
						{
							ID:            "webhook-2",
							URL:           "https://example.com/test-connector/webhook",
							CreatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
							UpdatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
							Description:   "description",
							EnabledEvents: []string{"book.transfer.completed"},
							Secret:        "secret",
							IsDisabled:    false,
						},
					},
				},
			)

			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				200,
				errors.New("deletion failed"),
			)

			resp, err := plg.uninstall(ctx, models.UninstallRequest{
				ConnectorID: "test-connector",
			})
			Expect(err).To(MatchError("failed to delete web hooks: deletion failed : "))
			Expect(resp).To(Equal(models.UninstallResponse{}))
		})

		It("should successfully uninstall and delete webhooks", func(ctx SpecContext) {
			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				200,
				nil,
			).SetArg(2, client.ListWebhookResponseWrapper[[]*client.EventSubscription]{
				WebhookEndpoints: []*client.EventSubscription{
					{
						ID:            "webhook-1",
						URL:           "https://example.com/test-connector/webhook",
						CreatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
						UpdatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
						Description:   "description",
						EnabledEvents: []string{"book.transfer.completed"},
						Secret:        "secret",
						IsDisabled:    false,
					},
					{
						ID:            "webhook-2",
						URL:           "https://example.com/other-connector/webhook",
						CreatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
						UpdatedAt:     now.Add(-time.Duration(5) * time.Minute).UTC().Format(time.RFC3339),
						Description:   "description",
						EnabledEvents: []string{"book.transfer.completed"},
						Secret:        "secret",
						IsDisabled:    false,
					},
				},
			})

			mockHTTPClient.EXPECT().Do(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(
				200,
				nil,
			)

			resp, err := plg.uninstall(ctx, models.UninstallRequest{
				ConnectorID: "test-connector",
			})
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(models.UninstallResponse{}))
		})
	})
})
