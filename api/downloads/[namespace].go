package downloads

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"math"
	"net/http"
	"strconv"
	labybadges "template-go-vercel"
)

const (
	GetModificationEndpoint = "https://flintmc.net/api/client-store/get-modification/-/%s"
)

var (
	FallbackLanguage = language.German
)

func Downloads(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	lang := r.URL.Query().Get("language")
	format := false
	round := false

	fmt.Println(r.URL.Query())

	if lang == "" {
		lang = FallbackLanguage.String()
	}

	if val, err := strconv.ParseBool(r.URL.Query().Get("format")); err == nil {
		format = val
	}
	if val, err := strconv.ParseBool(r.URL.Query().Get("round")); err == nil {
		round = val
	}

	response, err := http.Get(fmt.Sprintf(GetModificationEndpoint, namespace))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		_, _ = fmt.Fprint(w, "failed to request flintmc modification data")
		return
	}
	bytes, _ := io.ReadAll(response.Body)
	var addon = labybadges.Addon{}
	_ = json.Unmarshal(bytes, &addon)

	languageTag, _ := language.Parse(lang)
	printer := message.NewPrinter(languageTag)
	if round {
		downloadDigits := len(strconv.Itoa(addon.Downloads))
		divisor := math.Pow(10, math.Max(1, float64(downloadDigits-2)))
		rounded := int(math.Round(math.Floor(float64(addon.Downloads)/divisor)) * divisor)
		fmt.Println(rounded)
		if format {
			_, _ = printer.Fprintf(w, "%d", rounded)
			return
		}
		_, _ = fmt.Fprint(w, rounded)
		return
	}

	if format {
		fmt.Println("format")
		fmt.Println(printer.Printf("%d", addon.Downloads))
		_, _ = printer.Fprintf(w, "%d", addon.Downloads)
		return
	}

	_, _ = fmt.Fprint(w, addon.Downloads)
}
