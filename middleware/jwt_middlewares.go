package middleware

import (
	"net/http"

	"crossplatform_chatbot/utils/token"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware validates the JWT token
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) { // returning anonymous function as gin.HandlerFunc
		err := token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// extract user id from token
		userID, err := token.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// set the id in context
		c.Set("user_id", userID)
		c.Next() // used within middleware to continue to the remaining middleware and handlers in the request processing chain.
	}
}

// Checks if the user's role matches any of the provided roles.
/*func AuthorizeRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, err := token.ExtractTokenRole(c)
		if err != nil || userRole != role { // checks if userRole from token is the same as input role
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access"}) // 403
			c.Abort()
			return
		}
		c.Next()
	}
}*/
