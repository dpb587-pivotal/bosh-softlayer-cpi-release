# Testing the CPI
=================

## Normal testing

Generally the CPI can be tested by using the CPI release and deploying that using bosh lite using vagrant. See the [bosh-softlayer-cpi-release](http://github.com/maximilien/bosh-softlayer-cpi-release/README.md) for details.

Once the release is successfully deployed, then you can test changes to the CPI by doing:

1. update the CPI code under `src/bosh-softlayer-cpi` using:

```
cd src/bosh-softlayer-cpi
git pull
```

2. update your vagrant VM using:

```
vagrant provision
```

3. executing some bosh commands, e.g.,

To test create_stemcell method, using: `bosh upload stemcell path/to/your/stemcell`

To test create_vm method, using: `bosh deploy` using the [dummy release](https://github.com/pivotal-cf-experimental/dummy-boshrelease)

## Manual testing