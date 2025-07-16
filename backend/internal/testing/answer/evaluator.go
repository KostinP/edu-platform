package answer

import (
	"encoding/json"
	"strings"
)

type Choice struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

type ChoiceAnswer struct {
	Selected []int `json:"selected"` // индексы выбранных опций
}

func EvaluateAnswer(qType string, correctData []byte, userData []byte) (bool, string, bool) {
	switch qType {
	case "single_choice":
		var options []Choice
		var user ChoiceAnswer
		json.Unmarshal(correctData, &options)
		json.Unmarshal(userData, &user)
		if len(user.Selected) != 1 {
			return false, "Один ответ должен быть выбран", false
		}
		correctIndex := -1
		for i, opt := range options {
			if opt.Correct {
				correctIndex = i
				break
			}
		}
		return correctIndex == user.Selected[0], "", false

	case "multiple_choice":
		var options []Choice
		var user ChoiceAnswer
		json.Unmarshal(correctData, &options)
		json.Unmarshal(userData, &user)

		correctSet := map[int]bool{}
		userSet := map[int]bool{}
		for i, opt := range options {
			if opt.Correct {
				correctSet[i] = true
			}
		}
		for _, sel := range user.Selected {
			userSet[sel] = true
		}
		return mapsEqual(correctSet, userSet), "", false

	case "true_false":
		var correct bool
		var user bool
		json.Unmarshal(correctData, &correct)
		json.Unmarshal(userData, &user)
		return correct == user, "", false

	case "short_text":
		var correct string
		var user string
		json.Unmarshal(correctData, &correct)
		json.Unmarshal(userData, &user)
		return strings.TrimSpace(strings.ToLower(correct)) == strings.TrimSpace(strings.ToLower(user)), "", false

	default:
		return false, "Требуется ручная проверка", true
	}
}

func mapsEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
