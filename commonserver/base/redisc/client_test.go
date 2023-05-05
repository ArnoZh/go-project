// Package redisc .
 
package redisc

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	TestDB = 1
	// RemoteDBUrl = "127.0.0.1:6379"
	LocalDBUrl = ":6379"
	DBUrl      = LocalDBUrl
)

//func TestRedis(t *testing.T) {
//	client := DC()
//	_ = client.SelectDB(TestDB)
//	status := client.PoolStats()
//	fmt.Printf("%#v\n", status)
//
//	if _, err := client.LPush("1.list", 1, 2, 3, 4, 5, 6, 7); err != nil {
//		t.Errorf("LPush error %#v", err)
//	}
//	l, err := client.LRange("1.list", 0, -1)
//	if err != nil {
//		t.Errorf("LRange error %#v", err)
//	}
//
//	fmt.Printf("%#v\n", l)
//
//	res, err := client.LTrim("1.list", 0, 6)
//	fmt.Printf("LTrim %#v %#v\n", res, err)
//
//	l, err = client.LRange("1.list", 0, -1)
//	if err != nil {
//		t.Errorf("LRange error %#v", err)
//	}
//	fmt.Printf("%#v\n", l)
//
//}

type Post struct {
	PostID        int64
	LikeNum       int64
	Content       string
	CreateAt      int32
	LastCommentAt int32
	LastComment   string
}

//func TestRedisPipeline(t *testing.T) {
//	post := &Post{
//		ID:        111222,
//		LikeNum:       99,
//		Txt:       "afsfadf",
//		CreateAt:      111111111,
//		LastComAt: 222222222,
//		LastComment:   "hello",
//	}
//
//	nowMs := util.NowTs()
//	pipeline := DC().TxPipeline()
//	pipeline.Select(TestDB)
//	cmd := pipeline.Exists("p.111222")
//	_, err := pipeline.Exec()
//
//	if err != nil {
//		t.Errorf("pipeline exec error %v", err)
//	}
//
//	_, err = cmd.Result()
//	if err != nil {
//		t.Errorf("pipeline exec error %v", err)
//	}
//
//	cmd0 := pipeline.HMSetObject("test.111222", post, 0)
//	cmd1 := pipeline.HIncrBy("test.111222", "LikeNum", 1)
//	cmd2 := pipeline.ZAdd("test.rank", redis.Z{Score: float64(nowMs), Member: 111222})
//	cmd3 := pipeline.HGetAll("test.111222")
//
//	_, err = pipeline.Exec()
//	if err != nil {
//		t.Errorf("pipeline exec error %v", err)
//	}
//
//	t.Logf("%v", cmd0.Val())
//	t.Logf("%v", cmd1.Val())
//	t.Logf("%v", cmd2.Val())
//	t.Logf("%v", cmd3.Val())
//
//}

//func TestScanPost(t *testing.T) {
//	c, _ := NewClient(RemoteDBUrl, "", TestDB)
//	defer c.Close()
//	post := &Post{
//		ID:        1111,
//		LikeNum:       99,
//		Txt:       "afsfadf",
//		CreateAt:      111111111,
//		LastComAt: 222222222,
//		LastComment:   "hello",
//	}
//	err := c.HMSetObject("test.1111", post, 0)
//	if err != nil {
//		t.Errorf("test post set :%v", err)
//	}
//
//	ret, err := c.HMGetObject("test.1111", &Post{})
//	if err != nil {
//		t.Errorf("read struct error %v", err)
//	}
//	pNew, ok := ret.(*Post)
//	if !ok {
//		t.Errorf("error %#v %v\n", ret, ok)
//	} else {
//		fmt.Printf("read struct %#v \n", pNew)
//	}
//}

// NewID 生成帖子key
func NewID() int64 {
	id, err := DC().Incr("id.gen")
	if err != nil {
		return 0
	}
	return id
}

func BenchmarkWriteStruct(b *testing.B) {
	c, _ := NewClient(DBUrl, "", TestDB)
	defer c.Close()
	c.FlushDB()
	b.ResetTimer()
	suc, fail := 0, 0
	for i := 0; i < b.N; i++ {
		post := NewPost()
		err := c.HMSetObject(strconv.FormatInt(post.PostID, 10), post, 0)
		if err != nil {
			fail++
			b.Error(err)
			continue
		}
		suc++
	}
	if fail > 0 {
		b.Logf("total:%v suc:%v fail:%v\n", b.N, suc, fail)
	}
}

func BenchmarkReadStruct(b *testing.B) {
	c, _ := NewClient(DBUrl, "", TestDB)
	defer c.Close()
	b.ResetTimer()
	suc, fail := 0, 0
	for i := 0; i < b.N; i++ {
		ret, err := c.HMGetObject(strconv.FormatInt(int64(i), 10), &Post{})
		if err != nil {
			fail++
			b.Error(err)
			continue
		}
		_, ok := ret.(*Post)
		if !ok {
			fail++
			b.Error(err)
			continue
		}
		suc++
	}
	if fail > 0 {
		b.Logf("total:%v suc:%v fail:%v\n", b.N, suc, fail)
	}
}

func BenchmarkReadWriteStruct(b *testing.B) {
	c, _ := NewClient(DBUrl, "", TestDB)
	defer c.Close()
	c.FlushDB()
	b.ResetTimer()
	suc, fail := 0, 0
	for i := 0; i < b.N; i++ {
		post := NewPost()
		err := c.HMSetObject(strconv.FormatInt(post.PostID, 10), post, 0)
		if err != nil {
			fail++
			b.Error(err)
			continue
		}
		ret, err := c.HMGetObject(strconv.FormatInt(post.PostID, 10), &Post{})
		if err != nil {
			fail++
			b.Error(err)
			continue
		}
		_, ok := ret.(*Post)
		if !ok {
			fail++
			b.Error(err)
			continue
		}
		suc++
	}
	if fail > 0 {
		b.Logf("total:%v suc:%v fail:%v\n", b.N, suc, fail)
	}
}

// func BenchmarkPipelineMultiWrite(b *testing.B) {
// 	_ = DC().SelectDB(TestDB)
// 	DC().FlushDB()
// 	var wg sync.WaitGroup
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
//
// 			post := NewPost()
// 			key := strconv.FormatInt(post.ID, 10)
// 			pipeline := DC().TxPipeline()
// 			pipeline.HMSetObject(key, post, 0)
// 			pipeline.HIncrBy(key, "LikeNum", 1)
// 			cmd2 := pipeline.HGetAll(key)
//
// 			_, err := pipeline.Exec()
// 			if err != nil {
// 				b.Errorf("pipeline exec error %v", err)
// 			}
// 			res, err := cmd2.Result()
// 			if err != nil {
// 				b.Errorf("pipeline exec error %v", err)
// 			}
// 			_, err = ConvertTo(res, &Post{})
// 			if err != nil {
// 				b.Errorf("ConvertTo error %v", err)
// 			}
//
// 			// b.Logf("%#v", newPost)
// 		}()
// 	}
// 	wg.Wait()
// 	b.Logf("poolstatus:%v", DC().PoolStats())
// }

var idCounter = int64(0)

func BenchmarkPipelineWriteWithWorkerPool(b *testing.B) {
	client, _ := NewClient(DBUrl, "", TestDB)
	defer client.Close()
	client.FlushDB()
	b.ResetTimer()
	workNum := 60
	idCounter = 0
	c := make(chan int64, b.N)
	wg := &sync.WaitGroup{}
	wg.Add(workNum)
	for i := 0; i < workNum; i++ {
		go worker(c, b, wg, client)
	}

	for i := 0; i < b.N; i++ {
		c <- int64(i)
	}

	for i := 0; i < workNum; i++ {
		c <- int64(-1)
	}

	wg.Wait()
	close(c)
	// b.Logf("poolstatus:%+v, op:%v", DC().PoolStats(), idCounter)
}

func worker(c <-chan int64, b *testing.B, wg *sync.WaitGroup, client *Client) {
	defer wg.Done()
	for {
		select {
		case id := <-c:
			if id == -1 {
				return
			}

			post := NewPost()
			key := strconv.FormatInt(post.PostID, 10)
			pipeline := client.Pipeline()
			pipeline.HMSetObject(key, post, 0)
			pipeline.HIncrBy(key, "LikeNum", 1)
			cmd2 := pipeline.HGetAll(key)

			_, err := pipeline.Exec()
			if err != nil {
				b.Errorf("pipeline exec error %v", err)
			}
			res, err := cmd2.Result()
			if err != nil {
				b.Errorf("pipeline exec error %v", err)
			}
			_, err = ConvertTo(res, &Post{})
			if err != nil {
				b.Errorf("ConvertTo error %v", err)
			}
			atomic.AddInt64(&idCounter, 1)
		case <-time.After(5 * time.Second):
			return

		}
	}
}

func NewPost() *Post {
	post := &Post{
		PostID:        NewID(),
		LikeNum:       99,
		Content:       "$afsfadfsaHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHXXXXXXXXXXXXBSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSd#",
		CreateAt:      111111111,
		LastCommentAt: 222222222,
		LastComment:   "hello",
	}
	return post
}
