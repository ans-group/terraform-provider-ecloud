package lock

import (
	"github.com/ukfast/terraform-provider-ecloud/pkg/kvmutex"
)

var kvMutex *kvmutex.KVMutex = kvmutex.NewKVMutex()

func LockResource(resourceID string) func() {
	return kvMutex.Lock(resourceID)
}
func UnlockResource(resourceID string) {
	kvMutex.Unlock(resourceID)
}
