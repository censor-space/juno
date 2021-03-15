package quiz

type UserResult struct {
    UserId string `json:"user_id"`
    Name string  `json:"name"`
    Score int64 `json:"score"`
}


type QuetionResult struct {
    QuetionTitle string `json:"quetion_title"`
    Choice string `json:"choice"`
    Score int64 `json:"score"`
}
