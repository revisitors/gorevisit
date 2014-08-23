package gorevisit

func SimpleBlend(input *APIMsg) (*APIMsg, error) {
	imageContent, err := DataURIToDecodedContent(input.Content.Data)
	if err != nil {
		return input, err
	}

	// TODO: add transformation
	newImageBytes := imageContent.Data

	// FIXME: fix hard coded image type
	input.Content.Data = BytesToDataURI(newImageBytes, "image/jpg")

	return input, nil
}
