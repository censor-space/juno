package quiz

type Question struct {
	ImageURL            string `json:"image_url"`
	Title               string `json:"title"`
	Text                string `json:"text"`
	Choice1             string `json:"choice_1"`
	Choice2             string `json:"choice_2"`
	Choice3             string `json:"choice_3"`
	Choice4             string `json:"choice_4"`
	NotificationMessage string `json:"notification_message"`
}
