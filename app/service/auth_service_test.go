package service

// =============================================================================
// UNIT TEST — AuthService (register / login / logout / user CRUD).
//
// Pola: stub embed interface (nil), implement hanya method yang dipakai,
// tangkap argumen untuk assert. Logger singleton sudah di-init oleh
// TestMain di main_test.go.
//
// Jalankan:  go test ./app/service/ -run TestAuth -v
// =============================================================================

import (
	"context"
	"database/sql"
	"dc-transaction-service/app/model/entity"
	"dc-transaction-service/app/model/request"
	"dc-transaction-service/app/repository"
	"dc-transaction-service/config"
	"dc-transaction-service/util"
	"encoding/base64"
	"errors"
	"net/http"
	"testing"
)

// stubUserRepo — implement UserRepositoryInterface; method lain tetap
// nil via embed → panic jika service memanggil method tak terduga.
type stubUserRepo struct {
	repository.UserRepositoryInterface

	usernameExist    bool
	checkUsernameErr error
	checkUsernameArg string

	emailExist    bool
	checkEmailErr error
	checkEmailArg string

	user    entity.User
	userErr error
	userArg string

	insertedUser  entity.User
	insertUID     string
	insertCalled  bool
	insertUserErr error

	updatedUser  entity.User
	updateCalled bool
	updateRows   int64
	updateErr    error

	deleteArg    string
	deleteCalled bool
	deleteRows   int64
	deleteErr    error
}

func (s *stubUserRepo) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	s.checkUsernameArg = username
	return s.usernameExist, s.checkUsernameErr
}

func (s *stubUserRepo) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	s.checkEmailArg = email
	return s.emailExist, s.checkEmailErr
}

func (s *stubUserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	s.userArg = username
	return s.user, s.userErr
}

func (s *stubUserRepo) GetUserByUID(ctx context.Context, uid string) (entity.User, error) {
	s.userArg = uid
	return s.user, s.userErr
}

func (s *stubUserRepo) InsertUser(ctx context.Context, data entity.User) (string, error) {
	s.insertedUser = data
	s.insertCalled = true
	return s.insertUID, s.insertUserErr
}

func (s *stubUserRepo) UpdateUserByUID(ctx context.Context, data entity.User) (int64, error) {
	s.updatedUser = data
	s.updateCalled = true
	return s.updateRows, s.updateErr
}

func (s *stubUserRepo) DeleteUserByUID(ctx context.Context, uid string) (int64, error) {
	s.deleteArg = uid
	s.deleteCalled = true
	return s.deleteRows, s.deleteErr
}

func newTestAuthConfig() *config.ApplicationConfig {
	return &config.ApplicationConfig{
		SecretKey: base64.StdEncoding.EncodeToString([]byte("test-secret-key-for-jwt-signing")),
		Jwt: config.JwtConfig{
			TokenExp:    1,
			RefTokenExp: 24,
		},
	}
}

func TestAuthRegister_Success(t *testing.T) {
	repo := &stubUserRepo{usernameExist: false, insertUID: "ABCDEF1234567890"}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	uid, status, err := svc.AuthRegister(context.Background(), request.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Password:  "secret123",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, status)
	}
	if uid != "ABCDEF1234567890" {
		t.Fatalf("expected uid returned from repo, got %q", uid)
	}
	if !repo.insertCalled {
		t.Fatal("expected InsertUser to be called")
	}
	if repo.checkUsernameArg != "johndoe" {
		t.Fatalf("expected CheckUsernameExist called with johndoe, got %s", repo.checkUsernameArg)
	}
	if repo.insertedUser.RoleID != defaultRoleID {
		t.Fatalf("expected default role %d, got %d", defaultRoleID, repo.insertedUser.RoleID)
	}
	if repo.insertedUser.Password != util.HashSHA256("secret123") {
		t.Fatalf("expected password sha256 hash, got %s", repo.insertedUser.Password)
	}
	if len(repo.insertedUser.Password) != 64 {
		t.Fatalf("expected sha256 hex length 64, got %d", len(repo.insertedUser.Password))
	}
}

func TestAuthRegister_DuplicateUsername(t *testing.T) {
	repo := &stubUserRepo{usernameExist: true}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	_, status, err := svc.AuthRegister(context.Background(), request.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Password:  "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
	}
	if repo.insertCalled {
		t.Fatal("expected InsertUser NOT to be called")
	}
}

func TestAuthRegister_RepoError(t *testing.T) {
	repo := &stubUserRepo{checkUsernameErr: errors.New("db down")}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	_, status, err := svc.AuthRegister(context.Background(), request.RegisterRequest{
		FirstName: "John",
		Username:  "johndoe",
		Password:  "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
	}
}

func TestAuthLogin_Success(t *testing.T) {
	repo := &stubUserRepo{
		user: entity.User{
			UID:       "UID12345",
			FirstName: "John",
			LastName:  "Doe",
			Username:  "johndoe",
			RoleID:    3,
			Password:  util.HashSHA256("secret123"),
		},
	}
	cfg := newTestAuthConfig()
	svc := NewAuthServiceImpl(cfg, repo)

	user, token, status, err := svc.AuthLogin(context.Background(), request.LoginRequest{
		Username: "johndoe",
		Password: "secret123",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if user.UID != "UID12345" {
		t.Fatalf("expected UID12345, got %s", user.UID)
	}
	if token.Token == "" || token.RefreshToken == "" {
		t.Fatal("expected non-empty token and refresh token")
	}

	secretKey, _ := base64.StdEncoding.DecodeString(cfg.SecretKey)
	claims, err := util.ValidateToken(context.Background(), token.Token, string(secretKey))
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims.UID != "UID12345" || claims.Username != "johndoe" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestAuthLogin_WrongPassword(t *testing.T) {
	repo := &stubUserRepo{
		user: entity.User{
			UID:      "UID12345",
			Username: "johndoe",
			Password: util.HashSHA256("secret123"),
		},
	}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	_, _, status, err := svc.AuthLogin(context.Background(), request.LoginRequest{
		Username: "johndoe",
		Password: "wrongpass",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, status)
	}
}

func TestAuthLogin_UserNotFound(t *testing.T) {
	repo := &stubUserRepo{userErr: sql.ErrNoRows}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	_, _, status, err := svc.AuthLogin(context.Background(), request.LoginRequest{
		Username: "nobody",
		Password: "secret123",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, status)
	}
}

func TestAuthLogout(t *testing.T) {
	svc := NewAuthServiceImpl(newTestAuthConfig(), &stubUserRepo{})

	status, err := svc.AuthLogout(context.Background())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
}

func TestUserUpdate_Success(t *testing.T) {
	repo := &stubUserRepo{
		user: entity.User{
			UID:       "UID12345",
			FirstName: "John",
			Username:  "johndoe",
			RoleID:    3,
			Password:  util.HashSHA256("oldpass"),
		},
		updateRows: 1,
	}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserUpdate(context.Background(), "UID12345", request.UserUpdateRequest{
		FirstName: "Johnny",
		LastName:  "Doe",
		Username:  "johndoe",
		RoleID:    2,
		Password:  "newpass123",
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if !repo.updateCalled {
		t.Fatal("expected UpdateUserByUID to be called")
	}
	if repo.updatedUser.Password != util.HashSHA256("newpass123") {
		t.Fatal("expected password to be re-hashed with sha256")
	}
}

func TestUserUpdate_KeepPassword(t *testing.T) {
	oldHash := util.HashSHA256("oldpass")
	repo := &stubUserRepo{
		user: entity.User{
			UID:      "UID12345",
			Username: "johndoe",
			Password: oldHash,
		},
		updateRows: 1,
	}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserUpdate(context.Background(), "UID12345", request.UserUpdateRequest{
		FirstName: "Johnny",
		Username:  "johndoe",
		RoleID:    3,
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if repo.updatedUser.Password != oldHash {
		t.Fatal("expected old password hash to be kept")
	}
}

func TestUserUpdate_NotFound(t *testing.T) {
	repo := &stubUserRepo{userErr: sql.ErrNoRows}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserUpdate(context.Background(), "NOPE0000", request.UserUpdateRequest{
		FirstName: "Johnny",
		Username:  "johndoe",
		RoleID:    3,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, status)
	}
	if repo.updateCalled {
		t.Fatal("expected UpdateUserByUID NOT to be called")
	}
}

func TestUserUpdate_DuplicateUsername(t *testing.T) {
	repo := &stubUserRepo{
		user: entity.User{
			UID:      "UID12345",
			Username: "johndoe",
		},
		usernameExist: true,
	}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserUpdate(context.Background(), "UID12345", request.UserUpdateRequest{
		FirstName: "Johnny",
		Username:  "taken",
		RoleID:    3,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
	}
	if repo.updateCalled {
		t.Fatal("expected UpdateUserByUID NOT to be called")
	}
}

func TestUserDelete_Success(t *testing.T) {
	repo := &stubUserRepo{deleteRows: 1}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserDelete(context.Background(), "UID12345")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}
	if repo.deleteArg != "UID12345" {
		t.Fatalf("expected DeleteUserByUID called with UID12345, got %s", repo.deleteArg)
	}
}

func TestUserDelete_NotFound(t *testing.T) {
	repo := &stubUserRepo{deleteRows: 0}
	svc := NewAuthServiceImpl(newTestAuthConfig(), repo)

	status, err := svc.UserDelete(context.Background(), "NOPE0000")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if status != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, status)
	}
}
