# Change Log

## [v2.0.0-beta.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v2.0.0-beta.2) (2020-05-25)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.0.0-beta.1...v2.0.0-beta.2)

**Merged pull requests:**

- utilities: move package to public API [\#1395](https://github.com/grpc-ecosystem/grpc-gateway/pull/1395) ([johanbrandhorst](https://github.com/johanbrandhorst))

## [v2.0.0-beta.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v2.0.0-beta.1) (2020-05-25)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.0.0-alpha.2...v2.0.0-beta.1)

**Implemented enhancements:**

- Feature request: support of go protocol buffer v2 [\#1147](https://github.com/grpc-ecosystem/grpc-gateway/issues/1147)

**Fixed bugs:**

- Bazel tests are flaky [\#968](https://github.com/grpc-ecosystem/grpc-gateway/issues/968)

**Closed issues:**

- Merging swagger specs \(for multiple protos belonging to same package\) fails to emit summary for some RPCs [\#1387](https://github.com/grpc-ecosystem/grpc-gateway/issues/1387)
- Remove result/error envelope in streaming RPC response [\#1254](https://github.com/grpc-ecosystem/grpc-gateway/issues/1254)
- cleanup: rename 'swagger' references to 'openapi' [\#675](https://github.com/grpc-ecosystem/grpc-gateway/issues/675)
- Swagger: JSON definitions aren't CamelCased [\#375](https://github.com/grpc-ecosystem/grpc-gateway/issues/375)
- Support emitting default values in JSON [\#233](https://github.com/grpc-ecosystem/grpc-gateway/issues/233)

**Merged pull requests:**

- Generate changelog for v2.0.0-beta1 [\#1393](https://github.com/grpc-ecosystem/grpc-gateway/pull/1393) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Rename protoc-gen-swagger to protoc-gen-openapi [\#1390](https://github.com/grpc-ecosystem/grpc-gateway/pull/1390) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update docs for v2 [\#1388](https://github.com/grpc-ecosystem/grpc-gateway/pull/1388) ([johanbrandhorst](https://github.com/johanbrandhorst))
- runtime: remove DisallowUnknownFields\(\) [\#1386](https://github.com/grpc-ecosystem/grpc-gateway/pull/1386) ([johanbrandhorst](https://github.com/johanbrandhorst))
- runtime: make HTTPBodyMarshaler the default [\#1385](https://github.com/grpc-ecosystem/grpc-gateway/pull/1385) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Implement last-match-wins behaviour in mux [\#1384](https://github.com/grpc-ecosystem/grpc-gateway/pull/1384) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Fix various misspellings \(\#1381\) [\#1383](https://github.com/grpc-ecosystem/grpc-gateway/pull/1383) ([bvwells](https://github.com/bvwells))
- Remove usage of deprecated grpc.Errorf API \(\#1380\) [\#1382](https://github.com/grpc-ecosystem/grpc-gateway/pull/1382) ([bvwells](https://github.com/bvwells))
- Add missing documentation for openapiv2 proto definition \(\#1375\) [\#1378](https://github.com/grpc-ecosystem/grpc-gateway/pull/1378) ([bvwells](https://github.com/bvwells))
- runtime: make default marshaler emit default values [\#1377](https://github.com/grpc-ecosystem/grpc-gateway/pull/1377) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Turn off UseProtoNames in default marshaler [\#1376](https://github.com/grpc-ecosystem/grpc-gateway/pull/1376) ([johanbrandhorst](https://github.com/johanbrandhorst))
- cherry-picked the commit from master to v2 in customizingyourgateway.md [\#1374](https://github.com/grpc-ecosystem/grpc-gateway/pull/1374) ([iamrajiv](https://github.com/iamrajiv))
- Update google.golang.org/genproto commit hash to e9a78aa \(v2\) [\#1372](https://github.com/grpc-ecosystem/grpc-gateway/pull/1372) ([renovate[bot]](https://github.com/apps/renovate))
- PR and Issue template added to v2 branch [\#1371](https://github.com/grpc-ecosystem/grpc-gateway/pull/1371) ([iamrajiv](https://github.com/iamrajiv))
- Update google.golang.org/genproto commit hash to 08726f3 \(v2\) [\#1369](https://github.com/grpc-ecosystem/grpc-gateway/pull/1369) ([renovate[bot]](https://github.com/apps/renovate))
- Improve README.md [\#1363](https://github.com/grpc-ecosystem/grpc-gateway/pull/1363) ([amanjain97](https://github.com/amanjain97))
- all: replace all uses of golang/protobuf/proto [\#1358](https://github.com/grpc-ecosystem/grpc-gateway/pull/1358) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update module google/go-cmp to v0.4.1 \(v2\) [\#1357](https://github.com/grpc-ecosystem/grpc-gateway/pull/1357) ([renovate[bot]](https://github.com/apps/renovate))
- chore\(deps\): update golang docker tag to v1.14.3 \(v2\) [\#1356](https://github.com/grpc-ecosystem/grpc-gateway/pull/1356) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency io\_bazel\_rules\_go to v0.23.1 \(v2\) [\#1355](https://github.com/grpc-ecosystem/grpc-gateway/pull/1355) ([renovate[bot]](https://github.com/apps/renovate))
- remove unused PKGMAP [\#1353](https://github.com/grpc-ecosystem/grpc-gateway/pull/1353) ([johanbrandhorst](https://github.com/johanbrandhorst))
- internal/descriptor: move to generated apiconfig test [\#1352](https://github.com/grpc-ecosystem/grpc-gateway/pull/1352) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update google.golang.org/genproto commit hash to fc4c6c6 \(v2\) [\#1349](https://github.com/grpc-ecosystem/grpc-gateway/pull/1349) ([renovate[bot]](https://github.com/apps/renovate))
- Renovate changes from master [\#1348](https://github.com/grpc-ecosystem/grpc-gateway/pull/1348) ([johanbrandhorst](https://github.com/johanbrandhorst))
- docs: add example customizing unmarshal options \(\#1335\) [\#1341](https://github.com/grpc-ecosystem/grpc-gateway/pull/1341) ([srenatus](https://github.com/srenatus))
- chore\(deps\): update module golang/protobuf to v1.4.2 \(v2\) [\#1333](https://github.com/grpc-ecosystem/grpc-gateway/pull/1333) ([renovate[bot]](https://github.com/apps/renovate))
- Update google.golang.org/genproto commit hash to 8feb7f2 \(v2\) [\#1331](https://github.com/grpc-ecosystem/grpc-gateway/pull/1331) ([renovate[bot]](https://github.com/apps/renovate))
- all: correct use of go\_package [\#1329](https://github.com/grpc-ecosystem/grpc-gateway/pull/1329) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update dependency io\_bazel\_rules\_go to v0.23.0 \(v2\) [\#1328](https://github.com/grpc-ecosystem/grpc-gateway/pull/1328) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency bazel\_gazelle to v0.21.0 \(v2\) [\#1327](https://github.com/grpc-ecosystem/grpc-gateway/pull/1327) ([renovate[bot]](https://github.com/apps/renovate))
- Renovate: run go mod tidy after updates [\#1324](https://github.com/grpc-ecosystem/grpc-gateway/pull/1324) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update google.golang.org/genproto commit hash to 09dca8e \(v2\) [\#1322](https://github.com/grpc-ecosystem/grpc-gateway/pull/1322) ([renovate[bot]](https://github.com/apps/renovate))
- fixed typo in docs/\_docs/usage.md in v2 branch [\#1320](https://github.com/grpc-ecosystem/grpc-gateway/pull/1320) ([iamrajiv](https://github.com/iamrajiv))
- improved docs/\_docs/season\_of\_docs.md in v2 branch [\#1319](https://github.com/grpc-ecosystem/grpc-gateway/pull/1319) ([iamrajiv](https://github.com/iamrajiv))
- Add Bazel CI caching [\#1314](https://github.com/grpc-ecosystem/grpc-gateway/pull/1314) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update module antihax/optional to v1 \(v2\) [\#1311](https://github.com/grpc-ecosystem/grpc-gateway/pull/1311) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency com\_github\_bazelbuild\_buildtools to v3 \(v2\) [\#1310](https://github.com/grpc-ecosystem/grpc-gateway/pull/1310) ([renovate[bot]](https://github.com/apps/renovate))
- runtime: rewrite fieldmask logic with protoreflect [\#1309](https://github.com/grpc-ecosystem/grpc-gateway/pull/1309) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Update module google.golang.org/grpc to v1.29.1 \(v2\) [\#1307](https://github.com/grpc-ecosystem/grpc-gateway/pull/1307) ([renovate[bot]](https://github.com/apps/renovate))
- Update module golang/protobuf to v1.4.1 \(v2\) [\#1306](https://github.com/grpc-ecosystem/grpc-gateway/pull/1306) ([renovate[bot]](https://github.com/apps/renovate))
- Update golang Docker tag to v1.14.2 \(v2\) [\#1305](https://github.com/grpc-ecosystem/grpc-gateway/pull/1305) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency io\_bazel\_rules\_go to v0.22.4 \(v2\) [\#1301](https://github.com/grpc-ecosystem/grpc-gateway/pull/1301) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency com\_github\_bazelbuild\_buildtools to v0.29.0 \(v2\) [\#1300](https://github.com/grpc-ecosystem/grpc-gateway/pull/1300) ([renovate[bot]](https://github.com/apps/renovate))
- Update dependency bazel\_gazelle to v0.20.0 \(v2\) [\#1299](https://github.com/grpc-ecosystem/grpc-gateway/pull/1299) ([renovate[bot]](https://github.com/apps/renovate))
- Update google.golang.org/genproto commit hash to 43844f6 \(v2\) [\#1298](https://github.com/grpc-ecosystem/grpc-gateway/pull/1298) ([renovate[bot]](https://github.com/apps/renovate))
- Update golang.org/x/oauth2 commit hash to bf48bf1 \(v2\) [\#1297](https://github.com/grpc-ecosystem/grpc-gateway/pull/1297) ([renovate[bot]](https://github.com/apps/renovate))
- Add more instructions on the GitHub releases UI [\#1277](https://github.com/grpc-ecosystem/grpc-gateway/pull/1277) ([achew22](https://github.com/achew22))
- HttpBody mesage feature for stream RPC [\#1273](https://github.com/grpc-ecosystem/grpc-gateway/pull/1273) ([adasari](https://github.com/adasari))
- Consolidate error handling configuration [\#1265](https://github.com/grpc-ecosystem/grpc-gateway/pull/1265) ([johanbrandhorst](https://github.com/johanbrandhorst))

## [v1.14.5](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.5) (2020-05-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.0.0-alpha.2...v1.14.5)

**Implemented enhancements:**

- Standalone gateway package [\#1239](https://github.com/grpc-ecosystem/grpc-gateway/pull/1239) ([ssetin](https://github.com/ssetin))

**Fixed bugs:**

- Regeneration commands don't regenerate all files [\#1229](https://github.com/grpc-ecosystem/grpc-gateway/issues/1229)

**Closed issues:**

- bug in protoc-gen-swagger with put method and field\_mask [\#1271](https://github.com/grpc-ecosystem/grpc-gateway/issues/1271)
- How to change title and version in generated swagger json file [\#1269](https://github.com/grpc-ecosystem/grpc-gateway/issues/1269)
- Error when using oneof name in response body selector [\#1264](https://github.com/grpc-ecosystem/grpc-gateway/issues/1264)
- httpbody.proto has incorrect go\_package [\#1263](https://github.com/grpc-ecosystem/grpc-gateway/issues/1263)
- How to implement field validation for optional parameter [\#1256](https://github.com/grpc-ecosystem/grpc-gateway/issues/1256)
- HTTP GET with query params error [\#1245](https://github.com/grpc-ecosystem/grpc-gateway/issues/1245)
- how to get customezing header [\#1244](https://github.com/grpc-ecosystem/grpc-gateway/issues/1244)
- proto function not define in yaml then grpc-gateway\_out err [\#1233](https://github.com/grpc-ecosystem/grpc-gateway/issues/1233)
- ll [\#1232](https://github.com/grpc-ecosystem/grpc-gateway/issues/1232)
- Add WithUnmarshaler NewServeMux option. [\#1226](https://github.com/grpc-ecosystem/grpc-gateway/issues/1226)
- Feature request: Reject call if there are parameters/fields that don't exist in the protobuf [\#1210](https://github.com/grpc-ecosystem/grpc-gateway/issues/1210)
- Migrate away from using the protoc-gen-go/generator package [\#1209](https://github.com/grpc-ecosystem/grpc-gateway/issues/1209)
- Swagger doc generation got stuck for infinite time if a message refers itself. [\#1167](https://github.com/grpc-ecosystem/grpc-gateway/issues/1167)
- Remove "error" field from errors. [\#1098](https://github.com/grpc-ecosystem/grpc-gateway/issues/1098)

**Merged pull requests:**

- runtime: rewrite query.go in new protoreflect [\#1272](https://github.com/grpc-ecosystem/grpc-gateway/pull/1272) ([johanbrandhorst](https://github.com/johanbrandhorst))
- V2 [\#1267](https://github.com/grpc-ecosystem/grpc-gateway/pull/1267) ([Romeren](https://github.com/Romeren))
- runtime: replace StreamError with status.Status [\#1262](https://github.com/grpc-ecosystem/grpc-gateway/pull/1262) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Remove uses of protoc-gen-go/generator [\#1261](https://github.com/grpc-ecosystem/grpc-gateway/pull/1261) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Add Chef and Scaleway to ADOPTERS.md [\#1257](https://github.com/grpc-ecosystem/grpc-gateway/pull/1257) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Use public key to authenticate for docker registry [\#1253](https://github.com/grpc-ecosystem/grpc-gateway/pull/1253) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Add user testimonials and ADOPTERS.md [\#1251](https://github.com/grpc-ecosystem/grpc-gateway/pull/1251) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Minimize public API [\#1249](https://github.com/grpc-ecosystem/grpc-gateway/pull/1249) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Change build env to git repo v2 [\#1247](https://github.com/grpc-ecosystem/grpc-gateway/pull/1247) ([johanbrandhorst](https://github.com/johanbrandhorst))
- \#1098: remove error field from error messages [\#1242](https://github.com/grpc-ecosystem/grpc-gateway/pull/1242) ([srenatus](https://github.com/srenatus))
- v2: remove deprecated usage of github.com/golang/protobuf/protoc-gen-go/g… [\#1238](https://github.com/grpc-ecosystem/grpc-gateway/pull/1238) ([yinzara](https://github.com/yinzara))
- remove deprecated usage of github.com/golang/protobuf/protoc-gen-go/g… [\#1237](https://github.com/grpc-ecosystem/grpc-gateway/pull/1237) ([yinzara](https://github.com/yinzara))
- v2: set rpcMethodName in request.Context [\#1234](https://github.com/grpc-ecosystem/grpc-gateway/pull/1234) ([AyushG3112](https://github.com/AyushG3112))
- Fix generating helloworld.proto [\#1231](https://github.com/grpc-ecosystem/grpc-gateway/pull/1231) ([johanbrandhorst](https://github.com/johanbrandhorst))
- v2 fix query params for invoked in-process local request [\#1227](https://github.com/grpc-ecosystem/grpc-gateway/pull/1227) ([hb-chen](https://github.com/hb-chen))
- docs: fix some broken links and change godoc.org links to pkg.go.dev [\#1222](https://github.com/grpc-ecosystem/grpc-gateway/pull/1222) ([AyushG3112](https://github.com/AyushG3112))
- docs: change godoc.org links to pkg.go.dev [\#1220](https://github.com/grpc-ecosystem/grpc-gateway/pull/1220) ([johanbrandhorst](https://github.com/johanbrandhorst))
- v2: Fix gorelease base, update README [\#1218](https://github.com/grpc-ecosystem/grpc-gateway/pull/1218) ([johanbrandhorst](https://github.com/johanbrandhorst))
- Upgrade protobuf to 1.4.0 [\#1165](https://github.com/grpc-ecosystem/grpc-gateway/pull/1165) ([johanbrandhorst](https://github.com/johanbrandhorst))

## [v2.0.0-alpha.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v2.0.0-alpha.2) (2020-04-18)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.0.0-alpha.1...v2.0.0-alpha.2)

## [v2.0.0-alpha.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v2.0.0-alpha.1) (2020-04-18)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.4...v2.0.0-alpha.1)

## [v1.14.4](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.4) (2020-04-18)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.4-rc.1...v1.14.4)

**Closed issues:**

- custom unmarshal response error [\#1211](https://github.com/grpc-ecosystem/grpc-gateway/issues/1211)
- api annotation placeholder character case [\#1207](https://github.com/grpc-ecosystem/grpc-gateway/issues/1207)
- response\_body of http rule is ignored in case of server streaming [\#1202](https://github.com/grpc-ecosystem/grpc-gateway/issues/1202)
- Go gRPC Gateway - Type of one field in JSON is not same as in proto [\#1201](https://github.com/grpc-ecosystem/grpc-gateway/issues/1201)
- 1.14 breaks protoc-gen-go-json [\#1199](https://github.com/grpc-ecosystem/grpc-gateway/issues/1199)
- Use different field name in gateway [\#1197](https://github.com/grpc-ecosystem/grpc-gateway/issues/1197)
- No longer generate register...Server function in pb.gw.go? [\#1195](https://github.com/grpc-ecosystem/grpc-gateway/issues/1195)
- Support for source\_relative argument [\#1180](https://github.com/grpc-ecosystem/grpc-gateway/issues/1180)

**Merged pull requests:**

- Add Season of Docs project ideas page [\#1224](https://github.com/grpc-ecosystem/grpc-gateway/pull/1224) ([johanbrandhorst](https://github.com/johanbrandhorst))

## [v1.14.4-rc.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.4-rc.1) (2020-04-01)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.3...v1.14.4-rc.1)

**Implemented enhancements:**

- proto-gen-\* --version [\#649](https://github.com/grpc-ecosystem/grpc-gateway/issues/649)

**Closed issues:**

- protoc-gen-swagger: json\_names\_for\_fields=true does not respect json\_name for \*nested\* path parameters [\#1187](https://github.com/grpc-ecosystem/grpc-gateway/issues/1187)
- protoc-gen-grpc-gateway-ts plugin for generating Typescript types [\#1182](https://github.com/grpc-ecosystem/grpc-gateway/issues/1182)
- protoc-gen-swagger: support outputting enum parameters as integers [\#1177](https://github.com/grpc-ecosystem/grpc-gateway/issues/1177)
- protoc-gen-grpc-gateway should warn or fail if a selector does not exist [\#1175](https://github.com/grpc-ecosystem/grpc-gateway/issues/1175)
- Enum string values not supported by grpc gateway [\#1171](https://github.com/grpc-ecosystem/grpc-gateway/issues/1171)
- Swagger definition is broken if there's no default error response set [\#1162](https://github.com/grpc-ecosystem/grpc-gateway/issues/1162)
- Using golang/protobuf v1.4.0-rc.3 locks up generators [\#1158](https://github.com/grpc-ecosystem/grpc-gateway/issues/1158)

## [v1.14.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.3) (2020-03-11)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.2...v1.14.3)

**Closed issues:**

- Missing httprule breaks our module [\#1168](https://github.com/grpc-ecosystem/grpc-gateway/issues/1168)

## [v1.14.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.2) (2020-03-09)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.1...v1.14.2)

**Closed issues:**

- 1.14.x breaks clay [\#1161](https://github.com/grpc-ecosystem/grpc-gateway/issues/1161)
- Enum string values not supported by grpc gateway [\#1159](https://github.com/grpc-ecosystem/grpc-gateway/issues/1159)
- Path parameters visible in the body of the request [\#1157](https://github.com/grpc-ecosystem/grpc-gateway/issues/1157)
- 1.14.0 breaks existing pipelines [\#1154](https://github.com/grpc-ecosystem/grpc-gateway/issues/1154)

## [v1.14.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.1) (2020-03-05)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.14.0...v1.14.1)

## [v1.14.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.14.0) (2020-03-04)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.13.0...v1.14.0)

**Closed issues:**

- Swagger Codegen With Multiple Additional Bindings is resulting in missing Documentation [\#1149](https://github.com/grpc-ecosystem/grpc-gateway/issues/1149)
- google.api.http option not picked up/mapped correctly? [\#1148](https://github.com/grpc-ecosystem/grpc-gateway/issues/1148)
- cannot use multiple different error handlers in one gateway binary [\#1143](https://github.com/grpc-ecosystem/grpc-gateway/issues/1143)
- Is grpc-gateway a stateless application? [\#1139](https://github.com/grpc-ecosystem/grpc-gateway/issues/1139)
- Support for custom query parameters parsers [\#1128](https://github.com/grpc-ecosystem/grpc-gateway/issues/1128)

## [v1.13.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.13.0) (2020-02-11)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.12.2...v1.13.0)

**Implemented enhancements:**

- \[Feature Request\] Custom Type conversion [\#754](https://github.com/grpc-ecosystem/grpc-gateway/issues/754)

**Closed issues:**

- swagger: Prefix used in nested-field GET params does not respect json\_names\_for\_fields [\#1125](https://github.com/grpc-ecosystem/grpc-gateway/issues/1125)
- Generated spec missing grpc errors [\#1122](https://github.com/grpc-ecosystem/grpc-gateway/issues/1122)
- protoc-gen-swagger: Support for repeated custom message in url params [\#1119](https://github.com/grpc-ecosystem/grpc-gateway/issues/1119)
- protoc-gen-swagger: Support Response-specific examples  [\#1117](https://github.com/grpc-ecosystem/grpc-gateway/issues/1117)

## [v1.12.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.12.2) (2020-01-22)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.12.1...v1.12.2)

**Fixed bugs:**

- Gateway does not parse oneof types correctly when using camelCase [\#1113](https://github.com/grpc-ecosystem/grpc-gateway/issues/1113)

**Closed issues:**

- all SubConns are in TransientFailure, latest connection error: connection error: desc = \"transport: Error while dialing dial tcp \[::1\]:50051: connect: connection refused\" [\#1111](https://github.com/grpc-ecosystem/grpc-gateway/issues/1111)
- Streaming responses are put in `x-stream-definitions` rather than `/components/schema`, which tools do not support [\#1109](https://github.com/grpc-ecosystem/grpc-gateway/issues/1109)
- Why is metadataAnnotators not \[\]func\(metadata.MD, \*http.Request\) [\#1107](https://github.com/grpc-ecosystem/grpc-gateway/issues/1107)
- protoc\_gen\_swagger: Why `protobufListValue` type is represented as object? [\#1106](https://github.com/grpc-ecosystem/grpc-gateway/issues/1106)
- The metadata of local\_request function's return is always nil [\#1105](https://github.com/grpc-ecosystem/grpc-gateway/issues/1105)
- protoc-gen-swagger: support multiple files [\#1104](https://github.com/grpc-ecosystem/grpc-gateway/issues/1104)
- grpc-gateway is incompatible with --incompatible\_load\_proto\_rules\_from\_bzl [\#1101](https://github.com/grpc-ecosystem/grpc-gateway/issues/1101)
- Any example or doc of sending large data from Postman \(json\) to the server using client side streaming? [\#1100](https://github.com/grpc-ecosystem/grpc-gateway/issues/1100)
- Generate java grpc gateway code. [\#1097](https://github.com/grpc-ecosystem/grpc-gateway/issues/1097)
-  protoc-gen-swagger gen unsupport  http post  in query,  only support http get in query? [\#1096](https://github.com/grpc-ecosystem/grpc-gateway/issues/1096)
- Cannot generate swagger with generated protobuf [\#1094](https://github.com/grpc-ecosystem/grpc-gateway/issues/1094)
- genswagger fails on Ubuntu on Windows \(WSL\): value of type genswagger.alias is not assignable to type genswagger.alias [\#1092](https://github.com/grpc-ecosystem/grpc-gateway/issues/1092)
- protoc-gen-swagger: should use json\_name values for generating properties in swagger doc [\#1090](https://github.com/grpc-ecosystem/grpc-gateway/issues/1090)
- protoc-gen-swagger: json\_names\_for\_fields=true does not respect json\_name for path parameters [\#1084](https://github.com/grpc-ecosystem/grpc-gateway/issues/1084)
- extract protoc-gen-swagger to separate repo [\#1083](https://github.com/grpc-ecosystem/grpc-gateway/issues/1083)
-  How do I load balance call grpc service? [\#1081](https://github.com/grpc-ecosystem/grpc-gateway/issues/1081)

## [v1.12.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.12.1) (2019-11-06)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.12.0...v1.12.1)

**Closed issues:**

- Unable to create HTTP mapping with "/parent" [\#1079](https://github.com/grpc-ecosystem/grpc-gateway/issues/1079)

## [v1.12.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.12.0) (2019-11-04)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.11.3...v1.12.0)

**Implemented enhancements:**

- protoc-gen-swagger: support generating a Swagger definition with no schemes [\#1069](https://github.com/grpc-ecosystem/grpc-gateway/issues/1069)

**Fixed bugs:**

- "make test" fails with go mod error [\#895](https://github.com/grpc-ecosystem/grpc-gateway/issues/895)

**Closed issues:**

- jfbrandhorst/grpc-gateway-build-env image can't run on Windows [\#1073](https://github.com/grpc-ecosystem/grpc-gateway/issues/1073)
- EOF is received after one request [\#1071](https://github.com/grpc-ecosystem/grpc-gateway/issues/1071)
- grpc-ecosystem/grpc-gateway/third\_party/googleapis: warning: directory does not exist. [\#1068](https://github.com/grpc-ecosystem/grpc-gateway/issues/1068)
- third\_party/googleapis is missing from package [\#1065](https://github.com/grpc-ecosystem/grpc-gateway/issues/1065)
- handleForwardResponseOptions not called by DefaultHTTPError [\#1064](https://github.com/grpc-ecosystem/grpc-gateway/issues/1064)
- why marshal enum to json using string but received it with int . [\#1063](https://github.com/grpc-ecosystem/grpc-gateway/issues/1063)
- protoc-gen-swagger/genswagger does not build on go1.11 and earlier versions [\#1061](https://github.com/grpc-ecosystem/grpc-gateway/issues/1061)
- How to support custom output, implementation [\#1055](https://github.com/grpc-ecosystem/grpc-gateway/issues/1055)

## [v1.11.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.11.3) (2019-09-30)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.11.2...v1.11.3)

**Closed issues:**

- json Custom output support \(with examples\) [\#1051](https://github.com/grpc-ecosystem/grpc-gateway/issues/1051)
- Question: Override TransientFailure error with friendlier response [\#1047](https://github.com/grpc-ecosystem/grpc-gateway/issues/1047)
- Wrong codes generated when nested enum in path  [\#1017](https://github.com/grpc-ecosystem/grpc-gateway/issues/1017)

## [v1.11.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.11.2) (2019-09-20)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.11.1...v1.11.2)

**Implemented enhancements:**

- Support specifying servers in the swagger generator [\#891](https://github.com/grpc-ecosystem/grpc-gateway/issues/891)

**Fixed bugs:**

- Make protoc-gen-swagger build on 1.11 [\#1044](https://github.com/grpc-ecosystem/grpc-gateway/issues/1044)
- jsonpb panics when using numbers for parsing timestamps [\#1025](https://github.com/grpc-ecosystem/grpc-gateway/issues/1025)

**Closed issues:**

- Interceptors not called when using new RegisterHandler function [\#1043](https://github.com/grpc-ecosystem/grpc-gateway/issues/1043)
- How to use -grpc-gateway\_out sp that the result is in a specific folder? [\#1042](https://github.com/grpc-ecosystem/grpc-gateway/issues/1042)
- Is there any way to let json int32 can not accept string in grpc-gateway? [\#1029](https://github.com/grpc-ecosystem/grpc-gateway/issues/1029)
- Go integration tests are somewhat flaky [\#992](https://github.com/grpc-ecosystem/grpc-gateway/issues/992)

## [v1.11.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.11.1) (2019-09-02)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.11.0...v1.11.1)

**Fixed bugs:**

- protoc\_gen\_swagger openapiv2\_field definition ignores the type option [\#1002](https://github.com/grpc-ecosystem/grpc-gateway/issues/1002)

**Closed issues:**

- AnnotateIncomingContext not declared by package runtime [\#1023](https://github.com/grpc-ecosystem/grpc-gateway/issues/1023)
- Fuzzit CI job is failing unexpectedly [\#1019](https://github.com/grpc-ecosystem/grpc-gateway/issues/1019)
- Bazel Rule? [\#1010](https://github.com/grpc-ecosystem/grpc-gateway/issues/1010)

## [v1.11.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.11.0) (2019-08-30)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.10.0...v1.11.0)

**Fixed bugs:**

- protoc-gen-grpc-gateway fails silently after release 1.10 [\#1013](https://github.com/grpc-ecosystem/grpc-gateway/issues/1013)

**Closed issues:**

- protoc-gen-swagger does not generate parameters other than body and path parameters. [\#1012](https://github.com/grpc-ecosystem/grpc-gateway/issues/1012)

## [v1.10.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.10.0) (2019-08-28)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.6...v1.10.0)

**Implemented enhancements:**

- allow protobuf well known types in params [\#400](https://github.com/grpc-ecosystem/grpc-gateway/issues/400)

**Fixed bugs:**

- grpc-gateway don't work well when using github.com/golang/protobuf/ptypes/struct with streaming [\#999](https://github.com/grpc-ecosystem/grpc-gateway/issues/999)

**Closed issues:**

- Allow final url path parameter to be optional [\#1005](https://github.com/grpc-ecosystem/grpc-gateway/issues/1005)
- Update integration test dependencies [\#1004](https://github.com/grpc-ecosystem/grpc-gateway/issues/1004)
- Suggestion: Continuous Fuzzing [\#998](https://github.com/grpc-ecosystem/grpc-gateway/issues/998)
- Why grpc gateway does not call grpc callback directly? [\#952](https://github.com/grpc-ecosystem/grpc-gateway/issues/952)

## [v1.9.6](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.6) (2019-08-16)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.5...v1.9.6)

**Closed issues:**

- Returning a primitive type as a response instead of proto messages  [\#994](https://github.com/grpc-ecosystem/grpc-gateway/issues/994)
- protoc-gen-swagger: fix description of google/protobuf/struct.proto types [\#989](https://github.com/grpc-ecosystem/grpc-gateway/issues/989)
- Swagger generator does not convert parameters in URLs to camel case when `json\_names\_for\_fields` is enable.  [\#986](https://github.com/grpc-ecosystem/grpc-gateway/issues/986)
- The release upload job is broken [\#981](https://github.com/grpc-ecosystem/grpc-gateway/issues/981)
- Schema and field name questions from a front end developer [\#980](https://github.com/grpc-ecosystem/grpc-gateway/issues/980)
- undefined: runtime.AssumeColonVerbOpt [\#978](https://github.com/grpc-ecosystem/grpc-gateway/issues/978)
- I want to know how to transfer http+proto to grpc. [\#977](https://github.com/grpc-ecosystem/grpc-gateway/issues/977)
- Is it possible to use protoc-gen-swagger options in my own protos? [\#976](https://github.com/grpc-ecosystem/grpc-gateway/issues/976)

## [v1.9.5](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.5) (2019-07-22)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.4...v1.9.5)

**Fixed bugs:**

- Non-standard use of 412 HTTP Status Code [\#972](https://github.com/grpc-ecosystem/grpc-gateway/issues/972)

**Closed issues:**

- why response use enum's name [\#970](https://github.com/grpc-ecosystem/grpc-gateway/issues/970)

## [v1.9.4](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.4) (2019-07-09)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.3...v1.9.4)

**Closed issues:**

- Read the Http Post Body  [\#921](https://github.com/grpc-ecosystem/grpc-gateway/issues/921)
- Swagger document generation, required field is invalid [\#665](https://github.com/grpc-ecosystem/grpc-gateway/issues/665)

## [v1.9.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.3) (2019-06-28)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.2...v1.9.3)

**Fixed bugs:**

- EOF when calling Send for client streams [\#961](https://github.com/grpc-ecosystem/grpc-gateway/issues/961)

**Closed issues:**

- Please make a new release! [\#963](https://github.com/grpc-ecosystem/grpc-gateway/issues/963)
- application/x-www-form-urlencoded support. [\#960](https://github.com/grpc-ecosystem/grpc-gateway/issues/960)
- Bazel files are out of date [\#955](https://github.com/grpc-ecosystem/grpc-gateway/issues/955)
- Configurable AllowUnknownFields in jsonpb? [\#448](https://github.com/grpc-ecosystem/grpc-gateway/issues/448)

## [v1.9.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.2) (2019-06-17)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.1...v1.9.2)

**Fixed bugs:**

- 404s using colons in the middle of the last path segment [\#224](https://github.com/grpc-ecosystem/grpc-gateway/issues/224)

## [v1.9.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.1) (2019-06-13)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.9.0...v1.9.1)

**Closed issues:**

- grpc: received message larger than max [\#943](https://github.com/grpc-ecosystem/grpc-gateway/issues/943)
- json 1.1 api support for grpc-ecosystem to use queryparams with filter [\#938](https://github.com/grpc-ecosystem/grpc-gateway/issues/938)
- i import a new gateway.Endpoint without recompile [\#937](https://github.com/grpc-ecosystem/grpc-gateway/issues/937)
- all SubConns are in TransientFailure [\#936](https://github.com/grpc-ecosystem/grpc-gateway/issues/936)
- Merging swagger specs fails to use rpc comments \(again\) [\#923](https://github.com/grpc-ecosystem/grpc-gateway/issues/923)

## [v1.9.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.9.0) (2019-05-14)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.6...v1.9.0)

**Closed issues:**

- Errors in response streams do not go through the registered error handler [\#584](https://github.com/grpc-ecosystem/grpc-gateway/issues/584)

## [v1.8.6](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.6) (2019-05-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.5...v1.8.6)

**Fixed bugs:**

- can't specify an empty path? [\#414](https://github.com/grpc-ecosystem/grpc-gateway/issues/414)

**Closed issues:**

- JSON stream response not available [\#926](https://github.com/grpc-ecosystem/grpc-gateway/issues/926)
- why google/api/http.proto annotations.proto Field Numbers is 72295728 ? [\#925](https://github.com/grpc-ecosystem/grpc-gateway/issues/925)
- Documentation: 'base\_path' Swagger attribute confuses users [\#918](https://github.com/grpc-ecosystem/grpc-gateway/issues/918)
- go get: error loading module requirements go 1.11 [\#915](https://github.com/grpc-ecosystem/grpc-gateway/issues/915)
- gateway generation issue on windows [\#911](https://github.com/grpc-ecosystem/grpc-gateway/issues/911)

## [v1.8.5](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.5) (2019-03-15)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.4...v1.8.5)

**Closed issues:**

- Swagger get query param documentation shows repeated fields incorrectly [\#756](https://github.com/grpc-ecosystem/grpc-gateway/issues/756)

## [v1.8.4](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.4) (2019-03-13)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.3...v1.8.4)

**Closed issues:**

- Invalid swagger generated for bodies with repeated fields [\#906](https://github.com/grpc-ecosystem/grpc-gateway/issues/906)

## [v1.8.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.3) (2019-03-11)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.2...v1.8.3)

**Implemented enhancements:**

- Feature request from openapi 3: Allow apiKey in cookie [\#900](https://github.com/grpc-ecosystem/grpc-gateway/issues/900)

**Fixed bugs:**

- Error while defining enum comments [\#897](https://github.com/grpc-ecosystem/grpc-gateway/issues/897)

**Closed issues:**

- Its impossible to send response with non 200 status code [\#901](https://github.com/grpc-ecosystem/grpc-gateway/issues/901)

## [v1.8.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.2) (2019-03-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.1...v1.8.2)

**Implemented enhancements:**

- Update the build environment Dockerfile to Go 1.12 [\#885](https://github.com/grpc-ecosystem/grpc-gateway/issues/885)

**Fixed bugs:**

- Change in behavior of streaming request body \(1.4.1 vs 1.8.1\) [\#894](https://github.com/grpc-ecosystem/grpc-gateway/issues/894)
- Cannot download 1.8.0 with modules [\#886](https://github.com/grpc-ecosystem/grpc-gateway/issues/886)

**Closed issues:**

- Description and title ignored when field is not a scaler value type [\#892](https://github.com/grpc-ecosystem/grpc-gateway/issues/892)

## [v1.8.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.1) (2019-03-02)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.1-pre1...v1.8.1)

## [v1.8.1-pre1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.1-pre1) (2019-03-01)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.8.0...v1.8.1-pre1)

## [v1.8.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.8.0) (2019-03-01)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.7.0...v1.8.0)

**Implemented enhancements:**

- Support swagger annotations for default and required fields [\#851](https://github.com/grpc-ecosystem/grpc-gateway/issues/851)
- Support go modules [\#755](https://github.com/grpc-ecosystem/grpc-gateway/issues/755)

**Fixed bugs:**

- inconsistent identifier capitalization between protoc-gen-go and protoc-gen-grpc-gateway [\#683](https://github.com/grpc-ecosystem/grpc-gateway/issues/683)

**Closed issues:**

- Bazel incompatible changes [\#873](https://github.com/grpc-ecosystem/grpc-gateway/issues/873)
- Swagger has not existed for four years [\#872](https://github.com/grpc-ecosystem/grpc-gateway/issues/872)
- Improve README with AWS API Gateway findings [\#868](https://github.com/grpc-ecosystem/grpc-gateway/issues/868)
- swagger error [\#867](https://github.com/grpc-ecosystem/grpc-gateway/issues/867)
- A question about generating file protoc-gen-grpc-gateway/gengateway/template.go [\#864](https://github.com/grpc-ecosystem/grpc-gateway/issues/864)
- Repeated field documentation is overwritten by fields comments [\#863](https://github.com/grpc-ecosystem/grpc-gateway/issues/863)
- Using dep to depend on specific revision of golang/protobuf is causing transative dependency problems for users [\#829](https://github.com/grpc-ecosystem/grpc-gateway/issues/829)
- Mac OS X - Note about your tutorial [\#787](https://github.com/grpc-ecosystem/grpc-gateway/issues/787)
- Returning 302 redirect as response [\#607](https://github.com/grpc-ecosystem/grpc-gateway/issues/607)

## [v1.7.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.7.0) (2019-01-23)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.6.4...v1.7.0)

**Closed issues:**

- Error to build project with go module [\#846](https://github.com/grpc-ecosystem/grpc-gateway/issues/846)
- Result of gateway's Stream response is wrapped with "result" [\#579](https://github.com/grpc-ecosystem/grpc-gateway/issues/579)

## [v1.6.4](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.6.4) (2019-01-08)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.6.3...v1.6.4)

**Closed issues:**

- feature request: opt-out fieldmask behaviour in patch [\#839](https://github.com/grpc-ecosystem/grpc-gateway/issues/839)
- gRPC streaming keepAlive doesn't work with docker swarm [\#838](https://github.com/grpc-ecosystem/grpc-gateway/issues/838)

## [v1.6.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.6.3) (2018-12-21)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.6.2...v1.6.3)

**Closed issues:**

- Issue with google.protobuf.Empty representation in swagger file [\#831](https://github.com/grpc-ecosystem/grpc-gateway/issues/831)

## [v1.6.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.6.2) (2018-12-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.6.0...v1.6.2)

## [v1.6.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.6.0) (2018-12-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.6.1...v1.6.0)

## [v1.6.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.6.1) (2018-12-07)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.5.1...v1.6.1)

**Fixed bugs:**

- Cannot return HTTP header using "Grpc-Metadata-" prefix [\#782](https://github.com/grpc-ecosystem/grpc-gateway/issues/782)
- protoc-gen-swagger  throws the error: Any JSON doesn't have '@type' [\#771](https://github.com/grpc-ecosystem/grpc-gateway/issues/771)
- proto-gen-swagger: provide default description for HTTP 200 responses [\#766](https://github.com/grpc-ecosystem/grpc-gateway/issues/766)

**Closed issues:**

- Please release the repo, IOReaderFactory is not available on the latest release! [\#823](https://github.com/grpc-ecosystem/grpc-gateway/issues/823)
- Bazel CI breaks frequently [\#817](https://github.com/grpc-ecosystem/grpc-gateway/issues/817)
- Unable to add protobuf wrappers in url template option [\#808](https://github.com/grpc-ecosystem/grpc-gateway/issues/808)
- Class 'GPBMetadata\ProtocGenSwagger\Options\Annotations' not found [\#794](https://github.com/grpc-ecosystem/grpc-gateway/issues/794)
- REST gateway over RPCS? [\#789](https://github.com/grpc-ecosystem/grpc-gateway/issues/789)
- Why the rctx is substituted by a new empty context? [\#788](https://github.com/grpc-ecosystem/grpc-gateway/issues/788)
- grpc gateway intercepter [\#785](https://github.com/grpc-ecosystem/grpc-gateway/issues/785)
- "error" and "message" fields in error response [\#768](https://github.com/grpc-ecosystem/grpc-gateway/issues/768)
- Go1.11: http.CloseNotifier is deprecated [\#736](https://github.com/grpc-ecosystem/grpc-gateway/issues/736)
- Access to raw post body [\#652](https://github.com/grpc-ecosystem/grpc-gateway/issues/652)

## [v1.5.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.5.1) (2018-10-02)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.5.0...v1.5.1)

**Implemented enhancements:**

- protobuf well known types aren't represented in swagger output correctly [\#160](https://github.com/grpc-ecosystem/grpc-gateway/issues/160)

**Fixed bugs:**

- URLs using verb no longer work after upgrading to v1.5.0 [\#760](https://github.com/grpc-ecosystem/grpc-gateway/issues/760)
- protoc-gen-swagger doesn't generate any request objects for GET/DELETE [\#747](https://github.com/grpc-ecosystem/grpc-gateway/issues/747)

**Closed issues:**

- how to get proper fields name for method [\#745](https://github.com/grpc-ecosystem/grpc-gateway/issues/745)
- Make a new release [\#733](https://github.com/grpc-ecosystem/grpc-gateway/issues/733)
- how to provide interface type inside proto for grpc-gateway [\#723](https://github.com/grpc-ecosystem/grpc-gateway/issues/723)
- Is there any way we can remove fields from the response json in grpc-gateway? [\#710](https://github.com/grpc-ecosystem/grpc-gateway/issues/710)
- How to write tests for the gateway? [\#699](https://github.com/grpc-ecosystem/grpc-gateway/issues/699)
- protoc-gen-swagger: No comments for path parameters [\#694](https://github.com/grpc-ecosystem/grpc-gateway/issues/694)
- Can you differentiate between an empty map vs field not provided? [\#552](https://github.com/grpc-ecosystem/grpc-gateway/issues/552)
- import\_path option not working as intended [\#536](https://github.com/grpc-ecosystem/grpc-gateway/issues/536)

## [v1.5.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.5.0) (2018-09-09)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.4.1...v1.5.0)

**Fixed bugs:**

- forwarding binary metadata is broken [\#218](https://github.com/grpc-ecosystem/grpc-gateway/issues/218)

**Closed issues:**

- something wrong with service [\#748](https://github.com/grpc-ecosystem/grpc-gateway/issues/748)
- Support for repeated path parameters [\#741](https://github.com/grpc-ecosystem/grpc-gateway/issues/741)
- Uint64 is represented as type:"string" in the swagger docs. [\#735](https://github.com/grpc-ecosystem/grpc-gateway/issues/735)
- autoregister all provided services [\#732](https://github.com/grpc-ecosystem/grpc-gateway/issues/732)
- `go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway` fails on clean environment [\#731](https://github.com/grpc-ecosystem/grpc-gateway/issues/731)
- format tool [\#729](https://github.com/grpc-ecosystem/grpc-gateway/issues/729)
- how to do tls auth in grpc+grpc-gateway [\#727](https://github.com/grpc-ecosystem/grpc-gateway/issues/727)
- Let service choose it's own marshaller [\#725](https://github.com/grpc-ecosystem/grpc-gateway/issues/725)
- why gateway proxy can not distribute the http request to local server?  prompt 404 [\#722](https://github.com/grpc-ecosystem/grpc-gateway/issues/722)
- enc.SetIndent undefined \(type \*json.Encoder has no field or method SetIndent\) [\#717](https://github.com/grpc-ecosystem/grpc-gateway/issues/717)
- Travis CI fails on master branch [\#714](https://github.com/grpc-ecosystem/grpc-gateway/issues/714)
- google/protobuf/descriptor.proto: File not found. ? [\#713](https://github.com/grpc-ecosystem/grpc-gateway/issues/713)
- APIs with grpc-gateway \(S3,WebDav\) [\#709](https://github.com/grpc-ecosystem/grpc-gateway/issues/709)
- FR: Promote a field in the returned JSON message to a top-level returned value [\#707](https://github.com/grpc-ecosystem/grpc-gateway/issues/707)
- Does grpc-gateway support the HTTP 2.0 protocol? [\#703](https://github.com/grpc-ecosystem/grpc-gateway/issues/703)
- The swagger plugin couldn’t distinguish two rpcs if we use the resource name design style. [\#702](https://github.com/grpc-ecosystem/grpc-gateway/issues/702)
- Handling of optional parameters [\#697](https://github.com/grpc-ecosystem/grpc-gateway/issues/697)
- Vendor dependencies [\#689](https://github.com/grpc-ecosystem/grpc-gateway/issues/689)
- Output swagger seems incorrect [\#688](https://github.com/grpc-ecosystem/grpc-gateway/issues/688)
- how to use this in java? [\#685](https://github.com/grpc-ecosystem/grpc-gateway/issues/685)
- r [\#684](https://github.com/grpc-ecosystem/grpc-gateway/issues/684)
- url query parameters should support semicolon in value field [\#680](https://github.com/grpc-ecosystem/grpc-gateway/issues/680)
- how to install swagger-codegen@2.2.2? [\#670](https://github.com/grpc-ecosystem/grpc-gateway/issues/670)
- Impossible to use gogo/protobuf registered types in gRPC Status errors [\#576](https://github.com/grpc-ecosystem/grpc-gateway/issues/576)
- Path parameters can't have URL encoded values [\#566](https://github.com/grpc-ecosystem/grpc-gateway/issues/566)
- docs: show example of tracing over http-\>grpc boundary [\#348](https://github.com/grpc-ecosystem/grpc-gateway/issues/348)
- Response codes and descriptions in Swagger docs [\#304](https://github.com/grpc-ecosystem/grpc-gateway/issues/304)

## [v1.4.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.4.1) (2018-05-23)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.4.0...v1.4.1)

**Closed issues:**

- Next release ? [\#605](https://github.com/grpc-ecosystem/grpc-gateway/issues/605)

## [v1.4.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.4.0) (2018-05-20)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.3.1...v1.4.0)

**Implemented enhancements:**

- customize the error return [\#405](https://github.com/grpc-ecosystem/grpc-gateway/issues/405)
- Support map type in query string [\#316](https://github.com/grpc-ecosystem/grpc-gateway/issues/316)
- gRPC gateway Bazel build rules [\#66](https://github.com/grpc-ecosystem/grpc-gateway/issues/66)
- Support bytes fields in path parameter [\#5](https://github.com/grpc-ecosystem/grpc-gateway/issues/5)

**Closed issues:**

- the protoc\_gen\_swagger bazel rule generates non working import path. [\#633](https://github.com/grpc-ecosystem/grpc-gateway/issues/633)
- code.NotFound should return a 404 instead of a 405 [\#630](https://github.com/grpc-ecosystem/grpc-gateway/issues/630)
- field in query path not found [\#629](https://github.com/grpc-ecosystem/grpc-gateway/issues/629)
- how to use client pool in the gateway? [\#612](https://github.com/grpc-ecosystem/grpc-gateway/issues/612)
- pass http request uri to grpc server [\#587](https://github.com/grpc-ecosystem/grpc-gateway/issues/587)
- bidi streams have racy read caused by goroutine that closes over local variable [\#583](https://github.com/grpc-ecosystem/grpc-gateway/issues/583)
- Streamed response is not valid json \(or: is this the expected format?\) [\#581](https://github.com/grpc-ecosystem/grpc-gateway/issues/581)
- Import "google/api/annotations.proto" was not found or had errors. [\#574](https://github.com/grpc-ecosystem/grpc-gateway/issues/574)
- is there has a way to let grpc-gateway server support multiple endpoints [\#573](https://github.com/grpc-ecosystem/grpc-gateway/issues/573)
- would it be possible to avoid vendoring "third\_party/googleapis/" [\#572](https://github.com/grpc-ecosystem/grpc-gateway/issues/572)
- Is there anyway to output the access log of grpc gateway [\#556](https://github.com/grpc-ecosystem/grpc-gateway/issues/556)
- proto: no slice oenc for \*reflect.rtype = \[\]\*reflect.rtype [\#551](https://github.com/grpc-ecosystem/grpc-gateway/issues/551)
- autoreconf not found [\#549](https://github.com/grpc-ecosystem/grpc-gateway/issues/549)
- \[feature\]combine expvar into grpc-gateway [\#542](https://github.com/grpc-ecosystem/grpc-gateway/issues/542)
- Source code still imports "golang.org/x/net/context" [\#533](https://github.com/grpc-ecosystem/grpc-gateway/issues/533)
- Incorrect error message when execute protoc-gen-grpc-gateway to HTTP GET method with BODY [\#531](https://github.com/grpc-ecosystem/grpc-gateway/issues/531)
- add support for the google.api.HttpBody proto as a request [\#528](https://github.com/grpc-ecosystem/grpc-gateway/issues/528)
- Prefixed model names in generated swagger spec [\#525](https://github.com/grpc-ecosystem/grpc-gateway/issues/525)
- Better format for error.message in stream [\#519](https://github.com/grpc-ecosystem/grpc-gateway/issues/519)
- Getting this on go get . in the src directory: HelloService.pb.go:20:8 - no Go files in \go\src\google\api [\#518](https://github.com/grpc-ecosystem/grpc-gateway/issues/518)
- ci: set up codecov [\#513](https://github.com/grpc-ecosystem/grpc-gateway/issues/513)
- protoc-gen-swagger not using description field of info swagger object [\#511](https://github.com/grpc-ecosystem/grpc-gateway/issues/511)
- Cut a minor release for https://github.com/grpc-ecosystem/grpc-gateway/issues/495 [\#506](https://github.com/grpc-ecosystem/grpc-gateway/issues/506)
- bug: uncapitalized service name causes runtime error unknown function in service.pb.gw.go [\#484](https://github.com/grpc-ecosystem/grpc-gateway/issues/484)
- RESOURCE\_EXHAUSTED -\> 503 [\#431](https://github.com/grpc-ecosystem/grpc-gateway/issues/431)
- Adding authentication definitions to generated swagger files [\#428](https://github.com/grpc-ecosystem/grpc-gateway/issues/428)
- Move to stdlib context over x/net/context [\#326](https://github.com/grpc-ecosystem/grpc-gateway/issues/326)
- deprecate 1.6 and embrace \(\*http.Request\).Context by default [\#313](https://github.com/grpc-ecosystem/grpc-gateway/issues/313)

## [v1.3.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.3.1) (2017-12-23)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.3.0...v1.3.1)

**Implemented enhancements:**

- Support import\_path? [\#443](https://github.com/grpc-ecosystem/grpc-gateway/issues/443)

**Closed issues:**

- protoc-gen-swagger missing definition issue [\#504](https://github.com/grpc-ecosystem/grpc-gateway/issues/504)
- Backwards incompatible change to chunked encoding [\#495](https://github.com/grpc-ecosystem/grpc-gateway/issues/495)
- Map of list [\#493](https://github.com/grpc-ecosystem/grpc-gateway/issues/493)
- How to run `makefile` for this repo? [\#491](https://github.com/grpc-ecosystem/grpc-gateway/issues/491)
- all SubConns are in TransientFailure [\#490](https://github.com/grpc-ecosystem/grpc-gateway/issues/490)
- Appengine Standard Environment: "not an Appengine context" [\#487](https://github.com/grpc-ecosystem/grpc-gateway/issues/487)
- Enum Path Parameter to Swagger [\#486](https://github.com/grpc-ecosystem/grpc-gateway/issues/486)
- Should v1.3 be also tagged as v1.3.0? [\#483](https://github.com/grpc-ecosystem/grpc-gateway/issues/483)
- HTTP response is not correct json encoded if the grpc return stream of objects. [\#481](https://github.com/grpc-ecosystem/grpc-gateway/issues/481)
- Support JSON-RPCv2 [\#477](https://github.com/grpc-ecosystem/grpc-gateway/issues/477)
- Naming convention? [\#475](https://github.com/grpc-ecosystem/grpc-gateway/issues/475)
- Request context not being used [\#470](https://github.com/grpc-ecosystem/grpc-gateway/issues/470)
- Generate Swagger documentation [\#469](https://github.com/grpc-ecosystem/grpc-gateway/issues/469)
- Support Request | make: swagger-codegen: Command not found [\#468](https://github.com/grpc-ecosystem/grpc-gateway/issues/468)
- How do you generate a swagger yaml file instead of json? [\#467](https://github.com/grpc-ecosystem/grpc-gateway/issues/467)
- Add default support for proto over http [\#465](https://github.com/grpc-ecosystem/grpc-gateway/issues/465)
- Allow compiling the gateway code to a different go package [\#463](https://github.com/grpc-ecosystem/grpc-gateway/issues/463)
- support google.api.HttpBody [\#457](https://github.com/grpc-ecosystem/grpc-gateway/issues/457)
- \[swagger bug\] with google/protobuf/wrappers.proto [\#453](https://github.com/grpc-ecosystem/grpc-gateway/issues/453)
- The tensorflow serving support RESTful api：{"error":"json: cannot unmarshal object into Go value of type \[\]json.RawMessage","code":3} [\#444](https://github.com/grpc-ecosystem/grpc-gateway/issues/444)
- choose some return fields omit or  not omit by configure [\#439](https://github.com/grpc-ecosystem/grpc-gateway/issues/439)
- swagger title and version hardcoded [\#437](https://github.com/grpc-ecosystem/grpc-gateway/issues/437)
- Change the path though http header [\#424](https://github.com/grpc-ecosystem/grpc-gateway/issues/424)
- google/protobuf/descriptor.proto: File not found [\#422](https://github.com/grpc-ecosystem/grpc-gateway/issues/422)
- Output file will not compile if the .proto file does not contain a service with parameters in the url path [\#389](https://github.com/grpc-ecosystem/grpc-gateway/issues/389)
- Scaling support [\#381](https://github.com/grpc-ecosystem/grpc-gateway/issues/381)
- I cannot get the default value from client side [\#380](https://github.com/grpc-ecosystem/grpc-gateway/issues/380)
- Problem with Generated annotations.proto file [\#377](https://github.com/grpc-ecosystem/grpc-gateway/issues/377)
- Release 1.3.0 [\#357](https://github.com/grpc-ecosystem/grpc-gateway/issues/357)
- swagger: Unclear comments' parser behaviour [\#352](https://github.com/grpc-ecosystem/grpc-gateway/issues/352)
- Support semicolon syntax in go\_package protobuf option [\#341](https://github.com/grpc-ecosystem/grpc-gateway/issues/341)
- Add SOAP proxy [\#339](https://github.com/grpc-ecosystem/grpc-gateway/issues/339)
- Support combination of query params and body for POSTs with body: "\*" [\#234](https://github.com/grpc-ecosystem/grpc-gateway/issues/234)
- Interceptor [\#221](https://github.com/grpc-ecosystem/grpc-gateway/issues/221)

## [v1.3.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.3.0) (2017-11-03)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.3...v1.3.0)

## [v1.3](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.3) (2017-11-03)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.2.2...v1.3)

**Closed issues:**

- Extract basic auth from URL [\#480](https://github.com/grpc-ecosystem/grpc-gateway/issues/480)
- Lack of "google/protobuf/descriptor.proto" [\#476](https://github.com/grpc-ecosystem/grpc-gateway/issues/476)
- question: how to indicate whether call is through grpc gateway [\#456](https://github.com/grpc-ecosystem/grpc-gateway/issues/456)
- How to define this restful api using pb? [\#452](https://github.com/grpc-ecosystem/grpc-gateway/issues/452)
- how to output field as an array of json values? [\#449](https://github.com/grpc-ecosystem/grpc-gateway/issues/449)
- How do I override maxMsgSize? [\#445](https://github.com/grpc-ecosystem/grpc-gateway/issues/445)
- OpenAPI spec is generated with duplicated operation IDs. [\#442](https://github.com/grpc-ecosystem/grpc-gateway/issues/442)
- This process seems to generate conflicting code with go-micro [\#440](https://github.com/grpc-ecosystem/grpc-gateway/issues/440)
- any way to let int64 marshal to int not string? [\#438](https://github.com/grpc-ecosystem/grpc-gateway/issues/438)
- Support  streaming [\#435](https://github.com/grpc-ecosystem/grpc-gateway/issues/435)
- Update DO NOT EDIT header in generated files [\#433](https://github.com/grpc-ecosystem/grpc-gateway/issues/433)
- generate code use context not "golang.org/x/net/context" [\#430](https://github.com/grpc-ecosystem/grpc-gateway/issues/430)
- Replace \n with spaces in swagger definitions [\#426](https://github.com/grpc-ecosystem/grpc-gateway/issues/426)
- \[question\]Is there any example for  http headers process? [\#420](https://github.com/grpc-ecosystem/grpc-gateway/issues/420)
- Is there any way to support a multipart form request? [\#410](https://github.com/grpc-ecosystem/grpc-gateway/issues/410)
- Not able to pass allow\_delete\_body to protoc-gen-grpc-gateway. [\#402](https://github.com/grpc-ecosystem/grpc-gateway/issues/402)
- returned errors should conform to google.rpc.Status [\#399](https://github.com/grpc-ecosystem/grpc-gateway/issues/399)
- Is there any way to generate python gateway code? [\#398](https://github.com/grpc-ecosystem/grpc-gateway/issues/398)
- how to handle arbitrary \(json\) structs [\#395](https://github.com/grpc-ecosystem/grpc-gateway/issues/395)
- \[question\]can give a url with query sting demo? [\#394](https://github.com/grpc-ecosystem/grpc-gateway/issues/394)
- \[question\]the swagger url generated is what? [\#393](https://github.com/grpc-ecosystem/grpc-gateway/issues/393)
- \[Question\] How do I use semantic versions? [\#392](https://github.com/grpc-ecosystem/grpc-gateway/issues/392)
- \[question\]how to run examples? [\#391](https://github.com/grpc-ecosystem/grpc-gateway/issues/391)
- Why does gateway use ServerMetadata? [\#388](https://github.com/grpc-ecosystem/grpc-gateway/issues/388)
- Can't generate code with last version [\#384](https://github.com/grpc-ecosystem/grpc-gateway/issues/384)
- is it ready for production use? [\#382](https://github.com/grpc-ecosystem/grpc-gateway/issues/382)
- Support Google Flatbuffers [\#376](https://github.com/grpc-ecosystem/grpc-gateway/issues/376)
- calling Enum by string name in requests using gogo/protobuf results in error. [\#372](https://github.com/grpc-ecosystem/grpc-gateway/issues/372)
- Definitions containing URLs with trailing slashes won't compile [\#370](https://github.com/grpc-ecosystem/grpc-gateway/issues/370)
- Should metadata annotator include the headers from incoming matcher? [\#368](https://github.com/grpc-ecosystem/grpc-gateway/issues/368)
-  metadata.NewOutgoingContext is undefined [\#364](https://github.com/grpc-ecosystem/grpc-gateway/issues/364)
- Why does not gateway forward headers as-is? [\#311](https://github.com/grpc-ecosystem/grpc-gateway/issues/311)
- Question: Why passing context to RegisterMyServiceHandler is required?  [\#301](https://github.com/grpc-ecosystem/grpc-gateway/issues/301)
- Allow whitelisting of particular HTTP headers to map to metadata. [\#253](https://github.com/grpc-ecosystem/grpc-gateway/issues/253)
- Swagger definitions don't handle parameters that are not explicitly required in the url [\#159](https://github.com/grpc-ecosystem/grpc-gateway/issues/159)

## [v1.2.2](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.2.2) (2017-04-17)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.2.1...v1.2.2)

## [v1.2.1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.2.1) (2017-04-17)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.2.0...v1.2.1)

**Fixed bugs:**

- reflect upstream grpc metadata api change [\#358](https://github.com/grpc-ecosystem/grpc-gateway/issues/358)

**Closed issues:**

- Empty value omitted [\#355](https://github.com/grpc-ecosystem/grpc-gateway/issues/355)
- Release 1.2.0 [\#340](https://github.com/grpc-ecosystem/grpc-gateway/issues/340)
- Cut another release [\#278](https://github.com/grpc-ecosystem/grpc-gateway/issues/278)

## [v1.2.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.2.0) (2017-03-31)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.2.0.rc1...v1.2.0)

**Closed issues:**

- Problem with \*.proto as "no buildable Go source files" [\#338](https://github.com/grpc-ecosystem/grpc-gateway/issues/338)
- Invalid import during code generation [\#337](https://github.com/grpc-ecosystem/grpc-gateway/issues/337)

## [v1.2.0.rc1](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.2.0.rc1) (2017-03-24)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.1.0...v1.2.0.rc1)

**Implemented enhancements:**

- Support for Any types [\#80](https://github.com/grpc-ecosystem/grpc-gateway/issues/80)

**Fixed bugs:**

- Support for multi-segment elements [\#122](https://github.com/grpc-ecosystem/grpc-gateway/issues/122)

**Closed issues:**

- Go get breaks with autogenerated code [\#331](https://github.com/grpc-ecosystem/grpc-gateway/issues/331)
- Fresh install no longer generates necessary `google/api/annotations.pb.go` & `google/api/http.pb.go` files. [\#327](https://github.com/grpc-ecosystem/grpc-gateway/issues/327)
- Panic with query parameters [\#324](https://github.com/grpc-ecosystem/grpc-gateway/issues/324)
- Swagger-UI query parameters for enum types are sent as strings [\#320](https://github.com/grpc-ecosystem/grpc-gateway/issues/320)
- hide the object name in the response [\#317](https://github.com/grpc-ecosystem/grpc-gateway/issues/317)
- Package imported but not used [\#310](https://github.com/grpc-ecosystem/grpc-gateway/issues/310)
- Authorization headers aren't specified in Swagger.json [\#309](https://github.com/grpc-ecosystem/grpc-gateway/issues/309)
- Generating swagger version, contact name etc in generated docs [\#303](https://github.com/grpc-ecosystem/grpc-gateway/issues/303)
- Feature request: custom content type per service and rpc [\#302](https://github.com/grpc-ecosystem/grpc-gateway/issues/302)
- Reference: another RESTful api-gateway [\#299](https://github.com/grpc-ecosystem/grpc-gateway/issues/299)
- Integration with other languages is partially broken [\#298](https://github.com/grpc-ecosystem/grpc-gateway/issues/298)
- jsonpb convert int64 to integer instead of string [\#296](https://github.com/grpc-ecosystem/grpc-gateway/issues/296)
- default enum value is omitted [\#294](https://github.com/grpc-ecosystem/grpc-gateway/issues/294)
- Advice: could we simplify the flow as the below [\#292](https://github.com/grpc-ecosystem/grpc-gateway/issues/292)
- examples/browser test failure: TypeError: undefined is not a function \(evaluating 'window.location.protocol.startsWith\('chrome-extension'\)'\) [\#287](https://github.com/grpc-ecosystem/grpc-gateway/issues/287)
- ./entrypoint.go:25: undefined: api.RegisterYourServiceHandlerFromEndpoint [\#285](https://github.com/grpc-ecosystem/grpc-gateway/issues/285)
- Query params not handled in swagger file [\#284](https://github.com/grpc-ecosystem/grpc-gateway/issues/284)
- Please help: google/api/annotations.proto: File not found. [\#283](https://github.com/grpc-ecosystem/grpc-gateway/issues/283)
- Option to Allow Swagger for DELETEs with a body [\#279](https://github.com/grpc-ecosystem/grpc-gateway/issues/279)
- client declared and not used compilation error, after recent upgrade [\#276](https://github.com/grpc-ecosystem/grpc-gateway/issues/276)
- feature request / idea: generating JSONRPC2 client proxies from GRPC [\#272](https://github.com/grpc-ecosystem/grpc-gateway/issues/272)
- protoc-swagger-generator messes up the comments if there is rpc method that does not have rest [\#263](https://github.com/grpc-ecosystem/grpc-gateway/issues/263)
- Swagger Gen: underscores -\> lowerCamelCase field names and refs [\#261](https://github.com/grpc-ecosystem/grpc-gateway/issues/261)
- Timestamp as URL param causes bad request error [\#260](https://github.com/grpc-ecosystem/grpc-gateway/issues/260)
- "proto: no coders for int" printed whenever a gRPC error is returned over grpc-gateway. [\#259](https://github.com/grpc-ecosystem/grpc-gateway/issues/259)
- Compatibility with grpc.SupportPackageIsVersion4 [\#258](https://github.com/grpc-ecosystem/grpc-gateway/issues/258)
- How to use circuit breaker in this grpc gateway? [\#257](https://github.com/grpc-ecosystem/grpc-gateway/issues/257)
- cannot use example code to generate [\#255](https://github.com/grpc-ecosystem/grpc-gateway/issues/255)
- tests fail on go tip due to importing of main packages in test [\#250](https://github.com/grpc-ecosystem/grpc-gateway/issues/250)
- Add NGINX support [\#249](https://github.com/grpc-ecosystem/grpc-gateway/issues/249)
- Error when reverse proxy to gRPC server \(which is impl with Node.js\) [\#246](https://github.com/grpc-ecosystem/grpc-gateway/issues/246)
- Error output titlecase instead of lowercase [\#243](https://github.com/grpc-ecosystem/grpc-gateway/issues/243)
- Option field "\(google.api.http\)" is not a field or extension of message "ServiceOptions" [\#241](https://github.com/grpc-ecosystem/grpc-gateway/issues/241)
- Implement credentials handler in-box [\#238](https://github.com/grpc-ecosystem/grpc-gateway/issues/238)
- Proposal: Support WKT structs for URL params [\#237](https://github.com/grpc-ecosystem/grpc-gateway/issues/237)
- Example of /} in path template [\#232](https://github.com/grpc-ecosystem/grpc-gateway/issues/232)
- Serving swagger.json from runtime mux? [\#230](https://github.com/grpc-ecosystem/grpc-gateway/issues/230)
- ETCDclientv3 build error with the latest changes - github.com/grpc-ecosystem/grpc-gateway/runtime/marshal\_jsonpb.go:114: undefined: jsonpb.Unmarshaler [\#226](https://github.com/grpc-ecosystem/grpc-gateway/issues/226)
- Map in GET request [\#223](https://github.com/grpc-ecosystem/grpc-gateway/issues/223)
- HTTPS no longer works [\#220](https://github.com/grpc-ecosystem/grpc-gateway/issues/220)
- --swagger\_out plugin translates proto type int64 to string in Swagger specification [\#219](https://github.com/grpc-ecosystem/grpc-gateway/issues/219)
- Response body as a single field [\#217](https://github.com/grpc-ecosystem/grpc-gateway/issues/217)
- documentation of semantics of endpoint declarations [\#212](https://github.com/grpc-ecosystem/grpc-gateway/issues/212)
- gen-swagger does not generate PATCH method endpoints [\#211](https://github.com/grpc-ecosystem/grpc-gateway/issues/211)
- protoc-gen-grpc-gateway doesn't work correctly with option go\_package [\#207](https://github.com/grpc-ecosystem/grpc-gateway/issues/207)
- Browser Side Streaming Best Practices [\#206](https://github.com/grpc-ecosystem/grpc-gateway/issues/206)
- Does grpc-gateway support App Engine? [\#204](https://github.com/grpc-ecosystem/grpc-gateway/issues/204)
- "use of internal package" error, after moving to grpc-ecosystem [\#203](https://github.com/grpc-ecosystem/grpc-gateway/issues/203)
- Move to google.golang.org/genproto instead of shipping annotations.proto. [\#202](https://github.com/grpc-ecosystem/grpc-gateway/issues/202)
- Release v1.1.0 [\#196](https://github.com/grpc-ecosystem/grpc-gateway/issues/196)
- marshaler runtime.Marshaler does not handle io.EOF when decoding [\#195](https://github.com/grpc-ecosystem/grpc-gateway/issues/195)
- protobuf enumerated values now returned as strings instead of numbers. [\#186](https://github.com/grpc-ecosystem/grpc-gateway/issues/186)
- support annotating fields as required \(in swagger/oapi generation\)? [\#175](https://github.com/grpc-ecosystem/grpc-gateway/issues/175)
- architectural question: Can i codegen the client code for talking to the server ? [\#167](https://github.com/grpc-ecosystem/grpc-gateway/issues/167)
- Passing ENUM value as URL parameter throws error [\#166](https://github.com/grpc-ecosystem/grpc-gateway/issues/166)
- Support specifying which schemes should be output in swagger.json [\#161](https://github.com/grpc-ecosystem/grpc-gateway/issues/161)
- Use headers for routing [\#157](https://github.com/grpc-ecosystem/grpc-gateway/issues/157)
- ENUM in swagger.json makes client code failed to parse response from gateway [\#153](https://github.com/grpc-ecosystem/grpc-gateway/issues/153)
- Support map types [\#140](https://github.com/grpc-ecosystem/grpc-gateway/issues/140)
- generate OpenAPI/swagger documentation at run time? [\#138](https://github.com/grpc-ecosystem/grpc-gateway/issues/138)
- After the 1.7 release, update .travis.yaml to check the compiled proto output [\#137](https://github.com/grpc-ecosystem/grpc-gateway/issues/137)
- Getting parsed runtime.Pattern from server mux [\#127](https://github.com/grpc-ecosystem/grpc-gateway/issues/127)
- REST API without proxying [\#46](https://github.com/grpc-ecosystem/grpc-gateway/issues/46)

## [v1.1.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.1.0) (2016-07-23)
[Full Changelog](https://github.com/grpc-ecosystem/grpc-gateway/compare/v1.0.0...v1.1.0)

**Implemented enhancements:**

- Support oneof types of fields [\#82](https://github.com/grpc-ecosystem/grpc-gateway/issues/82)
- allow use of jsonpb for marshaling [\#79](https://github.com/grpc-ecosystem/grpc-gateway/issues/79)

**Closed issues:**

- Generating a gRPC stub using Gateway generates a gRPC internal error [\#198](https://github.com/grpc-ecosystem/grpc-gateway/issues/198)
- Build fails with error: use of internal package not allowed [\#197](https://github.com/grpc-ecosystem/grpc-gateway/issues/197)
- google/protobuf/descriptor.proto: File not found. [\#194](https://github.com/grpc-ecosystem/grpc-gateway/issues/194)
- please tag releases [\#189](https://github.com/grpc-ecosystem/grpc-gateway/issues/189)
- Support for path collapsing for embedded structs? [\#187](https://github.com/grpc-ecosystem/grpc-gateway/issues/187)
- \[ACTION Required\] Moving to grpc-ecosystem [\#179](https://github.com/grpc-ecosystem/grpc-gateway/issues/179)
- Ading grpc-timeout support [\#107](https://github.com/grpc-ecosystem/grpc-gateway/issues/107)
- Generation of one swagger file out of multiple protos? [\#99](https://github.com/grpc-ecosystem/grpc-gateway/issues/99)

## [v1.0.0](https://github.com/grpc-ecosystem/grpc-gateway/tree/v1.0.0) (2016-06-15)
**Implemented enhancements:**

- support protobuf-over-HTTP [\#124](https://github.com/grpc-ecosystem/grpc-gateway/issues/124)
- Static mapping from proto field names to golang field names [\#86](https://github.com/grpc-ecosystem/grpc-gateway/issues/86)
- Format Errors to JSON [\#25](https://github.com/grpc-ecosystem/grpc-gateway/issues/25)
- Emit API definition in Swagger schema format [\#9](https://github.com/grpc-ecosystem/grpc-gateway/issues/9)
- Method parameter in query string [\#6](https://github.com/grpc-ecosystem/grpc-gateway/issues/6)
- Integrate authentication [\#4](https://github.com/grpc-ecosystem/grpc-gateway/issues/4)

**Fixed bugs:**

- recent annotation change requires req.RemoteAddr to be populated [\#177](https://github.com/grpc-ecosystem/grpc-gateway/issues/177)
- Runtime panic with CloseNotify [\#115](https://github.com/grpc-ecosystem/grpc-gateway/issues/115)
- Gateway code generation broken when rpc method with a streaming response has an input paramter [\#35](https://github.com/grpc-ecosystem/grpc-gateway/issues/35)
- URL usage of nested messages causes nil pointer in proto3 [\#32](https://github.com/grpc-ecosystem/grpc-gateway/issues/32)
- Multiple .proto files generates invalid import statements. [\#22](https://github.com/grpc-ecosystem/grpc-gateway/issues/22)

**Closed issues:**

- remote peer address is lost in ctx - always resolves to localhost [\#173](https://github.com/grpc-ecosystem/grpc-gateway/issues/173)
- Bidirectional streams don't concurrently Send and Recv [\#169](https://github.com/grpc-ecosystem/grpc-gateway/issues/169)
- Error: failed to import google/api/annotations.proto [\#165](https://github.com/grpc-ecosystem/grpc-gateway/issues/165)
- Test datarace in controlapi [\#163](https://github.com/grpc-ecosystem/grpc-gateway/issues/163)
- not enough arguments in call to runtime.HTTPError [\#162](https://github.com/grpc-ecosystem/grpc-gateway/issues/162)
- String-values for Enums in request object are not recognized. [\#150](https://github.com/grpc-ecosystem/grpc-gateway/issues/150)
- Handling of import public "file.proto" [\#139](https://github.com/grpc-ecosystem/grpc-gateway/issues/139)
- Does grpc-gateway support http middleware? [\#132](https://github.com/grpc-ecosystem/grpc-gateway/issues/132)
- push to web clients using WS or SSE ? [\#131](https://github.com/grpc-ecosystem/grpc-gateway/issues/131)
- protoc-gen-swagger comment parsing for documentation gen [\#128](https://github.com/grpc-ecosystem/grpc-gateway/issues/128)
- generated code has a data race [\#123](https://github.com/grpc-ecosystem/grpc-gateway/issues/123)
- panic: net/http: CloseNotify called after ServeHTTP finished [\#121](https://github.com/grpc-ecosystem/grpc-gateway/issues/121)
- CloseNotify race with ServeHTTP [\#119](https://github.com/grpc-ecosystem/grpc-gateway/issues/119)
- echo service example does not compile [\#117](https://github.com/grpc-ecosystem/grpc-gateway/issues/117)
- go vet issues in template\_test.go [\#113](https://github.com/grpc-ecosystem/grpc-gateway/issues/113)
- undefined: proto.SizeVarint [\#103](https://github.com/grpc-ecosystem/grpc-gateway/issues/103)
- Closing the HTTP connection does not cancel the Context [\#101](https://github.com/grpc-ecosystem/grpc-gateway/issues/101)
- Logging [\#92](https://github.com/grpc-ecosystem/grpc-gateway/issues/92)
- Missing default values in JSON output? [\#91](https://github.com/grpc-ecosystem/grpc-gateway/issues/91)
- Better grpc error strings [\#87](https://github.com/grpc-ecosystem/grpc-gateway/issues/87)
- Fields aren't named in the same manner as golang/protobuf [\#84](https://github.com/grpc-ecosystem/grpc-gateway/issues/84)
- Header Forwarding from server. [\#73](https://github.com/grpc-ecosystem/grpc-gateway/issues/73)
- No pattern specified in google.api.HttpRule [\#70](https://github.com/grpc-ecosystem/grpc-gateway/issues/70)
- cannot find package "google/api" [\#67](https://github.com/grpc-ecosystem/grpc-gateway/issues/67)
- Generated .pb.go with services no longer works with latest version of grpc-go. [\#62](https://github.com/grpc-ecosystem/grpc-gateway/issues/62)
- JavaScript Proxy [\#61](https://github.com/grpc-ecosystem/grpc-gateway/issues/61)
- Add HTTP error code, error status to responseStreamChunk Error [\#58](https://github.com/grpc-ecosystem/grpc-gateway/issues/58)
- Reverse the code gen idea [\#44](https://github.com/grpc-ecosystem/grpc-gateway/issues/44)
- array of maps in json [\#43](https://github.com/grpc-ecosystem/grpc-gateway/issues/43)
- Examples break with 1.5 because of import of "main" examples package  [\#37](https://github.com/grpc-ecosystem/grpc-gateway/issues/37)
- Breaks with 1.5rc1 due to "internal" package name. [\#36](https://github.com/grpc-ecosystem/grpc-gateway/issues/36)
- Feature Request: Support for non-nullable nested messages. [\#20](https://github.com/grpc-ecosystem/grpc-gateway/issues/20)
- Is PascalFromSnake the right conversion to be doing? [\#19](https://github.com/grpc-ecosystem/grpc-gateway/issues/19)
- Infinite loop in generator when package name conflicts [\#17](https://github.com/grpc-ecosystem/grpc-gateway/issues/17)
- google.api.http options in multi-line format not supported [\#16](https://github.com/grpc-ecosystem/grpc-gateway/issues/16)
- Is there any plan to developing a C++ version? [\#15](https://github.com/grpc-ecosystem/grpc-gateway/issues/15)



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*
