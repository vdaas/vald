---
marp: true
title: OSS分散ベクトル検索エンジンValdと最新の取り組み
theme: vald
paginate: true
footer: © vdaas/vald
---

# OSS分散ベクトル検索エンジンValdと最新の取り組み

<div class="center">
Vald: Cloud Native Distributed Vector Database
</div>

<div class="center">
LINEヤフー株式会社 Matts966
</div>
<br>

<div class="center">

![alt text width:300](speakerdeck_vald.png)
<p>資料URL</p>
</div>

---

## 自己紹介

松井誠泰（GitHub: [Matts966](https://github.com/Matts966)）
- LINEヤフー株式会社
  - OSSのベクトルデータベースValdチームに転職して半年
- 趣味
  - 🍺 🍶 🥃 ☕️ 💻 📖 🚲

---

## お品書き

- Valdのご紹介
- 最新の取り組み
  - Grafana Foundation SDKを用いたGrafana Dashboard as Code
  - E2E V2: YAMLで宣言的にテストシナリオを記述
- ベクトルDB選定のすすめ

---

## Valdのご紹介

- クラウドネイティブな分散近似近傍ベクトルDB
- オープンソース
- CNCF Landscape
<!-- - 公式サイト: [vald.vdaas.org](https://vald.vdaas.org) -->

<div class="center">
  <img src="qr.png" />
  <p><a href="https://vald.vdaas.org">vald.vdaas.org</a></p>
</div>

[![bg contain right:60%](image-7.png)](https://landscape.cncf.io/?item=app-definition-and-development--database--vald)

---

## 最新の取り組み

1. **Grafana Foundation SDKを用いたGrafana Dashboard as Code**
2. E2E V2: YAMLで宣言的にテストシナリオを記述

---

## Grafanaボード管理の課題

- コンポーネント毎にボード・パネルをたくさん管理
  - コンポーネントがかなり違うものの共通化の余地あり
    - 繰り返し、条件分岐したい
  - パネル毎にアップグレード作業
- JSONでバージョン管理はしていたものの
  - シンプルなパネルでもGUIからexportされたJSONは大きくなってしまい、直接読み書きするのが難しい

![bg contain right:20%](image-3.png)

<!-- 
---

## Grafana Dashboard as Codeの選択肢 - JSONベース

| 方法                           | 概要                                     | 特徴・注意点                                |
| ------------------------------ | ---------------------------------------- | ------------------------------------------- |
| JSON管理（元の手法）           | GUIで作成後にJSON出力                    | 単純・最小構成向け、再利用や共通化は弱い    |
| Terraform Provider for Grafana | IaC統合（HCL）                           | JSON構造の記述が必要、Terraformに統合できる |
| Git Sync                       | GUI変更を自動でGit同期（Grafana 12以降） | GUI派に便利、繰り返しや再利用には不向き     |

---

## Grafana Dashboard as Codeの選択肢 - コードベース

| 方法                   | 概要                        | 特徴・注意点                                                   |
| ---------------------- | --------------------------- | -------------------------------------------------------------- |
| Grizzly                | CLIでリソースとして管理可能 | CLIが便利・Jsonnet使える                                       |
| Grafonnet              | Jsonnetで生成               | 繰り返し処理など対応                                           |
| Grabana                | Goで記述、宣言的            | 唯一JSON逆生成可能、開発は `grafana-foundation-sdk` に移行傾向 |
| grafana-foundation-sdk | 公式SDK（Go等）             | ⭐️**本日のお題**⭐️                                               | -->

<!--
---

## Grafana Dashboard as Code の選択肢

- JSON 出力してマニフェスト手動管理 or 自前で自動化
- Terraform Provider for Grafana
- Git Sync
- Grizzly
- Grafonnet
- Grabana
- grafana-foundation-sdk ← 今日の本題

---

## JSON出力とマニフェスト管理

- GUIでダッシュボードを作成 → JSONエクスポートしてバージョン管理
- CIで自動反映も可能
- シンプルだが繰り返しや共通化に弱い

---

## Terraform Provider for Grafana

- Terraformでダッシュボード・データソースなどを管理
- 他のIaCと統一できる
- JSON構造を記述する形のため、編集性はやや低い

---

## Git Sync

- Grafana 12で登場した機能、この５月に発表された
- GUI操作の結果をそのままgitに同期
- 職人的に凝ったグラフをたくさん作り、繰り返しが少ない運用では一番いいかも

---

## Grizzly

- `grr` CLIで `diff`, `apply` 操作でマニフェストを使った管理が自動化できる
- 複数のGrafanaオブジェクト（アラート等）も管理可能
- Jsonnetも使える

---

## Grafonnet

- Jsonnetライブラリでの構成
- 記述量が少なく複雑な構成に対応可能
- Jsonnetで書きたいならこれ

---

## Grabana

- Goでダッシュボード構築（宣言的に書ける）
- 唯一JSONからコードを逆生成できる
- ただし作者は `grafana-foundation-sdk` に注力しており[新機能対応がされていない](https://github.com/K-Phoen/grabana/issues/264)

🔗 [Three years of Grafana dashboards as code](https://blog.kevingomez.fr/2023/03/07/three-years-of-grafana-dashboards-as-code/) -->

---

## grafana-foundation-sdk の概要

- Grafana公式が提供する言語ごとのSDK
- GrafanaのAPIスキーマをベースに自動生成されている
- Go, TypeScript, Python, Java に対応

---

## 選定理由・メリット

- 繰り返しを簡単に表現できる
  - 同じようなダッシュボードをコンポーネントごとにつくっている場合などに、関数等で整理しやすい
- メトリクスを管理しているコードと同じ言語で書くことで、メトリクス名を参照でき、二重管理を避けられる
  - メトリクスの宣言→ダッシュボード作成まで自動化可能

---

## メリット

- メソッドチェーンで書けるので、補完に沿って書ける
- テキストなのでLLMの力を借りやすい
- GUIから出力できるJSONからGoへの自動変換が可能

![alt text](image-1.png)

---

## メリット

- 簡単にバージョンアップグレード
  - 公式がAPIスキーマから自動生成しているので
    - `go get` でタグを切り替えるだけで簡単に最新に追従できる
    - 網羅性が高い

```sh
go get github.com/grafana/grafana-foundation-sdk/go@v11.6.x+cog-v0.0.x
```

---

## メリット

- 公式から promql もビルダーが提供されていて、複雑な文字列、括弧の対応の管理を避けられる

<div class="center">

![width:800px](image-2.png)

</div>

---

## 注意点

![bg fit right:60%](image.png)

- [grafana/grafana-foundation-sdk#673](https://github.com/grafana/grafana-foundation-sdk/issues/673)
  - パネル配置にバグがあるため
  - 行や列の位置がズレるなど
  - 自分で整理するコードを書く必要あり
- 現状 [puzzle.go](https://github.com/vdaas/vald/blob/main/hack/grafana/gen/src/puzzle.go) としてValdレポジトリで公開

---

## 結果

<div class="center">

![alt text width:200](github_vald.png)

- コード量を1万行近く削減
- ほぼ同じボードを再現

![alt text](image-4.png)


<!-- 
| Board | Before | After |
|--------|--------|--------|
| Agent Memory | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/a902cf2a-69ba-44e0-aca5-cbe7e0566622" /><img width="1919" alt="image" src="https://github.com/user-attachments/assets/6215c298-bff5-416f-b1be-c7b90891c13c" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/cb076d4e-8aa2-4a80-a25d-30bded3d8c39" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/aea85f5e-bdbf-4689-9b93-41d0eeefc0c7" /><img width="1914" alt="image" src="https://github.com/user-attachments/assets/6a9f374b-815b-405c-8cb2-d3b6ee592a66" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/16056e2b-81d7-4d69-8605-9879fb514ff8" /> |
| Cluster Overview | <img width="1916" alt="image" src="https://github.com/user-attachments/assets/3e385e6d-5fbe-4186-9c4d-5ac688d6d357" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/b7e057a3-d026-4fe6-9d66-803fcfd13344" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/d7abad81-cb6c-401d-b1d4-f9e9d3f27164" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/0cb13301-868a-49e6-b169-267e98dcb9a0" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/f0bfa68b-9d6f-4dd9-bf88-e8c663e406fa" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/cd39f505-32e7-4662-b7de-9a34293b2c67" /> |
| Discoverer | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/1a1465d3-e2d4-4988-b6a2-266c9f54338a" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/44715b59-a728-4a72-b699-ed05db1fc596" /> |
| LB Gateway | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/3b311790-1ca4-4cb4-8952-2163b534f597" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/e44f51cf-5e3c-404b-91f9-f05f0267b259" /> |
| Agent | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/ace72cdb-5d3e-42d8-915a-25b5639a9f86" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/8a9c0b00-e3ee-4cc3-a915-e4ede2be0cb6" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/474bdc86-a650-4a35-b2d6-ad601662632f" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/f17ca7bc-c51b-41df-8cb4-9bcfb17f13fb" /> | 
| Benchmark Operator | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/5bf084f7-5dc2-4f71-ac23-b671adb27ba9" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/e444e6f9-9375-4723-9291-c554ca8a9bb1" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/835ddbe4-b50f-4e93-aebb-05eb62814d05" /><img width="1920" alt="image" src="https://github.com/user-attachments/assets/e84d56c2-1667-4c43-9151-e15d022c3150" /> |
| Helm Operator | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/ff04df50-7c83-40a1-bf11-757cf2101e03" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/29652b77-88d7-4f76-97e8-ed3f90bf156b" /> |
| Index Correction | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/12d22561-e9e9-4b40-b2a5-2c2427554f0c" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/7c66ddc9-0787-4da9-b3f8-3067875d70f2" /> |
| Index Manager | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/7400e4e1-8ec4-46ac-ace7-e2f33714d35f" /> | <img width="1920" alt="image" src="https://github.com/user-attachments/assets/36d5da3e-0a2c-4a26-b49d-f1b5341e28b8" /> |
-->
<!-- </div> -->

![bg fit right:60%](image-5.png)

---

## 最新の取り組み

1. Grafana Foundation SDKを用いたGrafana Dashboard as Code
2. **E2E V2: YAMLで宣言的にテストシナリオを記述**

---

## E2Eの悩み

- CRUDの処理がたくさんあるが、コードで書いていると冗長になりがち
- 違うコードベースに同じようなコードが散らばる
- データを取り出してアサートする流れも煩雑になりがち

---

## E2E V2: YAMLで宣言的にテストシナリオを記述

- 得られた成果
  - ジェネリクスを用いた汎用k8s, gRPCクライアント
  - 別環境でもYAMLをもとにk8s JobでE2Eが走る
  - パスで結果を取り出し、アサートできる
  - 並列実行・Loop処理を用いた負荷試験
- Future Work
  - PBT: Property Based Testing

![bg fit right:30%](image-9.png)

---

## ベクトルDB選定のすすめ

- CNCFにはハイブリッドサーチをサポートするOpenSearchもあり、検索用途で :+1:
- 推薦・検出などベクトル検索だけ必要で、パフォーマンス重視の方にはValdは :+1:

<div class="table">

|検索手法|エンジン|90 %ile (ms)|99 %ile (ms)|MRR|
|--|--|--|--|--|
|全文検索|OpenSearch|10.42|23.79|0.605|
|ハイブリッドサーチ|OpenSearch|21.56|28.823|**0.661**|
|ベクトル検索|OpenSearch|9.60|11.87|0.619|
|ベクトル検索|Vald|**1.93**|**2.363**|0.615|

</div>

[検索エンジン選定ガイド：ベクトル検索・全文検索からハイブリッドサーチまで](https://techblog.lycorp.co.jp/ja/20241225b)
[LINEヤフー Tech Blog](https://techblog.lycorp.co.jp/ja/20241225b)

---

# Contributions are Welcome!

<!-- [vald.vdaas.org](https://vald.vdaas.org)

![bg auto right](qr.png) -->

<div class="center">
  <img src="qr.png" />
  <p><a href="https://vald.vdaas.org">vald.vdaas.org</a></p>
</div>
