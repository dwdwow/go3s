package go3s

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestX(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*100), 5) // 1 request per second with a burst of 5
	for i := 0; i < 10; i++ {
		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Request", i, "processed at", time.Now())
	}
}
