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
	GetModificationEndpoint = "https://flintmc.net/api/client-store/get-modification/%s/%s"
	LabyBlue                = "#0a56a5"
	ShieldsEndpoint         = "https://img.shields.io/badge/%s-%s-%s?%s"
	ShieldLabelDownloads    = "Downloads"
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
	version := query.Get("version")
	color := query.Get("color")

	// shields.io does not pass the header through
	// acceptLanguageHeader := r.Header.Get("Accept-Language")

	// if acceptLanguageHeader != "" {
	//	lang = strings.Split(strings.Split(acceptLanguageHeader, ",")[0], "-")[0]
	// }

	lang := "de"

	if version == "" {
		version = "1.20"
	}
	if color == "" {
		color = LabyBlue
	}
	color = url.QueryEscape(color)
	response, err := http.Get(fmt.Sprintf(GetModificationEndpoint, version, namespace))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, "failed to request flintmc modification data")
		return
	}
	bytes, _ := io.ReadAll(response.Body)
	var addon = typing.Addon{}
	_ = json.Unmarshal(bytes, &addon)

	languageTag, _ := language.Parse(lang)
	printer := message.NewPrinter(languageTag)
	downloadDigits := len(strconv.Itoa(addon.Downloads))
	divisor := math.Pow(10, math.Max(1, float64(downloadDigits-2)))
	rounded := int(math.Round(math.Floor(float64(addon.Downloads)/divisor)) * divisor)

	result := typing.ShieldResponse{
		SchemaVersion: 1,
		Label:         "Downloads",
		Message:       strconv.Itoa(addon.Downloads),
		Color:         LabyBlue,
	}
	if style == "rounded" {
		result.Message = strconv.Itoa(rounded)
	} else if style == "formatted" {
		result.Message = printer.Sprintf("%d", addon.Downloads)
	} else if style == "formattedrounded" || style == "roundedformatted" {
		result.Message = printer.Sprintf("%d", rounded)
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
		result.Message,
		color,
		parameters,
	))
	if err != nil {
		_, _ = fmt.Fprint(w, "Failed to get shields icon")
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = io.Copy(w, svg.Body)
}
