# Changelog

## [0.0.48](https://github.com/malamtime/cli/compare/v0.0.47...v0.0.48) (2024-12-13)


### Bug Fixes

* **docs:** add performance explaination ([a467327](https://github.com/malamtime/cli/commit/a46732786e0f522f0cf08fd8ef17475040cc57c9))
* **docs:** fix config field in readme ([70ebee0](https://github.com/malamtime/cli/commit/70ebee03c944ce25533732ff94d97db963746c9c))
* **docs:** fix docs ([5473019](https://github.com/malamtime/cli/commit/547301994db88fb4f9a5a6c324d905698c019abc))
* **track:** increase the buffer length when parse file line by line ([6145772](https://github.com/malamtime/cli/commit/6145772df0fbf46cb317f8710d9f005630dc1b1c))

## [0.0.47](https://github.com/malamtime/cli/compare/v0.0.46...v0.0.47) (2024-12-13)


### Bug Fixes

* **sync:** make the sync could be force in `sync` command ([025ab3a](https://github.com/malamtime/cli/commit/025ab3a221c576f2ae8822d0211690d58ce6df7d))

## [0.0.46](https://github.com/malamtime/cli/compare/v0.0.45...v0.0.46) (2024-12-13)


### Features

* **sync:** add docs about `sync` command ([826f316](https://github.com/malamtime/cli/commit/826f31668dfc7140def0be123593d863964b3f36))
* **sync:** add sync command ([2236979](https://github.com/malamtime/cli/commit/2236979b0bbbc82a4a044db801e79a9307339a5e))


### Bug Fixes

* **docs:** remove unused docs ([a3e77fc](https://github.com/malamtime/cli/commit/a3e77fc2de237421a8cdc7571814b697c815e6c4))


### Miscellaneous Chores

* release 0.0.46 ([cc44364](https://github.com/malamtime/cli/commit/cc443646677dc7856a4b51c7a85e124de2a52d3c))

## [0.0.45](https://github.com/malamtime/cli/compare/v0.0.44...v0.0.45) (2024-12-13)


### Bug Fixes

* **docs:** update readme ([8c0f281](https://github.com/malamtime/cli/commit/8c0f281784c23afe273ed886ad9028e4d3bee48d))
* **gc:** remove unused empty line on gc ([eea8312](https://github.com/malamtime/cli/commit/eea8312943f629705ac50402ddbe463470c7ce64))


### Performance Improvements

* **model:** `GetPreCommandsTree` performance improve and add benchmark tests ([a0f8ae8](https://github.com/malamtime/cli/commit/a0f8ae8f6f4fb743e179c44b0eefc315063e6609))
* **model:** improve performance on `GetPreCommands` ([8d83ce2](https://github.com/malamtime/cli/commit/8d83ce26344f5645d1022722c37806ef3195e8b3))
* **model:** use bytes operators on postFile to improve performance ([6d15ca0](https://github.com/malamtime/cli/commit/6d15ca0d2794381cccb9f6625c1951cca64fe48d))

## [0.0.44](https://github.com/malamtime/cli/compare/v0.0.43...v0.0.44) (2024-12-13)


### Bug Fixes

* **ci:** ignore codecov generated files ([2109547](https://github.com/malamtime/cli/commit/21095477f36dcd357ad2548031e290ed92158f56))

## [0.0.43](https://github.com/malamtime/cli/compare/v0.0.42...v0.0.43) (2024-12-13)


### Bug Fixes

* **ci:** add uptrace on ci ([cac2a2e](https://github.com/malamtime/cli/commit/cac2a2e3a3cde61cad7497af89fa713ef8c77d38))
* **ci:** set timeout as 3m ([91ca363](https://github.com/malamtime/cli/commit/91ca3634747d6069f5a4a4f9488c7e319cc6fc89))
* **ci:** upgrade codecov action to v5 ([d02c5d6](https://github.com/malamtime/cli/commit/d02c5d6b260ecad2e66def25886e8ccc703d040b))
* **track:** fix mock config service ([bb944f9](https://github.com/malamtime/cli/commit/bb944f9396245d3976c76970708b1800dcc3c290))

## [0.0.42](https://github.com/malamtime/cli/compare/v0.0.41...v0.0.42) (2024-12-13)


### Features

* **trace:** add trace for cli ([bea74de](https://github.com/malamtime/cli/commit/bea74de6f45b3064eba5d1b163edb1cf7c159d62))


### Bug Fixes

* **docs:** update readme ([8bf7e8f](https://github.com/malamtime/cli/commit/8bf7e8fde4dba2578073701b03b351ef569cf309))


### Miscellaneous Chores

* release 0.0.42 ([5c56570](https://github.com/malamtime/cli/commit/5c56570dbc6e455bbec759130983cce373887016))

## [0.0.41](https://github.com/malamtime/cli/compare/v0.0.40...v0.0.41) (2024-12-10)


### Bug Fixes

* **docs:** add shelltime badge ([98e4529](https://github.com/malamtime/cli/commit/98e45296f07b7ab53cf9a4a7a479542750bfa757))

## [0.0.40](https://github.com/malamtime/cli/compare/v0.0.39...v0.0.40) (2024-12-10)


### Performance Improvements

* **api:** reduce tracking data size for performance improve ([c73c17a](https://github.com/malamtime/cli/commit/c73c17ab24f7f95da56576722f586e11d6fd43c4))

## [0.0.39](https://github.com/malamtime/cli/compare/v0.0.38...v0.0.39) (2024-12-09)


### Bug Fixes

* **handshake:** fix check method of handshake ([fd1af55](https://github.com/malamtime/cli/commit/fd1af55f433b0a93c7cfb5df169d58227ec2639d))

## [0.0.38](https://github.com/malamtime/cli/compare/v0.0.37...v0.0.38) (2024-12-09)


### Features

* **handshake:** add handshake support for smooth auth ([c9d8338](https://github.com/malamtime/cli/commit/c9d833819f0b5d26bdf7f0be30ebc90cb703b3f8))


### Bug Fixes

* **version:** remove legacy version info ([31d0c59](https://github.com/malamtime/cli/commit/31d0c5971c0d2e26a680c1f544080f29e6dfe487))


### Miscellaneous Chores

* release 0.0.38 ([412fdc5](https://github.com/malamtime/cli/commit/412fdc50cc56faa21859ae263c61fdc41634f6cd))

## [0.0.37](https://github.com/malamtime/cli/compare/v0.0.36...v0.0.37) (2024-12-07)


### Bug Fixes

* **api:** add testcase for model/api ([925ced5](https://github.com/malamtime/cli/commit/925ced5ad031d13bae467c7714b686cf48852a0d))
* **http:** add timeout on http send req ([4f7880c](https://github.com/malamtime/cli/commit/4f7880ca927f252a0f95e2084ff72e48197006f6))
* **http:** support multiple endpoint for debug ([1251f65](https://github.com/malamtime/cli/commit/1251f65ad95faf582128ac29cc9a4b2d737dde8e))

## [0.0.36](https://github.com/malamtime/cli/compare/v0.0.35...v0.0.36) (2024-12-04)


### Bug Fixes

* **build:** use default ldflags on build ([82cc3ca](https://github.com/malamtime/cli/commit/82cc3ca861aed8da3ffddf3b8d5d4b74a91137fe))

## [0.0.35](https://github.com/malamtime/cli/compare/v0.0.34...v0.0.35) (2024-12-04)


### Bug Fixes

* **ci:** update metadata in main ([28b2bbf](https://github.com/malamtime/cli/commit/28b2bbf8030701eb62f10d6c22cd191f42c52e8d))

## [0.0.34](https://github.com/malamtime/cli/compare/v0.0.33...v0.0.34) (2024-12-01)


### Bug Fixes

* **track:** fix first load tests on track ([0f30741](https://github.com/malamtime/cli/commit/0f3074188a0d8fc6b6bef4b9d25f78251008393e))

## [0.0.33](https://github.com/malamtime/cli/compare/v0.0.32...v0.0.33) (2024-11-24)


### Bug Fixes

* **track:** allow very first command sync to server ([676ece3](https://github.com/malamtime/cli/commit/676ece3340c166d222a43cda496c287846d27d78))

## [0.0.32](https://github.com/malamtime/cli/compare/v0.0.31...v0.0.32) (2024-11-21)


### Features

* **track:** add data masking for sensitive token ([bb2460a](https://github.com/malamtime/cli/commit/bb2460af3e1bc2c55d86f05d0880138d1a2f9e57))


### Bug Fixes

* **os:** correct os info in linux ([c9f5f98](https://github.com/malamtime/cli/commit/c9f5f98bdc38e5b6b01bd8c0246afd9a4f92f2b9))


### Miscellaneous Chores

* release 0.0.32 ([d298613](https://github.com/malamtime/cli/commit/d2986139e434947746811c39a73e7606f83faf8d))

## [0.0.31](https://github.com/malamtime/cli/compare/v0.0.30...v0.0.31) (2024-11-19)


### Bug Fixes

* **gc:** check original file exists before rename it ([3263d47](https://github.com/malamtime/cli/commit/3263d47df98ac91ae69d7e29df4d5e4c373ede71))

## [0.0.30](https://github.com/malamtime/cli/compare/v0.0.29...v0.0.30) (2024-11-16)


### Bug Fixes

* **api:** fix testcase on msgpack decode ([15e9c92](https://github.com/malamtime/cli/commit/15e9c9270eabe8f54bf76d0773ea97e438fd4854))
* **api:** support msgpack in api ([8656c15](https://github.com/malamtime/cli/commit/8656c155a73b96f34899ebb411a29b7bf4881abe))
* **ci:** use release action tag name instead of github ci one ([75fb9bf](https://github.com/malamtime/cli/commit/75fb9bfcb33496bfe36aba079295120c50796712))

## [0.0.29](https://github.com/malamtime/cli/compare/v0.0.28...v0.0.29) (2024-11-16)


### Bug Fixes

* **brand:** fix config folder combination ([7642c97](https://github.com/malamtime/cli/commit/7642c97e0bdfa1108de8408fb2784928016e1153))

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
