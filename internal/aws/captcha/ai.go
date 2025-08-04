package captcha

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/genai"
)

func SolveImage(base64Images []string, object string) (answer []int, err error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return
	}
	
	parts := []*genai.Part{
		genai.NewPartFromText(fmt.Sprintf("tell me the id of every image that contains \"%s\". example output [0,2,4]. MAKE SURE TO ONLY ANSWER WITH THE ARRAY AND NOTHING ELSE", object)),
	}
	
	for _, b64 := range base64Images {
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return answer, err
		}
		parts = append(parts, genai.NewPartFromBytes(data, "image/jpeg"))
	}
	
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}
	think := int32(0)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &think, // Disables thinking
			},
		},
	)
	if err != nil {
		return
	}
	
	fmt.Println("[+] Image Recognition Result:", result.Text())
	
	answer, err = JSONStringToSlice[int](result.Text())
	return
}

func JSONStringToSlice[T any](jsonStr string) ([]T, error) {
	var result []T
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}
