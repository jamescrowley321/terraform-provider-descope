# Changelog

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
