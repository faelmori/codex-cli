package utils

type Provider struct {
	Name    string
	BaseURL string
	EnvKey  string
}

var Providers = map[string]Provider{
	"openai": {
		Name:    "OpenAI",
		BaseURL: "https://api.openai.com/v1",
		EnvKey:  "OPENAI_API_KEY",
	},
	"openrouter": {
		Name:    "OpenRouter",
		BaseURL: "https://openrouter.ai/api/v1",
		EnvKey:  "OPENROUTER_API_KEY",
	},
	"gemini": {
		Name:    "Gemini",
		BaseURL: "https://generativelanguage.googleapis.com/v1beta/openai",
		EnvKey:  "GEMINI_API_KEY",
	},
	"ollama": {
		Name:    "Ollama",
		BaseURL: "http://localhost:11434/v1",
		EnvKey:  "OLLAMA_API_KEY",
	},
	"mistral": {
		Name:    "Mistral",
		BaseURL: "https://api.mistral.ai/v1",
		EnvKey:  "MISTRAL_API_KEY",
	},
	"deepseek": {
		Name:    "DeepSeek",
		BaseURL: "https://api.deepseek.com",
		EnvKey:  "DEEPSEEK_API_KEY",
	},
	"xai": {
		Name:    "xAI",
		BaseURL: "https://api.x.ai/v1",
		EnvKey:  "XAI_API_KEY",
	},
	"groq": {
		Name:    "Groq",
		BaseURL: "https://api.groq.com/openai/v1",
		EnvKey:  "GROQ_API_KEY",
	},
}
