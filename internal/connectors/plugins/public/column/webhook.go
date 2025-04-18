package column

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/formancehq/payments/internal/connectors/plugins/public/column/client"
	"github.com/formancehq/payments/internal/models"
)

const (
	HeadersSignature = "Column-Signature"
)

type webhookConfig struct {
	urlPath string
	fn      func(context.Context, client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error)
	secret  string
}

type defaultVerifier struct{}

type WebhookVerifier interface {
	verifyWebhookSignature(payload []byte, header, secret string) error
}

func (p *Plugin) initWebhookConfig() map[client.EventCategory]webhookConfig {
	p.webhookConfigs = map[client.EventCategory]webhookConfig{
		client.EventCategoryBookTransferCompleted: {
			urlPath: "/book/transfer/completed",
			fn:      p.translateBookTransfer,
		},
		client.EventCategoryBookTransferCanceled: {
			urlPath: "/book/transfer/canceled",
			fn:      p.translateBookTransfer,
		},
		client.EventCategoryBookTransferUpdated: {
			urlPath: "/book/transfer/updated",
			fn:      p.translateBookTransfer,
		},
		client.EventCategoryBookTransferHoldCreated: {
			urlPath: "/book/transfer/hold_created",
			fn:      p.translateBookTransfer,
		},
		client.EventCategoryWireTransferOutgoingCompleted: {
			urlPath: "/wire/outgoing_transfer/completed",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryWireTransferInitiated: {
			urlPath: "/wire/outgoing_transfer/initiated",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryWireTransferIncomingCompleted: {
			urlPath: "/wire/incoming_transfer/completed",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryWireTransferSubmitted: {
			urlPath: "/wire/outgoing_transfer/submitted",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryWireTransferRejected: {
			urlPath: "/wire/outgoing_transfer/rejected",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryWireTransferManualReview: {
			urlPath: "/wire/outgoing_transfer/manual_review",
			fn:      p.translateWireTransfer,
		},
		client.EventCategoryACHTransferSettled: {
			urlPath: "/ach/outgoing_transfer/settled",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferInitiated: {
			urlPath: "/ach/outgoing_transfer/initiated",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferSubmitted: {
			urlPath: "/ach/outgoing_transfer/submitted",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferCompleted: {
			urlPath: "/ach/outgoing_transfer/completed",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferManualReview: {
			urlPath: "/ach/outgoing_transfer/manual_review",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferReturned: {
			urlPath: "/ach/outgoing_transfer/returned",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferCanceled: {
			urlPath: "/ach/outgoing_transfer/canceled",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferReturnDishonored: {
			urlPath: "/ach/outgoing_transfer/return_dishonored",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferReturnContested: {
			urlPath: "/ach/outgoing_transfer/return_contested",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHTransferNOC: {
			urlPath: "/ach/outgoing_transfer/noc",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingScheduled: {
			urlPath: "/ach/incoming_transfer/scheduled",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingSettled: {
			urlPath: "/ach/incoming_transfer/settled",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingNSF: {
			urlPath: "/ach/incoming_transfer/nsf",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingCompleted: {
			urlPath: "/ach/incoming_transfer/completed",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingReturned: {
			urlPath: "/ach/incoming_transfer/returned",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingReturnDishonored: {
			urlPath: "/ach/incoming_transfer/return_dishonored",
			fn:      p.translateAchTransfer,
		},
		client.EventCategoryACHIncomingReturnContested: {
			urlPath: "/ach/incoming_transfer/return_contested",
			fn:      p.translateAchTransfer,
		},
		client.EventCategorySwiftOutgoingInitiated: {
			urlPath: "/swift/outgoing_transfer/initiated",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategoryInternationalWireCompleted: {
			urlPath: "/swift/outgoing_transfer/completed",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingManualReview: {
			urlPath: "/swift/outgoing_transfer/manual_review",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingSubmitted: {
			urlPath: "/swift/outgoing_transfer/submitted",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingPendingReturn: {
			urlPath: "/swift/outgoing_transfer/pending_return",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingReturned: {
			urlPath: "/swift/outgoing_transfer/returned",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingCancellationRequested: {
			urlPath: "/swift/outgoing_transfer/cancellation_requested",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingCancellationAccepted: {
			urlPath: "/swift/outgoing_transfer/cancellation_accepted",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingCancellationRejected: {
			urlPath: "/swift/outgoing_transfer/cancellation_rejected",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftOutgoingTrackingUpdated: {
			urlPath: "/swift/outgoing_transfer/tracking_updated",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingInitiated: {
			urlPath: "/swift/incoming_transfer/initiated",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingCompleted: {
			urlPath: "/swift/incoming_transfer/completed",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingPendingReturn: {
			urlPath: "/swift/incoming_transfer/pending_return",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingReturned: {
			urlPath: "/swift/incoming_transfer/returned",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingCancellationRequested: {
			urlPath: "/swift/incoming_transfer/cancellation_requested",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingCancellationAccepted: {
			urlPath: "/swift/incoming_transfer/cancellation_accepted",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingCancellationRejected: {
			urlPath: "/swift/incoming_transfer/cancellation_rejected",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategorySwiftIncomingTrackingUpdated: {
			urlPath: "/swift/incoming_transfer/tracking_updated",
			fn:      p.translateInternationalWireTransfer,
		},
		client.EventCategoryRealtimeTransferInitiated: {
			urlPath: "/realtime/outgoing_transfer/initiated",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeTransferManualReview: {
			urlPath: "/realtime/outgoing_transfer/manual_review",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeTransferManualReviewApproved: {
			urlPath: "/realtime/outgoing_transfer/manual_review_approved",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeTransferManualReviewRejected: {
			urlPath: "/realtime/outgoing_transfer/manual_review_rejected",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeTransferRejected: {
			urlPath: "/realtime/outgoing_transfer/rejected",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeIncomingTransferCompleted: {
			urlPath: "/realtime/incoming_transfer/completed",
			fn:      p.translateRealtimeTransfer,
		},
		client.EventCategoryRealtimeTransferCompleted: {
			urlPath: "/realtime/outgoing_transfer/completed",
			fn:      p.translateRealtimeTransfer,
		},
	}
	return p.webhookConfigs
}

func (p *Plugin) createWebhooks(ctx context.Context, req models.CreateWebhooksRequest) (models.CreateWebhooksResponse, error) {
	var others []models.PSPOther

	if req.FromPayload == nil {
		return models.CreateWebhooksResponse{}, models.ErrMissingFromPayloadInRequest
	}

	if req.WebhookBaseUrl == "" {
		return models.CreateWebhooksResponse{}, client.ErrWebhookUrlMissing
	}

	if !strings.HasPrefix(req.WebhookBaseUrl, "https://") {
		return models.CreateWebhooksResponse{}, fmt.Errorf("webhook URL must use HTTPS protocol")
	}

	for eventType, config := range p.webhookConfigs {
		url, err := url.JoinPath(req.WebhookBaseUrl, config.urlPath)
		if err != nil {
			return models.CreateWebhooksResponse{}, err
		}

		esr := client.CreateEventSubscriptionRequest{
			URL:           url,
			EnabledEvents: []string{string(eventType)},
		}
		resp, err := p.client.CreateEventSubscription(ctx, &esr)
		if err != nil {
			return models.CreateWebhooksResponse{}, fmt.Errorf("failed to create webhook subscription: %w", err)
		}
		config.secret = resp.Secret
		p.webhookConfigs[eventType] = config

		raw, err := json.Marshal(resp)
		if err != nil {
			return models.CreateWebhooksResponse{}, err
		}
		others = append(others, models.PSPOther{
			ID:    resp.ID,
			Other: raw,
		})
	}

	return models.CreateWebhooksResponse{
		Others: others,
	}, nil
}

func (p *Plugin) translateWebhook(ctx context.Context, req models.TranslateWebhookRequest) (models.TranslateWebhookResponse, error) {
	signatures, ok := req.Webhook.Headers[HeadersSignature]
	if !ok || len(signatures) == 0 {
		return models.TranslateWebhookResponse{}, client.ErrWebhookHeaderXSignatureMissing
	}

	config, ok := p.webhookConfigs[client.EventCategory(req.Name)]
	if !ok {
		return models.TranslateWebhookResponse{}, client.ErrWebhookNameUnknown
	}

	if err := p.verifier.verifyWebhookSignature(req.Webhook.Body, signatures[0], config.secret); err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	var webhook client.WebhookEvent[json.RawMessage]
	if err := json.Unmarshal(req.Webhook.Body, &webhook); err != nil {
		return models.TranslateWebhookResponse{}, fmt.Errorf("failed to unmarshal webhook: %w", err)
	}

	res, err := config.fn(ctx, webhook)
	if err != nil {
		return models.TranslateWebhookResponse{}, err
	}

	res.IdempotencyKey = webhook.ID
	return models.TranslateWebhookResponse{
		Responses: []models.WebhookResponse{res},
	}, nil
}

func (p *Plugin) translateBookTransfer(ctx context.Context, webhook client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error) {
	var transfer client.TransferResponse
	dataBytes, err := json.Marshal(webhook.Data)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal webhook data: %w", err)
	}
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return models.WebhookResponse{}, err
	}

	pspPayment, err := p.transferToPayment(&transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map webhook book transfer payment: %w", err)
	}

	return models.WebhookResponse{
		Payment: pspPayment,
	}, nil
}

func (p *Plugin) translateAchTransfer(ctx context.Context, webhook client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error) {
	var transfer client.ACHPayoutResponse
	dataBytes, err := json.Marshal(webhook.Data)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal webhook data: %w", err)
	}
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return models.WebhookResponse{}, err
	}

	paymentResponse, err := client.MapAchPayout(transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map ach transfer webhook response: %w", err)
	}

	pspPayment, err := p.payoutToPayment(paymentResponse)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map ach payout to payment: %w", err)
	}

	return models.WebhookResponse{
		Payment: pspPayment,
	}, nil
}

func (p *Plugin) translateRealtimeTransfer(ctx context.Context, webhook client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error) {
	var transfer client.RealtimeTransferResponse
	dataBytes, err := json.Marshal(webhook.Data)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal webhook data: %w", err)
	}
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return models.WebhookResponse{}, err
	}

	paymentResponse, err := client.MapRealtimePayout(transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map realtime transfer webhook response: %w", err)
	}

	pspPayment, err := p.payoutToPayment(paymentResponse)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map realtime payout to payment: %w", err)
	}

	return models.WebhookResponse{
		Payment: pspPayment,
	}, nil
}

func (p *Plugin) translateWireTransfer(ctx context.Context, webhook client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error) {
	var transfer client.WirePayoutResponse
	dataBytes, err := json.Marshal(webhook.Data)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal webhook data: %w", err)
	}
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return models.WebhookResponse{}, err
	}

	paymentResponse, err := client.MapWirePayout(transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map wire transfer webhook response: %w", err)
	}

	pspPayment, err := p.payoutToPayment(paymentResponse)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map wire payout to payment: %w", err)
	}

	return models.WebhookResponse{
		Payment: pspPayment,
	}, nil
}

func (p *Plugin) translateInternationalWireTransfer(ctx context.Context, webhook client.WebhookEvent[json.RawMessage]) (models.WebhookResponse, error) {
	var transfer client.InternationalWirePayoutResponse
	dataBytes, err := json.Marshal(webhook.Data)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to marshal webhook data: %w", err)
	}
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return models.WebhookResponse{}, err
	}

	paymentResponse, err := client.MapInternationalWirePayout(transfer)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map international wire webhook response: %w", err)
	}

	pspPayment, err := p.payoutToPayment(paymentResponse)
	if err != nil {
		return models.WebhookResponse{}, fmt.Errorf("failed to map international wire payout to payment: %w", err)
	}

	return models.WebhookResponse{
		Payment: pspPayment,
	}, nil
}

func (v *defaultVerifier) verifyWebhookSignature(payload []byte, header string, webhookSecret string) error {
	h := hmac.New(sha256.New, []byte(webhookSecret))
	h.Write(payload)
	computedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(computedSignature), []byte(header)) {
		return errors.New("signature verification failed")
	}

	return nil
}
