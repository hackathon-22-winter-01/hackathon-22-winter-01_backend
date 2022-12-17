package domain

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
)

// Card プレイヤーが使用するカードの情報
type Card struct {
	ID   uuid.UUID
	Type CardType
}

// CardType カードの種類
type CardType uint8

const (
	// YOLO
	// - レア度 : 1
	// - 自分の妨害がないレールを減らす
	// - すべて妨害がある場合は、何も起こらない
	CardTypeYolo CardType = iota

	// Galaxy Brain
	// - レア度 : 2
	// - カードをすべて捨てる ⇒ すべてのカードのステップ数を+1する
	CardTypeGalaxyBrain

	// Open Sourcerer
	// - レア度 : 3
	// - プレイヤーの体力を回復する
	// - 回復力 : 30
	CardTypeOpenSourcerer

	// Refactoring
	// - レア度 : 2
	// - 自身のレールに妨害を発生させる
	// - 妨害値 : 1
	// - 攻撃力 : 5
	CardTypeRefactoring

	// Pair Extraordinaire
	// - レア度 : 1
	// - 他プレイヤー1人をランダムに選択される
	// - レールに妨害を発生させる。
	// - 妨害値 : 2
	// - 攻撃力 : 30
	CardTypePairExtraordinaire

	// LGTM
	// - レア度 : 1
	// - 他プレイヤー1人をランダムに選択される
	// - レールに妨害を発生させる
	// - 妨害値 : 3
	// - 攻撃力 : 20
	CardTypeLgtm

	// Pull Shark
	// - レア度 : 2
	// - 他プレイヤー1人をランダムに選択される
	// - 敵のレールを増やす
	CardTypePullShark

	// Starstruck
	// - レア度 : 3
	// - 他プレイヤー1人をランダムに選択される
	// - レールに妨害を発生させる。
	// - 妨害値 : 5
	// - 攻撃力 : 50
	CardTypeStarstruck

	// Ooops!!!
	// - レア度 : 2
	// - 即時自動発動
	// - レールに妨害を発生させる
	// - 妨害値 : 1 ~ 3
	// - 攻撃力 : 妨害値 * 10
	CardTypeOoops

	// None
	// - カードを使用しないイベント用
	CardTypeNone
)

func NewCard(id uuid.UUID, typ CardType) *Card {
	return &Card{
		ID:   id,
		Type: typ,
	}
}

// DelayAndAttack 妨害値と攻撃力を返す
// いずれかがない場合はエラーを返す
func (t CardType) DelayAndAttack() (int, float32, error) {
	errCannotUse := errors.New("このカードを妨害に使用することはできません")

	switch t {
	case CardTypeYolo:
		return 0, 0, errCannotUse
	case CardTypeGalaxyBrain:
		return 0, 0, errCannotUse
	case CardTypeOpenSourcerer:
		return 0, 0, errCannotUse
	case CardTypeRefactoring:
		return 1, 5, nil
	case CardTypePairExtraordinaire:
		return 2, 30, nil
	case CardTypeLgtm:
		return 3, 20, nil
	case CardTypePullShark:
		return 0, 0, errCannotUse
	case CardTypeStarstruck:
		return 5, 50, nil
	case CardTypeOoops:
		delay := rand.Intn(3) + 1
		return delay, float32(delay) * 10, errCannotUse
	case CardTypeNone:
		return 0, 0, errCannotUse
	default:
		return 0, 0, errCannotUse
	}
}
