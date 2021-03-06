package action

import (
	bmscl "github.com/cloudfoundry-community/bosh-softlayer-tools/clients"
	sl "github.com/maximilien/softlayer-go/softlayer"

	. "bosh-softlayer-cpi/softlayer/common"
	slhw "bosh-softlayer-cpi/softlayer/hardware"
	slpool "bosh-softlayer-cpi/softlayer/pool"
	operations "bosh-softlayer-cpi/softlayer/pool/client/vm"
	slvm "bosh-softlayer-cpi/softlayer/vm"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

//go:generate counterfeiter -o fakes/fake_creator_provider.go . CreatorProvider
type CreatorProvider interface {
	Get(name string) VMCreator
}

//go:generate counterfeiter -o fakes/fake_deleter_provider.go . DeleterProvider
type DeleterProvider interface {
	Get(name string) VMDeleter
}

type creatorProvider struct {
	creators map[string]VMCreator
}

type deleterProvider struct {
	deleters map[string]VMDeleter
}

func NewCreatorProvider(softLayerClient sl.Client, baremetalClient bmscl.BmpClient, softLayerPoolClient operations.SoftLayerPoolClient, options ConcreteFactoryOptions, logger boshlog.Logger) CreatorProvider {
	agentEnvServiceFactory := NewSoftLayerAgentEnvServiceFactory(options.Registry, logger)

	vmFinder := slvm.NewSoftLayerFinder(
		softLayerClient,
		baremetalClient,
		agentEnvServiceFactory,
		logger,
	)

	virtualGuestCreator := slvm.NewSoftLayerCreator(
		vmFinder,
		softLayerClient,
		options.Agent,
		options.Softlayer.FeatureOptions,
		options.Registry,
		logger,
	)

	baremetalCreator := slhw.NewBaremetalCreator(
		vmFinder,
		softLayerClient,
		baremetalClient,
		options.Agent,
		logger,
	)

	poolCreator := slpool.NewSoftLayerPoolCreator(
		vmFinder,
		softLayerPoolClient,
		softLayerClient,
		options.Agent,
		options.Softlayer.FeatureOptions,
		options.Registry,
		logger,
	)

	return creatorProvider{
		creators: map[string]VMCreator{
			"virtualguest": virtualGuestCreator,
			"baremetal":    baremetalCreator,
			"pool":         poolCreator,
		},
	}
}

func (p creatorProvider) Get(name string) VMCreator {
	return p.creators[name]
}

func NewDeleterProvider(softLayerClient sl.Client, softLayerPoolClient operations.SoftLayerPoolClient, logger boshlog.Logger, vmFinder VMFinder) DeleterProvider {
	virtualGuestDeleter := slvm.NewSoftLayerVMDeleter(
		softLayerClient,
		logger,
		vmFinder,
	)

	poolDeleter := slpool.NewSoftLayerPoolDeleter(
		softLayerPoolClient,
		softLayerClient,
		logger,
	)

	return deleterProvider{
		deleters: map[string]VMDeleter{
			"virtualguest": virtualGuestDeleter,
			"pool":         poolDeleter,
		},
	}
}

func (p deleterProvider) Get(name string) VMDeleter {
	return p.deleters[name]
}
