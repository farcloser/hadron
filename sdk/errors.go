package sdk

import "errors"

var (
	// ErrNetworkCheck indicates failure checking if Docker network exists.
	ErrNetworkCheck = errors.New("failed to check network existence")

	// ErrNetworkCreate indicates failure creating Docker network.
	ErrNetworkCreate = errors.New("failed to create network")

	// ErrVolumeCheck indicates failure checking if Docker volume exists.
	ErrVolumeCheck = errors.New("failed to check volume existence")

	// ErrVolumeCreate indicates failure creating Docker volume.
	ErrVolumeCreate = errors.New("failed to create volume")

	// ErrContainerCheck indicates failure checking if Docker container exists.
	ErrContainerCheck = errors.New("failed to check container existence")
)
