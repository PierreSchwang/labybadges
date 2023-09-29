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
	"strconv"
)

const (
	GetModificationEndpoint = "https://flintmc.net/api/client-store/get-modification/-/%s"
	LabyBlue                = "#0a56a5"
)

type DownloadsResponse struct {
	Formatted           string `json:"formatted"`
	Rounded             string `json:"rounded"`
	FormattedAndRounded string `json:"formatted+rounded"`
	Raw                 string `json:"raw"`
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	style := r.URL.Query().Get("style")

	// shields.io does not pass the header through
	// acceptLanguageHeader := r.Header.Get("Accept-Language")

	// if acceptLanguageHeader != "" {
	//	lang = strings.Split(strings.Split(acceptLanguageHeader, ",")[0], "-")[0]
	// }

	lang := "de"

	response, err := http.Get(fmt.Sprintf(GetModificationEndpoint, namespace))
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
