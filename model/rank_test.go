package model

import (
	"fmt"
	"testing"
	"time"
)

func TestRank(t *testing.T) {
	fmt.Println("Hello world")
	RankMapInit(2)

	rankItems := []ArticleRankItem{
		{1, 3, nil},
		{2, 2, nil},
		{3, 1, nil},
		{4, 25, nil},
	}
	AddNewArticleList(1, rankItems)
	fmt.Println(GetTopicListByPageNum(1, 1, 2))

	<-time.After(3 * time.Second)
	fmt.Println(GetTopicListByPageNum(1, 1, 4))
}
