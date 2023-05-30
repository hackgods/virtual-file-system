package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user in the system.
type User struct {
	Username string
	Password string
	Role     string
}

// AuthService provides authentication services.
type AuthService struct {
	dbClient *mongo.Client
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(dbClient *mongo.Client) *AuthService {
	return &AuthService{dbClient: dbClient}
}

// Signup creates a new user account.
func (a *AuthService) Signup(username, password, role string) error {
	// Check if the username is already taken
	if a.isUsernameTaken(username) {
		return fmt.Errorf("username '%s' is already taken", username)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create the user document
	user := User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	// Insert the user document into the database
	if err := a.insertUser(user); err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	// Create the home directory for the user
	if err := createHomeDirectory(username); err != nil {
		return fmt.Errorf("failed to create home directory: %v", err)
	}

	return nil
}

// Login authenticates a user.
func (a *AuthService) Login(username, password string) (bool, error) {
	// Retrieve the user document from the database
	user, err := a.findUserByUsername(username)
	if err != nil {
		return false, fmt.Errorf("failed to find user: %v", err)
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil // Password does not match
	}

	return true, nil // Authentication successful
}

// Helper function to check if a username is already taken
func (a *AuthService) isUsernameTaken(username string) bool {
	collection := a.dbClient.Database("myfilesdb").Collection("users")

	filter := bson.M{"username": username}

	var result User
	err := collection.FindOne(nil, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false // Username is not taken
	} else if err != nil {
		log.Printf("Error checking username: %v", err)
		return true // Assume username is taken to be safe
	}

	return true // Username is taken
}

// Helper function to insert a user document into the database
func (a *AuthService) insertUser(user User) error {
	collection := a.dbClient.Database("myfilesdb").Collection("users")

	_, err := collection.InsertOne(nil, user)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to find a user by username
func (a *AuthService) findUserByUsername(username string) (*User, error) {
	collection := a.dbClient.Database("myfilesdb").Collection("users")

	filter := bson.M{"username": username}

	var user User
	err := collection.FindOne(nil, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Helper function to create the home directory for a user
func createHomeDirectory(username string) error {
	// Specify the base directory where the user directories will be created
	baseDir := "./storageData"

	// Create the user's home directory path
	homeDir := filepath.Join(baseDir, username)

	// Check if the home directory already exists
	if _, err := os.Stat(homeDir); !os.IsNotExist(err) {
		return fmt.Errorf("home directory '%s' already exists", homeDir)
	}

	// Create the user's home directory
	if err := os.Mkdir(homeDir, 0700); err != nil {
		return fmt.Errorf("failed to create home directory: %v", err)
	}

	return nil
}
