package image

import "fmt"

var harborVersionToImageLocator = map[string]ImageLocator{
	"v1.10.0": harborV1_10_0_ImageLocator{},
}

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

func NewImageGetterImpl(registry *string, harborVersion string) (ImageGetter, error) {
	locator, existed := harborVersionToImageLocator[harborVersion]
	if !existed {
		return nil, fmt.Errorf("failed to relate images with this harbor version %s ", harborVersion)
	}
	return &ImageGetterImpl{
		locator:       locator,
		registry:      registry,
		harborVersion: harborVersion,
	}, nil
}

func (i *ImageGetterImpl) CoreImage() *string {
	return GetImage(i.registry, i.locator.CoreImage())
}

func (i *ImageGetterImpl) ChartMuseumImage() *string {
	return GetImage(i.registry, i.locator.ChartMuseumImage())
}

func (i *ImageGetterImpl) ClairImage() *string {
	return GetImage(i.registry, i.locator.ClairImage())
}

func (i *ImageGetterImpl) ClairAdapterImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) JobServiceImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) NotaryServerImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) NotarySingerImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) NotaryDBMigratorImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) PortalImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) RegistryImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

func (i *ImageGetterImpl) RegistryControllerImage() *string {
	return GetImage(i.registry, i.locator.ClairAdapterImage())
}

// ImageLocator provider method to get harbor component image.
type ImageLocator interface {
	CoreImage() *string
	ChartMuseumImage() *string
	ClairImage() *string
	ClairAdapterImage() *string
	JobServiceImage() *string
	NotaryServerImage() *string
	NotarySingerImage() *string
	NotaryDBMigratorImage() *string
	PortalImage() *string
	RegistryImage() *string
	RegistryControllerImage() *string
}

func GetImage(registry *string, image *string) *string {
	var imageAddr string
	if registry == nil {
		imageAddr = fmt.Sprintf("%s", *image)
	} else {
		imageAddr = fmt.Sprintf("%s/%s", *registry, *image)
	}
	return &imageAddr
}

func StringToStringPtr(value string) *string {
	return &value
}
