package ecloud

func expandVirtualMachineSSHKeys(raw []interface{}) []string {
	sshKeys := make([]string, len(raw))
	for i, v := range raw {
		sshKey := v.(string)
		sshKeys[i] = sshKey
	}

	return sshKeys
}
