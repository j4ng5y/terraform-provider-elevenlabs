# elevenlabs_pvc_voice_sample

Manages a training sample for a Professional Voice Cloning (PVC) voice in ElevenLabs.

## Example Usage

```hcl
resource "elevenlabs_pvc_voice_sample" "sample_1" {
  voice_id   = elevenlabs_pvc_voice.enterprise_voice.id
  file_path  = "./training_samples/sample1.wav"
}
```

## Argument Reference

- `voice_id` (Required) - The ID of the PVC voice this sample belongs to.
- `file_path` (Required) - Local file path to the audio sample file.
- `transcription` (Optional) - The transcription of the audio sample.

## Attribute Reference

- `id` - The unique identifier for the PVC voice sample.
- `file_name` - The name of the uploaded file.
- `mime_type` - The MIME type of the audio file.
- `size_bytes` - The size of the file in bytes.
- `hash` - The hash of the file content.
- `state` - The processing state of the sample.
- `duration` - The duration of the audio sample in seconds.
- `sample_rate` - The sample rate of the audio file.
- `channels` - The number of audio channels.

## Import

You can find the ID in the ElevenLabs dashboard or retrieve it via the relevant data source in this provider.

PVC voice samples can be imported using the voice_id/sample_id format:

```bash
terraform import elevenlabs_pvc_voice_sample.sample_1 voice_123456789/sample_987654321
```