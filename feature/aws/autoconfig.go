package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/rs/zerolog/log"
)

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("aws load default config failed")

		return
	}

	client := imds.NewFromConfig(cfg)

	zlc := log.With().
		Str("cloud.provider", "aws")

	{
		output, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "instance-id",
		})
		if err != nil {
			log.Error().Err(err).Msg("aws get metadata failed")

			return
		}

		b, err := io.ReadAll(output.Content)
		if err != nil {
			log.Error().Err(err).Msg("aws read metadata content failed")

			return
		}

		zlc = zlc.Str("cloud.instance.id", string(b))
	}

	{
		region, err := client.GetRegion(context.TODO(), &imds.GetRegionInput{})
		if err != nil {
			log.Error().Err(err).Msg("aws get region failed")

			return
		}

		zlc = zlc.Str("cloud.region", region.Region)
	}

	{
		output, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "hostname",
		})
		if err != nil {
			log.Error().Err(err).Msg("aws get metadata failed")

			return
		}

		b, err := io.ReadAll(output.Content)
		if err != nil {
			log.Error().Err(err).Msg("aws read metadata content failed")

			return
		}

		zlc = zlc.Str("cloud.instance.name", string(b))
	}

	{
		output, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "placement/availability-zone",
		})
		if err != nil {
			log.Error().Err(err).Msg("aws get metadata failed")

			return
		}

		b, err := io.ReadAll(output.Content)
		if err != nil {
			log.Error().Err(err).Msg("aws read metadata content failed")

			return
		}

		zlc = zlc.Str("cloud.availability_zone", string(b))
	}

	{
		output, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "instance-type",
		})
		if err != nil {
			log.Error().Err(err).Msg("aws get metadata failed")

			return
		}

		b, err := io.ReadAll(output.Content)
		if err != nil {
			log.Error().Err(err).Msg("aws read metadata content failed")

			return
		}

		zlc = zlc.Str("cloud.machine.type", string(b))
	}

	{
		output, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "local-ipv4",
		})
		if err != nil {
			log.Error().Err(err).Msg("aws get metadata failed")

			return
		}

		b, err := io.ReadAll(output.Content)
		if err != nil {
			log.Error().Err(err).Msg("aws read metadata content failed")

			return
		}

		zlc = zlc.Strs("host.ip", []string{string(b)})
	}

	log.Logger = zlc.Logger()
}
