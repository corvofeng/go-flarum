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
		{1, 3},
		{2, 2},
		{3, 1},
		{4, 25},
	}
	AddNewArticleList(1, rankItems)
	fmt.Println(GetTopicListByPageNum(1, 1, 2))

	<-time.After(3 * time.Second)
	fmt.Println(GetTopicListByPageNum(1, 1, 4))
}
