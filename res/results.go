package res

//Res 查询结果
type Res struct {
	Results []Result `json:"results"`
}

//Result 查询结果的每个元素
type Result struct {
	Series []Serie `json:"series"`
}

//Serie 数据序列
type Serie struct {
	Name    string          `json:"name"`
	Columns []string        `json:"columns"`
	Values  [][]interface{} `json:"values"`
}
