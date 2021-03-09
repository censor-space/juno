package quiz

type UserResult struct {
    UserId string `json:"user_id"`
    Name string  `json:"name"`
    Score int64 `json:"score"`
}
