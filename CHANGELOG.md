# Changelog

## [1.3.0](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.2.2...v1.3.0) (2026-09-02)


### Features

* **sso-suite:** add hide_sso to sso_suite_settings ([#364](https://github.com/jamescrowley321/terraform-provider-descope/issues/364)) ([db2d050](https://github.com/jamescrowley321/terraform-provider-descope/commit/db2d0507d83210344fb0648d9e8ea27bb7803198))
* **sso-suite:** add hide_xaa to sso_suite_settings ([#367](https://github.com/jamescrowley321/terraform-provider-descope/issues/367)) ([7744a2a](https://github.com/jamescrowley321/terraform-provider-descope/commit/7744a2aa5bc4a11a66dd337f02f133c4b33d291b))
* **sso-suite:** replace hide_xaa with show_xaa ([#374](https://github.com/jamescrowley321/terraform-provider-descope/issues/374)) ([e067adb](https://github.com/jamescrowley321/terraform-provider-descope/commit/e067adbf69b8261f677593c16befe855f40181f0))
* **sso:** add allow_merge_users_with_multiple_tenants SSO attribute ([#332](https://github.com/jamescrowley321/terraform-provider-descope/issues/332)) ([49544bf](https://github.com/jamescrowley321/terraform-provider-descope/commit/49544bf891509d8f0c4bccfc71f54bcc3c364749))
* **sso:** add SSO Suite OIDC login ID attribute setting ([#373](https://github.com/jamescrowley321/terraform-provider-descope/issues/373)) ([973ed7f](https://github.com/jamescrowley321/terraform-provider-descope/commit/973ed7f80a93c966d9635aa08d7832efe24ffe05))


### Bug Fixes

* bump Go toolchain to 1.26.6 to patch stdlib vulnerabilities ([705c88f](https://github.com/jamescrowley321/terraform-provider-descope/commit/705c88ff8ebb2748b6218f0c6e13f7c61681815f))
* bump Go toolchain to 1.26.6 to patch stdlib vulnerabilities ([67fa8be](https://github.com/jamescrowley321/terraform-provider-descope/commit/67fa8bee54edc8e991be2878a1f095b600fe5d28))
* **deps:** update module github.com/descope/go-sdk to v1.29.0 ([#362](https://github.com/jamescrowley321/terraform-provider-descope/issues/362)) ([88089ab](https://github.com/jamescrowley321/terraform-provider-descope/commit/88089ab56e83ad96d71326f877915b3d3c48be59))
* **deps:** update module github.com/descope/go-sdk to v1.29.1 ([#366](https://github.com/jamescrowley321/terraform-provider-descope/issues/366)) ([de4876f](https://github.com/jamescrowley321/terraform-provider-descope/commit/de4876f444502d975dcd52557648841b0740554d))
* **deps:** update module github.com/descope/go-sdk to v1.30.0 ([#369](https://github.com/jamescrowley321/terraform-provider-descope/issues/369)) ([e499533](https://github.com/jamescrowley321/terraform-provider-descope/commit/e4995337d7a7c9cb4da9b94112f1dadb481f23ec))
* **deps:** update module github.com/descope/go-sdk to v1.31.0 ([#372](https://github.com/jamescrowley321/terraform-provider-descope/issues/372)) ([e1dd103](https://github.com/jamescrowley321/terraform-provider-descope/commit/e1dd103281aadff7ac70117cae03d64d32909fb6))
* **deps:** update module github.com/descope/go-sdk to v1.32.0 ([#376](https://github.com/jamescrowley321/terraform-provider-descope/issues/376)) ([1781c6d](https://github.com/jamescrowley321/terraform-provider-descope/commit/1781c6d05051a8c2c1c22ba351c3554f32351d61))
* **deps:** update module github.com/stretchr/testify to v1.12.0 ([#368](https://github.com/jamescrowley321/terraform-provider-descope/issues/368)) ([3caeac5](https://github.com/jamescrowley321/terraform-provider-descope/commit/3caeac5c7a4e1a1aae6b4fd8ca7cab418ed5eb07))
* **oauth:** validate system provider own-account attributes on initial create ([#360](https://github.com/jamescrowley321/terraform-provider-descope/issues/360)) ([3102f1c](https://github.com/jamescrowley321/terraform-provider-descope/commit/3102f1c6050db8e2f188bc775ce6edaf53218262))

## [1.2.2](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.2.1...v1.2.2) (2026-08-03)


### Bug Fixes

* **deps:** bump google.golang.org/grpc to v1.82.1 (GO-2026-6061) ([057096c](https://github.com/jamescrowley321/terraform-provider-descope/commit/057096c1c5b5690687f43be3e318ad12f6b16c39))
* **deps:** bump grpc to v1.82.1 (GO-2026-6061) ([45db743](https://github.com/jamescrowley321/terraform-provider-descope/commit/45db743af4ee551e7984747c66b748d8a9f46f9e))

## [1.2.1](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.2.0...v1.2.1) (2026-06-29)


### Miscellaneous Chores

* release v1.2.1 ([#172](https://github.com/jamescrowley321/terraform-provider-descope/issues/172)) ([ad08d0f](https://github.com/jamescrowley321/terraform-provider-descope/commit/ad08d0f352080d1a19a7010b8c8f979e71b11e34))

## [1.2.0](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.5...v1.2.0) (2026-06-29)


### Features

* **applications:** add reply_allowed_callback_urls to WS-Fed ([#307](https://github.com/jamescrowley321/terraform-provider-descope/issues/307)) ([0cd2fd7](https://github.com/jamescrowley321/terraform-provider-descope/commit/0cd2fd749dcc6d287e8d187d36dc2590f617bfd6))
* **authentication:** disallowed_characters + disallow_email_match ([#301](https://github.com/jamescrowley321/terraform-provider-descope/issues/301)) ([db5492b](https://github.com/jamescrowley321/terraform-provider-descope/commit/db5492b2c64b2857a4fed43fe5ca921b28aa66cb))
* **connectors:** add Outbound SCIM connector ([#302](https://github.com/jamescrowley321/terraform-provider-descope/issues/302)) ([097dace](https://github.com/jamescrowley321/terraform-provider-descope/commit/097dace86cedbc6048096f6f6e98565b737bee28))
* **isolation:** add tenant_user_isolation to descope_project settings ([#285](https://github.com/jamescrowley321/terraform-provider-descope/issues/285)) ([9d2f809](https://github.com/jamescrowley321/terraform-provider-descope/commit/9d2f8094bdad1b7bc9a0245d3a2ce78a0d5c7c76))
* **passkeys:** add configurable display_name attribute ([#315](https://github.com/jamescrowley321/terraform-provider-descope/issues/315)) ([242bb2c](https://github.com/jamescrowley321/terraform-provider-descope/commit/242bb2ce77b9f13a29b4fdcebf5832b77a4acf03))
* **password:** add setting for any case letter requirement ([#292](https://github.com/jamescrowley321/terraform-provider-descope/issues/292)) ([5b3e224](https://github.com/jamescrowley321/terraform-provider-descope/commit/5b3e224de0a7d96247672b52a15ccdb5db962adc))
* per-app roles and permissions for federated apps ([#298](https://github.com/jamescrowley321/terraform-provider-descope/issues/298)) ([64c4737](https://github.com/jamescrowley321/terraform-provider-descope/commit/64c473788f6f4369be1094c8cfb44ce9412a837f))
* SSO OIDC dedicated-client attributes for descope_project + force_pkce ([#313](https://github.com/jamescrowley321/terraform-provider-descope/issues/313)) ([74fb320](https://github.com/jamescrowley321/terraform-provider-descope/commit/74fb320c4670d9236bc44e111d73ecee72fba0b4))
* **sso:** add WS-Fed SSO application resource ([#282](https://github.com/jamescrowley321/terraform-provider-descope/issues/282)) ([822c470](https://github.com/jamescrowley321/terraform-provider-descope/commit/822c47070967e2f0395addf3ac2d51c4a7b610eb))


### Bug Fixes

* **deps:** update module github.com/descope/go-sdk to v1.15.0 ([#279](https://github.com/jamescrowley321/terraform-provider-descope/issues/279)) ([560b315](https://github.com/jamescrowley321/terraform-provider-descope/commit/560b315a4e76d5382249e7479e8894a2b177e20b))
* **deps:** update module github.com/descope/go-sdk to v1.23.0 ([#284](https://github.com/jamescrowley321/terraform-provider-descope/issues/284)) ([a947468](https://github.com/jamescrowley321/terraform-provider-descope/commit/a9474685a2c081a9ca4b25a46abe762ea43b5b58))
* **deps:** update module github.com/descope/go-sdk to v1.24.0 ([#318](https://github.com/jamescrowley321/terraform-provider-descope/issues/318)) ([b53ed03](https://github.com/jamescrowley321/terraform-provider-descope/commit/b53ed03bf0b539c24f701ace62ef1368646fc98a))

## [1.1.5](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.4...v1.1.5) (2026-03-30)


### Bug Fixes

* add ralph/claude runtime state to gitignore ([#136](https://github.com/jamescrowley321/terraform-provider-descope/issues/136)) ([867bec5](https://github.com/jamescrowley321/terraform-provider-descope/commit/867bec52ab0e7ca8ed7cfd01827c85253e1c96e4))

## [1.1.4](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.3...v1.1.4) (2026-03-29)


### Bug Fixes

* **test:** tolerate last-project deletion error in Destroy ([#129](https://github.com/jamescrowley321/terraform-provider-descope/issues/129)) ([10291ad](https://github.com/jamescrowley321/terraform-provider-descope/commit/10291ad6d4da4fc9bf4d2a7e5f64ca88da0338b4))

## [1.1.3](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.2...v1.1.3) (2026-03-29)


### Bug Fixes

* **ci:** add actions:write permission for workflow dispatch ([#127](https://github.com/jamescrowley321/terraform-provider-descope/issues/127)) ([0bf0e55](https://github.com/jamescrowley321/terraform-provider-descope/commit/0bf0e557d9504e71fe1df222b1326c54f50c7e3f))

## [1.1.2](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.1...v1.1.2) (2026-03-29)


### Bug Fixes

* **ci:** add --repo flag to gh workflow dispatch in release-please ([#125](https://github.com/jamescrowley321/terraform-provider-descope/issues/125)) ([eda0518](https://github.com/jamescrowley321/terraform-provider-descope/commit/eda051839c10a59c6c871f658b303b1f5e943c79))

## [1.1.1](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.1.0...v1.1.1) (2026-03-29)


### Bug Fixes

* **ci:** trigger goreleaser via workflow_dispatch from release-please ([#122](https://github.com/jamescrowley321/terraform-provider-descope/issues/122)) ([2eceaa5](https://github.com/jamescrowley321/terraform-provider-descope/commit/2eceaa55401cdf43c928ee3d926df783779dfb98))
* **ci:** use tag ref for SLSA provenance generator ([#124](https://github.com/jamescrowley321/terraform-provider-descope/issues/124)) ([1e90875](https://github.com/jamescrowley321/terraform-provider-descope/commit/1e908758aa61b17475bce6d4d9bdc31fc45fc51b))

## [1.1.0](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.0.1...v1.1.0) (2026-03-29)


### Features

* **docs:** comprehensive docs audit — fix requirements, add examples, standardize formatting ([#120](https://github.com/jamescrowley321/terraform-provider-descope/issues/120)) ([fd444e5](https://github.com/jamescrowley321/terraform-provider-descope/commit/fd444e5af18633887a3cb5fda4e349527249f665))

## [1.0.1](https://github.com/jamescrowley321/terraform-provider-descope/compare/v1.0.0...v1.0.1) (2026-03-29)


### Bug Fixes

* **release:** remove SBOMs from checksums and drop broken Cosign step ([#117](https://github.com/jamescrowley321/terraform-provider-descope/issues/117)) ([10e672b](https://github.com/jamescrowley321/terraform-provider-descope/commit/10e672b5be423bef1e888bd999a8f38ec51b7a44))
