package e2e

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

func TestSubscribeForm_Success(t *testing.T) {
	clearDB(t)
	ctx := setupChromeContext(t)

	var result string
	err := chromedp.Run(ctx, submit("http://localhost:8080/", "test5@gmail.com", "Kyiv", "daily", &result))
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	if !strings.Contains(result, "Subscription successful") {
		t.Fatalf("Expected success message got %v", result)
	}

}

func TestSubscribeForm_AlreadySubscribed(t *testing.T) {
	clearDB(t)
	ctx := setupChromeContext(t)
	var result string
	err := chromedp.Run(ctx, submit("http://localhost:8080/", "alreadysubscribed@gmail.com", "Kyiv", "daily", &result))
	if err != nil {
		t.Fatalf("first subscribe failed: %v", err)
	}

	if !strings.Contains(result, "Subscription successful") {
		t.Fatalf("Expected success message got %v", result)
	}

	err = chromedp.Run(ctx, submit("http://localhost:8080/", "alreadysubscribed@gmail.com", "Kyiv", "daily", &result))
	if err != nil {
		t.Fatalf("second subscribe failed: %v", err)
	}

	if !strings.Contains(result, "Error") {
		t.Fatalf("Expected error message got %v", result)
	}

}

func setupChromeContext(t *testing.T) context.Context {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Headless,
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	t.Cleanup(allocCancel)

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	t.Cleanup(ctxCancel)

	return ctx
}

func submit(urlstr, email, city, frequency string, result *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible("#subscription-form", chromedp.ByID),
		chromedp.SendKeys("#email", email),
		chromedp.SendKeys("#city", city),
		chromedp.SetValue("#frequency", frequency),
		chromedp.Click("#submit-btn", chromedp.ByID),
		chromedp.WaitNotPresent(`#response-message:empty`, chromedp.ByID),
		chromedp.Text(`#response-message`, result, chromedp.ByID),
	}
}

func clearDB(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/clear", nil)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
