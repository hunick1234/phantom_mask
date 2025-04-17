## 藥局口罩系統 API 規格說明

### 1. 查詢營業中藥局
**GET** `/api/pharmacies/open`

**描述**：列出在指定時間與星期幾營業的藥局。

**參數**：
- `time`（string，必填）：格式 `HH:MM`
- `day_of_week`（string，必填）：如 `Monday`、`Sunday`

**輸入範例**：
```
GET /api/pharmacies/open?time=09:00&day_of_week=Monday
```

**回傳格式**：
```json
[
  {
    "id": "ph001",
    "name": "Health Pharmacy",
    "address": "123 Main St",
    "open_hours": {
      "Monday": ["08:00", "18:00"]
    }
  }
]
```

---

### 2. 查詢藥局販售的口罩（排序）
**GET** `/api/pharmacies/{pharmacy_id}/masks`

**描述**：列出指定藥局販售的所有口罩，並依名稱或價格排序。

**參數**：
- `pharmacy_id`（path 參數）
- `sort_by`（query 參數，可選）：`name` 或 `price` 
    - 預設為 `name`

**輸入範例**：
```
GET /api/pharmacies/ph001/masks?sort_by=price
```

**回傳格式**：
```json
[
  {
    "mask_id": "m001",
    "name": "Blue Surgical Mask",
    "price": 5.0,
    "stock": 100
  }
]
```

---

### 3. 查詢口罩數量門檻的藥局（價格範圍）
**GET** `/api/pharmacies/filter_by_mask_count`

**描述**：篩選在指定價格範圍中，販售口罩數量多於/少於 x 件的藥局。

**參數**：
- `min_price`（float）
- `max_price`（float）
- `comparison`（string）：`more` 或 `less`
- `count`（int）

**輸入範例**：
```
GET /api/pharmacies/filter_by_mask_count?min_price=2.0&max_price=10.0&comparison=more&count=50
```

**回傳格式**：
```json
[
  {
    "id": "ph003",
    "name": "GreenCare Pharmacy",
    "mask_count": 12
  }
]
```

---

### 4. 交易金額最高前 x 名使用者
**GET** `/api/users/top_transactions`

**描述**：查詢指定期間內購買金額最高的前 x 名使用者。

**參數**：
- `start_date`（YYYY-MM-DD）
- `end_date`（YYYY-MM-DD）
- `top`（int）

**輸入範例**：
```
GET /api/users/top_transactions?start_date=2025-01-01&end_date=2025-01-31&top=5
```

**回傳格式**：
```json
[
  {
    "user_id": "u001",
    "name": "Alice",
    "total_amount": 150.0
  }
]
```

---

### 5. 查詢交易總數與金額
**GET** `/api/transactions/summary`

**描述**：查詢某期間內交易總口罩數量與金額。

**參數**：
- `start_date`（YYYY-MM-DD）
- `end_date`（YYYY-MM-DD）

**輸入範例**：
```
GET /api/transactions/summary?start_date=2025-01-01&end_date=2025-01-31
```

**回傳格式**：
```json
{
  "total_masks": 1000,
  "total_amount": 5000.0
}
```

---

### 6. 關鍵字搜尋藥局或口罩
**GET** `/api/search`

**描述**：以關鍵字搜尋藥局或口罩，並依關聯度排序。

**參數**：
- `q`（string，必填）
- `type`（string，可選）：`pharmacy` 或 `mask`

**輸入範例**：
```
GET /api/search?q=mask&type=pharmacy
```

**回傳格式**：
```json
[
  {
    "type": "mask",
    "id": "m010",
    "name": "Kids Mask",
    "relevance": 0.92
  },
  {
    "type": "pharmacy",
    "id": "ph008",
    "name": "Smile Pharmacy",
    "relevance": 0.88
  }
]
```

---

### 7. 處理購買交易（原子性）
**POST** `/api/purchase`

**描述**：使用者從藥局購買口罩，並完成扣庫存、交易紀錄等原子性資料更新。

**請求格式**：
```json
{
  "user_id": "u001",
  "pharmacy_id": "ph001",
  "mask_id": "m001",
  "quantity": 5
}
```

**輸入範例**：
```
POST /api/purchase
{
  "user_id": "u001",
  "pharmacy_id": "ph001",
  "mask_id": "m001",
  "quantity": 5
}
```

**回傳格式**：
```json
{
  "transaction_id": "t123",
  "status": "success",
  "remaining_stock": 45
}
```