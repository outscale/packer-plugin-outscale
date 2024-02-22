## 1.1.4 (February 22, 2024)

* Update go modules
* Update documentation links for Packer integration
* Update integration Organization name
* Fix documentation on variable use inconsistent ([GH-168](https://github.com/outscale/packer-plugin-outscale/issues/168))
* Fix Broken docs links ([GH-176](https://github.com/outscale/packer-plugin-outscale/issues/176))

## 1.1.3 (October 10, 2023)

* Packer can't reach with ssh a temporary VM created in a VPC on it's public IP ([GH-159](https://github.com/outscale/packer-plugin-outscale/issues/159))
* Update docs for local plugin install https://github.com/outscale/packer-plugin-outscale/pull/161

## 1.1.2 (August 3, 2023)

* Update Packer plugin SDK to fix go-cty package ([GH-138](https://github.com/outscale/packer-plugin-outscale/issues/138))

## 1.1.1 (Jully 17, 2023)

* add frieza to CI by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/150
* update code owners by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/149
* fix error retrieving blank pwd by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/148
* fix dataSource declaration missing in main by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/147

## 1.1.0 (May 22, 2023)

* fix builders links by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/135
* fix User Agent version and update go version in goReleaser  by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/140
* update links to fix official doc by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/139
* add new parameter ProductCode in CreateImage by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/142
* add doc for new parameter productCodes by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/143

## 1.0.7 (April 12, 2023)

* fix config servers while creating an osc client

## 1.0.6 (March 31, 2023)

* upgrade sdk go from v1 to v2 by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/95
* fix redirection links by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/116
* do not make snapshot public when global permission is true by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/117
* replace with working omi for testing purposes  by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/114
* add doc for global permission option by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/125
* add datasource to output all OMI informations by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/128
* fix osc sdk go v2 migration and pointers by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/115
* remove unecessary parameters by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/130
* Increase wait between getPassword calls by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/132
* add doc for data source by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/133

## 1.0.5 (January 17, 2023)

* Release v1.0.5 by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/109

## 1.0.4 (January 17, 2023)

* Handle globalPermission: allow image to be public by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/79
* Update Example to match the new version of packer  by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/80
* update go version for CI by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/90
* fix go fmt by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/93
* Fix lint and update workflow CI to use ubuntu and go latest version  by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/92
* add example fo test most-recent-omiFilter by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/85
* add logs for mostRecent test by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/97
* fix indent workflow yaml file by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/96
* fix read me format by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/100
* Bump github.com/zclconf/go-cty from 1.10.0 to 1.11.1 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/99
* Fix make generate and go-cty  by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/106
* remove handling of spot instance as we don't have such a use case in … by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/108
* add examples with the 3 builders by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/103

## 1.0.3 (September 20, 2022)

* Bump github.com/hashicorp/packer-plugin-sdk from 0.2.3 to 0.2.7 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/43
* Validate plugin from the packer-sdc plugin-validate command by @azr in https://github.com/outscale/packer-plugin-outscale/pull/45
* Bump github.com/zclconf/go-cty from 1.9.1 to 1.10.0 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/44
* Bump github.com/hashicorp/packer-plugin-sdk from 0.2.7 to 0.2.9 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/47
* Post migration fixes by @azr in https://github.com/outscale/packer-plugin-outscale/pull/48
* add .github/releases.yml by @azr in https://github.com/outscale/packer-plugin-outscale/pull/50
* Bump github.com/hashicorp/hcl/v2 from 2.10.1 to 2.11.1 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/49
* Bump github.com/hashicorp/packer-plugin-sdk from 0.2.9 to 0.2.10 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/51
* Bump github.com/hashicorp/packer-plugin-sdk from 0.2.10 to 0.2.11 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/52
* Migrate CI from Circle-Ci to  GitHb Action by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/53
* goreleaser: auto-generate changelog by @azr in https://github.com/outscale/packer-plugin-outscale/pull/55
* Add credential checking by @jerome-jutteau in https://github.com/outscale/packer-plugin-outscale/pull/59
* Fix docs links by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/60
* Add version into User-Agent by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/56
* Explicit message on missing params by @outscale-mdr in https://github.com/outscale/packer-plugin-outscale/pull/61
* Support for the HCP Packer registry by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/70
* Remove wrong setting of Profile by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/64
* bump packer sdk version to address legacy SSH key algorithms in SSH communicator by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/78
* Bump github.com/hashicorp/hcl/v2 from 2.11.1 to 2.14.0 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/74
* Bump github.com/hashicorp/packer-plugin-sdk from 0.2.12 to 0.3.2 by @dependabot in https://github.com/outscale/packer-plugin-outscale/pull/81
* Auto-merge dependabot requests when tests are OK by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/82
* Update release.yml by @nywilken in https://github.com/outscale/packer-plugin-outscale/pull/83
* update github token in worflow ci by @outscale-hmi in https://github.com/outscale/packer-plugin-outscale/pull/84

## 1.0.2 (October 7, 2021)

Improve the shutdown Behavior handler (#41)

* Update osc-sdk-go

This update is to fix the issue that DeleteOnVmDeletion was taken into
account when submitting the request to the API

* Reject misconfiguration

We cannot set shutdown_behavior to terminate and in the same time,
delete_on_vm_deletion to true for the launching device. The reason why
is because we stop the VM and then we perform the snapshot but in this
case, the volume is no longer available.

* Change the way we wait the VM

When the VM will be terminate, the state will never be 'stopped', in
this case we need to wait for 'deleted'

* Add tests for these fixes

- builder => reject when shutdown_behavior is terminate and
delete_on_vm_deletion is true
- acc => test that the snapshot is performed when the
delete_on_vm_deletion is terminate

* Fix compilation and test errors

* Refactoring of shutdown behaviors

## 1.0.1 (July 14, 2021)

* Fix github.com/outscale/osc-sdk-go/osc revision [GH-28]

## 1.0.0 (June 14, 2021)

* Allow use of owner-alias field for owners field in Outscale builders [GH-12]
* Update to packer-plugin-sdk v0.2.3 [GH-23]
* Update to github.com/aws/aws-sdk-go v1.38.57 [GH-27]

## 0.0.1 (April 20, 2021)

* Outscale Plugin break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md).
