package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/callmeskyy111/golang-jwt-auth/internal/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// service file depends on the project requirement

type Service struct{
	repo *Repo
	jwtSecret string
}

func NewService(r *Repo, jwtSecret string)*Service{
	return &Service{
		repo: r,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct{
	Token string `json:"token"`
	User PublicUser `json:"user"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput)(AuthResult, error){
	email:=strings.ToLower(strings.TrimSpace(input.Email))
	passw:=strings.ToLower(strings.TrimSpace(input.Password))

	// checks
	if email=="" || passw==""{
		return  AuthResult{},errors.New("email and password are required!")
	}
	if len(passw)<6{
		return AuthResult{},errors.New("password must be at least 6 chars long!")
	}

	// check if user already registered/exists with the email
	_,err:=s.repo.FindByEmail(ctx,email)
	if err == nil {
		return  AuthResult{}, errors.New("Email is already registerd.. Try with another one! ⚠️")
	}

	if !errors.Is(err,mongo.ErrNoDocuments){
		return AuthResult{},err
	}

	//! HASH THE PASSWORD 🛡️
	hashBytes, err:= bcrypt.GenerateFromPassword([]byte(passw),bcrypt.DefaultCost)

	if err!=nil{
		return AuthResult{}, fmt.Errorf("⚠️ Password-Hashing failed: %w",err)
	}

	now:=time.Now().UTC()
	u:=User{
		Email: email,
		PasswordHash: string(hashBytes),
		Role: "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdUser,err:= s.repo.Create(ctx,u)
	if err!=nil{
		return AuthResult{},err
	}

	token,err:=auth.CreateToken(s.jwtSecret,createdUser.ID.Hex(), createdUser.Role)
	if err!=nil{
		return AuthResult{},err
	}

	return AuthResult{
		Token: token,
		User:ToPublic(createdUser),
	},nil
}

func (s *Service) Login(ctx context.Context, input LoginInput)(AuthResult, error){
	email:=strings.ToLower(strings.TrimSpace(input.Email))
	passw:=strings.ToLower(strings.TrimSpace(input.Password))

	// checks
	if email=="" || passw==""{
		return  AuthResult{},errors.New("email and password are required!")
	}
	if len(passw)<6{
		return AuthResult{},errors.New("password must be at least 6 chars long!")
	}

	// Check if user is valid/already signed-in or not
	u,err:=s.repo.FindByEmail(ctx,email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return AuthResult{},errors.New("⚠️ Invalid credentials!")
		}
		return AuthResult{},err
	}

	if err:=bcrypt.CompareHashAndPassword([]byte(u.PasswordHash),[]byte(passw));err!=nil{
		return AuthResult{},errors.New("Invalid credentials or wrong password!")
	}

	// Create token (for register, token-creation is optional, here token-creation is necessary!)
	tkn,err:=auth.CreateToken(s.jwtSecret,u.ID.Hex(),u.Role)
	if err != nil {
		return AuthResult{},err
	}

	return AuthResult{
		Token: tkn,
		User: ToPublic(u),
	},nil
}