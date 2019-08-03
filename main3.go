package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/ontio/ontology/common/serialization"
	"os"
	"io"
	"sync"
	"time"

	"math/rand"
)

type RR struct {
	AccAddr string
	Bas []Balance
}
type Balance struct {
	ContractAdd string
	Amount int
}
func syncStudy()  {
	var wg sync.WaitGroup
	wg.Add(2)

	base58Addrs := []string{"aa", "bb"}
	contractAddrs := []string{"11", "22"}

	res := make([]RR, 2)

	for i:=0;i<2;i++ {
		go func(i int) {
			defer wg.Done()
			var wg2 sync.WaitGroup
			wg2.Add(2)
			acc := base58Addrs[i]
			res[i].AccAddr = acc
			res[i].Bas = make([]Balance, 2)
			for j:=0;j<2;j++ {
				go func(j int) {
					defer wg2.Done()
					conAddr := contractAddrs[j]
                    res[i].Bas[j].ContractAdd = conAddr
                    res[i].Bas[j].Amount = getBalance(acc, conAddr)
                    time.Sleep(2 * time.Second)
				}(j)
			}
			wg2.Wait()
		}(i)
	}
	wg.Wait()

	fmt.Println("res", res)
	fmt.Println("success")
}

func getBalance(acc string, contractAddr string) int {
	ma := map[string]int{"aa11":1, "aa22":2, "bb11":3,"bb22":4}
	return ma[acc+contractAddr]
}
func con(param interface{}){
	p , ok := param.([]interface{})
	fmt.Println(ok)
	fmt.Println(p)
}

func cc(p []interface{}) {
	switch p[0].(type) {
	case int:
		fmt.Println("p int:", p)
	case string:
		fmt.Println("p string:", p)
	default:
		fmt.Println("p default:", p)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i:= 0;i<10;i++ {
		x := rand.Uint64()
		fmt.Println("x:", x)
		y := rand.Intn(10)
		fmt.Println("y:", y)
	}
	//cc([]interface{}{1,"sss"})

	//path:="/api/v1/shardtxstate/:txhash/:notifyid"
	//path = "^" + path + "$"
	//matches := regexp.MustCompile(`:(\w+)`).FindAllStringSubmatch(path, -1)
	//params := make([]string, 0)
	//for _, v := range matches {
	//	params = append(params, v[1])
	//	path = strings.Replace(path, v[0], `(\w+)`, 1)
	//}
	//compiledPath, err := regexp.Compile(path)
	//if err != nil {
	//	panic(err)
	//}
	//
	//boo := compiledPath.MatchString("/api/v1/shardtxstate/fb64410d900a237ee63502af33f68944ef01006330213fdca8d29d1e7d06c44e")
	//fmt.Println("boo:", boo)
	//txState := xshard_state.CreateTxState(xshard_types.ShardTxID("111"))
	//fmt.Println(txState.ExecState)
	//txState2 := txState
	//txState2.ExecState = xshard_state.ExecYielded
	//
	//fmt.Println(txState.ExecState)

	//ss, err := hex.DecodeString("123g")
	//if err != nil {
	//	fmt.Println("err:", err)
	//}

	//ss := 123
	//con(ss)
	//syncStudy()
	//path := "/api/v1/smartcode/event/txhash/:hash/:sourcetxhash"
	//path2 := "/api/v1/smartcode/event/txhash/:hash/"
	//
	//initGet(path)
	//initGet(path2)



    //a :=  big.NewInt(1000000000000000000)
	//
    //fmt.Println(a.Uint64())
    //astr := fmt.Sprintf("%d", a.Uint64())
    //aa, err := strconv.Atoi(astr)
    //if err != nil {
    //	fmt.Println("err:", err)
	//}
    //fmt.Println("aa", aa)

    //id := int64(10)
    //fmt.Println(byte(id))

	//a := common.Uint256{}
	//fmt.Println("a:", a)
	//
	//b:= common.UINT256_EMPTY
	//fmt.Println("b:", b)
	//
	//fmt.Println(a == b)

	//u, err := common.Uint256FromHexString("0000000000000000000000000000000000000000000000000000000000000001")
	//if err !=nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("txHash: %x\n", u)
	//fmt.Printf("txHash: %v\n", u)


	//cache,_ := lru.NewARC(3)
	//cache.Add("sss", nil)
	//cache.Add("ss", nil)
	//cache.Add("s", nil)
	//fmt.Println(cache.Contains("sss"))
	//fmt.Println(cache.Contains("ss"))
	//m := make(map[string]bool)
	//m["sss"] = true
	//fmt.Println(m["sss"])
	//fmt.Println(m["ss"])
	//writeFile()
	//readFile()
}

func readFile() {
	f, err := getF("test.txt")
	if err != nil {
		return
	}
	r := bufio.NewReader(f)
	key := make([]byte, 4, 4)
	r.Read(key)
	h := binary.LittleEndian.Uint32(key)
	for i:=uint32(0);i<h;i++ {
		val, err := serialization.ReadVarBytes(r)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("val:", string(val))
	}
}

func writeFile(){
	f, err := getF("test.txt")
	if err != nil {
		return
	}
	key := make([]byte,4,4)
	_, err = f.ReadAt(key, 0)
	writer := bufio.NewWriter(f)
	var num uint32
	if err != nil {
		if err == io.EOF {
			num = 0
			serialization.WriteUint32(writer,0)
		} else {
			fmt.Println("err:", err)
			return
		}
	} else {
		num = binary.LittleEndian.Uint32(key)
	}


	//defer func() {
	//	writer.Flush()
	//	f.Close()
	//}()
	for i:=uint32(0);i<5 ;i++ {
		err = serialization.WriteVarBytes(writer, []byte("test"))
		if err != nil {
			fmt.Println("err:", err)
			return
		}
	}
	writer.Flush()
	num += 5

	key = make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(key, num)
	f.WriteAt(key,0)
	f.Sync()

	f.Close()
}



func getF(fileName string) (*os.File, error) {
	var f *os.File
	var err error
	if checkFileIsExist(fileName) {
		f, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("OpenFile err: %s\n", err)
			return nil, err
		}
	} else {
		f, err = os.Create(fileName)
		if err != nil {
			fmt.Errorf("Create err: %s\n", err)
			return nil, err
		}
	}
	return f, nil
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}