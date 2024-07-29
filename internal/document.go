
package internal

type Document struct {
	Id		  int	 `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Topic     string `json:"topic"`
	Watermark string `json:"watermark,omitempty"`
}

type Filter struct {
	Key string `json:"key"`

	// If value is empty, just return everything but sorted with the key
	Value string `json:"value,omitempty"`
}

type Status string

const (
	Pending    Status = "Pending"
	Started    Status = "Started"
	InProgress Status = "InProgress"
	Finished   Status = "Finished"
	Failed     Status = "Failed"
)
