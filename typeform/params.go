package typeform

type GetResponsesParams struct {
	FormId              string
	PageSize            int
	Since               Time
	Until               Time
	After               string
	Before              string
	IncludedResponseIds []string
	Completed           *bool
	Query               string
	Fields              []string
}
