package typeform

type GetResponsesAnswerField struct {
	Id   string    `json:"id"`
	Type FieldType `json:"type"`
	Ref  string    `json:"ref"`
}

type GetResponsesAnswer struct {
	Field   GetResponsesAnswerField `json:"field"`
	Type    AnswerType              `json:"type"`
	Text    string                  `json:"text"`
	Boolean bool                    `json:"boolean"`
	Email   string                  `json:"email"`
	Number  int                     `json:"number"`
	Choices Choices                 `json:"choices"`
	Date    string                  `json:"date"`
	Choice  Choice                  `json:"choice"`
	Url     string                  `json:"url"`
	FileUrl string                  `json:"file_url"`
	Payment Payment                 `json:"payment"`
}

type GetResponsesItemMetadata struct {
	UserAgent string `json:"user_agent"`
	Platform  string `json:"platform"`
	Referer   string `json:"referer"`
	NetworkId string `json:"network_id"`
	Browser   string `json:"browser"`
}

type GetResponsesItem struct {
	LandingId   string                   `json:"landing_id"`
	Token       string                   `json:"token"`
	ResponseId  string                   `json:"response_id"`
	LandedAt    string                   `json:"landed_at"`
	SubmittedAt string                   `json:"submitted_at"`
	Metadata    GetResponsesItemMetadata `json:"metadata"`
	Answers     []GetResponsesAnswer     `json:"answers"`
	Hidden      map[string]string
	Calculated  struct {
		Score int `json:"score"`
	}
}

type GetResponsesResponse struct {
	TotalItems int                `json:"total_items"`
	PageCount  int                `json:"page_count"`
	Items      []GetResponsesItem `json:"items"`
}
