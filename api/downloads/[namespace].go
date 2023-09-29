package downloads

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
	"strings"
)

const (
	GetModificationEndpoint = "https://flintmc.net/api/client-store/get-modification/-/%s"
)

var (
	FallbackLanguage = language.German
)

type DownloadsResponse struct {
	Formatted           string `json:"formatted"`
	Rounded             string `json:"rounded"`
	FormattedAndRounded string `json:"formatted+rounded"`
	Raw                 string `json:"raw"`
}

func Downloads(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	acceptLanguageHeader := r.Header.Get("Accept-Language")
	lang := "en"

	if acceptLanguageHeader != "" {
		lang = strings.Split(strings.Split(acceptLanguageHeader, ",")[0], "-")[0]
	}

	fmt.Println(r.URL.Query())

	if lang == "" {
		lang = FallbackLanguage.String()
	}

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
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(DownloadsResponse{
		Formatted:           printer.Sprintf("%d", addon.Downloads),
		Rounded:             strconv.Itoa(rounded),
		FormattedAndRounded: printer.Sprintf("%d", rounded),
		Raw:                 strconv.Itoa(addon.Downloads),
	})
}
