package sdk

import (
	"context"
	"errors"

	"github.com/farcloser/gluon/helpers"
	"github.com/farcloser/quark/sdk"
	"github.com/rs/zerolog/log"
)

var ErrNoSuchImage = errors.New("No such image")
var ErrFailedRetrievingCredentials = errors.New("Failed to retrieve credentials")
var ErrFailedLoadingGluon = errors.New("Failed to load gluon manifest")

func FromGluon(ctx context.Context, plan *sdk.Plan, name, filepath string) (*sdk.Image, error) {
	manifest, err := helpers.Load(filepath)
	if err != nil {
		return nil, ErrFailedLoadingGluon
	}

	for _, ref := range manifest.Credentials {
		secrets, err := sdk.GetSecret(ctx, ref, []string{"username", "password", "domain"})
		if err != nil {
			return nil, ErrFailedRetrievingCredentials
		}

		_, err = plan.Registry(secrets["domain"]).
			Username(secrets["username"]).
			Password(secrets["password"]).
			Build()
		if err != nil {
			log.Fatal().Err(err).Str("domain", secrets["domain"]).Msg("failed to build registry")
		}
	}

	for _, img := range manifest.Images {
		if img.Name == name {
			//nolint:wrapcheck
			return sdk.NewImage(img.Destination.Name).
				Domain(img.Destination.Domain).
				Version(img.Destination.Version).
				Digest(img.Destination.Digest).
				Build()
		}
	}

	return nil, ErrNoSuchImage
}
