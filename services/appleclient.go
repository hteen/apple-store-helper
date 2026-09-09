package services

import (
	"log"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

var (
	appleClientOnce sync.Once
	appleHTTP       *req.Client
	appleWarmMu     sync.Mutex
	appleWarmedAt   time.Time
	appleWarmKey    string
)

func shopHTTP() *req.Client {
	appleClientOnce.Do(func() {
		jar, _ := cookiejar.New(nil)
		appleHTTP = req.C().
			ImpersonateChrome().
			SetCookieJar(jar).
			SetTimeout(12 * time.Second)
	})
	return appleHTTP
}

func warmupShopPage(pageURL, acceptLanguage string) {
	appleWarmMu.Lock()
	defer appleWarmMu.Unlock()

	if pageURL == "" {
		return
	}
	if appleWarmKey == pageURL && time.Since(appleWarmedAt) < 20*time.Minute {
		return
	}

	resp, err := shopHTTP().R().
		SetHeader("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8").
		SetHeader("accept-language", acceptLanguage).
		Get(pageURL)
	if err != nil {
		log.Println("warmup shop page failed", err)
		return
	}
	if resp.GetStatusCode() >= 400 {
		log.Println("warmup shop page status", resp.GetStatusCode(), pageURL)
		return
	}

	appleWarmKey = pageURL
	appleWarmedAt = time.Now()
}

func resetShopSession() {
	appleWarmMu.Lock()
	defer appleWarmMu.Unlock()
	appleWarmedAt = time.Time{}
	appleWarmKey = ""
	if appleHTTP != nil {
		appleHTTP.ClearCookies()
	}
}
