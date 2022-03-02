package ecloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

//flattenInstanceDataVolumes flattens instance volumes into a set
func flattenInstanceDataVolumes(instanceVolumes []ecloudservice.Volume) *schema.Set {
	dataVolumes := schema.NewSet(schema.HashString, []interface{}{})
	for _, volume := range instanceVolumes {
		if volume.Type == ecloudservice.VolumeTypeOS || volume.IsShared == true {
			continue
		}

		dataVolumes.Add(volume.ID)
	}

	return dataVolumes
}

//rawVolumeExistsById returns true if value is in slice
func rawVolumeExistsById(rawVolumes []interface{}, value string) bool {
	for _, rawVolume := range rawVolumes {
		volume := rawVolume.(string)
		if volume == value {
			return true
		}
	}

	return false
}
