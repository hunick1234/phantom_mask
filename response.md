# 回應文件（Response）

## A. 必要資訊

### A.1. 任務完成度

- [x] 查詢指定時間與星期幾營業的藥局  
       ➤ API：`GET /api/pharmacies/open`
- [x] 查詢指定藥局販售的口罩，並依名稱或價格排序  
       ➤ API：`GET /api/pharmacies/{pharmacy_id}/masks`
- [x] 查詢價格範圍內，口罩數量多於或少於指定數量的藥局  
       ➤ API：`GET /api/pharmacies/filter_by_mask_count`
- [x] 查詢指定時間區間內購買金額最高的前 x 名使用者  
       ➤ API：`GET /api/users/top_transactions`
- [x] 查詢指定期間內總交易金額與口罩數量  
       ➤ API：`GET /api/transactions/summary`
- [x] 以關鍵字搜尋藥局或口罩，依關聯度排序  
       ➤ API：`GET /api/search`
- [x] 使用者購買口罩  
       ➤ API：`POST /api/users/me/purchase`

### A.2. API 文件

- 位置：`docs/api.md`
- postman collection：`docs/Pharmacy.postman_collection.json`

### A.3. 匯入資料指令

```go
$ go build -o seeder cmd/seeder/main.go
```
請依下列順序匯入初始資料：
```bash
$ ./seeder -t pharmacy -p data/pharmacies.json
$ ./seeder -t user -p data/users.json
```

---

## B. 加分項目
### B.1 Docker

- 提供 `Dockerfile` 與 `docker-compose.yml`
- 使用 `.env.example` 自動讀取環境變數
- 啟動專案與資料庫一鍵完成：

```bash
$ cp .env.example .env
$ docker compose up --build 

```

資料將自動透過 `init.sh` 導入。

---
