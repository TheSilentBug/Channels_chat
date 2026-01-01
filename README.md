## TCP Full-Duplex Peer-to-Peer Chat in Go (Using Channels)

---

## 🇬🇧 English Version

### 📌 Overview

This project demonstrates a **full-duplex TCP communication system** written in **Go**, where **two independent programs (peers)** can both **send and receive messages simultaneously**.

Each peer behaves as **both a TCP server and a TCP client**:

* Listens on a local port
* Tries to connect to the other peer
* Once connected, messages can flow **in both directions at the same time**

Concurrency and synchronization are handled using **Go channels and goroutines**.

---

### 🧠 Key Concepts Covered

* TCP networking with `net`
* Goroutines for concurrency
* Channels for message passing
* Full-duplex communication
* Graceful connection shutdown
* Race between `accept()` and `dial()`

---

### 🗂 Project Structure

```
peerA/
 ├── go.mod
 └── main.go

peerB/
 ├── go.mod
 └── main.go
```

Each peer is a **standalone Go application**.

---

### 🔁 How It Works

1. Each peer starts a **TCP listener**
2. At the same time, it repeatedly tries to **dial the other peer**
3. Whichever succeeds first (accept or dial) becomes the active connection
4. Four goroutines run in parallel:

   * Read from terminal (stdin)
   * Write messages to TCP
   * Read messages from TCP
   * Print received messages
5. Channels coordinate message flow and shutdown signals

---

### 🔀 Channel Architecture

| Channel                | Purpose                        |
| ---------------------- | ------------------------------ |
| `outgoing chan string` | Messages typed by the user     |
| `incoming chan string` | Messages received from TCP     |
| `done chan struct{}`   | Signals shutdown on disconnect |

---

### 🧵 Goroutines

| Goroutine     | Responsibility                     |
| ------------- | ---------------------------------- |
| `stdinReader` | Reads user input                   |
| `connWriter`  | Sends data over TCP                |
| `connReader`  | Receives data from TCP             |
| `main loop`   | Prints messages & handles shutdown |

---

### ▶️ How to Run

Open **two separate terminals**.

#### Terminal 1

```bash
cd peerA
go run .
```

#### Terminal 2

```bash
cd peerB
go run .
```

Now type messages in either terminal and press **Enter**.

---

### 📊 Communication Flow (Simplified)

```
User Input → outgoing channel → TCP write
TCP read → incoming channel → Console output
```

---

### ⚠️ Common Issues

* `connection refused` → peer is not listening yet
* Port conflict → ensure 8080 and 8081 are free
* IPv6 issues → use `127.0.0.1` instead of `localhost`

---

### 🚀 Possible Extensions

* Support multiple peers
* Add JSON message format
* Implement authentication
* Add message timestamps
* Use TLS for secure communication

---

### 📜 License

Free to use for learning, experimentation, and personal projects.

---

---

## 🇮🇷 نسخه فارسی

### 📌 معرفی پروژه

این پروژه یک نمونه‌ی **ارتباط TCP دوطرفه (Full-Duplex)** با زبان **Go** است که در آن **دو برنامه‌ی مستقل** می‌توانند **به‌صورت همزمان پیام ارسال و دریافت کنند**.

هر برنامه (Peer):

* هم **Server** است
* هم **Client**
* و پس از اتصال، ارتباط به‌صورت همزمان در هر دو جهت برقرار می‌شود

مدیریت همزمانی کاملاً با **goroutine** و **channel** انجام شده است.

---

### 🧠 مفاهیم آموزشی

* ارتباط TCP در Go
* برنامه‌نویسی همزمان (Concurrency)
* استفاده از Channel
* ارتباط دوطرفه (Full-Duplex)
* مدیریت قطع اتصال
* رقابت بین `Accept` و `Dial`

---

### 🗂 ساختار پروژه

```
peerA/
 ├── go.mod
 └── main.go

peerB/
 ├── go.mod
 └── main.go
```

هر کدام یک برنامه‌ی مستقل هستند.

---

### 🔁 منطق اجرا

1. هر Peer یک پورت را **Listen** می‌کند
2. همزمان تلاش می‌کند به Peer دیگر **وصل شود**
3. اولین اتصال موفق (چه Accept چه Dial) انتخاب می‌شود
4. چهار goroutine به‌صورت موازی اجرا می‌شوند:

   * خواندن ورودی کاربر
   * ارسال پیام روی TCP
   * دریافت پیام از TCP
   * نمایش پیام‌ها
5. ارتباط بین بخش‌ها با Channel انجام می‌شود

---

### 🔀 کانال‌ها

| نام کانال  | کاربرد                       |
| ---------- | ---------------------------- |
| `outgoing` | پیام‌های تایپ‌شده توسط کاربر |
| `incoming` | پیام‌های دریافتی از شبکه     |
| `done`     | اعلام قطع اتصال و خروج امن   |

---

### 🧵 Goroutineها

| Goroutine    | وظیفه             |
| ------------ | ----------------- |
| خواندن stdin | دریافت پیام کاربر |
| نویسنده TCP  | ارسال پیام        |
| خواننده TCP  | دریافت پیام       |
| حلقه اصلی    | نمایش پیام‌ها     |

---

### ▶️ نحوه اجرا

دو ترمینال جدا باز کنید.

#### ترمینال اول

```bash
cd peerA
go run .
```

#### ترمینال دوم

```bash
cd peerB
go run .
```

اکنون در هر کدام پیام بنویسید و Enter بزنید.

---

### 📊 فلو پیام‌ها

```
کاربر → outgoing → TCP
TCP → incoming → نمایش
```

---

### ⚠️ خطاهای رایج

* خطای `connection refused` یعنی طرف مقابل هنوز Listen نکرده
* پورت اشغال شده
* مشکل IPv6 (بهتر است از `127.0.0.1` استفاده شود)


بگو تا برات آماده کنم ✅
