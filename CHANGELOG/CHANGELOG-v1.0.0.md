# CHANGELOG v1.0.0

## Go 版本升级

- Go 版本从 **1.13** 升级至 **1.26.3**
- 新增 `godebug default=go1.26.3` 兼容性指令

## 安全漏洞修复

| CVE/问题 | 依赖 | 旧版本 | 新版本 |
|----------|------|--------|--------|
| CVE-2020-8911, CVE-2020-8912 | `aws-sdk-go` | v1.27.0 | 迁移至 v2 |
| CVE-2021-20329 | `mongo-driver` | v1.1.2 | 已消除（升级 `strfmt`） |
| CVE-2025-27144, CVE-2024-28180, WS-2023-0431 | `go-jose.v2` | v2.4.0 | 迁移至 `go-jose/v3` v3.0.5 |

## 废弃依赖迁移

### aws-sdk-go v1 → v2

- `github.com/aws/aws-sdk-go` v1.27.0 → `github.com/aws/aws-sdk-go-v2` v1.42.0（+ 独立模块 `credentials`、`service/s3`、`smithy-go`）
- 重写文件：
  - `pkg/client/internals3/s3.go` — `NewStaticCredentialsProvider`, `s3.NewFromConfig`, `UsePathStyle`
  - `pkg/repoiface/s3.go` — 所有 API 方法增加 `ctx` 参数，存储 `endpoint` 字段
  - `pkg/service/attachment/resource_control.go` — `ListObjectsV2` 替代 `ListObjectsWithContext`，`isNoSuchKey()` 使用 `smithy.APIError`
  - `pkg/service/attachment/handler.go` — `io.ReadAll` 替代 `ioutil.ReadAll`

### mongo-driver

- `go.mongodb.org/mongo-driver` v1.1.2 → 通过升级 `go-openapi/strfmt` 至 v0.26.3 消除该间接依赖
- 同时升级 `go-openapi/runtime` v0.19.7 → v0.32.3，`go-openapi/spec` v0.19.4 → v0.22.6，`go-openapi/validate` v0.19.5 → v0.26.0

### go-jose v2 → v3

- `gopkg.in/square/go-jose.v2` v2.4.0 → `github.com/go-jose/go-jose/v3` v3.0.5
- 修改文件：`pkg/util/jwtutil/jwt.go`（仅导入路径变更，API 兼容）

## k8s.io 全套升级

| 依赖 | 旧版本 | 新版本 |
|------|--------|--------|
| `k8s.io/api` | v0.18.4 | v0.36.2 |
| `k8s.io/apimachinery` | v0.18.4 | v0.36.2 |
| `k8s.io/client-go` | v0.18.4 | v0.36.2 |
| `k8s.io/apiextensions-apiserver` | v0.18.4 | v0.36.2 |
| `k8s.io/kubectl` | v0.18.4 | v0.36.2 |

### Helm 升级

- `helm.sh/helm/v3` 从 openpitrix fork（2020年）迁移至上游 **v3.21.1**
- 移除 fork 的 `replace` 指令

### etcd 迁移

- `go.etcd.io/etcd` v0.0.0-20200520232829 → `go.etcd.io/etcd/client/v3` v3.6.12 + `go.etcd.io/etcd/api/v3` v3.6.12
- 导入路径变更（7 个文件）：
  - `go.etcd.io/etcd/clientv3` → `go.etcd.io/etcd/client/v3`
  - `go.etcd.io/etcd/clientv3/concurrency` → `go.etcd.io/etcd/client/v3/concurrency`
  - `go.etcd.io/etcd/clientv3/namespace` → `go.etcd.io/etcd/client/v3/namespace`
  - `go.etcd.io/etcd/mvcc/mvccpb` → `go.etcd.io/etcd/api/v3/mvccpb`
  - `go.etcd.io/etcd/contrib/recipes` → `go.etcd.io/etcd/client/v3/experimental/recipes`

### 移除废弃 K8s API 导入

- `pkg/service/helm/install.go` 重写：移除 `apps/v1beta1`、`apps/v1beta2`、`apiextensions/v1beta1`、`k8s.io/kubernetes/pkg/apis/apps`，仅保留 v1 API

## Go 1.26.3 兼容性修复

- 修复 86+ 个 `go vet` 非恒定格式字符串错误：`logger.XXX(nil, variable)` → `logger.XXX(nil, "%s", variable)`
- 修复 `fmt.Errorf(variable)` 非恒定格式字符串（2 处）
- 修复字符串拼接格式调用 `logger.XXX(nil, "prefix: "+variable)` → `logger.XXX(nil, "prefix: %s", variable)`（4 处）
- 修复 int 到 string 转换 bug：`string(i)` → `strconv.Itoa(i)`（`pkg/service/cluster/handler.go`）
- 移除不可达代码：`pkg/plugins/vmbased/frame_interface.go`

## 依赖源迁移

| 原依赖 | 新源 | 新版本 |
|--------|------|--------|
| `kubesphere.io/im` | `codefuture.io/im`（GitHub: `codefuture-io/im`） | v0.2.0 |
| `openpitrix.io/iam` | `openpitrix.io/iam`（GitHub: `codefuture-io/iam`） | v0.1.0 |
| `openpitrix.io/notification` | `openpitrix.io/notification`（GitHub: `codefuture-io/notification`） | v0.2.2 |

- `im` 模块路径变更：`kubesphere.io/im/pkg/pb` → `codefuture.io/im/pkg/pb`（6 个文件）
- `iam`、`notification` 模块路径不变，通过 `replace` 重定向至 GitHub

## replace 指令清理

- 移除 `github.com/ugorji/go` replace（不再需要，`go/codec` 为独立模块）
- 移除 `github.com/docker/docker` → `docker/engine` replace（不再需要）
- `github.com/gocraft/dbr` 更新至最新 commit `v0.0.0-20190714181702-8114670a83bd`，移除 replace 指令

## 文件变更统计

33 个文件，970 行新增，1134 行删除
