```markdown
🎟️ Voucher API

A RESTful API to manage **brands**, **vouchers**, **customers**, and **redemption transactions** — built with **Go** and **PostgreSQL**. This API empowers you to create, fetch, and redeem vouchers efficiently while tracking customer points and transactions.

---

📚 Table of Contents

- [⚙️ Prerequisites](#️-prerequisites)
- [📥 Installation](#installation)
- [🔧 Configuration](#configuration)
- [🗃️ Database Setup](#database-setup)
- [🚀 Running the Application](#running-the-application)
- [🧪 Testing](#testing)
- [📡 API Endpoints](#api-endpoints)
- [📁 Project Structure](#project-structure)
- [🤝 Contributing](#contributing)
- [📄 License](#license)

---

⚙️ Prerequisites

Ensure the following are installed on your machine:

- [Go](https://golang.org/dl/) `v1.21+`
- [PostgreSQL](https://www.postgresql.org/) `v12+`
- Git
- Terminal (PowerShell or Bash)

---

📥 Installation

1. **Clone the repository**

```powershell
git clone https://github.com/yourusername/voucher-api.git
cd voucher-api
```

2. **Install dependencies**

```powershell
go mod tidy
```

---

## 🔧 Configuration

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgres://postgres:123@localhost:5432/voucher?sslmode=disable
PORT=8080
```

✔️ Ensure your PostgreSQL server is running and accessible based on the `DATABASE_URL`.

---

## 🗃️ Database Setup

1. **Create the Database**

```sql
CREATE DATABASE voucher;
```

2. **Run Migrations**

```sql
\c voucher
-- Then execute table creation queries listed in the README
```

3. **(Optional) Seed Initial Data**

```sql
INSERT INTO brands (name, description) VALUES ('Indomaret', 'Retail') ON CONFLICT DO NOTHING;
INSERT INTO vouchers (...) VALUES (...) ON CONFLICT DO NOTHING;
```

---

## 🚀 Running the Application

Start the API server:

```powershell
go run cmd/main.go
```

API will be available at: [http://localhost:8080](http://localhost:8080)

---

## 🧪 Testing

**Run Unit Tests:**

```powershell
go test ./test/... -v
```

**Check Code Coverage:**

```powershell
go test ./test/... -covermode=count -coverpkg=./internal/handler,./internal/service -coverprofile=cover.out
go tool cover -func=cover.out
go tool cover -html=cover.out -o cover.html
```

🧩 Open `cover.html` in your browser to visualize test coverage.

---

## 📡 API Endpoints

All endpoints accept and return `JSON`.

| Method | Endpoint                            | Description                         | Status Codes           |
|--------|-------------------------------------|-------------------------------------|-------------------------|
| POST   | `/brand`                            | Create a new brand                  | 201, 400, 500           |
| POST   | `/voucher`                          | Create a new voucher                | 201, 400, 500           |
| GET    | `/voucher?id={id}`                 | Get a voucher by ID                 | 200, 400, 404, 500      |
| GET    | `/voucher/brand?id={brandId}`      | Get all vouchers for a brand        | 200, 400, 500           |
| POST   | `/customer`                         | Create a new customer               | 201, 400, 500           |
| POST   | `/transaction/redemption`          | Redeem voucher(s)                   | 201, 400, 404, 500      |
| GET    | `/transaction/redemption?transactionId={id}` | Get transaction by ID   | 200, 400, 404, 500      |

---

## 🧾 Example Requests

**Create a Brand:**

```powershell
Invoke-WebRequest -Uri http://localhost:8080/brand -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name":"Indomaret","description":"Retail"}'
```

**Create a Redemption:**

```powershell
Invoke-WebRequest -Uri http://localhost:8080/transaction/redemption -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"customer_id":1,"vouchers":[{"voucher_id":1,"quantity":2}]}'
```

---

## 🚨 Error Responses

| Code | Description                      | Example JSON Response                            |
|------|----------------------------------|--------------------------------------------------|
| 400  | Bad request / validation error   | `{"error":"brand name is required"}`             |
| 404  | Resource not found               | `{"error":"customer not found"}`                 |
| 500  | Internal server/database error   | `{"error":"Failed to create brand: database error"}` |

---

## 📁 Project Structure

```
voucher-api/
├── cmd/                  # Entry point (main.go)
├── internal/
│   ├── config/           # Load .env and settings
│   ├── database/         # DB connection
│   ├── handler/          # API handlers
│   ├── model/            # Struct definitions
│   ├── repository/       # DB operations
│   └── service/          # Business logic
├── test/                 # Unit tests
├── .env                  # Environment variables
├── go.mod / go.sum       # Go module files
└── README.md             # You are here 📖
```

---

## 🤝 Contributing

1. Fork this repo
2. Create a branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -m "Add feature"`
4. Push: `git push origin feature/my-feature`
5. Open a pull request 🔥

> Don't forget to write tests for new features!

---

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for full details.

