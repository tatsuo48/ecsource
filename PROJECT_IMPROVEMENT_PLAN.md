# ecsource プロジェクト改善計画

## 概要
このドキュメントは、ecsourceプロジェクトの品質、保守性、拡張性を向上させるための包括的な改善計画を示します。

## 現状分析

### 強み
- ✅ シンプルで明確な目的（ECS Task Definitionのリソース表示）
- ✅ 自動リリース管理（tagpr + goreleaser）
- ✅ Homebrew対応
- ✅ 負の値を赤色で表示する視覚的フィードバック
- ✅ 2つの入力形式に対応（wrapped/unwrapped JSON）

### 改善が必要な領域
- ❌ **テストが存在しない** - 品質保証の欠如
- ⚠️ **単一ファイル構造** - 保守性とスケーラビリティの制限
- ⚠️ **エラーハンドリング** - strconv.ParseIntのエラーを無視
- ⚠️ **コード可読性** - main関数が長く複雑
- ⚠️ **CLI実装** - 手動の引数解析、flagパッケージ未使用
- ⚠️ **ドキュメント不足** - コードコメントが少ない

---

## 改善計画

### フェーズ1: コード品質と信頼性の向上（高優先度）

#### 1.1 テストの追加
**優先度: 🔴 高**

**目標:**
- ユニットテストカバレッジ80%以上
- CI/CDパイプラインでテスト自動実行

**実施内容:**
```
- main_test.go の作成
- JSON解析テスト
  - 正常系：wrapped/unwrapped形式
  - 異常系：不正なJSON、空ファイル
- リソース計算テストテスト
  - CPU、Memory、MemoryReservationの合計計算
  - leftover計算の正確性
  - 負の値の検出
- テーブル出力テスト
  - 赤色表示のテスト（負の値）
- GitHub Actionsワークフローの追加
  - go test ./... の自動実行
  - カバレッジレポート生成
```

**メリット:**
- バグの早期発見
- リファクタリング時の安全性向上
- 将来の機能追加時の回帰防止

#### 1.2 エラーハンドリングの改善
**優先度: 🔴 高**

**問題箇所:**
```go
// main.go:141-142
taskCPU, _ := strconv.ParseInt(taskDefinition.CPU, 10, 64)
taskMemory, _ := strconv.ParseInt(taskDefinition.Memory, 10, 64)
```

**実施内容:**
- strconv.ParseIntのエラーを適切に処理
- 無効な値の場合はエラーメッセージを表示
- エラーメッセージの統一と改善

**修正例:**
```go
taskCPU, err := strconv.ParseInt(taskDefinition.CPU, 10, 64)
if err != nil {
    log.Fatalf("failed to parse task CPU '%s': %v", taskDefinition.CPU, err)
}
```

---

### フェーズ2: コード構造とアーキテクチャの改善（中優先度）

#### 2.1 パッケージ構造のリファクタリング
**優先度: 🟡 中**

**現状:** すべてがmain.goに集約（180行）

**提案構造:**
```
ecsource/
├── main.go                 # エントリーポイントのみ（CLI引数処理）
├── internal/
│   ├── parser/
│   │   └── parser.go      # JSON解析ロジック
│   ├── calculator/
│   │   └── calculator.go  # リソース計算ロジック
│   ├── renderer/
│   │   └── table.go       # テーブル表示ロジック
│   └── models/
│       └── task.go        # データ構造定義
└── internal/parser/parser_test.go
    internal/calculator/calculator_test.go
    internal/renderer/table_test.go
```

**実施内容:**
1. データ構造を `internal/models` に移動
2. JSON解析を `internal/parser` に分離
3. リソース計算を `internal/calculator` に分離
4. テーブル表示を `internal/renderer` に分離
5. 各パッケージのユニットテスト作成

**メリット:**
- 責任の分離（Single Responsibility Principle）
- テストが書きやすくなる
- 将来の機能追加が容易
- コードの再利用性向上

#### 2.2 main関数のリファクタリング
**優先度: 🟡 中**

**実施内容:**
- main関数を短く（20行以内）
- 処理を意味のある関数に分割
  - `parseArgs()` - CLI引数処理
  - `loadTaskDefinition(filepath)` - ファイル読込と解析
  - `displayResourceTable(td)` - テーブル表示

**修正例:**
```go
func main() {
    args := parseArgs()

    taskDef, err := loadTaskDefinition(args.filepath)
    if err != nil {
        log.Fatal(err)
    }

    if err := displayResourceTable(taskDef); err != nil {
        log.Fatal(err)
    }
}
```

---

### フェーズ3: 機能拡張と利便性向上（中〜低優先度）

#### 3.1 CLIインターフェースの改善
**優先度: 🟡 中**

**実施内容:**
- `flag` パッケージまたは `cobra` ライブラリの導入
- より柔軟なオプション追加
  ```
  -f, --file         ファイルパス
  -o, --output       出力形式 (table, json, csv)
  -q, --quiet        エラー以外の出力を抑制
  -c, --color        色の有効/無効切り替え
  --no-header        ヘッダーを非表示
  ```

**使用例:**
```bash
ecsource -f task.json -o json
ecsource --file task.json --output csv --no-header
```

#### 3.2 出力フォーマットの追加
**優先度: 🟢 低**

**実施内容:**
- JSON出力対応
- CSV出力対応
- デフォルトはテーブル（現状維持）

**メリット:**
- スクリプトでの利用が容易
- 他ツールとの連携性向上

#### 3.3 複数ファイル処理
**優先度: 🟢 低**

**実施内容:**
- 複数ファイルを同時に処理
- ワイルドカード対応
- ディレクトリ内の全JSONファイル処理

**使用例:**
```bash
ecsource task1.json task2.json task3.json
ecsource tasks/*.json
ecsource --dir ./ecs-definitions
```

#### 3.4 追加情報の表示
**優先度: 🟢 低**

**実施内容:**
- Task Family名の表示
- Network Mode表示
- Essential コンテナの表示
- 使用率パーセンテージ表示
  ```
  CPU使用率: 0/512 (0%)
  Memory使用率: 1024/1024 (100%)
  ```

---

### フェーズ4: ドキュメントとメンテナンス（低優先度）

#### 4.1 コードドキュメントの追加
**優先度: 🟢 低**

**実施内容:**
- すべての公開関数にGoDocコメント追加
- パッケージレベルのドキュメント追加
- 複雑なロジックへのインラインコメント

#### 4.2 README.mdの拡充
**優先度: 🟢 低**

**実施内容:**
- より詳細な使用例
- トラブルシューティングセクション
- FAQセクション
- Contributing ガイドライン

#### 4.3 依存関係の更新
**優先度: 🟢 低**

**実施内容:**
- Go 1.23への更新（現在1.19）
- tablewriterの最新版へ更新
- Renovateの設定確認（既に存在）

---

## 実装ロードマップ

### スプリント1（1-2週間）: 基盤強化
- [ ] テストフレームワークのセットアップ
- [ ] 基本的なユニットテスト作成（カバレッジ50%）
- [ ] エラーハンドリングの改善
- [ ] GitHub Actionsでテスト自動実行

### スプリント2（1-2週間）: リファクタリング
- [ ] パッケージ構造の設計と実装
- [ ] main関数のリファクタリング
- [ ] テストカバレッジ80%達成
- [ ] 既存機能の動作確認

### スプリント3（1週間）: CLI改善
- [ ] flagパッケージの導入
- [ ] ヘルプメッセージの改善
- [ ] バージョン表示の改善

### スプリント4（2週間）: 機能拡張
- [ ] JSON/CSV出力対応
- [ ] 複数ファイル処理
- [ ] 追加情報の表示

### スプリント5（1週間）: ドキュメント
- [ ] GoDocコメント追加
- [ ] README更新
- [ ] Contributing ガイド作成

---

## 成功指標

### 品質指標
- ✅ ユニットテストカバレッジ: 80%以上
- ✅ すべてのエラーパスに適切なハンドリング
- ✅ Go lint/vet でワーニングゼロ

### 保守性指標
- ✅ 平均関数長: 30行以下
- ✅ 循環的複雑度: 10以下
- ✅ パッケージ構造が明確（責任分離）

### ユーザビリティ指標
- ✅ ヘルプメッセージの改善
- ✅ エラーメッセージの明確化
- ✅ 新しい出力フォーマット対応

---

## リスクと対策

### リスク1: 破壊的変更
**対策:**
- セマンティックバージョニングの厳守
- 既存のCLIインターフェースは維持
- 新機能はオプトイン方式

### リスク2: リファクタリング時のバグ混入
**対策:**
- テストファースト開発
- 既存機能のテスト完備
- 段階的なリファクタリング

### リスク3: 開発工数の見積もり
**対策:**
- スプリント単位で優先順位付け
- 高優先度から着手
- 各スプリント後にレビュー

---

## 次のステップ

1. **即座に着手可能:**
   - エラーハンドリングの修正（2-3時間）
   - 基本的なテストの追加（1日）

2. **短期（1-2週間）:**
   - テストカバレッジ向上
   - パッケージ構造のリファクタリング

3. **中長期（1-2ヶ月）:**
   - 機能拡張
   - ドキュメント整備

---

## 参考資料

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [tablewriter Documentation](https://github.com/olekukonko/tablewriter)

---

**最終更新:** 2025-10-21
**作成者:** Claude Code
