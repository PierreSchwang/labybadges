package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"labybadges/typing"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetModificationEndpoint = "https://flintmc.net/api/client-store/get-modification/%s"
	LabyBlue                = "#0a56a5"
	ShieldsEndpoint         = "https://img.shields.io/badge/%s-%s-%s?%s"
	ShieldLabelDownloads    = "Downloads"
	LocalLanguage           = "de"
)

type DownloadsResponse struct {
	Formatted           string `json:"formatted"`
	Rounded             string `json:"rounded"`
	FormattedAndRounded string `json:"formatted+rounded"`
	Raw                 string `json:"raw"`
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	namespace := query.Get("namespace")
	style := query.Get("style")
	color := query.Get("color")

	if color == "" {
		color = LabyBlue
	}

	color = url.QueryEscape(color)
	response, err := http.Get(fmt.Sprintf(GetModificationEndpoint, namespace))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = fmt.Fprint(w, "failed to request flintmc modification data")
		return
	}
	bytes, _ := io.ReadAll(response.Body)
	var addon = typing.Addon{}
	if err = json.Unmarshal(bytes, &addon); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = fmt.Fprint(w, "failed to unmarshal json response")
		return
	}

	languageTag, _ := language.Parse(LocalLanguage)
	printer := message.NewPrinter(languageTag)

	formattedMessage := strconv.Itoa(addon.Downloads)
	if style == "formatted" {
		formattedMessage = printer.Sprintf("%d", addon.Downloads)
	} else if style == "rounded" || style == "formattedrounded" || style == "roundedformatted" {
		downloadDigits := len(strconv.Itoa(addon.Downloads))
		divisor := math.Pow(10, math.Max(1, float64(downloadDigits-2)))
		rounded := int(math.Round(math.Floor(float64(addon.Downloads)/divisor)) * divisor)
		if style == "rounded" {
			formattedMessage = strconv.Itoa(rounded)
		} else {
			formattedMessage = printer.Sprintf("%d", rounded)
		}
	}

	// Build the shields.io url
	parameters := ""
	for s := range query {
		if s == "namespace" || s == "style" || s == "version" {
			continue
		}
		parameters += s + "=" + url.QueryEscape(query.Get(s))
	}

	svg, err := http.Get(fmt.Sprintf(
		ShieldsEndpoint,
		ShieldLabelDownloads,
		formattedMessage,
		color,
		parameters,
	))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = fmt.Fprint(w, "Failed to get shields icon")
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = io.Copy(w, svg.Body)
}
