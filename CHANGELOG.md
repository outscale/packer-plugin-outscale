# 📜 Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### 💥 Breaking
- (nothing yet)

### ✨ Added
- (nothing yet)

### 🛠️ Changed / Refactoring
- (nothing yet)

### 📝 Documentation
- (nothing yet)

### ⚰️ Deprecated
- (nothing yet)

### 🗑️ Removed
- (nothing yet)

### 🐛 Fixed
- (nothing yet)

### 🔒 Security
- (nothing yet)

### 📦 Dependency updates
- (nothing yet)

### 🌱 Others
- (nothing yet)

## [1.6.0-beta.1] - 2026-02-11

### ✨ Added
* 🚧 chore: migrate to sdk v3 by @jobs62 in https://github.com/outscale/packer-plugin-outscale/pull/220

### 🛠️ Changed / Refactoring
* 🔧 chore: increase default ssh timeout by @ryohkhn in https://github.com/outscale/packer-plugin-outscale/pull/222

### 📝 Documentation
* 📝 docs: Update the README by @outscale-rce in https://github.com/outscale/packer-plugin-outscale/pull/215

### 🌱 Others
* Update examples by @outscale-toa in https://github.com/outscale/packer-plugin-outscale/pull/214
* :construction_worker: install ansible for playbook provisoner by @outscale-toa in https://github.com/outscale/packer-plugin-outscale/pull/216
* :construction_worker: Fix example tests by @outscale-toa in https://github.com/outscale/packer-plugin-outscale/pull/217
* Fix examples workflows by @outscale-toa in https://github.com/outscale/packer-plugin-outscale/pull/218

## [1.5.0] - 2025-06-25

* **Feature:** Add new parameters `omi_boot_modes`  and ` boot_mode` in `CreateImage`
* **Docs:** Add documentation for new parameters `omi_boot_modes`  and ` boot_mode`

### :recycle: Refactored
- **Refactor:** Various code improvements

### :arrow_up: Updated
- **Update:** Go modules

---

## [1.4.0] - 2025-01-30

### :bug: Fixed
- **Fix:** `keypairName` parameter
- **Fix:** Prevalidation of OMI name filters

### :recycle: Refactored
- **Refactor:** Various code improvements

### :arrow_up: Updated
- **Update:** Go modules

---

## [1.3.0] - 2024-09-19

### :bug: Fixed
- **Fix:** BSU builder
- **Fix:** Security group filters
- **Fix:** Check call validity when `global_permission` is `true`

### :hammer_and_wrench: Changed
- **Improve:** Endpoint management

### :memo: Documentation
- **Docs:** Update documentation and add HCL examples

### :arrow_up: Updated
- **Update:** Go modules

---

## [1.2.0] - 2024-04-08

### :sparkles: Added
- **Feature:** `skip_create_omi` option

---

## [1.1.5] - 2024-02-23

### :bug: Fixed
- **Fix:** GoReleaser GitHub action

### :bookmark: Release
- **Release:** Fixed GoReleaser issue (v1.1.4 tag was invalid)

---

## [1.1.4] - 2024-02-22

### :hammer_and_wrench: Changed
- **Update:** Go modules
- **Update:** Integration organization name

### :memo: Documentation
- **Docs:** Update documentation links for Packer integration
- **Docs:** Fix inconsistent variable usage ([GH-168](https://github.com/outscale/packer-plugin-outscale/issues/168))
- **Docs:** Fix broken documentation links ([GH-176](https://github.com/outscale/packer-plugin-outscale/issues/176))

---

## [1.1.3] - 2023-10-10
* **Fix:** Packer can't reach with SSH a temporary VM created in a VPC on its public IP ([GH-159](https://github.com/outscale/packer-plugin-outscale/issues/159))
* **Docs:** Update documentation for local plugin install [PR-161](https://github.com/outscale/packer-plugin-outscale/pull/161)

## [1.1.2] - 2023-08-03
* **Update:** Packer plugin SDK to fix `go-cty` package ([GH-138](https://github.com/outscale/packer-plugin-outscale/issues/138))

## [1.1.1] - 2023-07-17
* **Add:** Frieza to CI by @outscale-hmi in [PR-150](https://github.com/outscale/packer-plugin-outscale/pull/150)
* **Update:** Code owners by @outscale-hmi in [PR-149](https://github.com/outscale/packer-plugin-outscale/pull/149)
* **Fix:** Error retrieving blank password by @outscale-hmi in [PR-148](https://github.com/outscale/packer-plugin-outscale/pull/148)
* **Fix:** `dataSource` declaration missing in main by @outscale-hmi in [PR-147](https://github.com/outscale/packer-plugin-outscale/pull/147)

## [1.1.0] - 2023-05-22
* **Fix:** Builders links by @outscale-hmi in [PR-135](https://github.com/outscale/packer-plugin-outscale/pull/135)
* **Fix:** User-Agent version and update Go version in GoReleaser by @outscale-hmi in [PR-140](https://github.com/outscale/packer-plugin-outscale/pull/140)
* **Update:** Links to fix official doc by @outscale-hmi in [PR-139](https://github.com/outscale/packer-plugin-outscale/pull/139)
* **Feature:** Add new parameter `ProductCode` in `CreateImage` by @outscale-hmi in [PR-142](https://github.com/outscale/packer-plugin-outscale/pull/142)
* **Docs:** Add documentation for new parameter `productCodes` by @outscale-hmi in [PR-143](https://github.com/outscale/packer-plugin-outscale/pull/143)

## [1.0.7] - 2023-04-12
* **Fix:** Config servers while creating an OSC client

## [1.0.6] - 2023-03-31
* **Upgrade:** SDK Go from v1 to v2 by @outscale-hmi in [PR-95](https://github.com/outscale/packer-plugin-outscale/pull/95)
* **Fix:** Redirection links by @outscale-hmi in [PR-116](https://github.com/outscale/packer-plugin-outscale/pull/116)
* **Fix:** Do not make snapshot public when `global_permission` is `true` by @outscale-hmi in [PR-117](https://github.com/outscale/packer-plugin-outscale/pull/117)

## [1.0.5] - 2023-01-17
* **Release:** v1.0.5 by @outscale-mdr in [PR-109](https://github.com/outscale/packer-plugin-outscale/pull/109)

## [1.0.4] - 2023-01-17
* **Feature:** Handle `globalPermission`: allow image to be public by @outscale-hmi in [PR-79](https://github.com/outscale/packer-plugin-outscale/pull/79)
* **Docs:** Update example to match the new version of Packer by @outscale-hmi in [PR-80](https://github.com/outscale/packer-plugin-outscale/pull/80)

## [1.0.3] - 2022-09-20
* **Upgrade:** Bump `github.com/hashicorp/packer-plugin-sdk` from `0.2.3` to `0.2.7` by @dependabot in [PR-43](https://github.com/outscale/packer-plugin-outscale/pull/43)

## [1.0.2] - 2021-10-07
* **Fix:** Improve shutdown behavior handler ([GH-41](https://github.com/outscale/packer-plugin-outscale/issues/41))

## [1.0.1] - 2021-07-14
* **Fix:** `github.com/outscale/osc-sdk-go/osc` revision [GH-28](https://github.com/outscale/packer-plugin-outscale/issues/28)

## [1.0.0] - 2021-06-14
* **Feature:** Allow use of `owner-alias` field for `owners` field in Outscale builders [GH-12](https://github.com/outscale/packer-plugin-outscale/issues/12)

## [0.0.1] - 2021-04-20
* **Outscale Plugin** break out from Packer core. Changes prior to break out can be found in [Packer's CHANGELOG](https://github.com/hashicorp/packer/blob/master/CHANGELOG.md).
