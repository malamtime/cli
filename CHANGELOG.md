# Changelog

## [0.0.28](https://github.com/malamtime/cli/compare/v0.0.27...v0.0.28) (2024-11-16)


### Bug Fixes

* **brand:** rename to shelltime.xyz ([79099e4](https://github.com/malamtime/cli/commit/79099e4a207ae703f58bf52298122163cf07b71e))
* **brand:** rename to shelltime.xyz ([1336fd9](https://github.com/malamtime/cli/commit/1336fd9856ee11a8cc60f6f25535b8043a5553ac))
* **cli:** add version field ([1077519](https://github.com/malamtime/cli/commit/10775191428a76ac2d2c7ac69675dc724e980c15))

## [0.0.27](https://github.com/malamtime/cli/compare/v0.0.26...v0.0.27) (2024-11-16)


### Bug Fixes

* add os and osVersion to tracking data ([68f5f21](https://github.com/malamtime/cli/commit/68f5f214daf20b8a4ca0708633c6935a3ba2f4e9))

## [0.0.26](https://github.com/malamtime/cli/compare/v0.0.25...v0.0.26) (2024-11-10)


### Bug Fixes

* **ci:** add tag to binary ([8b583ed](https://github.com/malamtime/cli/commit/8b583ed08df7754a81707c3c786d37ead17f58e7))

## [0.0.25](https://github.com/malamtime/cli/compare/v0.0.24...v0.0.25) (2024-10-20)


### Bug Fixes

* **track:** fix cursor writer ([fb0cd4d](https://github.com/malamtime/cli/commit/fb0cd4de5ddc1b750e65056bc7175deacfa79449))

## [0.0.24](https://github.com/malamtime/cli/compare/v0.0.23...v0.0.24) (2024-10-19)


### Bug Fixes

* **logger:** close logger only on program close ([187ec6f](https://github.com/malamtime/cli/commit/187ec6ffc007f7bdca92cdf68b40e3fa8064a9eb))

## [0.0.23](https://github.com/malamtime/cli/compare/v0.0.22...v0.0.23) (2024-10-15)


### Bug Fixes

* **db:** ignore empty line on db ([67cdcd5](https://github.com/malamtime/cli/commit/67cdcd5d4ed6b6d381a29e54e2a26e3fff0697c2))
* **db:** use load-once to avoid buffer based parse ([97af089](https://github.com/malamtime/cli/commit/97af0892d0cf634d54daa9888b7be80b78545560))
* **tests:** skip logger settings in testing ([9c9e8ed](https://github.com/malamtime/cli/commit/9c9e8ed7b6af83c5495cac593d70b918157eecd0))

## [0.0.22](https://github.com/malamtime/cli/compare/v0.0.21...v0.0.22) (2024-10-14)


### Bug Fixes

* **db:** fix line parser ([4f01d88](https://github.com/malamtime/cli/commit/4f01d8843035aafee5076f417a45cd615c2e6d8f))
* **log:** enable log on each command ([1374164](https://github.com/malamtime/cli/commit/1374164d736e4fa1e672bc8a57f8fc92be18a841))

## [0.0.21](https://github.com/malamtime/cli/compare/v0.0.20...v0.0.21) (2024-10-13)


### Bug Fixes

* **db:** handle cursor file data not found issue ([6bc2adb](https://github.com/malamtime/cli/commit/6bc2adb9df2eba9e9a91011581957df0c294d0a6))

## [0.0.20](https://github.com/malamtime/cli/compare/v0.0.19...v0.0.20) (2024-10-13)


### Bug Fixes

* **gc:** fix closest node check ([f93a8fb](https://github.com/malamtime/cli/commit/f93a8fb20c73c580f050ed1d207342733ac1ac1a))
* **gc:** fix gc command remove incorrectly unfinished pre commands issue ([2aa7663](https://github.com/malamtime/cli/commit/2aa7663b7c32f38ce615358cd6c17840e75920b4))

## [0.0.19](https://github.com/malamtime/cli/compare/v0.0.18...v0.0.19) (2024-10-13)


### Bug Fixes

* **api:** not parse api response if it's ok ([5d40712](https://github.com/malamtime/cli/commit/5d407126338dd43960726ca31f71cf6e72ef2809))
* **ci:** disable race in testing ([09eaa83](https://github.com/malamtime/cli/commit/09eaa8375c4535eaff3024903f81501f0dcacab7))
* **docs:** add testing badge to readme ([45e8e28](https://github.com/malamtime/cli/commit/45e8e28c3dfd3a085f924b88613d7699073a8a00))
* **gc:** clean pre, post and cursor in gc command ([b56676a](https://github.com/malamtime/cli/commit/b56676a47d8af279901924e4fa8a2378dd290a14))
* **gc:** fix gc command issue and add tests ([9a8ed39](https://github.com/malamtime/cli/commit/9a8ed392cde80e37947535cf258339d1ebb72c23))
* **track:** fix issue that could be sync data more than once ([c822898](https://github.com/malamtime/cli/commit/c8228988092065b5bccfee323b7703a8717f6946))


### Performance Improvements

* **mod:** remove unused mod ([44d2121](https://github.com/malamtime/cli/commit/44d21217c31bb271103b6d17d7dc38c1673f2c74))
* **tracking:** check the pair of pre and post command and sync to server ([3c3a13b](https://github.com/malamtime/cli/commit/3c3a13b894559f502410180d695e775d19d9c77f))
* **track:** use append file to improve performance ([6e361b4](https://github.com/malamtime/cli/commit/6e361b436c1c9389c48ce17d9e6c6531eb85491a))

## [0.0.18](https://github.com/malamtime/cli/compare/v0.0.17...v0.0.18) (2024-10-11)


### Bug Fixes

* **db:** ignore db close error ([0a2c41f](https://github.com/malamtime/cli/commit/0a2c41f11f9d2e4bb9e552afb8c93d6284121d0b))
* **http:** change auth method for http client ([ef398d9](https://github.com/malamtime/cli/commit/ef398d90c756507a7fde26f2f714dca56c7eb25d))

## [0.0.17](https://github.com/malamtime/cli/compare/v0.0.16...v0.0.17) (2024-10-05)


### Bug Fixes

* **log:** ignore panic and put it in log file silently ([7ca505b](https://github.com/malamtime/cli/commit/7ca505bee629da3383fdbcce827409b0c1d20d9a))

## [0.0.16](https://github.com/malamtime/cli/compare/v0.0.15...v0.0.16) (2024-10-05)


### Bug Fixes

* **ci:** fix missing id ([19d5f13](https://github.com/malamtime/cli/commit/19d5f13ef8dfcab7146edf86504e51e53bbaabc8))

## [0.0.15](https://github.com/malamtime/cli/compare/v0.0.14...v0.0.15) (2024-10-05)


### Bug Fixes

* **ci:** enable zip ([5c3464b](https://github.com/malamtime/cli/commit/5c3464b865be3e31e1b8c5cb67ddb2a6c9a7e2ef))

## [0.0.14](https://github.com/malamtime/cli/compare/v0.0.13...v0.0.14) (2024-10-05)


### Bug Fixes

* **ci:** enable sign and nortary ([3a73076](https://github.com/malamtime/cli/commit/3a73076edcc7e3adb3cf508d4b7e69b9711d638a))

## [0.0.13](https://github.com/malamtime/cli/compare/v0.0.12...v0.0.13) (2024-10-05)


### Bug Fixes

* **ci:** fix release config ([da230cf](https://github.com/malamtime/cli/commit/da230cf7aba05bc5b2a6a23c75bd356b21bb73f2))

## [0.0.12](https://github.com/malamtime/cli/compare/v0.0.11...v0.0.12) (2024-10-05)


### Bug Fixes

* **ci:** lock releaser version ([c6e2a21](https://github.com/malamtime/cli/commit/c6e2a21305cf3e3a4022266b03300fd87425a122))

## [0.0.11](https://github.com/malamtime/cli/compare/v0.0.10...v0.0.11) (2024-10-05)


### Features

* **gc:** add gc command ([533044f](https://github.com/malamtime/cli/commit/533044fb10f6eeb4631d670dbca82ca8ae04dc5d))


### Bug Fixes

* **db:** fix db api and log permission issue ([e37dbc4](https://github.com/malamtime/cli/commit/e37dbc456e0c91febd0049c7ed0361aa2e728491))
* **gc:** fix syntax error and update parameters of gc command ([329566e](https://github.com/malamtime/cli/commit/329566ec096d3b37e65f40ca08c7e5ca23a1b9f0))


### Performance Improvements

* **db:** migrate to NutsDB since sqlite is too hard to compile ([f765e7b](https://github.com/malamtime/cli/commit/f765e7bc6500eec805b9e117913984873f8c0a4c))


### Miscellaneous Chores

* release 0.0.11 ([b96d666](https://github.com/malamtime/cli/commit/b96d6663cee287058475d9caed17fb775072368e))

## [0.0.10](https://github.com/malamtime/cli/compare/v0.0.9...v0.0.10) (2024-10-05)


### Bug Fixes

* **ci:** ignore coverage.txt for release ([39b7994](https://github.com/malamtime/cli/commit/39b7994cdb89d68cf87fa075ad43034cec2bd2a4))

## [0.0.9](https://github.com/malamtime/cli/compare/v0.0.8...v0.0.9) (2024-10-05)


### Bug Fixes

* **ci:** use quill for binary sign and nortary ([f885733](https://github.com/malamtime/cli/commit/f8857334545c24f61472ad08da2da94c1b1df0b0))

## [0.0.8](https://github.com/malamtime/cli/compare/v0.0.7...v0.0.8) (2024-10-05)


### Bug Fixes

* **ci:** disable windows build ([e15c7ab](https://github.com/malamtime/cli/commit/e15c7ab33d26c76306a4d85f09f1a948534c0e59))

## [0.0.7](https://github.com/malamtime/cli/compare/v0.0.6...v0.0.7) (2024-10-05)


### Bug Fixes

* **ci:** enable cgo for sqlite ([d3d662e](https://github.com/malamtime/cli/commit/d3d662e26e8b06006ae919ec1979807afa117620))

## [0.0.6](https://github.com/malamtime/cli/compare/v0.0.5...v0.0.6) (2024-10-05)


### Bug Fixes

* **ci:** enable cgo for sqlite ([374e5c3](https://github.com/malamtime/cli/commit/374e5c3da4965b181d51dd1d1408ce2dae6db5ab))

## [0.0.5](https://github.com/malamtime/cli/compare/v0.0.4...v0.0.5) (2024-10-05)


### Performance Improvements

* **track:** save to fs first for performance purpose and sync it later ([57c5606](https://github.com/malamtime/cli/commit/57c56066f92ca22b289ddf233657107b72afdbc0))

## [0.0.4](https://github.com/malamtime/cli/compare/v0.0.3...v0.0.4) (2024-10-04)


### Bug Fixes

* **logger:** add logger for app ([0a9821f](https://github.com/malamtime/cli/commit/0a9821ff2c876f6a0f621df3c15ec28197009ed3))
* **track:** fix track method and logger ([3e52212](https://github.com/malamtime/cli/commit/3e5221280586649bb674ff190743a8412798dbff))

## [0.0.3](https://github.com/malamtime/cli/compare/v0.0.2...v0.0.3) (2024-10-04)


### Bug Fixes

* **track:** add tip on track command ([000a367](https://github.com/malamtime/cli/commit/000a367e73fcfdff566f3eeb2a7e9e5d5242ad18))

## [0.0.2](https://github.com/malamtime/cli/compare/v0.0.1...v0.0.2) (2024-10-04)


### Bug Fixes

* **ci:** fix release command ([1ea4e73](https://github.com/malamtime/cli/commit/1ea4e730c5ab3abe0220be1a208fd295da6d3c2b))

## 0.0.1 (2024-10-04)


### Features

* add track command ([f85660a](https://github.com/malamtime/cli/commit/f85660a63f83c69229fe1d4a4b534a1c76f49b58))
* **app:** add basic command and ci ([976973f](https://github.com/malamtime/cli/commit/976973fde38bd054cbdcff9de26b80b73c855892))


### Bug Fixes

* **ci:** fix ci branch name ([5b08dd8](https://github.com/malamtime/cli/commit/5b08dd85d0d818cd1d7d5686ffeb03303d4b00ae))
* **track:** add more params for track ([24de070](https://github.com/malamtime/cli/commit/24de070375e2acec0aaec479387d50baa42b2561))


### Miscellaneous Chores

* release 0.0.1 ([5517712](https://github.com/malamtime/cli/commit/5517712672634a2d7fb5e1438028b1f3a58beb02))
