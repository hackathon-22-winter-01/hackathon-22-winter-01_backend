package domain

// Card プレイヤーが使用するカードの情報
type Card struct {
	ID   string
	Type CardType
}

// CardType カードの種類
type CardType uint8
