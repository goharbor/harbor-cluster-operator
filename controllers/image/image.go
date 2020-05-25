package image

import "fmt"

var harborVersionToImageLocator = map[string]ImageLocator{
	"v1.10.0": harborV1_10_0_ImageLocator{},
}

type ImageGetter struct {
	ImageLocator
	register      *string
	harborVersion string
}

func NewImageGetter(registry *string, harborVersion string) *ImageGetter {
	imageGetter := &ImageGetter{
		register:      registry,
		harborVersion: harborVersion,
	}
	imageGetter.ImageLocator = harborVersionToImageLocator[harborVersion]
	return imageGetter
}

func (i ImageGetter) CoreImage() string {
	return GetImage(i.register, i.ImageLocator.CoreImage())
}

func (i ImageGetter) ChartMuseumImage() string {
	return GetImage(i.register, i.ImageLocator.ChartMuseumImage())
}

func (i ImageGetter) ClairImage() string {
	return GetImage(i.register, i.ImageLocator.ClairImage())
}

func (i ImageGetter) ClairAdapterImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) JobServiceImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) NotaryServerImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) NotarySingerImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) NotaryDBMigratorImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) PortalImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) RegistryImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
}

func (i ImageGetter) RegistryControllerImage() string {
	return GetImage(i.register, i.ImageLocator.ClairAdapterImage())
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
	if registry == nil {
		return fmt.Sprintf("%s", image)
	} else {
		return fmt.Sprintf("%s/%s", registry, image)
	}
}
