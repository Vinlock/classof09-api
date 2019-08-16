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
	MultipleChoiceField           = "multiple_choice"
)

type AnswerType string

const (
	TextAnswer    AnswerType = "text"
	BooleanAnswer AnswerType = "boolean"
	EmailAnswer   AnswerType = "email"
	NumberAnswer  AnswerType = "number"
)

type Choices struct {
	Labels []string
}

type Choice struct {
	Label string
}

type Time struct {
	Time  time.Time
	Valid bool
}
