package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// Service for authentication
type Service interface {
	Verify(ctx context.Context, idToken string) error
	Create(ctx context.Context, uid, email, name, password string) (*auth.UserRecord, error)
	GetCustomToken(ctx context.Context, uid string) (string, error)
	VerifyCustomToken(ctx context.Context, customToken string) (string, error)
	Delete(ctx context.Context, uid string) error
	GetUserID(ctx context.Context, idToken string) (string, error)
	GetUserByID(ctx context.Context, uid string) (*auth.UserRecord, error)
	GetUserByToken(ctx context.Context, idToken string) (*auth.UserRecord, error)
	DeleteUserByID(ctx context.Context, uid string) error
}

type svc struct {
	auth   *auth.Client
	apiKey string
}

// New authentication service
func New(ctx context.Context, credentialsFile, apiKey string) (Service, error) {
	firebase, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	auth, err := firebase.Auth(ctx)
	if err != nil {
		return nil, err
	}

	// TODO : Validate the apiKey

	return svc{auth: auth, apiKey: apiKey}, nil
}

// Verify validates the token
func (s svc) Verify(ctx context.Context, idToken string) error {
	if _, err := s.getUserID(ctx, idToken); err != nil {
		return err
	}
	return nil
}

// Create creates a test-user (other users are created in UI)
func (s svc) Create(ctx context.Context, uid, email, name, password string) (*auth.UserRecord, error) {
	user := &auth.UserToCreate{}
	user.Email(email)
	user.UID(uid)
	user.DisplayName(name)
	created, err := s.auth.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Set testClaims - so we know that the user was created from a test
	testClaims := map[string]interface{}{"TestAt": time.Now(), "Test": true}
	if err := s.auth.SetCustomUserClaims(ctx, created.UID, testClaims); err != nil {
		return nil, err
	}

	return created, nil
}

// Delete user (can only delete test-users)
func (s svc) Delete(ctx context.Context, uid string) error {
	user, err := s.auth.GetUser(ctx, uid)
	if err != nil {
		return err
	}
	if test, ok := user.CustomClaims["Test"]; !test.(bool) || !ok {
		return errors.New("Cannot delete this user")
	}
	if err := s.auth.DeleteUser(ctx, uid); err != nil {
		return fmt.Errorf("Cannot delete user: %v", err)
	}
	return nil
}

func (s svc) GetCustomToken(ctx context.Context, uid string) (string, error) {
	return s.auth.CustomToken(ctx, uid)
}

func (s svc) VerifyCustomToken(ctx context.Context, customToken string) (string, error) {
	idToken, err := s.signInWithCustomToken(customToken)
	if err != nil {
		return "", err
	}
	return idToken.Token, nil
}

// Verify validates the token
func (s svc) GetUserID(ctx context.Context, idToken string) (string, error) {
	return s.getUserID(ctx, idToken)
}

// GetUserByToken returns the user specified by token
func (s svc) GetUserByToken(ctx context.Context, idToken string) (*auth.UserRecord, error) {
	uid, err := s.getUserID(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return s.auth.GetUser(ctx, uid)
}

// GetUserByID returns user specified by ID
func (s svc) GetUserByID(ctx context.Context, uid string) (*auth.UserRecord, error) {
	return s.auth.GetUser(ctx, uid)
}

func (s svc) DeleteUserByID(ctx context.Context, uid string) error {
	return s.auth.DeleteUser(ctx, uid)
}

func (s svc) getUserID(ctx context.Context, idToken string) (string, error) {
	if len(idToken) == 0 {
		return "", errors.New("token is empty")
	}
	token, err := s.auth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

type idToken struct {
	Token        string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

func (s svc) signInWithCustomToken(customToken string) (*idToken, error) {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", s.apiKey)

	data := []byte(fmt.Sprintf(`{"token":"%s","returnSecureToken":true}`, customToken))
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Cannot get ID-Token: response: %v", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var idToken idToken
	if err := json.Unmarshal(respBody, &idToken); err != nil {
		return nil, err
	}

	return &idToken, nil
}
