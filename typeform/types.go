package typeform

import "time"

type FieldType string

const (
	DropdownField       FieldType = "dropdown"
	LegalField          FieldType = "legal"
	YesNoField          FieldType = "yes_no"
	ShortTextField      FieldType = "short_text"
	EmailField          FieldType = "email"
	NumberField         FieldType = "number"
	RatingField         FieldType = "rating"
	LongTextField       FieldType = "long_text"
	OpinionScaleField   FieldType = "opinion_scale"
	PictureChoiceField  FieldType = "picture_choice"
	DateField           FieldType = "date"
	MultipleChoiceField FieldType = "multiple_choice"
)

type AnswerType string

const (
	TextAnswer    AnswerType = "text"
	BooleanAnswer AnswerType = "boolean"
	EmailAnswer   AnswerType = "email"
	NumberAnswer  AnswerType = "number"
	UrlAnswer     AnswerType = "url"
	FileUrlAnswer AnswerType = "file_url"
	PaymentAnswer AnswerType = "payment"
	ChoiceAnswer  AnswerType = "choice"
	ChoicesAnswer AnswerType = "choices"
)

type Choices struct {
	Labels []string `json:"labels"`
}

type Choice struct {
	Label string `json:"label"`
}

type Payment struct {
	Amount  string `json:"amount"`
	Last4   string `json:"last4"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

type Time struct {
	Time  time.Time
	Valid bool
}
