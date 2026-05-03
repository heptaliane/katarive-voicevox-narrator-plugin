package main

import (
	"bytes"
	"context"
	"io"
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

var SupportedEncoding []pb.AudioEncoding = []pb.AudioEncoding{
	pb.AudioEncoding_AUDIO_ENCODING_WAV,
}

type VoiceVoxNarratorService struct {
	pb.UnimplementedNarratorServiceServer

	Speaker speaker.VoiceVoxHandler
	Logger  hclog.Logger
}

func (n *VoiceVoxNarratorService) Narrate(
	ctx context.Context,
	req *pb.NarrateRequest,
) (*pb.NarrateResponse, error) {
	n.Logger.Debug("Start generating narration", "output", req.GetPath())
	// TODO: implementation
	return &pb.NarrateResponse{}, nil
}
func (n *VoiceVoxNarratorService) GetNarratorServiceMetadata(
	ctx context.Context,
	req *pb.GetNarratorServiceMetadataRequest,
) (*pb.GetNarratorServiceMetadataResponse, error) {
	return &pb.GetNarratorServiceMetadataResponse{
		Name:              NAME,
		Version:           VERSION,
		SupportedEncoding: SupportedEncoding,
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

	sh, err := speaker.NewHttpVoiceVoxHandler(ctx, server)
	if err != nil {
		logger.Error(
			"Failed to establish VoiceVox server",
			"error", err,
			"server", server,
		)
		os.Exit(1)
	}
	logger.Info("Connection with VoiceVox server is established.")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: katarive.Handshake,
		Plugins: map[string]plugin.Plugin{
			"narrator": &katarive.NarratorPlugin{
				Impl: &VoiceVoxNarratorService{
					Speaker: sh,
					Logger:  logger,
				},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})
}
