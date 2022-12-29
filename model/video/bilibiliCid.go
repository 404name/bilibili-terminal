package video

//https://www.zhaoyanbai.com/articles/BilibiliGolangDownloader

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/404name/termui-demo/global"
)

var qnMap = map[string]struct {
	QN         int
	NeedCookie bool
	Detail     string
}{
	"4K":      {120, true, "超清 4K"},
	"1080P60": {116, true, "高清1080P60"},
	"1080P+":  {112, true, "高清1080P+"},
	"1080P":   {80, true, "1080P"},
	"720P60":  {74, true, "高清720P60"},
	"720P":    {64, false, "高清720P"},
	"480P":    {32, false, "清晰480P"},
	"360P":    {16, false, "流畅360P"},
}

type bilibiliCid struct {
	Bvid     string
	Cid      string
	Title    string
	QN       int
	PlayURLs []string
}

var (
	bvid        string
	page        int // 视频合集的分p, 如果不指定，默认为0表示全下
	dir         string
	qn          int
	sessionData string
	prePos      int64 // 用于记录下载百分比的
)

func init() {
	var _qn string

	flag.StringVar(&bvid, "b", "BV1xx411c79H", "bvid with BV prefix in the url")
	flag.IntVar(&page, "p", 0, "single video page number of video list")
	flag.StringVar(&dir, "d", "./", "which directory to save")
	flag.StringVar(&_qn, "q", "360p", "video quality")
	flag.StringVar(&sessionData, "s", "", "value of your bilibili cookie[\"SESSDATA\"]")
	flag.Parse()

	_qn = strings.ToUpper(_qn)

	if bvid == "" {
		log.Fatalf("you must specify bvid")
	}

	if v, ok := qnMap[_qn]; ok {
		qn = v.QN

		if v.NeedCookie && sessionData == "" {
			log.Fatalf("need set value of bilibili cookie[\"SESSDATA\"] when you download %v", v.Detail)
		}
	} else {
		log.Fatalf("invalid qn value")
	}

}

// func main() {
// 	if bvid == "" {
// 		log.Fatalf("invalid bvid\n")
// 	}

// 	global.Log.Errorf("%v@%v\n", bvid, qn)

// 	videos := getCidList(bvid, qn)
// 	for i, v := range videos {
// 		if page == 0 || page == (i+1) {
// 			v.download(dir)
// 		}
// 	}
// }

type downloader struct {
	io.ReadCloser
	Total   int64
	Current int64
}

func (d *downloader) Read(p []byte) (n int, err error) {
	if Player.Ready {
		d.ReadCloser.Close()
		return
	}

	n, err = d.ReadCloser.Read(p)

	d.Current += int64(n)
	nowPos := d.Current * 100 / d.Total

	// 每2%进度刷新一下，共刷新50次
	if (nowPos > prePos || d.Current == d.Total) && !Player.Ready {
		prePos = nowPos
		VideoDownloadBarRender(int(d.Current), int(d.Total))
	}

	if d.Current == d.Total {
		global.Log.Errorln("关闭Reader")
		d.ReadCloser.Close()
	}
	return
}

func (c *bilibiliCid) download(dir string) {
	prePos = 0
	for i, URL := range c.PlayURLs {
		u, err := url.Parse(URL)
		if err != nil {
			global.Log.Errorf("ERR: %v", err)
			continue
		}

		name := fmt.Sprintf("%v_%v_%v", c.Bvid, c.Title, path.Base(path.Base(u.Path)))
		global.Log.Infof("Downloading[%d]: name:%v\n\turl:%v\n", i, name, URL)

		client := &http.Client{}
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			log.Println(err)
			return
		}
		setUserAgent(req)
		setCookie(req)
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Range", "bytes=0-")                               // Range 的值要为 bytes=0- 才能下载完整视频
		req.Header.Set("Referer", "https://www.bilibili.com/video/"+bvid) // 必需添加
		req.Header.Set("Origin", "https://www.bilibili.com")
		req.Header.Set("Connection", "keep-alive")

		rsp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer rsp.Body.Close()

		// global.Log.Errorf("save to: %v", dir)
		out, err := os.Create(dir)
		if err != nil {
			global.Log.Errorf("err: %v", err)
			continue
		}
		defer out.Close()

		dr := &downloader{
			rsp.Body,
			rsp.ContentLength,
			0,
		}
		io.Copy(out, dr)
		global.Log.Infof("视频下载完成")
	}
}

func (c *bilibiliCid) getPlayURLs() {
	url := fmt.Sprintf("https://api.bilibili.com/x/player/playurl?bvid=%v&cid=%v&qn=%v&fourk=1", c.Bvid, c.Cid, c.QN)
	fmt.Println(url)

	pl := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Quality int `json:"quality"`
			Durl    []struct {
				Order     int    `json:"order"`
				URL       string `json:"url"`
				BackupURL string `json:"backup_url"`
			} `json:"durl"`
		} `json:"data"`
	}{}

	data := rawGetURL(url, setCookie)
	fmt.Println(data)
	json.Unmarshal([]byte(data), &pl)

	for i, p := range pl.Data.Durl {
		global.Log.Infof("PlayList[%d]: quality %v order %v url %v %v", i, pl.Data.Quality, p.Order, p.URL, p.BackupURL)
		c.PlayURLs = append(c.PlayURLs, p.URL)
	}
}

func getCidList(bvid string, qn int) []bilibiliCid {
	cl := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Cid       int64  `json:"cid"`
			Page      int    `json:"page"`
			Part      string `json:"part"`
			Duration  int    `json:"duration"`
			Vid       string `json:"vid"`
			Dimension struct {
				Width  int `json:"width"`
				Height int `json:"height"`
				Rotate int `json:"rotate"`
			} `json:"Dimension"`
		} `json:"data"`
	}{}

	data := getURL("https://api.bilibili.com/x/player/pagelist?bvid=" + bvid)

	json.Unmarshal([]byte(data), &cl)

	if len(cl.Data) == 0 {
		global.Log.Errorf("ERR: get cid list failed")
	}

	var cids []bilibiliCid
	for i, d := range cl.Data {
		c := bilibiliCid{}
		c.Cid = strconv.FormatInt(d.Cid, 10)
		c.Title = d.Part
		c.Bvid = bvid
		c.QN = qn
		c.getPlayURLs()
		cids = append(cids, c)
		global.Log.Infof("CidList[%d]: %v %v %v %v", i, d.Cid, d.Part, d.Dimension.Width, d.Dimension.Height)
	}

	return cids
}

func setUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
}
func setCookie(req *http.Request) {
	cookie := http.Cookie{Name: "SESSDATA", Value: sessionData, Expires: time.Now().Add(30 * 24 * 60 * 60 * time.Second)}
	global.Log.Infof("got bilibili cookie, SESSDATA:%v", sessionData)
	req.AddCookie(&cookie)
}

func getURL(url string) string {
	return rawGetURL(url, nil)
}

func rawGetURL(url string, headerSet func(*http.Request)) (s string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	setUserAgent(req)

	if headerSet != nil {
		headerSet(req)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		global.Log.Errorf("http return %v\n", res.StatusCode)
		return
	}

	rsp, _ := ioutil.ReadAll(res.Body)

	s = string(rsp)

	return
}
