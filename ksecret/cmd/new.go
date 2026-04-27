package cmd

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/borghives/kosmos-go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GeneratePayload(cmd *cobra.Command) string {
	payload, _ := cmd.Flags().GetString("payload")
	if payload == "" {
		fmt.Println("Generating random string for payload.")
		payload = GenerateRandomString(32)
	}
	return payload
}

func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)

	// Read random data from the OS into the byte slice
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate random string: %v", err)
	}

	for i, b := range bytes {
		// Use modulo to map the random byte to our charset
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// Define the "list" context command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new secret",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Action: Creating a new secret...\n")
		secretName, _ := cmd.Flags().GetString("name")
		if secretName == "" {
			log.Fatalf("Secret name is required")
		}

		force, _ := cmd.Flags().GetBool("force")

		fmt.Printf("Secret name: %s\n", secretName)

		// 1. Build the request to list secrets
		manager := kosmos.SummonSecretManager()

		//ignore error if forced
		err := manager.CreateSecret(secretName)
		if err != nil && !force {
			if status.Code(err) != codes.AlreadyExists {
				log.Fatalf("Failed Creating secret: %s, %v", secretName, err)
			}
			fmt.Printf("Secret already exists: %s\n", secretName)
			return
		}

		payload := GeneratePayload(cmd)

		err = manager.AddSecretVersion(secretName, payload)
		if err != nil {
			log.Fatalf("Failed Adding secret version: %s, %v", secretName, err)
		}
		fmt.Printf("Successfully Created secret: %s\n", secretName)
	},
}

func init() {
	newCmd.Flags().StringP("name", "n", "", "Secret name")
	newCmd.Flags().StringP("payload", "", "", "Secret payload")
	newCmd.Flags().BoolP("force", "f", false, "Force new payload if secret exists")

}
