package models

type RAGIndexRequest struct {
	Model string `json:"model"`
}

type RAGDocumentIndexUsage struct {
	UsedBytes int64 `json:"used_bytes"`
}

type RAGDocumentIndexResponse struct {
	ID                      string                `json:"id"`
	Model                   string                `json:"model"`
	Status                  string                `json:"status"`
	ProgressPercentage      float64               `json:"progress_percentage"`
	DocumentModelIndexUsage RAGDocumentIndexUsage `json:"document_model_index_usage"`
}

type RAGDocumentIndexesResponse struct {
	Indexes []RAGDocumentIndexResponse `json:"indexes"`
}
