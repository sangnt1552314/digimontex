package models

type DigimonSearchQueryParams struct {
	Name      string `json:"name"`
	Exact     string `json:"exact"`
	Level     string `json:"level"`
	Attribute string `json:"attribute"`
	XAntibody string `json:"xAntibody"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

type Digimon struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Href  string `json:"href"`
	Image string `json:"image"`
}

type DigimonResponse struct {
	Content  []Digimon `json:"content"`
	Pageable struct {
		CurrentPage    int    `json:"currentPage"`
		ElementsOnPage int    `json:"elementsOnPage"`
		TotalElements  int    `json:"totalElements"`
		TotalPages     int    `json:"totalPages"`
		PreviousPage   string `json:"previousPage"`
		NextPage       string `json:"nextPage"`
	} `json:"pageable"`
}

type DigimonDetail struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	XAntibody bool   `json:"xAntibody"`
	Images    []struct {
		Href        string `json:"href"`
		Transparent bool   `json:"transparent"`
	} `json:"images"`
	Levels []struct {
		ID    int    `json:"id"`
		Level string `json:"level"`
	} `json:"levels"`
	Types []struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"types"`
	Attributes []struct {
		ID        int    `json:"id"`
		Attribute string `json:"attribute"`
	} `json:"attributes"`
	Fields []struct {
		ID    int    `json:"id"`
		Field string `json:"field"`
		Image string `json:"image"`
	} `json:"fields"`
	ReleaseDate  string `json:"releaseDate"`
	Descriptions []struct {
		Origin      string `json:"origin"`
		Language    string `json:"language"`
		Description string `json:"description"`
	} `json:"descriptions"`
	Skills []struct {
		ID int `json:"id"`
		Skill string `json:"skill"`
		Translation string `json:"translation"`
		Description string `json:"description"`
	} `json:"skills"`
	PriorEvolutions []struct {
        ID        int    `json:"id"`
        Digimon   string `json:"digimon"`
        Condition string `json:"condition"`
        Image     string `json:"image"`
        URL       string `json:"url"`
    } `json:"priorEvolutions"`
	NextEvolutions []struct {
        ID        int    `json:"id"`
        Digimon   string `json:"digimon"`
        Condition string `json:"condition"`
        Image     string `json:"image"`
        URL       string `json:"url"`
    } `json:"nextEvolutions"`
}
