# internal/k8s Client 抽象化統合 — Design

**Goal:** `internal/k8s` 配下に並存する複数の "Client" 抽象化・非直感的な命名(`ObjectAPI`等)・重複実装を局所化/再利用化/Generics化し、パッケージ全体をシンプルかつコントロール可能にする。

**Scope:** `internal/k8s/{client,resource,reconciler,portforward}` および root package (`internal/k8s/*.go`) のみ。
CRD型定義パッケージ (`internal/k8s/vald/benchmark/api/v1`, `internal/k8s/vald/mirror/api/v1`, `internal/k8s/vald/operator/api/v1`, `internal/k8s/vald/operator/api/valdrelease`) は対象外。

## 現状分析

"Client" を名乗る抽象化が実質4種類、意味が重複したまま並存している:

| 抽象化 | 実体 | 生産コードでの利用 |
|---|---|---|
| `k8s.Client`(root) | controller-runtime `client.Client` のエイリアス | reconciler内で `mgr.GetClient()` 経由 |
| `internal/k8s/client.Client` | `cli.WithWatch` を手書きで包んだ独自interface、`client.New()`で生成 | 生産コード3箇所のみ (index-operator、benchmark-job config、read-replica rotator) — いずれも「manager外でCRUD+LabelSelectorが必要」な場面 |
| `resource.ObjectAPI` | `= kclient.Client` (実は `k8s.Client` と同一実体を別名で再宣言) | `ObjectClient[T,PT]`, `Syncer`, `GetObject`/`ListObjects` |
| `client.ClientSet` (生client-go clientset) | 型付きclientset + rest.Config | `resource.Pod/Deployment/...`等12種の型別ファクトリ(resources.go/status.go/rollout.go/token.go/dynamic.go) + `portforward`パッケージ。**tests/v2/e2e/crud/\* からのみ使用され、生産コード(pkg/)からの呼び出しはゼロ** |

追加調査で判明した具体的な無駄:
- `internal/k8s/client.Client` の `MatchingLabels` メソッドは定義されているが、実際の呼び出しは全て free var (`k8s.MatchingLabels`/`client.MatchingLabels`、controller-runtime の `cli.MatchingLabels` を指す)経由であり、interfaceのメソッドとして呼ばれている箇所はゼロ(死んだ surface)。
- `LabelSelector` は4箇所で実際に使われている、唯一価値のある独自メソッド。
- `resource/status.go` の `possibleStatuses()` と `checkResourceState()` は同じ12 Kindの集合を別々の type switch で列挙しており、Kind追加時に両方の更新が必要(同期ミスのリスク)。

`internal/k8s/reconciler` (`ListReconciler[L]`/`ObjectReconciler[T]`) と `resource.Base[T,PT]`/`ResourceInterface[T,L,C]`/`baseClient[T,L,C]` は既にGenerics化済みで質が高く、大きな手直しは不要と判断。

## 変更方針

### 1. Client抽象化の統合 (3種 → 実質1種+1拡張)

- `resource.ObjectAPI` を削除。`resource/object.go`, `objectclient.go`, `syncer.go` は root `internal/k8s` パッケージを import し、`k8s.Client` を直接使用する(import cycle無し: root package は `resource`/`client` のどちらにも依存していないため安全)。
- `internal/k8s/client.Client` interface を次の形に再定義する:
  ```go
  type Client interface {
      k8s.Client // controller-runtime の Get/List/Create/Delete/Update/Patch/Apply/Watch を埋め込みで継承
      LabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error)
  }
  ```
  `MatchingLabels` メソッドは削除(死んだsurface)。
- 具体型 `client` struct は `cli.WithWatch` を無名埋め込みにし、Get/List/Create/Delete/Update/Patch/Apply/Watch の8つの手書きforwardingメソッド(約90行)を削除する。promoted methodでinterfaceを自動的に満たす。
- `Patcher`/`NewPatcher`/`PodPredicates` はそれぞれ独立した用途(単一の実呼び出し元)を持つ既存の妥当な抽象化のため変更しない。

### 2. resource/status.go の重複統合

`possibleStatuses(obj)` + `checkResourceState(obj)` の2つのtype switchを、Kindごとに `{possibleStatuses []ResourceStatus, evaluate func(T) (ResourceStatus, string, error)}` を1組にまとめた単一テーブルに統合する。Kind追加時の更新箇所を1箇所に減らす。

### 3. e2e専用レイヤーの扱い

`client.ClientSet` + 型別ファクトリ一式 (`resources.go`/`status.go`/`rollout.go`/`token.go`/`dynamic.go`/`portforward`パッケージ) は配置を変えず `internal/k8s` に残す。上記1・2の重複削減以外の広範な書き直しは行わない(12 Kindの型別ファクトリは `baseClient[T,L,C]` で既にGenerics化されており、各ワークロード種別のPodTemplateアクセスパスがそれぞれ異なる構造体経路にあるため、これ以上のGenerics化は反射(reflection)無しには困難と判断)。tests/v2/e2e/crud/* が利用している公開API形状(resource.Pod/Deployment/.../WaitForStatus/RolloutRestart/CreateServiceAccountToken、client.ClientSet/Dynamic/RESTMapper/NewDynamicClient、portforward.New等)はそのまま尊重し、シグネチャ変更は行わない。

### 4. 呼び出し元・モックの追随

以下を新しい `client.Client` interface (embedding後) に追随させる:
- `pkg/index/operator/service/operator.go`
- `pkg/tools/benchmark/job/config/config.go`
- `pkg/index/job/readreplica/rotate/service/rotator.go`
- `internal/test/mock/k8s/client.go` (`ValdK8sClientMock`) — embeddingに伴い、mock構造も `k8s.Client` mock埋め込み + `LabelSelector` のみのmockに簡素化。`MatchingLabels` mockメソッドは削除。

## Success Criteria

- `go build ./...` / `go vet ./...` が通る。
- `go test ./internal/k8s/... ./pkg/index/operator/... ./pkg/tools/benchmark/job/... ./pkg/index/job/readreplica/...` が全てグリーン。
- `resource.ObjectAPI` という名前がコードベースから消えている。
- `internal/k8s/client/client.go` に Get/List/Create/Delete/Update/Patch/Apply/Watch の手書きforwardingメソッドが存在しない(embeddingで代替)。
- `resource/status.go` に `possibleStatuses` と `checkResourceState` の2重type switchが存在しない(単一テーブルに統合済み)。

## Out of Scope

- CRD型定義パッケージ (`vald/benchmark`, `vald/mirror`, `vald/operator/api`, `valdrelease`) の deepcopy/values.gen.go 等。
- e2eテスト専用レイヤーの `internal/k8s` からの移設。
- 12種の型別ファクトリ(`resource.Pod`/`Deployment`/...)自体の設計変更。

## 実装結果 (2026-07-11)

上記方針どおり実装済み。加えて `vald-reviewer`/`code-reviewer` によるレビューを実施し、以下を反映:

- `resource/status.go` の `evaluatorFor[R any]` を安全な型アサーション(`obj.(R)` の `ok` 判定付き)に変更し、テーブルのkeyとevaluatorの対応が万一ずれてもpanicせず `StatusUnknown`/`Unsupported resource type` にフォールバックするようにした。
- `internal/k8s/client/client.go` の `client` struct に、`New()` が `Client` interfaceを返すため `cli.WithWatch` から昇格する余分なメソッド(`DeleteAllOf`/`Status()`/`Scheme()`等)が外部に漏れない旨のコメントを追加。

## 追加スコープ: 本番/e2e Client統合 (2026-07-11 続き)

初回実装完了後、ユーザーから追加分析の依頼があり、`internal/k8s/resource` パッケージ内に「本番系(controller-runtimeベース)」と「e2e専用系(typed clientsetベース)」という性質の異なる2つのClient抽象化が混在している点を指摘した。

検証の結果、**typed clientset(`client.ClientSet`)はCRD(VolumeSnapshot, 自前CRD等)を一切扱えない**ため、本番側をe2e系(`resource.ResourceInterface`/`baseClient`)に寄せる統合は不可能と判明。一方、**scheme-aware な `k8s.Client` はCRDと組み込み型を同一APIで扱える**ため、本番側の3実呼び出し元(index-operator, benchmark-job config, read-replica rotator)がこれまで手書きで行っていたCRUDを、Genericsベースの `resource.ObjectClient[T,PT]` / 新設 `resource.ListClient[T,L,PT]` に統合する方向で実装した。

### 実装内容

1. **`resource.ListClient[T,L,PT]`を新設**(`internal/k8s/resource/objectclient.go`): `ObjectClient[T,PT]`(Get/UpdateStatus/Wait)を埋め込み、List/Create/Update/Delete/Watchを追加。List/Create/Update/Deleteは`ObjectClient`が保持する`k8s.Client`(=controller-runtimeの`client.Client`、`mgr.GetClient()`と同じ型)経由、WatchのみCRUDより広い`cli.WithWatch`型の`watchAPI`フィールド経由(`k8s.Client`にWatchメソッドが無いため)。
2. **`internal/k8s/client.Client`に`Raw() cli.WithWatch`アクセサを追加**(純粋加算、既存メソッドは無変更): 独自のname/namespaceベース`Get`を持つ`client.Client`は、controller-runtime標準のObjectKeyベース`Get`とシグネチャが非互換のため`k8s.Client`引数として渡せない。`Raw()`で埋め込みの標準クライアントを取り出し、`ObjectClient`/`ListClient`のコンストラクタに渡す。
3. **本番3呼び出し元を移行**:
   - `pkg/tools/benchmark/job/config/config.go`: CRD(`ValdBenchmarkJob`)の単発Getを`resource.ObjectClient`経由に。
   - `pkg/index/operator/service/operator.go`: Deployment/JobのList/Createを`resource.ListClient`経由に(`o.client.LabelSelector`は変更なし)。
   - `pkg/index/job/readreplica/rotate/service/rotator.go`: VolumeSnapshot(CRD)/PVC/DeploymentのList/Create/Update/Delete/Watchを`resource.ListClient`経由に。
4. **テスト書き換え**: `operator_test.go`/`rotator_test.go`のList/Create検証を、testifyモックの`.On("List"/"Create",...)`から実際のcontroller-runtime fakeクライアント(`fake.NewClientBuilder`)+ `interceptor.NewClient`によるCreate呼び出し回数計測に置き換え(実際のnamespace/labelフィルタリングが動作する分、以前より厳密な検証になった)。

### レビュー反映
- `code-reviewer`指摘により、Create呼び出し検証をbool(呼ばれたか否か)からcount(正確に1回)に強化。
- `ListClient`の`watch`フィールドを`watchAPI`にリネーム(`Watch`メソッドとの混同を避けるため)。
- `Create`のオプション型を`k8s.CreateOption`から`kclient.CreateOption`に統一(Update/Deleteと同じ`kclient.*`スタイルに揃えた)。
- テストの`event.Type != "ADDED"`を`event.Type != watch.Added`に変更(文字列リテラルではなく定数を使用)。

### 明示的に見送った項目(将来的な検討事項)
- ~~`internal/k8s.Client`(root、controller-runtimeエイリアス)と`internal/k8s/client.Client`(独自リッチラッパー)の名前衝突解消のためのリネーム~~ → 「追加スコープ3」で`StandaloneClient`にリネームして解消済み。
- `ListClient`の`func() *X { return new(X) }`という繰り返しのリストコンストラクタ引数を、`Objectable[T]`と同様のポインタ制約トリックで排除する案 — 7箇所程度の軽微な重複であり、今回は見送り。

## 追加スコープ2: resources.go の Client レイヤー化 (2026-07-11 続き)

ユーザーから、`internal/k8s/resource/resources.go`(e2e専用のtyped clientsetファクトリ群)を「Kubernetes Actionごとにクライアントを作るのではなく、リソースごとに一括構築する」Clientレイヤーに昇格してはどうか、という提案があった。

### 調査結果
- `tests/v2/e2e/crud/kubernetes_test.go`の`processKubernetes`は、Kubernetes Actionを実行するたびに`resource.Pod(r.k8s, k.Namespace)`等を都度呼び出しており、e2eシナリオ内でアクション数だけ再構築が発生する。
- ただし`c.GetClientSet().CoreV1().Pods(namespace)`等の構築自体はネットワークI/Oを伴わない軽量なstruct生成であり、**実行時性能への影響は軽微**と判断(ユーザーもこの点を確認し、目的は「コードの一元化/可読性」であるとの回答)。
- `config.KubernetesConfig.Namespace`はアクション単位のフィールドであり、スキーマ上はアクションごとに異なるnamespaceを指定可能(実際の`tests/v2/e2e/assets/*.yaml`は全て単一namespaceだが、将来的な変更で異なる可能性を排除できない)。そのため「起動時に単一namespaceで全Kind分を一括構築」という単純な設計は、将来のサイレントな挙動変化リスクがあるため採用しなかった。

### 実装内容
`internal/k8s/resource/clients.go`に`Clients`型を新設:
- `NewClients(c client.ClientSet) *Clients` — 軽量なコンストラクタ(内部マップは遅延初期化)。
- `Pod(namespace)`/`Deployment(namespace)`/.../`Endpoints(namespace)` — namespaceをキーにしたマップで遅延構築+キャッシュ(`cachedNamespaced[T comparable]`ジェネリックヘルパーで共通化)。アクションごとのnamespaceを完全に尊重しつつ、同一(Kind, namespace)の組み合わせの再構築を回避。
- `MutatingWebhookConfiguration()`/`ValidatingWebhookConfiguration()` — cluster-scopedなため単一インスタンスをキャッシュ。
- `resources.go`自体(12種のファクトリ関数)は無変更 — 既存の公開APIを壊さない。

`tests/v2/e2e/crud/strategy_test.go`の`runner`に`clients *resource.Clients`フィールドを追加し、`r.k8s`構築後に`resource.NewClients(r.k8s)`で一度だけ構築。`kubernetes_test.go`の`processKubernetes`は`resource.X(r.k8s, k.Namespace)`から`r.clients.X(k.Namespace)`に置き換え。

### テスト
`internal/k8s/resource/clients_test.go`(TDD): 全10種のnamespace付きKindをテーブル駆動で同一namespace呼び出しでの同一インスタンス返却/異なるnamespaceでの別インスタンスを検証、cluster-scoped側の単一インスタンスキャッシュ、`-race`付き並行アクセス安全性を検証。e2e側(`//go:build e2e`)は実クラスタが無いため`go build -tags e2e`/`go vet -tags e2e`によるコンパイル確認のみ。

### レビュー反映(2巡目)
`vald-reviewer`/`code-reviewer`の指摘により以下を修正:
- **Law 5違反**: `clients.go`/`clients_test.go`が標準の`sync`パッケージを直接importしていたのを`internal/sync`(エイリアス経由)に修正。
- **Law 3違反**: `rotator_test.go`に追加した`panic(err)`(テストのテーブル構築中、`t`がスコープ内にあるにも関わらずpanicを使用)を`t.Fatalf(...)`に修正。
- `cachedNamespaced[T comparable]`は実際には値の比較を一切行わない(map値としてのみ使用)ため、`T any`に緩和しコメントで明記。
- 単一インスタンスキャッシュ(`MutatingWebhookConfiguration`/`ValidatingWebhookConfiguration`)がキャッシュヒット時も無条件で書き込みロックを取得していた不整合を、`cachedSingleton[T comparable]`ジェネリックヘルパー(RLock高速パス付き)に統一し、namespace付きキャッシュと同じダブルチェックロッキングのスタイルに揃えた。
- テストカバレッジをPod/Deploymentの2種類のみから全10種のnamespace付きKindへテーブル駆動で拡張。

## 追加スコープ3: 残存する非直感的命名の解消 (2026-07-11 続き)

コミット後、ユーザーから再度同じゴールでの深掘り依頼があり、再調査したところ以下2点が残存していた:

1. **`APIPodMetrics`/`APINodeMetrics`系エイリアス**(`internal/k8s/types.go`) — ユーザーが当初から挙げていた「APIXXXのような非直感的な操作」に該当する最後の識別子。唯一の利用者である`pkg/discoverer/k8s/service`が独自ドメイン型`PodMetrics`/`NodeMetrics`を持つため名前衝突回避のために`API`接頭辞を付けていたが、利用箇所が1パッケージのみだったため、エイリアスを削除し`pkg/discoverer/k8s/service`が`k8s.io/metrics/pkg/apis/metrics/v1beta1`を直接importする形に変更(実装の局所化)。
2. **`internal/k8s.Client`(root、controller-runtimeエイリアス) と `internal/k8s/client.Client`(独自リッチラッパー)の名前衝突** — 両レビューで指摘済みだったが前回まで見送っていた項目。`internal/k8s/client.Client`を`StandaloneClient`にリネーム(「manager外で動作するコードのためのscheme-awareクライアント」であることを名前で明示)。実際に`internal/k8s/client`パッケージをimportしている全ファイルを`grep`で洗い出し、`client.Client`型を実使用する7箇所(rotator.go x2, operator.go, options.go, config.go, job.go, option.go)とmockの型アサーションを更新。他パッケージの同名だが無関係な`client.Client`(gRPC vald client、controller-runtime自身のclient.Client等)には一切手を触れていない。

### レビュー反映
- `vald-reviewer`: 両パートともVald Law違反なし。リネームが全呼び出し元に漏れなく適用されていることをリポジトリ全体のgrepで独立検証済み。
- `code-reviewer`: ブロッキング指摘なし。`StandaloneClient`という命名は妥当と判断。コメントアウト済みテストコード内の`client.Client`表記(`operator_test.go`)も`client.StandaloneClient`に更新し、設計ドキュメントの「見送り」記載も本セクションで解消済みに更新。

## 追加スコープ4: ListClient のクロージャ引数排除 & ClientSet フォールバック重複解消 (2026-07-11 続き)

さらなる深掘り依頼を受けて再調査したところ、以下2点が見つかった:

1. **`ListClient`の`newList func() L`クロージャ引数** — `resource.NewListClient[T, L](api, func() L { return new(L) })`という呼び出し形が、全実呼び出し元(Deployment/Job/VolumeSnapshot/PersistentVolumeClaim/ConfigMap)で`new(X)`を返すだけの同一パターンだった。`ObjectClient`が既に`Objectable[T] = interface{ *T; Object }`という「PがTへのポインタでObjectを実装する」制約でこれを型引数のみから解決していたのと同じ発想で、`ListClient`にも`ListPtr[L] = interface{ *L; ObjectListType }`制約を導入し、`PL(new(L))`で同じことをクロージャなしで実現。型引数はコンストラクタ呼び出し時点では`NewListClient[T, L](api)`の2つの明示引数から残り2つ(`PT`, `PL`)が構造的単一ポインタ項の制約により推論されるため、呼び出し側は変更前と同じ2型引数で済む。ただし構造体フィールド宣言(`operator.go`/`rotator.go`の`*ListClient[...]`型フィールド)ではGoが型引数を推論しないため、4引数のフル表記が必要になる制約は事前に確認済み。
2. **`client.NewClientSet`の3箇所に重複したin-clusterフォールブックロック** — `kubeConfig`未指定時・`ClientConfig()`失敗時・`clientSetFromConfig()`失敗時の3箇所で「in-cluster設定にフォールバックし、失敗したら元エラーとJoinする」処理が逐語的に重複していた。`fallbackToInCluster(origErr error) (ClientSet, error)`ヘルパーに抽出し、3箇所とも呼び出しに置き換え。`NewClientSet`の戻り値も named return(`c ClientSet, err error`)から素の`(ClientSet, error)`に変更(named returnを使う理由がなくなったため)。

### 実装内容
- `internal/k8s/resource/objectclient.go`: `ListPtr[L any] interface { *L; ObjectListType }`を追加。`ListClient[T, L, PT, PL]`の`List`/`Watch`は`PL(new(L))`で新規リストを生成。`operator.go`/`rotator.go`の呼び出し元(19箇所)から`func() L { return new(L) }`クロージャを削除。
- `internal/k8s/client/clientset.go`: `fallbackToInCluster`ヘルパーを抽出し、`NewClientSet`内の3箇所の重複ブロックを置き換え。

### レビュー反映
- `vald-reviewer`: Vald Law違反なし。`PL(new(L))`が実際に使われている全5 Kind(Deployment/Job/VolumeSnapshot/PersistentVolumeClaim/ConfigMap)で旧クロージャと挙動同一であることを検証。`fallbackToInCluster`のエラーJoin挙動が元の3ブロックすべての分岐で同一であることも検証。
- `code-reviewer`: ブロッキング指摘なし。以下の提案を追加で反映:
  - `internal/k8s/client`にテストファイルが皆無だった(pre-existing gap)ため、`clientset_test.go`を新設し`fallbackToInCluster`のnilエラー時/既存エラーJoin時の両ケースをテスト(`rest.InClusterConfig()`がテスト環境では常に`rest.ErrNotInCluster`を返す決定性を利用)。
  - `operator.go`と`rotator.go`の両方に重複していた`*resource.ListClient[k8s.Deployment, k8s.DeploymentList, *k8s.Deployment, *k8s.DeploymentList]`という4引数フル表記を、`resource.DeploymentListClient`型エイリアスとして`objectclient.go`に集約し、両ファイルのフィールド宣言を置き換え。
  - (残存・対応不要と判断) `ListPtr`/`Objectable`は構造的制約であるため、「LがTの正しいリスト型であること」自体は型システムでは強制できない(旧設計から変わらない既知の限界であり、今回のリファクタでは悪化していない)。

## 追加スコープ5: reconciler パッケージの残存クロージャ排除 (2026-07-11 続き)

ユーザーから「更に高度なGenerics化と実装の局所化、構造パスの最適化、コールスタックの最適化を実施してください」という依頼を受けて再調査したところ、以下2点が見つかった:

1. **`reconciler.NewListReconciler`の`newList func() L`クロージャ**(Round4で対応済みの`resource.ListClient`と同型のパターンが`reconciler`パッケージ側に残存) — `resource.ListPtr[L]`制約を再利用し、型パラメータを`[L any, PL resource.ListPtr[L]]`に分割。`ListOption`/`listReconciler`/全`WithXxx`ヘルパー/`Reconcile`を追従。呼び出し元(`pkg/discoverer/k8s/service/discover.go`4箇所、`pkg/operator/benchmark/service/operator.go`1箇所)からクロージャを削除。多くの呼び出しは`WithOnReconcile`等の引数からGoの型推論で`L`/`PL`が解決されるため、明示型引数は不要(コンパイラの`unnecessary type arguments`ヒントで確認)。
2. **`reconciler.NewObjectReconciler`の`newObj func() T`クロージャ** — 同一パターン。ただし今回は既存の`resource.Objectable[T]`制約(`internal/k8s/resource/objectclient.go`で定義済み、`resource.ObjectClient`が既に使用中)をそのまま再利用できたため、新規制約の追加は不要だった。型パラメータを`[T any, PT resource.Objectable[T]]`に分割し、`ObjectOption`/`objectReconciler`/全`WithObjectXxx`ヘルパー/`Reconcile`/`NewReconciler`/`For`を追従。唯一の実呼び出し元`pkg/index/operator/service/operator.go`からクロージャを削除。

このほか、`internal/k8s/reconciler.go`/`option.go`の非公開フィールド`merticsAddr`(typo)を`metricsAddr`に修正(Round4のレビューで指摘されていた軽微な事項)。

コールスタック/構造パスの観点では、`baseReconciler`(共通ロジック集約済み)・`resource.Base[T,PT]`(DeepCopy基盤)・`resource.List[T,PT]`(全リスト型の統合)・`client.StandaloneClient`(embedding済み)・`client/dynamic.go`は、いずれも既存ラウンドで既に妥当な設計に到達しており追加の圧縮余地はないと判断した。reconcile loop自体はKubernetes APIへのネットワークI/Oが律速するため、インターフェース間接呼び出しの削減による実質的な性能改善は見込めず、可読性とのトレードオフに見合わないと判断し対象外とした。

### レビュー反映
- `vald-reviewer`: PASS。`newList`/`newObj`クロージャ排除後も`PL(new(L))`/`PT(new(T))`が全実使用箇所(metricsv1beta1.NodeMetricsList/PodMetricsList、k8s.PodList/NodeList/JobList、k8s.Pod)で旧クロージャと挙動同一であることを検証。Vald Law違反なし。
- `code-reviewer`: 承認、ブロッキング指摘なし。2型パラメータ化は既存の`resource.ListClient`/`resource.ObjectClient`と同じ設計一貫性を持つと評価。

(このコミット後、`vald-reviewer`/`code-reviewer`のteammate版がSendMessageツールを持たず応答を返せないまま待機し続ける問題が判明。以降のレビューは`fork`(オーケストレーターの会話コンテキストを継承し、完了時に自動通知が届く)で実施する方式に切り替えた。)

## 追加スコープ6: pkg/gateway/mirror での ObjectClient パターン未適用箇所の解消 (2026-07-11 続き)

再度「更に高度なGenerics化と実装の局所化、構造パスの最適化、コールスタックの最適化」の依頼を受けて、これまで未調査だった`internal/k8s/portforward`・ルートパッケージ残り・`pkg/gateway/mirror`/`pkg/agent`/`pkg/manager/index`を調査した結果、以下1点が見つかった:

**`pkg/gateway/mirror/service/discovery.go`の`updateMirrorTargetPhase`が手書きの`ObjectKey`構築+`Status().Update()`を使用** — `resource.ObjectClient[T,PT]`(`Get(ctx, name, namespace)` + `UpdateStatus(ctx, obj)`)として既に一般化済みのパターンと完全に一致していたが未適用だった。`target.MirrorTarget`(`= mirrv1.ValdMirrorTarget`、`internal/k8s/vald/mirror/api/v1`で定義済みのCRD型で`DeepCopyObject`実装済み)は`resource.Objectable[T]`制約を満たす。

### 実装内容
`discovery`構造体に`mirrorTargets *resource.ObjectClient[target.MirrorTarget, *target.MirrorTarget]`フィールドを追加し、`NewDiscovery`内で`d.ctrl`確定直後に`resource.NewObjectClient[target.MirrorTarget](d.ctrl.GetManager().GetClient())`で一度だけ構築(`mgr.GetClient()`は`manager.New()`完了時点で即座に取得可能な値であり、`Start()`前でも安全に構築できることを確認)。`updateMirrorTargetPhase`内の手書き`c.Get(ctx, k8s.ObjectKey{...}, mt)` + `c.Status().Update(ctx, mt)`を`d.mirrorTargets.Get(ctx, name, d.namespace)` + `d.mirrorTargets.UpdateStatus(ctx, mt)`に置き換え。

`createMirrorTargetResource`の単発`Create`呼び出しは、`resource.ObjectClient`にCreateメソッドがなく(List系操作を持つ`ListClient`が必要になり不釣り合い)、スコープ外の重複でもないため変更しなかった(スコープ外の改善を追加しない方針)。

### テスト
新規テストは追加していない: 既存の`Test_discovery_syncWithAddr`(`pkg/gateway/mirror/service/discovery_test.go`)が`NewDiscovery`経由で`d.ctrl`に`k8smock.ControllerMock{GetManagerFunc: k8smock.NewDefaultManagerMock}`を注入した上で`syncWithAddr`→`updateMirrorTargetPhase`のパスを実際に実行しており、リファクタ後も全ケースPASSすることで新旧プラミングの挙動同一性を検証済み。

### レビュー反映
- `vald-reviewer`: PASS。新旧`Get`/`UpdateStatus`の意味的同一性を`internal/k8s/resource/object.go`/`objectclient.go`の実装から検証。`d.ctrl.GetManager().GetClient()`をStart前に呼ぶ安全性を`internal/k8s/reconciler.go`の`New()`実装(managerが同期的に構築済み)から検証。`k8s.New(...)`失敗時の`return nil, err`追加(nilのまま`GetManager()`を呼ぶパニック防止)も正しい修正と評価。Vald Law違反なし。
- `code-reviewer`: 承認、ブロッキング指摘なし。`pkg/tools/benchmark/job/config/config.go`の既存`resource.NewObjectClient`呼び出しと一貫。非ブロッキングの所感: フィールドの構築が`updateMirrorTargetPhase`初回呼び出し時の遅延構築から`NewDiscovery`内の即時構築に変わったが、fail-fast方向の改善であり実質的にリスクなし。

## 追加スコープ7: multi-agentワークフローによる網羅的再調査と11項目の実装 (2026-07-12)

ユーザーから再度「GenericsでさらにSimplify可能か」という依頼を受け、ultracodeモードのmulti-agentワークフロー(8クラスタ並列サーベイ→重複排除→新規性/実現可能性/価値の3レンズ敵対的検証、各候補はoverlayコンパイル実験で事前検証)を実施。30候補を抽出し、3レンズ全会一致10件+部分支持3件を得た(残り17件はセッション上限で未検証、後続ラウンドで検証再開予定)。ユーザー承認のもと以下11項目を実装:

1. **`Base.DeepCopyInto`削除**(`base.go`) — 契約上常にshadowされる死んだsurfaceであり、埋め込み型が手書き`DeepCopyInto`を書き忘れた場合に昇格メソッド自身が`DeepCopyIntoer[T]`制約を自己充足してコンパイルが通り、実行時に`self().DeepCopyInto`→昇格メソッド→…の無限再帰でstack overflowするfootgunだった。削除により書き忘れは「missing method DeepCopyInto」のコンパイルエラーになる。理由を説明するコメントを残置。`base_test.go`の`TestBase_DeepCopyInto`と`b.DeepCopyInto(nil)`行を削除。
2. **`List.DeepCopyInto`の死んだ`TypeMeta`再代入削除**(`list.go`) — 直前の`*out = *in`でコピー済み。
3. **`Create/Update/Delete`を`ListClient`→`ObjectClient`へ移動**(`objectclient.go`) — 単一オブジェクト操作は`ObjectClient`の責務。`ListClient`は`*ObjectClient`埋め込みのmethod promotionで既存呼び出し元(rotator.go/operator.go)は無変更。
4. **`Watch`のlist引数排除**(`objectclient.go`) — controller-runtimeはlistをGVK解決にのみ使いItemsを無視するため、`PL(new(L))`内部生成で`List`と対称のシグネチャに。`rotator.go`の「このsnapshotだけをwatchしている」ように読める誤解誘発リスト構築(2箇所約8行)を削除。
5. **`newWorkloadClient`のsetter導出**(`resources.go`) — 5 Kind×2クロージャ(getter/setter)を、getterが返すフィールドポインタへの`*b.podTemplate(obj) = *pt`書き込みでsetterを導出する形に統一し5クロージャに半減。Goは型パラメータへのフィールドアクセスを許さない(golang/go#48522)ためgetter 1個/Kindが証明可能な下限。
6. **`withExt[I,R,...]`ヘルパー**(`resources.go`) — 7つの拡張インターフェースディスパッチ(UpdateStatus/DeleteCollection/UpdateEphemeralContainers/UpdateResize/GetScale/UpdateScale/ApplyScale)の「型アサーション+okチェック+ゼロ値+ErrUnimplemented(名前文字列)」定型を1ヘルパーに集約。型引数は全呼び出しで推論。
7. **`cachedNamespaced`/`cachedSingleton`のファクトリ直接渡し化**(`clients.go`) — 12個の束縛クロージャ`func() X { return Pod(cs.c, ns) }`を、ヘルパーが`*Clients`を受け取り`construct(cs.c, namespace)`を呼ぶ形にしてファクトリ関数値(`Pod`等)の直接渡しに。Round 4/5で排除したのと同じクロージャ引数パターンの残存だった。
8. **`extractItems`削除→`apimeta.ExtractList`**(`status.go`) — 27行の手書きreflection(FieldByName+Addr/valueフォールバック)を上流の維持されているヘルパーに委譲。`checkResourceState`は動的型でテーブルディスパッチするだけなので非generic(`any`引数)化。`status_test.go`の型不一致テストは「Itemsを持たない戻り型で抽出自体が失敗する」ケースに書き換え(ガード目的=抽出失敗の即時浮上は維持)。`internal/errors`の`ErrItemsFieldNotFoundOrNotASlice`/`ErrItemIsNotOfType`はdead化するがスコープ外のため残置。
9. **`waitLoop`統合**(`objectclient.go`/`status.go`) — `ObjectClient.Wait`と`WaitForStatus`が各々手書きしていたticker+timer+selectのポーリング骨格を非genericヘルパー`waitLoop(ctx, onTimeout, step)`に統合。`WaitForStatus`はname/labelSelector分岐をループ外に持ち上げtick毎の再分岐も解消。※`ObjectClient.Wait`自体は生産呼び出しゼロのdead surfaceであることが検証で判明したが、ユーザーが将来のoperator実装でのstatus待機用に温存を選択したため、削除ではなく統合を実施。
10. **`registerKind[R]`**(`status.go`) — `kindStatusTable`のmap literal(11 Kind分のkey/evaluator手動ペア)を`init()`内の`registerKind(possible, evaluator)`呼び出しに置換。テーブルキーとevaluator引数型が単一の型パラメータRから導出されるため、copy-pasteによるkey/evaluatorずれがコンパイル時に閉じる。
11. **`RolloutRestart`の型パラメータ縮小4→1**(`rollout.go`) — 本体は`SetPodAnnotations`1メソッドしか呼ばないため、制約を`WorkloadControllerResourceClient[T,L,C]`から既存の`podAnnotationInterface[T]`に縮小。Round 2の「e2e公開API凍結」に文言上抵触するが、呼び出し元`tests/v2/e2e/crud/kubernetes_test.go`はメソッドベース推論により無変更で済むことをoverlay vet(e2eタグ)で検証済みのため、凍結の趣旨(e2eを壊さない)は維持されるとユーザーが判断し実施を承認。

検証: `go build`(internal/k8s+pkg全体)/`go build -tags e2e`/`go vet`/`go vet -tags e2e`/`go test -race -count=1`全て緑、`gofmt`差分なし。12ファイル、+241/-382行(純減141行)。

### 見送り・棄却
- `ObjectClient.Wait`削除(2/3支持): ユーザー判断で温存(上記9参照)。検証で判明した注意点として、`defaultWaitInterval`/`defaultWaitTimeout`は`WaitForStatus`と共有のため、将来削除する場合も2定数は温存または`status.go`へ移動が必要。
- `WaitForStatus`の型パラメータC排除(1/2): Round 2のe2e公開API凍結に抵触(noveltyレンズがNG)のため見送り。
- 「4型パラメータ`ListClient`のgeneric type aliasによる削減」: サーベイ自身が「Go 1.26でも不可能」という否定的結論を報告(型パラメータの制約で*TがObjectable[T]を満たすことを表現できない)。

### レビュー反映
- `vald-reviewer`: PASS。重点9項目を全て独立検証 — CRD 4パッケージのテスト通過、`SetPodTemplate`新旧1対1対応(5 Kind全て)、`WaitForStatus`セマンティクス保持、`registerKind`の11 Kind分のpossible集合を旧map literalと全件突合して完全一致、controller-runtime `typedWatch`のソースを直接確認しWatchのlist引数がGVK解決専用であることを実証、`go vet -tags e2e`でRolloutRestart縮小後の呼び出し元無変更を確認。Vald Law 1-5違反なし。
- `code-reviewer`: 承認、ブロッキング指摘なし。withExtの型引数推論が全7呼び出しで有効、waitLoopの挙動同一性(ticker/timer生成順・エラー伝播順)、registerKindのinit()初期化順の安全性、SetPodTemplateのロック下in-place変異の安全性を確認。非ブロッキング所見3件(WaitForStatus List路の短絡は実到達不能な差・dead errors残置は記載済み・base.goのコメント位置は妥当)はいずれも対応不要。
