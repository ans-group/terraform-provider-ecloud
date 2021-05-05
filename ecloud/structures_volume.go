package ecloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

//flattenInstanceDataVolumes flattens instance volumes into a set
func flattenInstanceDataVolumes(instanceVolumes []ecloudservice.Volume) *schema.Set {
	dataVolumes := schema.NewSet(schema.HashString, []interface{}{}) 
	for _, volume := range instanceVolumes {
		if volume.Type == ecloudservice.VolumeTypeOS {
			continue
		}

		dataVolumes.Add(volume.ID)
	}

	return dataVolumes
}

func rawVolumeExistsById(rawVolumes []interface{}, value string) bool {
	for _, rawVolume := range rawVolumes {
		volume := rawVolume.(map[string]string)
		if volume["id"] == value {
			return true
		}
	}

	return false
}
