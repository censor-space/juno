package quiz


type QuetionResult struct {
    QuetionTitle string `json:"quetion_title"`
    Choice string `json:"choice"`
    Score int64 `json:"score"`
}

type UserResult struct {
    UserId string `json:"user_id"`
    Name string  `json:"name"`
    Score int64 `json:"score"`
}

type UserResults []UserResult

func (urs UserResults) Len() int {
    return len(urs)
}

func (urs UserResults) Less(i, j int) bool {
    return urs[i].Score < urs[j].Score
}

func (urs UserResults) Swap(i, j int) {
    urs[i], urs[j] = urs[j], urs[i]
}
