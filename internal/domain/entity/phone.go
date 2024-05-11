package domain

type Phone struct {
	Value  string `bson:"value"`
	IsChat bool   `bson:"is_chat"`
	IsMain bool   `bson:"is_main"`
}
