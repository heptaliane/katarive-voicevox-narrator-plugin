package main

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	katarive "github.com/heptaliane/katarive-go-sdk"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"

	"github.com/heptaliane/katarive-voicevox-narrator-plugin/internal/speaker"
)

const NAME string = "voicevox"
const VERSION string = "v1"
const DEFAULT_SERVER string = "http://localhost:50021"
const ENV_SERVER string = "KATARIVE_VOICEVOX_SERVER"
const CACHE_DIR string = ".cache/katarive-voicevox"

var SupportedEncoding []pb.AudioEncoding = []pb.AudioEncoding{
	pb.AudioEncoding_AUDIO_ENCODING_WAV,
	pb.AudioEncoding_AUDIO_ENCODING_MP3,
	pb.AudioEncoding_AUDIO_ENCODING_M4A,
}

type VoiceVoxNarratorService struct {
	pb.UnimplementedNarratorServiceServer

	Speaker speaker.NarrationGenerator
	Logger  hclog.Logger
}

func (n *VoiceVoxNarratorService) Narrate(
	ctx context.Context,
	req *pb.NarrateRequest,
) (*pb.NarrateResponse, error) {
	n.Logger.Debug("Start generating narration", "output", req.GetPath())
	err := n.Speaker.Do(
		ctx,
		req.GetPath(),
		req.GetText(),
		speaker.WithEncoding(req.GetEncoding()),
		speaker.WithSpeakerId(int(req.GetSpeakerId())),
	)
	return &pb.NarrateResponse{}, err
}
func (n *VoiceVoxNarratorService) GetNarratorServiceMetadata(
	ctx context.Context,
	req *pb.GetNarratorServiceMetadataRequest,
) (*pb.GetNarratorServiceMetadataResponse, error) {
	return &pb.GetNarratorServiceMetadataResponse{
		Name:              NAME,
		Version:           VERSION,
		SupportedEncoding: SupportedEncoding,
		Speakers:          n.Speaker.Speakers(),
	}, nil
}

// Check NarratorServiceServer implementation
var _ pb.NarratorServiceServer = new(VoiceVoxNarratorService)

func GetenvWithDefault(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func main() {
	ctx := context.Background()

	server := GetenvWithDefault(ENV_SERVER, DEFAULT_SERVER)

	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: os.Stderr,
	})

	handler, err := speaker.NewHttpVoiceVoxHandler(ctx, server)
	if err != nil {
		logger.Error(
			"Failed to establish VoiceVox server",
			"error", err,
			"server", server,
		)
		os.Exit(1)
	}
	logger.Info("Connection with VoiceVox server is established.")

	os.MkdirAll(CACHE_DIR, 0755)
	narrator := &speaker.ChunkedNarrationGenerator{
		Handler:  handler,
		Chunker:  new(speaker.LineBreakTextChunker),
		CacheDir: CACHE_DIR,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: katarive.Handshake,
		Plugins: map[string]plugin.Plugin{
			"narrator": &katarive.NarratorPlugin{
				Impl: &VoiceVoxNarratorService{
					Speaker: narrator,
					Logger:  logger,
				},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})
}
