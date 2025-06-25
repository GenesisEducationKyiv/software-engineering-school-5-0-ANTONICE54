package e2e

import (
	"context"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestSubscribeForm_Success(t *testing.T) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Headless,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

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
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Headless,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

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
