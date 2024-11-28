package http

import (
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type DashboardView struct {
	InviteLink DashboardNetworkInviteLinkView
}

type DashboardNetworkInviteLinkView struct {
	Code models.JoinNetworkCode
	URL  string
}

func (v DashboardView) WithInviteLink(
	inviteLink DashboardNetworkInviteLinkView,
) DashboardView {
	v.InviteLink = inviteLink
	return v
}

func (v DashboardView) ShowInviteLink() bool {
	return v.InviteLink.URL != "" && !v.InviteLink.Code.Expired()
}

func (v DashboardNetworkInviteLinkView) Expiry() string {
	return models.Duration(time.Until(v.Code.ExpiresAt)).String()
}

func (v DashboardNetworkInviteLinkView) String() string {
	return v.URL
}

func NewDashboardView() DashboardView {
	return DashboardView{}
}

func NewDashNetworkInviteLinkView(
	url string,
	code models.JoinNetworkCode,
) DashboardNetworkInviteLinkView {
	return DashboardNetworkInviteLinkView{
		Code: code,
		URL:  url,
	}
}
