# hackathon-22-winter-01

## 必要環境

- Go 1.19

## 開発者メモ

- [openapi](./docs/openapi/v1.yml)を変更したら`go generate ./...`を実行

### 依存図

←: 依存 ⇠: 実装

```mermaid
classDiagram

class domain { <<struct>> }
class ws { <<interface, struct>> }
class repository { <<interface>> }
class repoimpl { <<struct>> }
class handler { <<struct>> }
class oapi { <<interface>> }

domain <|-- ws
domain <|-- repository
repository <|.. repoimpl
ws <|-- handler
repository <|-- handler
handler ..|> oapi
oapi <|-- main
```
