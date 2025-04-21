# Phantom Mask API 文件

---

## 1. 查詢營業中藥局

**GET** `/api/pharmacies/open`

| 參數          | 類型   | 是否必填 | 說明                                    |
| ------------- | ------ | -------- | --------------------------------------- |
| `time`        | string | ✅       | 時間，格式為 `HH:MM`                    |
| `day_of_week` | string | ✅       | 星期幾（Mon,Tue,Wed,Thur,Fru,,Sat,Sun） |

📦 範例：

```
GET /api/pharmacies/open?time=10:00&day_of_week=Mon
```

- response:

```json
[
  {
    "id": "1",
    "name": "DFW Wellness",
    "address": "123 Main St",
    "opening_hours": {
      "Mon": ["08:00", "12:00"],
      "Tue": ["14:00", "18:00"]
    }
  }
]
```

---

## 2. 查詢藥局販售的口罩

**GET** `/api/pharmacies/{pharmacy_id}/masks`

| 參數      | 類型   | 是否必填 | 說明              |
| --------- | ------ | -------- | ----------------- |
| `sort_by` | string | ✅       | `name` 或 `price` |

📦 範例：

```
GET /api/pharmacies/1/masks?sort_by=price
```

- response:

```json
[
  {
    "id": "10",
    "name": "MaskT (green) (10 per pack)",
    "price": 41.86,
    "stock": 50
  },
  {
    "id": "11",
    "name": "Second Smile (black) (3 per pack)",
    "price": 5.84,
    "stock": 80
  }
]
```

---

## 3. 查詢符合條件的藥局（依口罩數量）

**GET** `/api/pharmacies/filter_by_mask_count`

| 參數         | 類型   | 是否必填 | 說明             |
| ------------ | ------ | -------- | ---------------- |
| `min_price`  | float  | ✅       | 價格下限         |
| `max_price`  | float  | ✅       | 價格上限         |
| `comparison` | string | ✅       | `more` 或 `less` |
| `count`      | int    | ✅       | 數量門檻         |

📦 範例：

```
GET /api/pharmacies/filter_by_mask_count?min_price=5&max_price=30&comparison=more&count=3
```

- response:

```json
[
  {
    "id": "3",
    "name": "GreenCare Pharmacy",
    "mask_count": 5
  }
]
```

---

## 4. 查詢交易金額前 N 名使用者

**GET** `/api/users/top_transactions`

| 參數         | 類型   | 是否必填 | 說明                 |
| ------------ | ------ | -------- | -------------------- |
| `start_date` | string | ✅       | 起始日（YYYY-MM-DD） |
| `end_date`   | string | ✅       | 結束日               |
| `top`        | int    | ✅       | 顯示幾名使用者       |

📦 範例：

```
GET /api/users/top_transactions?start_date=2025-01-01&end_date=2025-01-31&top=3
```

- response:

```json
[
  [
    {
      "user_id": "2",
      "name": "Ada Larson",
      "total_amount": 143.12
    },
    {
      "user_id": "4",
      "name": "Lester Arnold",
      "total_amount": 136.84
    }
  ]
]
```

---

## 5. 查詢總交易口罩數與金額

**GET** `/api/transactions/summary`

| 參數         | 類型   | 是否必填 | 說明                  |
| ------------ | ------ | -------- | --------------------- |
| `start_date` | string | ✅       | 起始日（YYYY-MM-DD）  |
| `end_date`   | string | ✅       | 結束日 （YYYY-MM-DD） |

📦 範例：

```
GET /api/transactions/summary?start_date=2025-01-01&end_date=2025-01-31
```

- response:

```json
{
  "total_amount": 279.96,
  "total_masks": 10
}
```

---

## 6. 關鍵字搜尋藥局與口罩

**GET** `/api/search`

| 參數 | 類型   | 是否必填 | 說明       |
| ---- | ------ | -------- | ---------- |
| `q`  | string | ✅       | 搜尋關鍵字 |

📦 範例：

```
GET /api/search?q=blue
```

- response:

```json
[
  {
    "id": "1",
    "type": "mask",
    "name": "Second Smile (blue) (10 per pack)"
  },
  {
    "id": "5",
    "type": "pharmacy",
    "name": "blue Mask Warehouse"
  }
]
```

---

## 7. 購買口罩（需登入）

**POST** `/api/users/me/purchase`

🔒 需附上 Authorization: `Bearer <JWT>`

| 欄位          | 類型 | 是否必填 | 說明     |
| ------------- | ---- | -------- | -------- |
| `pharmacy_id` | uint | ✅       | 藥局 ID  |
| `mask_id`     | uint | ✅       | 口罩 ID  |
| `quantity`    | int  | ✅       | 購買數量 |

📦 範例：

```json
POST /api/users/me/purchase
Authorization: Bearer <token>

{
  "pharmacy_id": 2,
  "mask_id": 5,
  "quantity": 3
}
```

- response:

```json
{
  "transaction_id": 101,
  "status": "success"
}
```

---

## 8. 登入以取得 JWT（測試用）

**POST** `/api/users/login`

| 欄位      | 類型 | 是否必填 | 說明              |
| --------- | ---- | -------- | ----------------- |
| `user_id` | uint | ✅       | 測試用的使用者 ID |

📦 範例：

```
POST /api/users/login
{
  "user_id": 1
}
```

- response:

```json
{
  "token": "xxx.yyy.zzz"
}
```

---
