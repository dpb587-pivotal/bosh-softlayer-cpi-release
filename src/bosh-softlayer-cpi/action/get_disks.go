package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-softlayer-cpi/api"
	"bosh-softlayer-cpi/softlayer/virtual_guest_service"
)

type GetDisks struct {
	vmService instance.Service
}

func NewGetDisks(
	vmService instance.Service,
) GetDisks {
	return GetDisks{
		vmService: vmService,
	}
}

func (gd GetDisks) Run(vmCID VMCID) (disks []string, err error) {
	disks, err = gd.vmService.AttachedDisks(vmCID.Int())
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Finding disks for vm '%s'", vmCID.String())
	}

	return disks, nil
}