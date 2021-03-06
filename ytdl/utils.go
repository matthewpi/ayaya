package ytdl

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func reverseStringSlice(s []string) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func interfaceToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func (c *Client) httpGet(cx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(cx)
	// Youtube responses depend on language and user agent
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:70.0) Gecko/20100101 Firefox/70.0")

	return c.HTTPClient.Do(req)
}

func (c *Client) httpGetAndCheckResponse(cx context.Context, url string) (*http.Response, error) {
	resp, err := c.httpGet(cx, url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return resp, nil
}

func (c *Client) httpGetAndCheckResponseReadBody(cx context.Context, url string) ([]byte, error) {
	resp, err := c.httpGetAndCheckResponse(cx, url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
