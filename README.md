## Stocky Backend

### Specification

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
