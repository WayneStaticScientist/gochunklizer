package user

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func VerifyToken(token string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(os.Getenv("NEXT_PUBLIC_API_URL") + "/user/verifyToken?t=" + token)
	if err != nil {
		return fmt.Errorf("error making GET request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("there was error verifying user: %s", err)
	}
	return nil
}
