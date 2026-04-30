package main

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	katarive "github.com/heptaliane/katarive-go-sdk"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
)

type VoiceVoxNarratorService struct {
	pb.UnimplementedNarratorServiceServer
	Logger hclog.Logger
}

func (n *VoiceVoxNarratorService) Narrate(
	ctx context.Context,
	req *pb.NarrateRequest,
) (*pb.NarrateResponse, error) {
	// TODO: implement this
	return nil, nil
}
func (n *VoiceVoxNarratorService) GetNarratorServiceMetadata(
	ctx context.Context,
	req *pb.GetNarratorServiceMetadataRequest,
) (*pb.GetNarratorServiceMetadataResponse, error) {
	// TODO: implement this
	return nil, nil
}

// Check NarratorServiceServer implementation
var _ pb.NarratorServiceServer = new(VoiceVoxNarratorService)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: os.Stderr,
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: katarive.Handshake,
		Plugins: map[string]plugin.Plugin{
			"narrator": &katarive.NarratorPlugin{
				Impl: &VoiceVoxNarratorService{
					Logger: logger,
				},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})
}
