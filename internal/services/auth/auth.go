package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go_grpc/auth/internal/domain/models"
	"github.com/go_grpc/auth/internal/lib/jwt"
	"github.com/go_grpc/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	appSaver    AppSaver
	tokenTTL    time.Duration
}

type AppSaver interface {
	SaveApp(ctx context.Context, userId int64, name string, secret string) (int64, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid App Id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	appSaver AppSaver,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		appSaver:    appSaver,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.Any("err", err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("comparing password hash")
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Debug("fetching app", slog.Int("app_id", appID))
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("err", err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists", slog.Any("err", err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", slog.Any("err", err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// previously SaveApp was not present in the service, but Ensure appSaver exists and used in CreateApp

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found", slog.Any("err", err))

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

// CreateApp creates a new application with provided name and secret.
func (a *Auth) CreateApp(ctx context.Context, userId int64, name string, secret string) (int64, string, error) {
	const op = "Auth.CreateApp"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("creating app")

	id, err := a.appSaver.SaveApp(ctx, userId, name, secret)
	if err != nil {
		if errors.Is(err, storage.ErrAppExists) {
			return 0, "", fmt.Errorf("%s: %w", op, storage.ErrAppExists)
		}

		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("app created", slog.Int64("app_id", id))

	return id, name, nil
}
