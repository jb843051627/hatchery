//go:build ignore
// 复现脚本: 仅供 go run ./tmp_repro_run.go 手动复现超时问题, 不参与 go build/test.

package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const base = "http://127.0.0.1:18080"
const dbPath = "D:/Develop/Workspace-trae/go-workspace/hatchery/hatchery-bug-005-test/tmp_repro.db"

type reading struct {
	IncubatorID int64   `json:"incubator_id"`
	SensorType  string  `json:"sensor_type"`
	Value       float64 `json:"value"`
	RecordedAt  string  `json:"recorded_at"`
}

func mkReadings(incID int64, n int) []reading {
	rs := make([]reading, n)
	t0 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := range rs {
		rs[i] = reading{incID, "temperature", float64(i), t0.Add(time.Duration(i) * time.Second).Format(time.RFC3339)}
	}
	return rs
}

func createIncubator() int64 {
	cr := struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
		Status   string `json:"status"`
	}{"inc-repro", "loc", 1000, "active"}
	b, _ := json.Marshal(cr)
	resp, err := http.Post(base+"/api/incubators", "application/json", bytes.NewReader(b))
	if err != nil {
		fmt.Println("FATAL create incubator:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var m map[string]int64
	json.NewDecoder(resp.Body).Decode(&m)
	return m["id"]
}

func dbCount() int {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(15000)")
	if err != nil {
		fmt.Println("open db:", err)
		return -1
	}
	defer db.Close()
	var c int
	if err := db.QueryRow("SELECT COUNT(*) FROM readings").Scan(&c); err != nil {
		fmt.Println("count:", err)
		return -1
	}
	return c
}

func main() {
	incID := createIncubator()
	fmt.Printf("[setup] created incubator id=%d\n", incID)

	// ---- Phase A: 测速 (无取消) ----
	N_A := 5000
	bodyA, _ := json.Marshal(mkReadings(incID, N_A))
	start := time.Now()
	resp, err := http.Post(base+"/api/readings/batch", "application/json", bytes.NewReader(bodyA))
	if err != nil {
		fmt.Println("FATAL phaseA:", err)
		return
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
	elapsedA := time.Since(start)
	perRow := elapsedA / time.Duration(N_A)
	fmt.Printf("[phase A] N=%d elapsed=%v (~%v/row incl. http overhead)\n", N_A, elapsedA, perRow)

	// ---- Phase B: 复现 —— 客户端中途断连(模拟 30s 超时),看服务端是否继续写完 ----
	targetDur := 3 * time.Second
	N_B := int(targetDur / perRow)
	if N_B < 2000 {
		N_B = 2000
	}
	if N_B > 200000 {
		N_B = 200000
	}
	bodyB, _ := json.Marshal(mkReadings(incID, N_B))
	fmt.Printf("[phase B] N=%d (target ~%v of server work), client will disconnect after 200ms\n", N_B, targetDur)

	c0 := dbCount()
	fmt.Printf("[phase B] DB rows BEFORE = %d\n", c0)

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, "POST", base+"/api/readings/batch", bytes.NewReader(bodyB))
	req.Header.Set("Content-Type", "application/json")
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
		fmt.Println(">>> [disconnect] client cancelled ctx at 200ms (simulating 30s timeout / connection drop)")
	}()
	start = time.Now()
	resp, err = http.DefaultClient.Do(req)
	fmt.Printf("[phase B] client Do returned: err=%v, elapsed=%v (client error is EXPECTED — it gave up)\n", err, time.Since(start))
	if resp != nil {
		io.ReadAll(resp.Body)
		resp.Body.Close()
	}

	// 给服务端足够时间把剩下的写完
	wait := targetDur + 6*time.Second
	fmt.Printf("[phase B] waiting %v to let the server keep writing after disconnect...\n", wait)
	time.Sleep(wait)

	c1 := dbCount()
	written := c1 - c0
	fmt.Printf("[phase B] DB rows AFTER  = %d  (delta = %d, expected if cancel ignored = %d)\n", c1, written, N_B)
	fmt.Println("---- RESULT ----")
	switch {
	case written >= N_B:
		fmt.Printf(">>> REPRO CONFIRMED: server wrote ALL %d rows AFTER the client disconnected at 200ms.\n", written)
		fmt.Println(">>> The request context (r.Context()) was cancelled, but BatchIngest kept inserting row-by-row until done.")
	case written > 0:
		fmt.Printf(">>> PARTIAL: server wrote %d/%d rows after disconnect.\n", written, N_B)
	default:
		fmt.Println(">>> no rows written after disconnect (cancel may have prevented entry into BatchIngest).")
	}
}
