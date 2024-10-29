package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Generates a JWT token for the user
func GenerateToken(userID string, role string) (string, error) {
	// Retrieves the token lifespan from environment variables (.env) and converts it to int
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	// Defines the claims for the JWT token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["role"] = role                                                           // Assuming you want to include role information
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix() // Sets expiration time

	// Creates a new JWT token with the specified claims and signing method (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with the API secret key (defined in .env) and returns it
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// // ValidateAndGenerateToken checks if the user exists and generates a token if not
// func ValidateAndGenerateToken(db database.Database, userProfile *linebot.UserProfileResponse, userID string) (bool, *string, error) {
// 	var dbUser models.User

// 	// Check if the user exists in the database
// 	err := db.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// Create a new user if not found
// 			dbUser = models.User{
// 				UserID:       userID,
// 				UserName:     userProfile.DisplayName,
// 				FirstName:    "", // LINE doesn't provide first and last names
// 				LastName:     "",
// 				LanguageCode: userProfile.Language,
// 			}

// 			// Save the new user
// 			if err := db.GetDB().Create(&dbUser).Error; err != nil {
// 				return false, nil, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token for the new user
// 			tokenStr, err := GenerateToken(userID, "user")
// 			if err != nil {
// 				return false, nil, fmt.Errorf("error generating JWT: %w", err)
// 			}

// 			return false, &tokenStr, nil
// 		}
// 		return false, nil, fmt.Errorf("error retrieving user: %w", err)
// 	}

// 	return true, nil, nil // User already exists
// }

// Checks if the provided token is valid
func TokenValid(c *gin.Context) error {
	// Extracts the token from the request
	tokenString := ExtractToken(c)

	// Parses the token and validates the signature method (check the api secret key)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Extracts the token from the request
func ExtractToken(c *gin.Context) string {
	// Checks if the token is present in the query parameters (part of url in HTTP GET request)
	token := c.Query("token")
	if token != "" {
		return token
	}

	// Extracts the token from the Authorization header if it is in the "Bearer" format
	bearerToken := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(bearerToken), "bearer ") {
		return bearerToken[7:] // Extract the token part after "Bearer "
	}
	return ""

	/* // Splitting via " ", extract token from the second string (discarded)
	bearerToken := c.Request.Header.Get("Authorization")
	bearerTokenSlice := strings.Split(bearerToken, " ") // split the header into parts using " " (the header should be "Bearer/Token token_value")
	if len(bearerTokenSlice) == 2 {                     // Check if the result slice has 2 elements exactly
		return bearerTokenSlice[1] // The second element (token) is returned
	}
	return ""*/
}

// Extracts the user ID from the token

func ExtractTokenID(c *gin.Context) (int, error) {
	// Extracts the token from the request
	tokenString := ExtractToken(c)

	// Parses the token and validates its signature method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	// Extracts the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return int(uid), nil
	}
	return 0, nil
}

// Extracts the user role from the token
/*
func ExtractTokenRole(c *gin.Context) (string, error) {
	// Extracts the token from the request
	tokenString := ExtractToken(c)

	// Parses the token and validates its signature method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	// Extracts the user role from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok {
			return "", fmt.Errorf("role not found in token")
		}
		return role, nil
	}
	return "", nil
}*/
