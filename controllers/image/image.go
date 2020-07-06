package image

import "fmt"

// ImageGetter will proxy the ImageLocator
type ImageGetter interface {
	ImageLocator
}

// ImageGetterImpl contains the concrete ImageLocator instance,
// if registry is not null, all the methods in ImageGetter will be wrapped to add the registry prefix.
type ImageGetterImpl struct {
	locator       ImageLocator
	registry      *string
	harborVersion string
}

func NewImageGetter(registry *string, harborVersion string) (ImageGetter, error) {
	// The version should be validated at the spec level to make sure it's in the supported list
	// or keep the current returns
	var locator ImageGetter
	switch harborVersion {
	case "1.10.0":
		locator = &harborV1_10_0_ImageLocator{}
	}
	if locator == nil {
		return nil, fmt.Errorf("failed to get relate images with this harbor version %s ", harborVersion)
	}
	return &ImageGetterImpl{
		locator:       locator,
		registry:      registry,
		harborVersion: harborVersion,
	}, nil
}

func (i *ImageGetterImpl) CoreImage() string {
	return GetImage(i.registry, i.locator.CoreImage())
}

func (i *ImageGetterImpl) ChartMuseumImage() string {
	return GetImage(i.registry, i.locator.ChartMuseumImage())
}

func (i *ImageGetterImpl) ClairImage() string {
	return GetImage(i.registry, i.locator.ClairImage())
}

func (i *ImageGetterImpl) ClairAdapterImage() string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) JobServiceImage() string {
	return GetImage(i.registry, i.locator.JobServiceImage())
}

func (i *ImageGetterImpl) NotaryServerImage() string {
	return GetImage(i.registry, i.locator.NotaryServerImage())
}

func (i *ImageGetterImpl) NotarySingerImage() string {
	return GetImage(i.registry, i.locator.NotarySingerImage())
}

func (i *ImageGetterImpl) NotaryDBMigratorImage() string {
	return GetImage(i.registry, i.locator.NotaryDBMigratorImage())
}

func (i *ImageGetterImpl) PortalImage() string {
	return GetImage(i.registry, i.locator.PortalImage())
}

func (i *ImageGetterImpl) RegistryImage() string {
	return GetImage(i.registry, i.locator.RegistryImage())
}

func (i *ImageGetterImpl) RegistryControllerImage() string {
	return GetImage(i.registry, i.locator.RegistryControllerImage())
}

// ImageLocator provider method to get harbor component image.
type ImageLocator interface {
	CoreImage() string
	ChartMuseumImage() string
	ClairImage() string
	ClairAdapterImage() string
	JobServiceImage() string
	NotaryServerImage() string
	NotarySingerImage() string
	NotaryDBMigratorImage() string
	PortalImage() string
	RegistryImage() string
	RegistryControllerImage() string
}

func GetImage(registry *string, image string) string {
	var imageAddr string
	if registry == nil {
		imageAddr = fmt.Sprintf("%s", image)
	} else {
		imageAddr = fmt.Sprintf("%s/%s", *registry, image)
	}
	return imageAddr
}

func String(value string) *string {
	return &value
}
