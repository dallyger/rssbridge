package forgejo

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"vnbr.de/rssbridge/internal/util"
)

type Client struct {
	Context *util.ScrapeCtx
	Domain  string
}

func (c *Client) Request(method string, endpoint string) (status int, body []byte, err error) {

	url := endpoint
	if !strings.Contains(url, "://"+c.Domain) {
		url = fmt.Sprintf("https://%s/api/v1/%s", c.Domain, endpoint)
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Context.Token))
	req.Header.Set("User-Agent", util.UserAgent())
	req.Header.Set("X-Forwarded-For", c.Context.InboundIP)
	req.Header.Set("X-Forwarded-Host", c.Context.InboundHost)
	req.Header.Set("X-Forwarded-Proto", c.Context.InboundProto)

	resp, err := http.DefaultClient.Do(req)
	slog.Debug("forgejo: sent request", "url", url, "status", resp.StatusCode)

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
