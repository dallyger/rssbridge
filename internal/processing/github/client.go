package github

import (
	"fmt"
	"io"
	"net/http"

	"vnbr.de/rssbridge/internal/util"
)

func Request(method string, url string, token string, ctx *util.ScrapeCtx) (status int, body []byte, err error) {

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", util.UserAgent())
	req.Header.Set("X-Forwarded-For", ctx.InboundIP)
	req.Header.Set("X-Forwarded-Host", ctx.InboundHost)
	req.Header.Set("X-Forwarded-Proto", ctx.InboundProto)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	status = resp.StatusCode
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	return
}
