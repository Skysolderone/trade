package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"trade/db"
	"trade/model"
)

// Strategy1Response 策略一API响应
type Strategy1Response struct {
	Results []Strategy1ResultWithDetails `json:"results"`
}

// Strategy1ResultWithDetails 策略一结果（包含详细记录）
type Strategy1ResultWithDetails struct {
	model.Strategy1Result
	Details []model.Strategy1DetailRecord `json:"details"`
}

// Strategy2Response 策略二API响应
type Strategy2Response struct {
	Results []Strategy2ResultWithDetails `json:"results"`
}

// Strategy2ResultWithDetails 策略二结果（包含详细记录）
type Strategy2ResultWithDetails struct {
	model.Strategy2Result
	Details []model.Strategy2DetailRecord `json:"details"`
}

// StartServer 启动Web服务器
func StartServer(port int) {
	// 注册路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/strategy1", getStrategy1Results)
	http.HandleFunc("/api/strategy2", getStrategy2Results)

	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// indexHandler 首页处理
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

// getStrategy1Results 获取策略一结果
func getStrategy1Results(w http.ResponseWriter, r *http.Request) {
	// 设置CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 获取查询参数
	symbol := r.URL.Query().Get("symbol")
	interval := r.URL.Query().Get("interval")

	// 构建查询
	query := db.Pog.Model(&model.Strategy1Result{})
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}

	// 查询结果
	var results []model.Strategy1Result
	if err := query.Order("created_at DESC").Find(&results).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取详细记录
	var response Strategy1Response
	for _, result := range results {
		var details []model.Strategy1DetailRecord
		db.Pog.Where("result_id = ?", result.ID).Order("year ASC").Find(&details)

		response.Results = append(response.Results, Strategy1ResultWithDetails{
			Strategy1Result: result,
			Details:         details,
		})
	}

	json.NewEncoder(w).Encode(response)
}

// getStrategy2Results 获取策略二结果
func getStrategy2Results(w http.ResponseWriter, r *http.Request) {
	// 设置CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 获取查询参数
	symbol := r.URL.Query().Get("symbol")
	interval := r.URL.Query().Get("interval")
	hourStr := r.URL.Query().Get("hour")

	// 构建查询
	query := db.Pog.Model(&model.Strategy2Result{})
	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	if interval != "" {
		query = query.Where("interval = ?", interval)
	}
	if hourStr != "" {
		hour, err := strconv.Atoi(hourStr)
		if err == nil {
			query = query.Where("hour = ?", hour)
		}
	}

	// 查询结果
	var results []model.Strategy2Result
	if err := query.Order("hour ASC").Find(&results).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取详细记录
	var response Strategy2Response
	for _, result := range results {
		var details []model.Strategy2DetailRecord
		db.Pog.Where("result_id = ?", result.ID).Order("date ASC").Find(&details)

		response.Results = append(response.Results, Strategy2ResultWithDetails{
			Strategy2Result: result,
			Details:         details,
		})
	}

	json.NewEncoder(w).Encode(response)
}
