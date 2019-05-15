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

	i := 0
	getNextDiskName := func() string {
		var name string
		for {
			i++
			name = fmt.Sprintf("Hard disk %d", i)
			if !rawDiskExistsByProperty(rawDisks, "name", name) {
				break
			}
		}

		return name
	}

	for _, rawDisk := range rawDisks {
		disk := rawDisk.(map[string]interface{})

		name := disk["name"].(string)
		if len(name) < 1 {
			name = getNextDiskName()
		}

		disks = append(disks, ecloudservice.CreateVirtualMachineRequestDisk{
			Name:     name,
			Capacity: disk["capacity"].(int),
		})
	}

	return disks
}

func flattenVirtualMachineDisks(currentRawDisks []interface{}, vmDisks []ecloudservice.VirtualMachineDisk) interface{} {
	var flattenedDisks []map[string]interface{}

	getDiskByUUID := func(vmDisks []ecloudservice.VirtualMachineDisk, uuid string) *ecloudservice.VirtualMachineDisk {
		for _, vmDisk := range vmDisks {
			if vmDisk.UUID == uuid {
				return &vmDisk
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

		vmDisk := getDiskByUUID(vmDisks, currentDisk["uuid"].(string))
		if vmDisk != nil {
			flattenedDisks = append(flattenedDisks, map[string]interface{}{
				"uuid":     (*vmDisk).UUID,
				"capacity": (*vmDisk).Capacity,
			})
		}
	}

	// Next, find UUID for disks we do not have a UUID for, using the current capacity. This will use our current
	// array of flattenedDisks to determine which disks we shouldn't use
	for _, currentRawDisk := range currentRawDisks {
		currentDisk := currentRawDisk.(map[string]interface{})
		if len(currentDisk["uuid"].(string)) > 0 {
			continue
		}

		for _, vmDisk := range vmDisks {
			// TODO - check vmDisk.Type once available
			if len(vmDisk.UUID) > 0 && !diskExistsByProperty(flattenedDisks, "uuid", vmDisk.UUID) && currentDisk["capacity"].(int) == vmDisk.Capacity {
				flattenedDisks = append(flattenedDisks, map[string]interface{}{
					"uuid":     (vmDisk).UUID,
					"capacity": (vmDisk).Capacity,
				})
			}
		}
	}

	// Finally, add any new disks
	for _, vmDisk := range vmDisks {
		// TODO - check vmDisk.Type once available
		if len(vmDisk.UUID) > 0 && !diskExistsByProperty(flattenedDisks, "uuid", vmDisk.UUID) {
			flattenedDisks = append(flattenedDisks, map[string]interface{}{
				"uuid":     (vmDisk).UUID,
				"capacity": (vmDisk).Capacity,
			})
		}
	}

	return flattenedDisks
}

func resourceVirtualMachineUpdateDisk(old, new interface{}) []ecloudservice.PatchVirtualMachineRequestDisk {
	var disks []ecloudservice.PatchVirtualMachineRequestDisk

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

		if rawDiskExistsByProperty(newRawDisks, "uuid", oldDisk["uuid"]) {
			continue
		}

		disks = append(disks, ecloudservice.PatchVirtualMachineRequestDisk{
			UUID:  oldDisk["uuid"].(string),
			State: ecloudservice.PatchVirtualMachineRequestDiskStateAbsent,
		})
	}

	return disks
}

func diskExistsByProperty(disks []map[string]interface{}, name string, value interface{}) bool {
	for _, disk := range disks {
		if disk[name] == value {
			return true
		}
	}

	return false
}

func rawDiskExistsByProperty(rawDisks []interface{}, name string, value interface{}) bool {
	for _, rawDisk := range rawDisks {
		disk := rawDisk.(map[string]interface{})
		if disk[name] == value {
			return true
		}
	}

	return false
}
