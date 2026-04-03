package otp

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

// Send sends the OTP code via Text.lk SMS API to the given phone number.
// Expects environment variables:
// - TEXTLK_API_TOKEN: API token for Text.lk
// - TEXTLK_SENDER_ID: Approved sender ID in Text.lk (e.g., TextLKDemo or your brand)
func Send(phone, code string) error {
    apiToken := os.Getenv("TEXTLK_API_TOKEN")
    sender := os.Getenv("TEXTLK_SENDER_ID")
    if apiToken == "" {
        return errors.New("TEXTLK_API_TOKEN not set")
    }
    if sender == "" {
        sender = "TextLKDemo"
    }

    // Ensure phone is in international format without +, e.g., 9471xxxxxxx
    payload := map[string]string{
        "api_token": apiToken,
        "recipient": phone,
        "sender_id": sender,
        "type":      "plain",
        "message":   fmt.Sprintf("Your verification code is %s", code),
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest(http.MethodPost, "https://app.text.lk/api/http/sms/send", bytes.NewReader(body))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    respBody, _ := io.ReadAll(resp.Body)

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("text.lk send failed: status=%d body=%s", resp.StatusCode, string(respBody))
    }
    return nil
}
