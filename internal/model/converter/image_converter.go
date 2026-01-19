package converter

import "pura-agung-kertajaya-backend/internal/model"

func ToImageVariants(images map[string]string) model.ImageVariants {
	if images == nil {
		return model.ImageVariants{}
	}
	return model.ImageVariants{
		Blur:   images["blur"],
		Avatar: images["avatar"],
		Xs:     images["xs"],
		Sm:     images["sm"],
		Md:     images["md"],
		Lg:     images["lg"],
		Xl:     images["xl"],
		TwoXl:  images["2xl"],
		Fhd:    images["fhd"],
	}
}
