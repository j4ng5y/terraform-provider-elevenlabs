# elevenlabs_convai_knowledge_base_rag_index

Manages a RAG index for a ConvAI knowledge base document.

## Example Usage

```hcl
resource "elevenlabs_convai_knowledge_base" "kb" {
  name    = "KB"
  content = "Hello"
}

resource "elevenlabs_convai_knowledge_base_rag_index" "rag" {
  documentation_id = elevenlabs_convai_knowledge_base.kb.id
  model            = "e5_mistral_7b_instruct"
}
```

## Argument Reference

- `documentation_id` (Required) - Knowledge base document ID to index.
- `model` (Required) - Embedding model to use (`e5_mistral_7b_instruct` or `multilingual_e5_large_instruct`).

## Attribute Reference

- `id` - RAG index ID for the document/model.
- `status` - Index status (`created`, `processing`, `failed`, `succeeded`, etc.).
- `progress_percentage` - Indexing progress percentage.
- `used_bytes` - Storage consumed by the index.

## Import

RAG indexes can be imported using the composite ID:

```bash
terraform import elevenlabs_convai_knowledge_base_rag_index.rag kb-123:rag-123
```
