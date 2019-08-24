package controllers

import (
	"bytes"
	"io"
	"net/http"
	"regexp"

	"NetTracer/models"
	"NetTracer/util"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseText(c *gin.Context) {

	//1.通过是否有相同的ip和端口来判断出本地地址和远程地址谁提供的服务
	for i := 0; ; i++ { //处理多个上传文件
		loadfile := "file" + strconv.Itoa(i)
		file, _, err := c.Request.FormFile(loadfile)
		if err != nil {
			break
		}
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, file)
		pending_content := strings.Split(buf.String(), "\n")

		s_map := make(map[string]map[string]int)                          //存放服务端协议+IP地址+端口
		var tmp_slice []string                                            //临时存放无法分辨出提供服务的连接
		var s_slice []string                                              //存放无法分辨出提供服务的连接
		re := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}:\d+\b`) //正则匹配ip
		tmp_protocol := ""                                                //临时存放当前连接协议
		t_c := 0

		for _, line := range pending_content {
			if strings.Contains(line, "0.0.0.0") || !strings.Contains(line, "ESTABLISHED") { //跳过地址为0.0.0.0和未建立稳定连接的行
				continue
			}
			//  TCP    1.1.1.1:56891       2.2.2.2:3389      ESTABLISHED       InHost
			line_data := re.FindAllString(line, -1)
			//[1.1.1.1:56891,2.2.2.2:3389]
			if len(line_data) == 2 {
				t_c = t_c + 1
				if strings.Contains(line, "TCP") {
					tmp_protocol = "TCP"
				} else if strings.Contains(line, "UDP") {
					tmp_protocol = "UDP"
				}
				c0_map, ok0 := s_map[tmp_protocol+"#"+line_data[0]]
				c1_map, ok1 := s_map[tmp_protocol+"#"+line_data[1]]
				if ok0 { //本地地址是服务地址
					array_line_data := strings.Split(line_data[1], ":")
					ct, ot := c0_map[array_line_data[0]] //判断客户端IP是否相同
					if ot {
						ct = ct + 1
						c0_map[array_line_data[0]] = ct
					} else {
						c0_map[array_line_data[0]] = 1
					}
					continue
				} else if ok1 { //外部地址是服务地址
					array_line_data := strings.Split(line_data[0], ":")
					ct, ot := c1_map[array_line_data[0]] //判断客户端IP是否相同
					if ot {
						ct = ct + 1
						c1_map[array_line_data[0]] = ct
					} else {
						c1_map[array_line_data[0]] = 1
					}
					continue
				} else { //无法判断或者s_map为空的情况下
					s_length := len(tmp_slice)
					if s_length == 0 { //初始状态或为空时,添加到临时切片中
						tmp_slice = append(tmp_slice, tmp_protocol+"#"+line_data[0]+"#"+line_data[1])
					} else {
						for i := 0; i < s_length; i++ {
							if strings.Contains(tmp_slice[i], line_data[0]) { //临时切片中有相同的本地地址，则是服务地址
								array_line_data := strings.Split(line_data[1], ":")
								c_map := make(map[string]int) //存放客户端IP和发起连接数量
								c_map[array_line_data[0]] = 1
								s_map[tmp_protocol+"#"+line_data[0]] = c_map
								array_tmp_slice := strings.Split(tmp_slice[i], "#")
								array_client := strings.Split(array_tmp_slice[2], ":")
								client_ip := array_client[0]
								ct, ot := c_map[client_ip] //判断客户端IP是否相同
								if ot {
									ct = ct + 1
									c_map[client_ip] = ct
								} else {
									c_map[client_ip] = 1
								}
								tmp_slice = append(tmp_slice[:i], tmp_slice[i+1:]...) //删除
								break
							} else if strings.Contains(tmp_slice[i], line_data[1]) { //临时切片中有相同的外部地址，则是服务地址
								array_line_data := strings.Split(line_data[0], ":")
								c_map := make(map[string]int)
								c_map[array_line_data[0]] = 1
								s_map[tmp_protocol+"#"+line_data[1]] = c_map
								array_tmp_slice := strings.Split(tmp_slice[i], "#")
								array_client := strings.Split(array_tmp_slice[1], ":")
								client_ip := array_client[0]
								ct, ot := c_map[client_ip] //判断客户端IP是否相同
								if ot {
									ct = ct + 1
									c_map[client_ip] = ct
								} else {
									c_map[client_ip] = 1
								}
								tmp_slice = append(tmp_slice[:i], tmp_slice[i+1:]...) //删除
								break
							} else {
								if i+1 == s_length { //如果遍历完临时切片也没发现有相同的服务
									tmp_slice = append(tmp_slice, tmp_protocol+"#"+line_data[0]+"#"+line_data[1])
								}
							}
						}
					}
				}
			}
		}
		//2.通过常见的端口来判断提供服务IP(可自定义更改)
		well_known := []string{"17", "21", "22", "23", "25", "53", "69", "80", "81", "86", "110", "123", "135", "139", "143", "161", "389", "443", "445", "587", "636", "1311", "1433", "1434", "1720", "2301", "2381", "3306", "3389", "4443", "47001", "5060", "5061", "5432", "5500", "5900", "5901", "5985", "5986", "7080", "8080", "8081", "8082", "8089", "8000", "8180", "8443"}
		//加载额外的自定义服务端口
		ext_port := util.ReadPortConfig()
		well_known = append(well_known, ext_port...)

		well_known = util.RemoveRepeatedElement(well_known)
		fmt.Println(well_known)
		for _, v := range tmp_slice {
			tv := strings.Split(v, "#")
			p1 := strings.Split(tv[1], ":")[1]
			ip1 := strings.Split(tv[1], ":")[0]
			p2 := strings.Split(tv[2], ":")[1]
			ip2 := strings.Split(tv[2], ":")[0]
			for k, p := range well_known {
				if p == p1 {
					c_map := make(map[string]int)
					c_map[ip2] = 1
					s_map[tv[0]+"#"+tv[1]] = c_map
					break
				} else if p == p2 {
					c_map := make(map[string]int)
					c_map[ip1] = 1
					s_map[tv[0]+"#"+tv[2]] = c_map
					break
				} else {
					if (k + 1) == len(well_known) {
						s_slice = append(s_slice, v)
					}
				}
			}
		}
		//3.将s_map和s_slice中的结果存到数据库中
		for k, v := range s_map {
			k_arr1 := strings.Split(k, "#")
			k_arr2 := strings.Split(k_arr1[1], ":")
			s_ip := k_arr2[0]
			s_port := k_arr2[1]
			s_port_i, err := strconv.Atoi(s_port)
			util.CheckErr(err)
			protocol := k_arr1[0]
			for k2, v2 := range v {
				c_ip := k2
				c_port := 0
				c_count := v2
				err = models.InsertNetstat(s_ip, s_port_i, c_ip, c_port, c_count, protocol)
				util.CheckErr(err)
			}
		}
		for _, v := range s_slice {
			v_arr := strings.Split(v, "#")
			s_ip := strings.Split(v_arr[1], ":")[0]
			s_port := strings.Split(v_arr[1], ":")[1]
			s_port_i, err := strconv.Atoi(s_port)
			util.CheckErr(err)
			c_ip := strings.Split(v_arr[2], ":")[0]
			c_port := strings.Split(v_arr[2], ":")[1]
			c_port_i, err := strconv.Atoi(c_port)
			util.CheckErr(err)
			c_count := 1
			protocol := v_arr[0]
			err = models.InsertNetstat(s_ip, s_port_i, c_ip, c_port_i, c_count, protocol)
			util.CheckErr(err)
		}
	}
	c.String(http.StatusOK, "SUCCESS")
}
