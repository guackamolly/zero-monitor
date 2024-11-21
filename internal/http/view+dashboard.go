package http

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type DashboardView struct {
	InviteLink DashboardNetworkInviteLinkView
}

type DashboardNetworkInviteLinkView struct {
	URL       string
	ExpiresAt time.Time
}

func (v DashboardView) WithInviteLink(
	inviteLink DashboardNetworkInviteLinkView,
) DashboardView {
	v.InviteLink = inviteLink
	return v
}

func (v DashboardView) ShowInviteLink() bool {
	return v.InviteLink.URL != "" && v.InviteLink.ExpiresAt.After(time.Now())
}

func (v DashboardNetworkInviteLinkView) Expiry() string {
	return models.Duration(time.Until(v.ExpiresAt)).String()
}

func NewDashboardView() DashboardView {
	return DashboardView{}
}

func NewDashNetworkInviteLinkView(
	url string,
	expiresAt time.Time,
) DashboardNetworkInviteLinkView {
	return DashboardNetworkInviteLinkView{
		URL:       url,
		ExpiresAt: expiresAt,
	}
}
