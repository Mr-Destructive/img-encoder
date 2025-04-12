package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

func Main(Context openruntimes.Context) openruntimes.Response {
	var imageURL string

	switch Context.Req.Method {
	case "GET":
		query := Context.Req.Query
		if url, ok := query["url"]; ok && url != "" {
			imageURL = url
		} else {
			return Context.Res.Json(map[string]string{
				"error": "Missing 'url' query parameter.",
			})
		}

	case "POST":
		var body struct {
			URL string `json:"url"`
		}

		if err := Context.Req.BodyJson(&body); err != nil || body.URL == "" {
			return Context.Res.Json(map[string]string{
				"error": "Invalid JSON body. Expected { \"url\": \"...\" }",
			})
		}

		imageURL = body.URL

	default:
		return Context.Res.Json(map[string]string{
			"error": "Method not allowed. Use GET with ?url=... or POST with JSON body { \"url\": \"...\" }",
		})
	}

	client := http.Client{}
	resp, err := client.Get(imageURL)
	Context.Log(resp.StatusCode)
	Context.Log(resp)
	Context.Log(err)
	if err != nil || resp.StatusCode != 200 {
		return Context.Res.Json(map[string]string{
			"error": fmt.Sprintf("Failed to fetch image: %v", err),
		})
	}
	defer resp.Body.Close()

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Context.Res.Json(map[string]string{
			"error": "Failed to read image data.",
		})
	}

	base64Str := base64.StdEncoding.EncodeToString(imgBytes)
	return Context.Res.Binary([]byte(base64Str))
}
