package labeler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kien-wati/go-openai"
	"golang.org/x/exp/rand"
)

const (
	AZURE_OPENAI_MODEL       = "gpt4o"
	HELICONE_API_KEY         = "sk-helicone-proxy-cyfmjkq-qehutty-tk6tzsi-ebaovky-58988fca-3f84-4131-81f0-ed52d216b425"
	HELICONE_TEAM_NAME       = "DIKES"
	HELICONE_OPENAI_API_BASE = "https://gpt4-wati.openai.azure.com"
	HELICONE_BASE_URL        = "http://34.87.179.186:8787/openai"
)

type TemplateLabeler struct {
	oaClient *openai.Client
}

func NewTemplateLabeler() *TemplateLabeler {
	resource := "/deployments/gpt4o"
	heliconConfig := openai.DefaultHeliconeConfig(HELICONE_API_KEY, HELICONE_BASE_URL+resource, HELICONE_OPENAI_API_BASE, HELICONE_TEAM_NAME)
	oaClient := openai.NewClientWithConfig(heliconConfig)

	return &TemplateLabeler{
		oaClient: oaClient,
	}
}

func (s *TemplateLabeler) CreateLabelForTemplate(ctx context.Context, content string) ([]string, error) {
	sampleContent := `Sapna, the global enterprise governance, risk & compliance market size is expected to reach USD 75 billion by 2028, according to Fortune. 
		To enter this market, you can qualify as a US CPA, but that is an endlessly long route.
		You need not take that route to begin earning USD 10-20/hour. You can simply work as a remote paralegal or virtual assistant for a CPA, corporate secretary, law firm, or US companies and startups. 
		Tap on "Tell Work" to know what kind of international work you can do as CA/CS/CMA/ even B.Com graduate. 
		Even if you only charge USD 10 per hour (& I have deliberately chosen the rock bottom rate) & work for a maximum of 200 hours a month (or eight hours per day for 25 days), you make USD 2000. 
		Eventually, the goal will be to start charging USD 50 per hour as you build up a phenomenal track record & client base, so you can get closer to an income of USD 10,000/month! Join our FREE bootcamp on 3rd August @ 7 PM IST to know how to get such work.`

	sampleResponse := "english|usa|event invitation|motivational content|call-to-action content"

	marketingCategory := `new product announcement, product updates,
	event invitation, survey, feedback request, informative content,
	motivational content, appreciation content, cross-sell content, call-to-action content, subscription content,
	promotional content, reward program, seasonal content, holiday content, other`

	// Call OpenAI API to generate labels
	ccReq := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: `You are a marketing template labeler. Please labels the following template with the following format: <template_language>|<predicted_country>|<marketing_category>.`,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("You will use the following list for marketing category: %s", marketingCategory),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: sampleContent,
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: sampleResponse,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		Temperature: 0.7,
		N:           1,
	}

	ccResp, err := s.oaClient.CreateChatCompletion(ctx, ccReq)
	if err != nil {
		return nil, err
	}

	var labels []string
	if len(ccResp.Choices) != 0 {
		labels = strings.Split(ccResp.Choices[0].Message.Content, "|")
	}

	return labels, nil
}

func (s *TemplateLabeler) CreateLabelForContact(ctx context.Context) ([]string, error) {
	categoryLabels := []string{
		"new product announcement",
		"product updates",
		"event invitation",
		"survey",
		"feedback request",
		"informative content",
		"motivational content",
		"appreciation content",
		"cross-sell content",
		"call-to-action content",
		"subscription content",
		"promotional content",
		"reward program",
		"seasonal content",
		"holiday content",
		"other",
	}

	languageLabels := []string{"english", "arabic", "vietnamese", "chinese"}

	countryLabels := []string{"usa", "india", "uk", "canada",
		"australia", "singapore", "malaysia", "philippines",
		"indonesia", "thailand", "vietnam", "china", "japan",
		"south korea", "uae", "saudi arabia", "qatar", "kuwait", "bahrain",
		"oman", "egypt", "nigeria", "south africa", "kenya", "ghana", "tanzania", "uganda", "zambia", "zimbabwe", "mozambique", "angola", "morocco", "algeria", "tunisia", "libya", "sudan", "ethiopia", "somalia", "yemen", "syria", "jordan", "lebanon"}

	rand.Seed(uint64(time.Now().UnixNano()))

	// Randomly pick 1 country
	country := countryLabels[rand.Intn(len(countryLabels))]

	// Randomly pick 1 language
	language := languageLabels[rand.Intn(len(languageLabels))]

	// Randomly pick 3 categories
	rand.Shuffle(len(categoryLabels), func(i, j int) { categoryLabels[i], categoryLabels[j] = categoryLabels[j], categoryLabels[i] })
	categories := categoryLabels[:3]

	// Combine selected labels into one string array
	finalLabels := []string{country, language}
	finalLabels = append(finalLabels, categories...)
	return finalLabels, nil
}

// func CreateTemplateLabelerModel(ctx context.Context) error {
// 	bs, err := os.ReadFile("./datasets/training_prepared.jsonl") //read the content of file
// 	if err != nil {
// 		return err
// 	}

// 	// Upload JSONL file to OpenAI
// 	uploadedFile, err := oaClient.CreateFileBytes(ctx, openai.FileBytesRequest{
// 		Bytes:   bs,
// 		Purpose: openai.PurposeFineTune,
// 	})
// 	if err != nil {
// 		fmt.Printf("Upload JSONL file error: %v\n", err)
// 		return err
// 	}

// 	fineTuneReq := openai.FineTuningJobRequest{
// 		TrainingFile: uploadedFile.ID,
// 		Model:        openai.GPT4oMini20240718,
// 	}

// 	fineTuningJob, err := oaClient.CreateFineTuningJob(ctx, fineTuneReq)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	fineTuningJob, err = oaClient.RetrieveFineTuningJob(ctx, fineTuningJob.ID)
// 	if err != nil {
// 		fmt.Printf("Getting fine tune model error: %v\n", err)
// 		return err
// 	}
// 	fmt.Println(fineTuningJob.FineTunedModel)

// 	return nil
// }
