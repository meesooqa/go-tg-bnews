package telegram

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/term"
)

type noSignUp struct{}

// AuthFlow implements auth.UserAuthenticator prompting the terminal for input
type AuthFlow struct {
	noSignUp
	PhoneNumber string // optional, will be prompted if empty
}

// NewTelegramAuthFlow creates a new Telegram authentication flow
func NewTelegramAuthFlow() auth.Flow {
	return auth.NewFlow(AuthFlow{
		noSignUp:    noSignUp{},
		PhoneNumber: os.Getenv("PHONE"),
	}, auth.SendCodeOptions{})
}

// SignUp is called when the user needs to sign up for a new account
func (noSignUp) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

// AcceptTermsOfService is called when the user needs to accept terms of service
func (noSignUp) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// SignUp is called when the user needs to sign up for a new account
func (AuthFlow) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("signing up not implemented in AuthFlow")
}

// AcceptTermsOfService is called when the user needs to accept terms of service
func (AuthFlow) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// Code prompts the user for the authentication code sent to their phone
func (AuthFlow) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")

	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(code), nil
}

// Phone prompts the user for their phone number
func (o AuthFlow) Phone(_ context.Context) (string, error) {
	if o.PhoneNumber != "" {
		return o.PhoneNumber, nil
	}

	fmt.Print("Enter phone in international format (e.g. +1234567890): ")

	phone, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(phone), nil
}

// Password prompts the user for their 2FA password
func (AuthFlow) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")

	bytePwd, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytePwd)), nil
}
