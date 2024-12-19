package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/zaketn/sso/internal/domain/models"
	"github.com/zaketn/sso/internal/lib/jwt"
	"github.com/zaketn/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("the credentials is invalid")
	ErrInvalidAppId       = errors.New("the app id is invalid")
	ErrUserExists         = errors.New("user already exists")
)

// New returns new instance of the Auth structure
func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    usrSaver,
		userProvider: usrProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appId int) (token string, err error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)
	log.Info("attempting to log user")

	usr, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PassHash, []byte(password)); err != nil {
		log.Error("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = jwt.NewToken(usr, app, a.tokenTTL)
	if err != nil {
		log.Error("failed generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (userId int64, err error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("the user already exists")

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
	)

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("failed to get admin info")

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		log.Error("failed to get admin info")

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
