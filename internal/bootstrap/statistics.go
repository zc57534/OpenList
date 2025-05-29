package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type MachineInfo struct {
	SystemName    string
	Version       string
	KernelVersion string
	Architecture  string
	CpuCount      int
	TotalMemory   uint64
}

const urlSystemInfo = "https://dapi.alistgo.com/v0/system/info"

func DealMachineInfoStatistics() {
	//获取系统名称、系统版本
	osNsme, version := getOsNameAndVersion()
	//获取内核版本
	kernelVersion := getKernelVersion()
	//获取系统架构
	arch := getArch()
	//获取cpu核数
	cpuCores := getCpuCores()
	//总内存大小
	totalMemory := getTotalMemory()
	machineInfo := MachineInfo{
		SystemName:    osNsme,
		Version:       version,
		KernelVersion: kernelVersion,
		Architecture:  arch,
		CpuCount:      cpuCores,
		TotalMemory:   totalMemory,
	}
	fmt.Println(machineInfo)
	//上报数据
	res, err := request("POST", urlSystemInfo, machineInfo)
	if err != nil {
		fmt.Printf("上报系统信息出错，err: %v, res: %v", err, res)
	}
	fmt.Printf("上报系统信息返回的结果是：%v", res)
}

// 获取系统名称
func getOsNameAndVersion() (osName string, version string) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		fmt.Println("无法读取系统信息:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		//提取 ID 字段（如 "ubuntu"、"centos"）
		if strings.HasPrefix(line, "ID=") {
			osName = strings.Trim(line[3:], "\"")
			break
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(line[11:], "\"")
		}
	}
	fmt.Printf("系统名称: %s, 系统版本：%s\n", osName, version)
	return
}

func executeCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// 获取内核版本
func getKernelVersion() string {
	kv, err := executeCommand("uname", "-r")
	if err != nil {
		return ""
	}
	return kv
}

// 获取系统架构
func getArch() string {
	arch, err := executeCommand("uname", "-m")
	if err != nil {
		return ""
	}
	return arch
}

// 获取cpu核数
func getCpuCores() int {
	cpuCoresStr, err := executeCommand("nproc")
	if err != nil {
		return 0
	}
	cpuCores, err := strconv.Atoi(cpuCoresStr)
	if err != nil {
		fmt.Printf("解析CPU核数失败， err: %v\n", err)
		return 0
	}
	return cpuCores
}

// 获取总内存大小
func getTotalMemory() (mem uint64) {
	output, err := executeCommand("free", "-b")
	if err != nil {
		return mem
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				mem, err = strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					fmt.Printf("<UNK>Mem<UNK> err: %v\n", err)
				}
				return mem
			}
		}
	}
	fmt.Printf("无法解析 free 输出, err: %v\n", err)
	return 0
}

// 发送请求
func request(method string, url string, params interface{}) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// 构建请求的 JSON 数据
	jsonData, _ := json.Marshal(params)
	// 创建请求对象
	req, _ := http.NewRequest(
		method, // POST/GET
		url,
		bytes.NewBuffer(jsonData),
	)
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 处理非 200 状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败，状态码:", resp.StatusCode)
		return nil, nil
	}

	// 解析 JSON 响应
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	fmt.Println("响应结果:", result)
	return result, nil
}
