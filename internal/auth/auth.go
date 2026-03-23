package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	tea "github.com/charmbracelet/bubbletea"
)

type Context struct {
	Credential       *azidentity.DefaultAzureCredential
	SubscriptionID   string
	TenantID         string
	UserPrincipal    string
	SubscriptionName string
}

type Subscription struct {
	ID        string
	Name      string
	TenantID  string
	IsDefault bool
}

type SubscriptionsMsg struct {
	Subscriptions []Subscription
}

type AuthReadyMsg struct {
	Ctx *Context
}

type AuthErrMsg struct {
	Err error
}

type accountShowOutput struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	Name     string `json:"name"`
	User     struct {
		Name string `json:"name"`
	} `json:"user"`
}

type accountListOutput struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
}

func ListSubscriptions() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "az", "account", "list", "--output", "json").Output()
		if err != nil {
			return AuthErrMsg{Err: fmt.Errorf("az account list: %w", err)}
		}
		var accounts []accountListOutput
		if err := json.Unmarshal(out, &accounts); err != nil {
			return AuthErrMsg{Err: fmt.Errorf("parse az account list: %w", err)}
		}
		subs := make([]Subscription, len(accounts))
		for i, a := range accounts {
			subs[i] = Subscription{
				ID:        a.ID,
				Name:      a.Name,
				TenantID:  a.TenantID,
				IsDefault: a.IsDefault,
			}
		}
		if len(subs) == 0 {
			return AuthErrMsg{Err: fmt.Errorf("no subscriptions; run az login")}
		}
		return SubscriptionsMsg{Subscriptions: subs}
	}
}

func ResolveWithSubscription(subID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "az", "account", "show", "--subscription", subID).Output()
		if err != nil {
			return AuthErrMsg{Err: fmt.Errorf("az account show: %w", err)}
		}
		var acc accountShowOutput
		if err := json.Unmarshal(out, &acc); err != nil {
			return AuthErrMsg{Err: fmt.Errorf("parse az account show: %w", err)}
		}
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return AuthErrMsg{Err: fmt.Errorf("DefaultAzureCredential: %w", err)}
		}
		return AuthReadyMsg{Ctx: &Context{
			Credential:       cred,
			SubscriptionID:   acc.ID,
			TenantID:         acc.TenantID,
			UserPrincipal:    acc.User.Name,
			SubscriptionName: acc.Name,
		}}
	}
}
