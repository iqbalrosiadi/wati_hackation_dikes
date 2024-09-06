package main

import (
	"context"

	"github.com/kien-wati/go-openai"
)

const (
	AZURE_OPENAI_MODEL       = "gpt4o"
	HELICONE_API_KEY         = "sk-helicone-proxy-cyfmjkq-qehutty-tk6tzsi-ebaovky-58988fca-3f84-4131-81f0-ed52d216b425"
	HELICONE_TEAM_NAME       = "DIKES"
	HELICONE_OPENAI_API_BASE = "https://gpt4-wati.openai.azure.com"
	HELICONE_BASE_URL        = "http://34.87.179.186:8787/openai"
)

var (
	oaClient *openai.Client
)

func init() {
	// For chat completion
	resource := "/deployments/gpt4o"
	heliconConfig := openai.DefaultHeliconeConfig(HELICONE_API_KEY, HELICONE_BASE_URL+resource, HELICONE_OPENAI_API_BASE, HELICONE_TEAM_NAME)
	oaClient = openai.NewClientWithConfig(heliconConfig)
}

func main() {
	test := "أهلًا Othmane👋 تبحرتي ولا باقي 🌞 إلا باقي حاول تستافد من هاد شهر لباقي من صيف و تلاقى مع لحباب و صحاب 🫂 و أحسن حاجة تقدر دير دابا هي تكمل Application ديالك فبرنامج Virtual Assistant و ضمن بلاصتك معنا و دوز هاد شهر هاني مهني و نتا عارف راسك فشهر 9 أتبدا أحسن Program لغادي يعاونوك باش تبدا مسارك المهني كمساعد إفتراضي ف 8 أسابيع فقط 🔥 و قول ماكنفكروش ليك 😉 و دابا ما عليك غير تبرك على Button below باش تكمل التسجيل ديالك فبرنامج 👇 نهارك مبروك 😁"
	resp, err := CreateLabelForTemplate(test)
	if err != nil {
		panic(err)
	}
	println(resp)
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

func CreateLabelForTemplate(content string) (string, error) {
	sampleContent := `Sapna, the global enterprise governance, risk & compliance market size is expected to reach USD 75 billion by 2028, according to Fortune. To enter this market, you can qualify as a US CPA, but that is an endlessly long route. You need not take that route to begin earning USD 10-20/hour. You can simply work as a remote paralegal or virtual assistant for a CPA, corporate secretary, law firm, or US companies and startups. Tap on "Tell Work" to know what kind of international work you can do as CA/CS/CMA/ even B.Com graduate. Even if you only charge USD 10 per hour (& I have deliberately chosen the rock bottom rate) & work for a maximum of 200 hours a month (or eight hours per day for 25 days), you make USD 2000. Eventually, the goal will be to start charging USD 50 per hour as you build up a phenomenal track record & client base, so you can get closer to an income of USD 10,000/month! Join our FREE bootcamp on 3rd August @ 7 PM IST to know how to get such work.`
	sampleResponse := "english|bootcamp promotion|remote paralegal/virtual assistant|career development"

	// Call OpenAI API to generate labels
	ccReq := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a template labeler. Please label the following template.",
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

	ccResp, err := oaClient.CreateChatCompletion(context.Background(), ccReq)
	if err != nil {
		return "", err
	}

	return ccResp.Choices[0].Message.Content, nil
}
