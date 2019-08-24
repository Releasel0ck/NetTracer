package models

import (
	"NetTracer/util"
	"database/sql"

	"sort"
	"strconv"
	"strings"

	"github.com/awalterschulze/gographviz"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var graph *gographviz.Graph
var service_ip string
var service_port int
var client_ip string
var client_port int
var connect_count int
var protocol string
var tmp map[string]struct{}

type sortMapS struct {
	key   string
	value int
}

var sortMapL []sortMapS

func CreateDB(dbname string) bool {
	dbFile := `./db/` + dbname + `.db`
	if util.CheckFileIsExist(dbFile) { //判断文件是否存在
		return false
	} else {
		status := util.CreateFile(dbFile) //创建文件是否成功
		if status {
			db, err := sql.Open("sqlite3", dbFile) //连接数据库
			if err == nil {
				DB = db
				netstat_table := `
	CREATE TABLE [netstat](
  		[service_ip] CHAR(15), 
 		[service_port] INT, 
 		[client_ip] CHAR(15), 
	    [client_port] INT, 
        [connect_count] INT, 
        [protocol] CHAR(3));
    CREATE INDEX service_ip_index ON netstat (service_ip);
    CREATE INDEX client_ip_index ON netstat (client_ip);
	`
				DB.Exec(netstat_table) //初始化数据表
				return true
			} else {
				return false
			}
		} else {
			return status
		}
	}
}

//连接数据库
func ConnectDB(dbname string) (*sql.DB, error) {
	dbFile := `./db/` + dbname + `.db`
	db, err := sql.Open("sqlite3", dbFile)
	if err == nil {
		DB = db
		return db, err
	}
	return nil, err
}

//插入新数据
func InsertNetstat(s_ip string, s_port int, c_ip string, c_port int, c_count int, protocol string) error {
	stmt, err := DB.Prepare("INSERT INTO netstat(service_ip, service_port, client_ip,client_port,connect_count,protocol) values(?,?,?,?,?,?)")
	if err != nil {
		return err
	} else {
		_, err = stmt.Exec(s_ip, s_port, c_ip, c_port, c_count, protocol)
		if err != nil {
			return err
		}
		return nil
	}
}

//按连接数据量排序并返回前6个
func RankByMostCon() string {
	if DB == nil {
		return ""
	}

	rankString := ""
	rows, err := DB.Query("select service_ip,service_port,connect_count from netstat")
	portInfo := make(map[string]int)
	ipInfo := make(map[string]int)
	for rows.Next() {
		err = rows.Scan(&service_ip, &service_port, &connect_count)
		if err != nil {
			return ""
		}
		//util.CheckErr(err)
		service_port_str := strconv.Itoa(service_port)
		_, ok := portInfo[service_port_str]
		if ok {
			portInfo[service_port_str] = portInfo[service_port_str] + connect_count
		} else {
			portInfo[service_port_str] = connect_count
		}
		_, ok2 := ipInfo[service_ip]
		if ok2 {
			ipInfo[service_ip] = ipInfo[service_ip] + connect_count
		} else {
			ipInfo[service_ip] = connect_count
		}
	}
	for k, v := range portInfo {
		sortMapL = append(sortMapL, sortMapS{k, v})
	}
	sort.Slice(sortMapL, func(i, j int) bool {
		return sortMapL[i].value > sortMapL[j].value // 降序
	})
	for i := 0; i < 6 && i < len(sortMapL); i++ {
		rankString = rankString + "#" + sortMapL[i].key + ":" + strconv.Itoa(sortMapL[i].value)
	}
	rankString = rankString + "$"
	sortMapL = sortMapL[0:0]
	for k, v := range ipInfo {
		sortMapL = append(sortMapL, sortMapS{k, v})
	}
	sort.Slice(sortMapL, func(i, j int) bool {
		return sortMapL[i].value > sortMapL[j].value // 降序
	})
	for i := 0; i < 6 && i < len(sortMapL); i++ {
		rankString = rankString + "#" + sortMapL[i].key + ":" + strconv.Itoa(sortMapL[i].value)
	}
	sortMapL = sortMapL[0:0]
	return rankString
}

//按端口查询
func QueryByPort(port string) string {
	rows, err := DB.Query("select * from netstat where service_port=?", port)
	util.CheckErr(err)
	graph = gographviz.NewGraph()
	graph.SetName("G")
	uid := 0
	for rows.Next() {
		err = rows.Scan(&service_ip, &service_port, &client_ip, &client_port, &connect_count, &protocol)
		util.CheckErr(err)
		PraseQueryPort()
		uid = uid + 1
	}
	if uid == 0 {
		return "null"
	}
	output := graph.String()
	output_arr := strings.Split(output, "\n")
	new_output_arr := util.RemoveRepeatedElement(output_arr)
	new_output := strings.Join(new_output_arr, "\n")
	return new_output
}

//按服务端ip查询
func QueryBySIP(ip string) string {
	rows, err := DB.Query("select * from netstat where service_ip=?", ip)
	util.CheckErr(err)
	graph = gographviz.NewGraph()
	graph.SetName("G")
	uid := 0
	for rows.Next() {
		err = rows.Scan(&service_ip, &service_port, &client_ip, &client_port, &connect_count, &protocol)
		util.CheckErr(err)
		PraseQueryIP()
		uid = uid + 1
	}
	if uid == 0 {
		return "null"
	}
	output := graph.String()
	output_arr := strings.Split(output, "\n")
	new_output_arr := util.RemoveRepeatedElement(output_arr)
	new_output := strings.Join(new_output_arr, "\n")
	return new_output
}

//按客户端IP查询
func QueryByCIP(ip string) string {
	rows, err := DB.Query("select * from netstat where client_ip=?", ip)
	util.CheckErr(err)
	graph = gographviz.NewGraph()
	graph.SetName("G")
	uid := 0
	for rows.Next() {
		err = rows.Scan(&service_ip, &service_port, &client_ip, &client_port, &connect_count, &protocol)
		util.CheckErr(err)
		PraseQueryIP()
		uid = uid + 1
	}
	if uid == 0 {
		return "null"
	}
	output := graph.String()
	output_arr := strings.Split(output, "\n")
	new_output_arr := util.RemoveRepeatedElement(output_arr)
	new_output := strings.Join(new_output_arr, "\n")
	return new_output
}

//按ip查询
func QueryByIP(ip string) string {
	rows, err := DB.Query("select * from netstat where service_ip=? or client_ip=?", ip, ip)
	util.CheckErr(err)
	graph = gographviz.NewGraph()
	graph.SetName("G")
	uid := 0
	for rows.Next() {
		err = rows.Scan(&service_ip, &service_port, &client_ip, &client_port, &connect_count, &protocol)
		util.CheckErr(err)
		PraseQueryIP()
		uid = uid + 1
	}
	if uid == 0 {
		return "null"
	}
	output := graph.String()
	output_arr := strings.Split(output, "\n")
	new_output_arr := util.RemoveRepeatedElement(output_arr)
	new_output := strings.Join(new_output_arr, "\n")
	return new_output
}

//处理port查询结果，生成dot语言
func PraseQueryPort() {
	if client_port != 0 {
		service_node := "s" + service_ip
		service_node_style := make(map[string]string)
		service_node_style["color"] = "#D3D3D3" //LightGray
		service_node_style["label"] = service_ip

		s_port_node := service_ip + strconv.Itoa(service_port)
		s_port_node_style := make(map[string]string)
		s_port_node_style["color"] = "#D3D3D3" //LightGray
		s_port_node_style["shape"] = "circle"
		s_port_node_style["fontsize"] = "10"
		s_port_node_style["label"] = `"` + strconv.Itoa(service_port) + `"`

		c_port_node := client_ip + strconv.Itoa(client_port)
		c_port_node_style := make(map[string]string)
		c_port_node_style["color"] = "#D3D3D3"
		c_port_node_style["shape"] = "circle"
		c_port_node_style["fontsize"] = "10"
		c_port_node_style["label"] = `"` + strconv.Itoa(client_port) + `"`

		client_node := "c" + client_ip
		client_node_style := make(map[string]string)
		client_node_style["color"] = "#D3D3D3"
		client_node_style["label"] = client_ip

		line_to_service_style := make(map[string]string)
		line_to_service_style["label"] = strconv.Itoa(service_port)
		line_to_service_style["dir"] = "none"
		line_to_service_style["style"] = "dashed"
		line_to_service_style["color"] = "#D3D3D3"

		line_from_client_style := make(map[string]string)
		line_from_client_style["label"] = strconv.Itoa(client_port)
		line_from_client_style["dir"] = "none"
		line_from_client_style["style"] = "dashed"

		graph.AddNode("G", service_node, service_node_style)
		graph.AddNode("G", s_port_node, s_port_node_style)
		graph.AddEdge(s_port_node, service_node, true, line_to_service_style)
		graph.AddNode("G", client_node, client_node_style)
		graph.AddNode("G", c_port_node, c_port_node_style)
		graph.AddEdge(c_port_node, client_node, true, line_from_client_style)
		graph.AddEdge(c_port_node, s_port_node, true, line_to_service_style)
	} else {
		service_node := "s" + service_ip
		service_node_style := make(map[string]string)
		service_node_style["color"] = "#CDCD00"
		service_node_style["label"] = service_ip

		s_port_node := service_ip + strconv.Itoa(service_port)
		s_port_node_style := make(map[string]string)
		s_port_node_style["color"] = "#D3D3D3" //LightGray
		s_port_node_style["shape"] = "circle"
		s_port_node_style["fontsize"] = "10"
		s_port_node_style["label"] = `"` + strconv.Itoa(service_port) + `"`

		client_node := "c" + client_ip
		client_node_style := make(map[string]string)
		client_node_style["color"] = "#00CD00"
		client_node_style["label"] = client_ip

		line_to_service_style := make(map[string]string)
		line_to_service_style["label"] = protocol
		line_to_service_style["color"] = "#D3D3D3"

		line_from_client_style := make(map[string]string)
		line_from_client_style["label"] = `"` + strconv.Itoa(connect_count) + `"`
		line_from_client_style["color"] = "#D3D3D3"

		graph.AddNode("G", service_node, service_node_style)
		graph.AddNode("G", s_port_node, s_port_node_style)
		graph.AddEdge(s_port_node, service_node, true, line_to_service_style)
		graph.AddNode("G", client_node, client_node_style)
		graph.AddEdge(client_node, s_port_node, true, line_from_client_style)
	}

}

//处理ip查询结果，生成dot语言
func PraseQueryIP() {
	if client_port != 0 {
		service_node := "s1" + service_ip
		service_node_style := make(map[string]string)
		service_node_style["color"] = "#D3D3D3" //LightGray
		service_node_style["label"] = service_ip

		s_port_node := service_ip + strconv.Itoa(service_port)
		s_port_node_style := make(map[string]string)
		s_port_node_style["color"] = "#D3D3D3" //LightGray
		s_port_node_style["shape"] = "circle"
		s_port_node_style["fontsize"] = "10"
		s_port_node_style["label"] = `"` + strconv.Itoa(service_port) + `"`

		c_port_node := client_ip + strconv.Itoa(client_port)
		c_port_node_style := make(map[string]string)
		c_port_node_style["color"] = "#D3D3D3"
		c_port_node_style["shape"] = "circle"
		c_port_node_style["fontsize"] = "10"
		c_port_node_style["label"] = `"` + strconv.Itoa(client_port) + `"`

		client_node := "c1" + client_ip
		client_node_style := make(map[string]string)
		client_node_style["color"] = "#D3D3D3"
		client_node_style["label"] = client_ip

		line_to_service_style := make(map[string]string)
		line_to_service_style["label"] = strconv.Itoa(service_port)
		line_to_service_style["dir"] = "none"
		line_to_service_style["style"] = "dashed"
		line_to_service_style["color"] = "#D3D3D3"

		line_from_client_style := make(map[string]string)
		line_from_client_style["label"] = strconv.Itoa(client_port)
		line_from_client_style["dir"] = "none"
		line_from_client_style["style"] = "dashed"

		graph.AddNode("G", service_node, service_node_style)
		graph.AddNode("G", s_port_node, s_port_node_style)
		graph.AddEdge(s_port_node, service_node, true, line_to_service_style)
		graph.AddNode("G", client_node, client_node_style)
		graph.AddNode("G", c_port_node, c_port_node_style)
		graph.AddEdge(c_port_node, client_node, true, line_from_client_style)
		graph.AddEdge(c_port_node, s_port_node, true, line_to_service_style)
	} else {
		service_node := "s2" + service_ip
		service_node_style := make(map[string]string)
		service_node_style["color"] = "#CDCD00"
		service_node_style["label"] = service_ip

		s_port_node := service_ip + strconv.Itoa(service_port)
		s_port_node_style := make(map[string]string)
		s_port_node_style["color"] = "#D3D3D3" //LightGray
		s_port_node_style["shape"] = "circle"
		s_port_node_style["fontsize"] = "10"
		s_port_node_style["label"] = `"` + strconv.Itoa(service_port) + `"`

		client_node := "c2" + client_ip
		client_node_style := make(map[string]string)
		client_node_style["color"] = "#00CD00"
		client_node_style["label"] = client_ip

		line_to_service_style := make(map[string]string)
		line_to_service_style["label"] = protocol
		line_to_service_style["color"] = "#D3D3D3"

		line_from_client_style := make(map[string]string)
		line_from_client_style["label"] = `"` + strconv.Itoa(connect_count) + `"`
		line_from_client_style["color"] = "#D3D3D3"

		graph.AddNode("G", service_node, service_node_style)
		graph.AddNode("G", s_port_node, s_port_node_style)
		graph.AddEdge(s_port_node, service_node, true, line_to_service_style)
		graph.AddNode("G", client_node, client_node_style)
		graph.AddEdge(client_node, s_port_node, true, line_from_client_style)
	}
}
