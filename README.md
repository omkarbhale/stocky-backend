# Stocky Backend

## Database schema

### **User and Rewards**

| Table       | Fields                                                                            | Relationships                      |
| ----------- | --------------------------------------------------------------------------------- | ---------------------------------- |
| **users**   | `id`, `name`, `created_at`, `updated_at`                                          | —                                  |
| **rewards** | `id`, `user_id`, `symbol_id`, `quantity`, `timestamp`, `created_at`, `updated_at` | belongs to **User** and **Symbol** |

- One **user** → many **rewards**
- Each **reward** references a **symbol**

---

### **Symbols and Price History**

| Table                      | Fields                                                                      | Relationships         |
| -------------------------- | --------------------------------------------------------------------------- | --------------------- |
| **symbols**                | `id`, `name`, `created_at`, `updated_at`                                    | —                     |
| **symbol_price_histories** | `id`, `symbol_id`, `price`, `time_hour`, `date`, `created_at`, `updated_at` | belongs to **Symbol** |

- One **symbol** → many **price history** records
- Price history stores **hourly prices** per day (0–23 hours)

---

### **Ledger (Company Accounting)**

| Table            | Fields                                                                                              | Relationships                              |
| ---------------- | --------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| **accounts**     | `id`, `name`, `type`, `description`, `created_at`, `updated_at`                                     | —                                          |
| **transactions** | `id`, `description`, `created_at`, `updated_at`                                                     | —                                          |
| **entries**      | `id`, `transaction_id`, `account_id`, `type` (`debit/credit`), `amount`, `created_at`, `updated_at` | belongs to **Transaction** and **Account** |

- One **transaction** → many **entries**
- Debits and credits are balanced per transaction
- Each **entry** references a specific **account**

---

### **Entity Relationships**

```
User ───< Reward >─── Symbol
│
└──< SymbolPriceHistory

Account ───< Entry >─── Transaction
```

---

## Reward Flow

When user **Omkar** is rewarded **10 INFY shares**:

1. The company buys shares worth ₹10,000 at the latest market price.
2. 4% transaction fee = ₹400.
3. A double-entry ledger transaction is recorded:

| Account          | Type   | Amount  |
| ---------------- | ------ | ------- |
| StockInvestments | Debit  | ₹10,000 |
| TransactionFees  | Debit  | ₹400    |
| Cash             | Credit | ₹10,400 |

4. User sees only the **reward** in their portfolio; fees stay internal.

---

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL

## Notes

- Transaction fee is fixed at 4%.

## API Specification (YAML)

```
openapi: 3.0.0
info:
  title: Stocky
  version: 1.0.0
  description: ""
servers:
  - url: http://localhost:8080
paths:
  /historical-inr/1:
    get:
      summary: Get valuation till yesterday
      tags:
        - Portfolio
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
  /stats/1:
    get:
      summary: Stats today for user
      tags:
        - Portfolio
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
  /portfolio/1:
    get:
      summary: "Bonus: Current Holdings per symbol"
      tags:
        - Portfolio
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
  /reward/:
    post:
      summary: Create reward
      tags:
        - Rewards
      responses:
        "201":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                userId:
                  type: integer
                symbolId:
                  type: integer
                quantity:
                  type: integer
                time:
                  type: string
                  format: date-time
  /today-stocks/1:
    get:
      summary: All rewards today for user
      tags:
        - Rewards
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
  /symbol/:
    get:
      summary: Get all symbols (just names and ID)
      tags:
        - Symbols
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
  /user/:
    post:
      summary: Create user
      tags:
        - Users
      responses:
        "201":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
    get:
      summary: Get all users
      tags:
        - Users
      responses:
        "200":
          description: Example
          headers:
            Date:
              schema:
                type: string
            Content-Length:
              schema:
                type: integer

```
