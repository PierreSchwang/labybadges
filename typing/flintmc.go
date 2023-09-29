package typing

type Addon struct {
	Id               int    `json:"id"`
	Namespace        string `json:"namespace"`
	Name             string `json:"name"`
	Featured         bool   `json:"featured"`
	Verified         bool   `json:"verified"`
	Organization     int    `json:"organization"`
	Author           string `json:"author"`
	Downloads        int    `json:"downloads"`
	DownloadString   string `json:"download_string"`
	ShortDescription string `json:"short_description"`
	Rating           struct {
		Count  int `json:"count"`
		Rating int `json:"rating"`
	} `json:"rating"`
	Changelog            string        `json:"changelog"`
	RequiredLabymodBuild int           `json:"required_labymod_build"`
	Releases             int           `json:"releases"`
	LastUpdate           int           `json:"last_update"`
	Licence              string        `json:"licence"`
	VersionString        string        `json:"version_string"`
	Meta                 []interface{} `json:"meta"`
	Dependencies         []interface{} `json:"dependencies"`
	Permissions          []string      `json:"permissions"`
	BrandImages          []struct {
		Type string `json:"type"`
		Hash string `json:"hash"`
	} `json:"brand_images"`
	Tags []int `json:"tags"`
}
