package ecloud

import "log"

func expandCreateInstanceRequestImageData(rawData map[string]interface{}) map[string]interface{} {
	imageData := make(map[string]interface{})

	for k, v := range rawData {
		imageData[k] = v
	}

	log.Printf("[INFO] Image data: [%+v]", imageData)

	return imageData
}
