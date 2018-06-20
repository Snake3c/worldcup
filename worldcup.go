package main

import (
	// "bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"net/http"
	"os"
	// "regexp"
	// "strconv"
	"encoding/json"
	// "strings"
	// "time"
)

const (
	// STARTURL string = "http://2018.sina.com.cn/ballgame/player.shtml?id=" //妹子图模块列表页url
	API_player     string = "http://api.sports.sina.com.cn/?p=sports&s=sport_client&a=index&_sport_t_=football&_sport_s_=opta&_sport_a_=getPlayer&id="
	API_player_end string = "&type=108&dpc=1"

	API_team     string = "http://api.sports.sina.com.cn/?p=sports&s=sport_client&a=index&_sport_t_=football&_sport_s_=opta&_sport_a_=teamStatics&show_players=1&type=108&season=2017&id="
	API_team_end string = "&dpc=1"

	API_query_tp     string = "http://api.sports.sina.com.cn/?p=sports&s=sport_client&a=index&_sport_t_=football&_sport_s_=opta&_sport_a_=getTeamPlayers&type=108&season=2017&id="
	API_query_tp_end string = "&dpc=1"
)

// p: sports
// s: sport_client
// a: index
// _sport_t_: football
// _sport_s_: opta
// _sport_a_: teamStatics
// show_players: 1
// type: 108
// callback: cb_teamStatsCur_f23b83f6_2343_4ee2_8419_8ed143d954c1
// season: 2017
// id: 941
// dpc: 1

// "result":{
//         "status":{
//             "code":0,
//             "msg":""
//         },
//         "data":{
//             "id":"51247",
//             "team_id":"536",
//             "type_id":"4",
//             "league_id":"2017",
//             "team_name":"Russia",
//             "player_name":"Vladimir Granat",
//             "player_name_cn":"格拉纳特",
//             "position":"2",
//             "real_position":"Central Defender",
//             "real_position_side":"Centre",
//             "birth_date":"1987-05-22",
//             "weight":"78",
//             "height":"184",
//             "jersey_num":"14",
//             "country":"Russia",
//             "club":"Rubin Kazan",
//             "team_name_cn":"俄罗斯",
//             "country_cn":"俄罗斯",
//             "position_cn":"后卫",
//             "club_cn":"喀山鲁宾",
//             "age":"31",
//             "pic":"http://www.sinaimg.cn/ty/opta/players/51247.jpg",
//             "sl_team_id":"941",
//             "sl_type_id":"108",
//             "team_logo":"http://www.sinaimg.cn/lf/sports/logo85/941.png",
//             "team_colors":"http://www.sinaimg.cn/lf/sports/kits/2015/home/941.png"
//         }
//     }

type InfoRsp struct {
	Id             string `json:"id"`
	Player_name    string `json:"player_name"`
	Player_name_cn string `json:"player_name_cn"`
	Birth_date     string `json:"birth_date"`
	Weight         string `json:"weight"`
	Height         string `json:"height"`
	Jersey_num     string `json:"jersey_num"`
	Country_cn     string `json:"country_cn"`
	Position_cn    string `json:"position_cn"`
	Club_cn        string `json:"club_cn"`
	Age            string `json:"age"`
	Pic            string `json:"pic"`
}

type TeamRsp struct {
	Sl_team_id   int    `json:"sl_team_id"`
	Team_name_cn string `json:"team_name_cn"`
}

type resultData struct {
	Player []InfoRsp `json:"player"`
	Team   TeamRsp   `json:"team"`
}

type statusRsp struct {
	Data resultData `json:"data"`
}

type playerInfo struct {
	Result statusRsp `json:"result"`
}

type QueryPlayer struct {
	Result QueryPlayerData `json:"result"`
}
type QueryPlayerData struct {
	Data []InfoRsp `json:"data"`
}

func init() {

	countrys := []string{"俄罗斯", "比利时", "德国", "英格兰", "西班牙", "法国", "波兰", "冰岛", "葡萄牙", "塞尔维亚", "瑞士", "克罗地亚", "丹麦", "瑞典", "伊朗", "日本", "韩国", "沙特阿拉伯", "澳大利亚",
		"巴西", "乌拉圭", "阿根廷", "哥伦比亚", "秘鲁", "尼日利亚", "埃及", "突尼斯", "塞内加尔", "摩洛哥", "墨西哥", "哥斯达黎加", "巴拿马"}

	for _, country := range countrys {

		//此处路径替换成自己的本地路径
		dir := "/Users/leoyu/worldcup/world_HD/" + country + "/"
		creatPath(dir)
	}

	// return

	// i := 19724 //沙特阿拉伯 id = 19725
	//从 0开始遍历球队，大概遍历到 4000多就只剩下沙特阿拉伯了
	i := 0
	for {

		// url := API + fmt.Sprintf("%d", i) + API_2
		url := API_team + fmt.Sprintf("%d", i) + API_team_end

		beego.Error(url)

		i++

		_, html, _ := getHtml(url)

		// beego.Error(html)

		playerInfo := playerInfo{}

		err := json.Unmarshal([]byte(html), &playerInfo)
		if err != nil {
			beego.Error(err.Error())
			continue
		}

		if playerInfo.Result.Data.Team.Team_name_cn != "" {
			teamid := fmt.Sprintf("%d", playerInfo.Result.Data.Team.Sl_team_id)

			url := API_query_tp + teamid + API_query_tp_end

			_, html, _ := getHtml(url)

			queryp := QueryPlayer{}

			err := json.Unmarshal([]byte(html), &queryp)
			if err != nil {
				beego.Error(err.Error())
				continue
			}

			for _, p := range queryp.Result.Data {

				filename := p.Player_name_cn + "_" + p.Country_cn + "_" + p.Position_cn + ".jpg"
				// p_url := p.Pic
				dir := "/Users/leoyu/worldcup/world_HD/" + playerInfo.Result.Data.Team.Team_name_cn + "/"

				// beego.Error(p.Id)
				// beego.Error(p.Pic)
				p_url := "http://www.sinaimg.cn/lf/sports/wc_2018/player/" + p.Id + ".jpg"

				req := httplib.Get(p_url)
				// err := req.ToFile(dir)

				b, err := req.Bytes()

				if err != nil {
					beego.Error(err.Error())
				}

				err = saveImageTo(dir, filename, b)

				if err != nil {
					beego.Error(err.Error())
				}
				beego.Error(p.Country_cn, " ", p.Player_name_cn, " successed")

			}
		}
	}
}

func saveImageToII(path string, filename string, imgByte []byte) error {
	pInfo, pErr := os.Stat(path)
	if pErr != nil || pInfo.IsDir() == false {
		errDir := os.Mkdir(path, os.ModePerm)
		if errDir != nil {
			fmt.Println(errDir)
			os.Exit(-1)
			return errDir
		}
	}
	fn := path + "/" + filename
	_, fErr := os.Stat(fn)
	var fh *os.File
	if fErr != nil {
		fh, _ = os.Create(fn)
	} else {
		fh, _ = os.Open(fn)
	}
	defer fh.Close()
	fh.Write(imgByte)

	return nil
}

func saveImageTo(Imagepath string, fileName string, image []byte) error {

	//判断当前逻辑分库对应的文件夹是否存在
	_, err := os.Stat(Imagepath)
	if os.IsNotExist(err) {
		beego.Error("exsit")
		//不存在，则根据逻辑分库创建文件夹
		err = creatPath(Imagepath)

		if err != nil {
			beego.Error(err.Error())
			return err
		}
	}

	// //4、保存图片信息
	fileName = Imagepath + fileName
	err = ioutil.WriteFile(fileName, image, 0666)
	if err != nil {

		return err
	}

	return nil
}

func creatPath(Imagepath string) error {
	_, err := os.Stat(Imagepath)
	if os.IsExist(err) {
		return nil
	}
	//不存在，则创建
	err = os.MkdirAll(Imagepath, 0666)
	if err != nil {
		beego.Error("创建图片保存路径失败!")
		return err
	}
	return nil
}

func main() {
}

//下载html
func getHtml(url string) (error, string, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	html, err1 := ioutil.ReadAll(response.Body)
	return err, string(html), err1
}
