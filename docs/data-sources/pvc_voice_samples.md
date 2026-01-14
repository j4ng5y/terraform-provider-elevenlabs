# elevenlabs_pvc_voice_samples

Lists training samples for a Professional Voice Cloning (PVC) voice in ElevenLabs.

## Example Usage

```hcl
data "elevenlabs_pvc_voice_samples" "enterprise_samples" {
  voice_id = elevenlabs_pvc_voice.enterprise_voice.id
}
```

## Argument Reference

- `voice_id` (Required) - The ID of the PVC voice to list samples for.

## Attribute Reference

- `samples` - A list of PVC voice samples, each containing:
  - `sample_id` - The unique identifier for the sample.
  - `file_name` - The name of the sample file.
  - `mime_type` - The MIME type of the sample.
  - `size_bytes` - The size of the sample in bytes.
  - `hash` - The hash of the sample content.
  - `state` - The processing state of the sample.
  - `transcription` - The transcription of the sample.
  - `duration` - The duration of the sample in seconds.
  - `sample_rate` - The sample rate of the audio.
  - `channels` - The number of audio channels.