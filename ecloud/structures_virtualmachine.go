package ecloud

import (
	"fmt"

	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func expandVirtualMachineSSHKeys(raw []interface{}) []string {
	sshKeys := make([]string, len(raw))
	for i, v := range raw {
		sshKey := v.(string)
		sshKeys[i] = sshKey
	}

	return sshKeys
}

func expandCreateVirtualMachineRequestDisks(rawDisks []interface{}) []ecloudservice.CreateVirtualMachineRequestDisk {
	var disks []ecloudservice.CreateVirtualMachineRequestDisk
	var i int
	for _, v := range rawDisks {
		i++
		disk := v.(map[string]interface{})
		disks = append(disks, ecloudservice.CreateVirtualMachineRequestDisk{
			Name:     fmt.Sprintf("Hard disk %d", i),
			Capacity: disk["capacity"].(int),
		})
	}

	return disks
}

func flattenVirtualMachineDisks(currentRawDisks []interface{}, new []ecloudservice.VirtualMachineDisk) interface{} {
	var newDisks []map[string]interface{}

	newDiskExists := func(uuid string) bool {
		for _, v := range newDisks {
			if v["uuid"].(string) == uuid {
				return true
			}
		}

		return false
	}

	getDiskByUUID := func(disks []ecloudservice.VirtualMachineDisk, uuid string) *ecloudservice.VirtualMachineDisk {
		for _, disk := range disks {
			if disk.UUID == uuid {
				return &disk
			}
		}

		return nil
	}

	// First, find all disks that we have a UUID for
	for _, currentRawDisk := range currentRawDisks {
		currentDisk := currentRawDisk.(map[string]interface{})
		if len(currentDisk["uuid"].(string)) < 1 {
			continue
		}

		newDisk := getDiskByUUID(new, currentDisk["uuid"].(string))
		if newDisk != nil {
			newDisks = append(newDisks, map[string]interface{}{
				"uuid":     (*newDisk).UUID,
				"capacity": (*newDisk).Capacity,
			})
		}
	}

	// Next, find UUID for disks we do not have a UUID for, using the current capacity. This will use our current
	// array of newDisks to determine which disks we shouldn't use
	for _, currentRawDisk := range currentRawDisks {
		currentDisk := currentRawDisk.(map[string]interface{})
		if len(currentDisk["uuid"].(string)) > 0 {
			continue
		}

		for _, newDisk := range new {
			if !newDiskExists(newDisk.UUID) && currentDisk["capacity"].(int) == newDisk.Capacity {
				newDisks = append(newDisks, map[string]interface{}{
					"uuid":     (newDisk).UUID,
					"capacity": (newDisk).Capacity,
				})
			}
		}
	}

	// Finally, add any new disks
	for _, newDisk := range new {
		if !newDiskExists(newDisk.UUID) {
			newDisks = append(newDisks, map[string]interface{}{
				"uuid":     (newDisk).UUID,
				"capacity": (newDisk).Capacity,
			})
		}
	}

	return newDisks
}

func resourceVirtualMachineUpdateDisk(old, new interface{}) []ecloudservice.PatchVirtualMachineRequestDisk {
	var disks []ecloudservice.PatchVirtualMachineRequestDisk

	rawDiskExists := func(rawDisks []interface{}, uuid string) bool {
		for _, rawDisk := range rawDisks {
			disk := rawDisk.(map[string]interface{})
			if disk["uuid"].(string) == uuid {
				return true
			}
		}

		return false
	}

	oldRawDisks := old.([]interface{})
	newRawDisks := new.([]interface{})

	// First, find present disks. At this point, we take all new disks and mark them as present, with the new capacity
	for _, newRawDisk := range newRawDisks {
		newDisk := newRawDisk.(map[string]interface{})

		disks = append(disks, ecloudservice.PatchVirtualMachineRequestDisk{
			Capacity: newDisk["capacity"].(int),
			UUID:     newDisk["uuid"].(string),
			State:    ecloudservice.PatchVirtualMachineRequestDiskStatePresent,
		})
	}

	// Next, find absent disks. This will check whether each old disk exists in the new state, and mark them as
	// absent if not
	for _, oldRawDisk := range oldRawDisks {
		oldDisk := oldRawDisk.(map[string]interface{})

		if rawDiskExists(newRawDisks, oldDisk["uuid"].(string)) {
			continue
		}

		disks = append(disks, ecloudservice.PatchVirtualMachineRequestDisk{
			UUID:  oldDisk["uuid"].(string),
			State: ecloudservice.PatchVirtualMachineRequestDiskStateAbsent,
		})
	}

	return disks
}
