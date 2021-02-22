/********************************************

@Author :yinjinlin<yinjinlin_uplook@163.com>
@Time : 2021/2/22 11:30
@Description:

*********************************************/
package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"strings"
)

// 自定义变量
var (
	maxSearchNum = 2000  // 单次搜索最大返回数量,防止陷入深度搜索，拖垮搜索服务器
	numShareds   = 3     // 每个索引的分片数
	numReplicas  = 1     // 每个索引每个分片对应的副本数
	indexPre     string  // 索引前缀
	indexPrefix  = "kc-" // ik_smart=插件，中文分析器,对中文分词支持优秀 ik_max_word=插件，加强版ik_smart,对短语分词更细腻

	client *elasticsearch.Client
)

// 定义结构体
type ElasticClient struct {
	index     string // 索引名称
	indexType string // 索引类型
	indexPre  string // 索引前缀

	term    []map[string]map[string]interface{} // 聚合刷选条件
	terms   []map[string]map[string]interface{} // 聚合筛选条件In
	should  []map[string]map[string]interface{} // 聚合筛选条件Or
	mustNot []map[string]map[string]interface{} // 聚合排除筛选条件

	keyword []map[string]map[string]interface{} // 文本刷选
	ranges  []map[string]map[string]interface{} // 范围筛选条件
	sort    []map[string]interface{}            // 排序
	limit   map[string]int                      // 分页
	geo     map[string]interface{}              // 地理位置
	factor  map[string]interface{}              // 欢迎度

	lastQuery string // 打印请求json
	// 客户端的请求与响应
	client   *elasticsearch.Client // es request client
	response *esapi.Response       // es response client

}

// 初始化一些配置信息
func (e *ElasticClient) Init() {
	// 基础配置信息
	hostAddr := "xxxx"
	indexPre := "xxxx"
	username := "xxxx"
	password := "xxxx"

	fmt.Println(hostAddr)
	fmt.Println(indexPre)
	fmt.Println(username)
	fmt.Println(password)

	var err error
	// 客户端端配置,
	cfg := elasticsearch.Config{
		//Addresses []string // 要使用的Elasticsearch节点列表
		//Username  string   // HTTP基本身份验证的用户名
		//Password  string   // HTTP基本认证密码.
	}
	client, err = elasticsearch.NewClient(cfg)

	if err != nil {
		err = errors.New(fmt.Sprintf("es NewCient error,error is %s", err.Error()))
		log.Panic(err)
		return
	}
}

// 文档字段值精确值筛选
// @docField             字段名
// @interface{} value    字段值 类型可以是整型、字符串
func (e *ElasticClient) SetFilter(docField string, value interface{}) *ElasticClient {
	maps := map[string]map[string]interface{}{
		"term": {docField: value},
	}
	if len(e.term) == 0 {
		e.term = make([]map[string]map[string]interface{}, 0)
		e.term = append(e.term, maps)
	} else {
		e.term = append(e.term, maps)
	}

	return e
}

// 关键字搜索 (字符串过滤)
func (e *ElasticClient) Search(keyword string, docField []string) *ElasticClient {
	maps := map[string]map[string]interface{}{
		"query_string": {
			"query":  strings.TrimSpace(keyword),
			"fields": docField,
		},
	}
	if len(e.keyword) == 0 {
		e.keyword = make([]map[string]map[string]interface{}, 0)
		e.keyword = append(e.keyword, maps)
	} else {
		e.keyword = append(e.keyword, maps)
	}
	return e
}

// 文档字段刷选大于
func (e *ElasticClient) SetFilterGt(docField string, value interface{}) *ElasticClient {
	maps := map[string]map[string]interface{}{
		"range": {docField: map[string]interface{}{
			"gt": value,
		}},
	}
	if len(e.ranges) == 0 {
		e.ranges = make([]map[string]map[string]interface{}, 0)
		e.ranges = append(e.ranges, maps)
	} else {
		e.ranges = append(e.ranges, maps)
	}
	return e
}

// 文档字段筛选 大于等于
func (e *ElasticClient) SetFilterGte(docField string, value interface{}) *ElasticClient {
	maps := map[string]map[string]interface{}{
		"range": {docField: map[string]interface{}{
			"gte": value,
		}},
	}
	if len(e.ranges) == 0 {
		e.ranges = make([]map[string]map[string]interface{}, 0)
		e.ranges = append(e.ranges, maps)
	} else {
		e.ranges = append(e.ranges, maps)
	}
	return e
}

// 文档字段筛选 小于
func (e *ElasticClient) SetFilterLt(docField string, value interface{}) *ElasticClient {
	maps := map[string]map[string]interface{}{
		"range": {docField: map[string]interface{}{
			"lt": value,
		}},
	}
	if len(e.ranges) == 0 {
		e.ranges = make([]map[string]map[string]interface{}, 0)
		e.ranges = append(e.ranges, maps)
	} else {
		e.ranges = append(e.ranges, maps)
	}
	return e
}

// 文档字段值范围筛选不在范围内
// @docField            文档字段名称
// @starRange           字段值范围开始
// @endRange            字段值范围结束
// @params              可变参数 左边界【gt:大于、gte:大于等于】 右边界【lt:小于、lte:小于等于】
func (e *ElasticClient) SetFilterNotRange(docField string, starRange interface{}, endRange interface{}, params ...interface{}) *ElasticClient {
	var leftLimiter string
	var rightLimter string
	if len(params) == 0 {
		leftLimiter = "gte"
		rightLimter = "lte"
	} else {
		if len(params) == 1 {
			leftLimiter = params[0].(string)
			rightLimter = "lte"
		} else {
			leftLimiter = params[0].(string)
			rightLimter = params[1].(string)
		}
	}

	maps := map[string]map[string]interface{}{
		"range": {docField: map[string]interface{}{
			leftLimiter: starRange,
			rightLimter: endRange,
		}},
	}
	if len(e.mustNot) == 0 {
		e.mustNot = make([]map[string]map[string]interface{}, 0)
		e.mustNot = append(e.mustNot, maps)
	} else {
		e.mustNot = append(e.mustNot, maps)
	}
	return e
}

// 文档字段值范围筛选
// @docField        文档字段名称
// @minValue        字段值范围开始
// @maxValue        字段值范围结束
// @params          可变参数 左边界【gt:大于、gte:大于等于】 右边界【lt:小于、lte:小于等于】
func (e *ElasticClient) SetFilterRange(docField string, minV interface{}, maxV interface{}, params ...interface{}) *ElasticClient {
	var leftLimiter string
	var rightLimter string
	if len(params) == 0 {
		leftLimiter = "gte"
		rightLimter = "lte"
	} else {
		if len(params) == 1 {
			leftLimiter = params[0].(string)
			rightLimter = "lte"
		} else {
			leftLimiter = params[0].(string)
			rightLimter = params[1].(string)
		}
	}
	maps := map[string]map[string]interface{}{
		"range": {docField: map[string]interface{}{
			leftLimiter: minV,
			rightLimter: maxV,
		}},
	}
	if len(e.ranges) == 0 {
		e.ranges = make([]map[string]map[string]interface{}, 0)
		e.ranges = append(e.ranges, maps)
	} else {
		e.ranges = append(e.ranges, maps)
	}
	return e
}

// 设置分页
func (e *ElasticClient) SetLimit(start, limit int) *ElasticClient {
	searchNum := start * limit
	if searchNum > maxSearchNum {
		start -= 1
	}
	e.limit = make(map[string]int)
	e.limit["from"] = start
	e.limit["size"] = limit
	return e
}

// 设置排序 默认升序
func (e *ElasticClient) SetSortMode(docField string, sort ...string) *ElasticClient {
	var sortMode string
	if len(sort) == 0 {
		sortMode = "desc"
	} else {
		sortMode = sort[0]
	}
	maps := map[string]interface{}{
		docField: map[string]string{"order": sortMode},
	}

	if len(e.sort) == 0 {
		e.sort = make([]map[string]interface{}, 0)
		e.sort = append(e.sort, maps)
		e.sort = append(e.sort, map[string]interface{}{"_score": map[string]string{"order": "desc"}})
	} else {
		e.sort = append(e.sort, maps)
	}
	return e
}

func (e *ElasticClient) Query() map[string]interface{} {
	requestBoby := make(map[string]interface{})
	// 1: 封装requesBoby
	// 如调用getSearch

	byteBody, _ := json.Marshal(requestBoby)
	// 2: 利用searchRequest 配置请求
	req := esapi.SearchRequest{
		Index:        []string{e.index}, // 索引名称
		DocumentType: []string{e.indexType},
		Body:         bytes.NewReader(byteBody),
	}
	// 3: 请求
	response, err := req.Do(context.Background(), e.client)
	if err != nil {
		log.Println(err)
		return make(map[string]interface{})
	}
	e.response = response
	// 4：关闭请求
	defer response.Body.Close()

	// 5:对请求得到的数据处理
	var maps = make(map[string]map[string]interface{})
	responseStr :=  response.String()[strings.Index(response.String(),"{"):]
	_ = json.Unmarshal([]byte(responseStr), &maps)
	var resMap = make(map[string]interface{})
	if len(maps) == 0 {
		return resMap
	}
	if value,ok := maps["hits"]["hits"];ok {
		resMap["result"] = value
	}
	if value,ok := maps["hits"]["total"];ok {
		resMap["total"]=value
	}
	return resMap
}

func (e *ElasticClient) getSearch() (map[string]map[string]map[string]map[string]map[string]interface{}, bool) {
	// 组合map过滤
	search := make(map[string]map[string]map[string]map[string]map[string]interface{})
	search["query"] = make(map[string]map[string]map[string]map[string]interface{})
	search["query"]["function_score"] = make(map[string]map[string]map[string]interface{})
	search["query"]["function_score"]["query"] = make(map[string]map[string]interface{})
	search["query"]["function_score"]["query"]["bool"] = make(map[string]interface{})
	search["query"]["function_score"]["query"]["bool"]["must"] = make([][]map[string]map[string]interface{}, 0)
	search["query"]["function_score"]["query"]["bool"]["must_not"] = make([][]map[string]map[string]interface{}, 0)
	search["query"]["function_score"]["query"]["bool"]["should"] = make([][]map[string]map[string]interface{}, 0)

	var capBool bool
	// 精确值条件
	if len(e.term) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["must"] = append(search["query"]["function_score"]["query"]["bool"]["must"].([][]map[string]map[string]interface{}), e.term)
	}
	// 范围in条件
	if len(e.terms) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["must"] = append(search["query"]["function_score"]["query"]["bool"]["must"].([][]map[string]map[string]interface{}), e.terms)
	}
	// 精确值排除条件
	if len(e.mustNot) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["must_not"] = append(search["query"]["function_score"]["query"]["bool"]["must_not"].([][]map[string]map[string]interface{}), e.mustNot)
	}
	// 范围条件
	if len(e.ranges) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["must"] = append(search["query"]["function_score"]["query"]["bool"]["must"].([][]map[string]map[string]interface{}), e.ranges)
	}
	// 地理位置，筛选附近公里数内的数据
	if len(e.geo) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["filter"] = map[string]interface{}{
			"geo_bounding_box": e.geo["geo_bounding_box"],
		}
	}
	if len(e.should) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["should"] = append(search["query"]["function_score"]["query"]["bool"]["should"].([][]map[string]map[string]interface{}), e.should)
	}
	if len(e.keyword) > 0 {
		capBool = true
		search["query"]["function_score"]["query"]["bool"]["must"] = append(search["query"]["function_score"]["query"]["bool"]["must"].([][]map[string]map[string]interface{}), e.keyword)
	}
	return search, capBool
}

// 实例化引用
func NewElasticClient(index string) (*ElasticClient, error) {
	var err error
	return &ElasticClient{
		index:     "upelastic", // 可以做到配置文件中
		indexType: "_doc",
		client:    client,
	}, err

}
