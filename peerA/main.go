package main // Main package: entry point of the Go application

import (
	"bufio"   // For buffered I/O (reading from stdin, writing to TCP)
	"fmt"     // For formatted input/output (printing logs)
	"net"     // For TCP networking
	"os"      // For accessing OS features (stdin)
	"strings" // For string manipulation (TrimSpace)
	"time"    // For timeouts and retry intervals
)

/*
Constant configuration values

مقادیر ثابت پیکربندی برنامه:
- آدرس گوش‌دادن محلی
- آدرس اتصال به peer مقابل
- فاصله تلاش مجدد برای اتصال
- timeout نوشتن روی TCP
*/
const (
	localListenAddr  = "0.0.0.0:8080"         // Peer A listens on this address | آدرس Listen این برنامه
	remoteDialAddr   = "127.0.0.1:8081"       // Peer A dials Peer B | آدرس Peer مقابل
	dialRetryEvery   = 700 * time.Millisecond // Delay between dial retries | فاصله تلاش مجدد اتصال
	connWriteTimeout = 5 * time.Second        // TCP write timeout | تایم‌اوت نوشتن روی TCP
)

func main() {
	// Startup logs | پیام‌های شروع برنامه
	fmt.Println("PeerA starting...")
	fmt.Println("Local listen:", localListenAddr)
	fmt.Println("Remote dial :", remoteDialAddr)
	fmt.Println("Type and press Enter to send. Ctrl+C to exit.")

	/*
		Channels definition

		outgoing: messages typed by user (to be sent)
		incoming: messages received from TCP
		done:     signal channel for shutdown

		تعریف کانال‌ها:
		- outgoing: پیام‌های خروجی کاربر
		- incoming: پیام‌های دریافتی از شبکه
		- done: اعلام پایان و قطع اتصال
	*/
	outgoing := make(chan string, 32)
	incoming := make(chan string, 32)
	done := make(chan struct{})

	// Start TCP listener | شروع گوش‌دادن روی TCP
	ln, err := net.Listen("tcp", localListenAddr)
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer ln.Close() // Ensure listener is closed on exit | بستن listener هنگام خروج

	/*
		acceptCh is used to receive an incoming connection asynchronously

		کانالی برای دریافت اتصال ورودی به‌صورت همزمان
	*/
	acceptCh := make(chan net.Conn, 1)
	go acceptOnce(ln, acceptCh) // Run accept in a goroutine | اجرای Accept در goroutine

	/*
		Establish connection:
		- Either accept incoming connection
		- Or dial the remote peer

		برقراری اتصال:
		- یا اتصال ورودی را می‌پذیرد
		- یا به peer مقابل وصل می‌شود
	*/
	conn := establishConn(acceptCh, remoteDialAddr)
	if conn == nil {
		fmt.Println("Failed to establish connection.")
		return
	}
	defer conn.Close() // Close connection on exit | بستن اتصال هنگام خروج

	fmt.Println("Connected to:", conn.RemoteAddr())

	// Start concurrent I/O goroutines | شروع goroutineهای ورودی/خروجی
	go stdinReader(outgoing, done)      // Read user input | خواندن ورودی کاربر
	go connWriter(conn, outgoing, done) // Write to TCP | ارسال پیام روی TCP
	go connReader(conn, incoming, done) // Read from TCP | دریافت پیام از TCP

	/*
		Main event loop:
		- Prints incoming messages
		- Exits on done signal

		حلقه اصلی:
		- نمایش پیام‌های دریافتی
		- خروج امن در صورت بسته‌شدن اتصال
	*/
	for {
		select {
		case msg := <-incoming:
			fmt.Println(msg)
		case <-done:
			fmt.Println("Connection closed. Bye.")
			return
		}
	}
}

/*
acceptOnce waits for a single incoming TCP connection
and sends it into acceptCh.

این تابع منتظر یک اتصال TCP ورودی می‌ماند
و آن را داخل کانال acceptCh ارسال می‌کند
*/
func acceptOnce(ln net.Listener, acceptCh chan<- net.Conn) {
	conn, err := ln.Accept() // Block until a connection arrives | انتظار برای اتصال
	if err != nil {
		return
	}
	acceptCh <- conn // Send accepted connection | ارسال اتصال پذیرفته‌شده
}

/*
establishConn races between:
- accepting an incoming connection
- dialing the remote peer

این تابع بین دو حالت رقابت ایجاد می‌کند:
- دریافت اتصال ورودی
- تلاش برای اتصال به peer مقابل
*/
func establishConn(acceptCh <-chan net.Conn, remote string) net.Conn {
	for {
		select {
		case c := <-acceptCh:
			// Incoming connection wins | اتصال ورودی برنده می‌شود
			return c
		default:
			// Try dialing remote peer | تلاش برای اتصال به peer مقابل
			c, err := net.Dial("tcp", remote)
			if err == nil {
				return c
			}
			time.Sleep(dialRetryEvery) // Wait before retry | صبر قبل از تلاش مجدد
		}
	}
}

/*
stdinReader reads user input from terminal
and sends it to outgoing channel.

این تابع ورودی کاربر را از ترمینال می‌خواند
و داخل کانال outgoing قرار می‌دهد
*/
func stdinReader(outgoing chan<- string, done <-chan struct{}) {
	sc := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-done:
			return // Stop on shutdown | توقف هنگام خروج
		default:
		}

		if !sc.Scan() {
			return // End of input | پایان ورودی
		}
		line := strings.TrimSpace(sc.Text()) // Remove extra spaces | حذف فاصله‌های اضافی
		if line == "" {
			continue // Ignore empty lines | نادیده گرفتن خطوط خالی
		}
		outgoing <- "A: " + line // Send message | ارسال پیام
	}
}

/*
connWriter writes messages from outgoing channel
to the TCP connection.

این تابع پیام‌ها را از outgoing گرفته
و روی اتصال TCP می‌نویسد
*/
func connWriter(conn net.Conn, outgoing <-chan string, done chan struct{}) {
	w := bufio.NewWriter(conn)
	for {
		select {
		case <-done:
			return // Stop on shutdown | توقف در صورت خروج
		case msg := <-outgoing:
			_ = conn.SetWriteDeadline(time.Now().Add(connWriteTimeout)) // Set write timeout | تنظیم تایم‌اوت
			_, err := w.WriteString(msg + "\n")                         // Write message | نوشتن پیام
			if err != nil {
				closeDone(done)
				return
			}
			if err := w.Flush(); err != nil { // Flush buffer | ارسال نهایی داده
				closeDone(done)
				return
			}
		}
	}
}

/*
connReader reads messages from TCP connection
and sends them to incoming channel.

این تابع پیام‌ها را از TCP می‌خواند
و داخل incoming قرار می‌دهد
*/
func connReader(conn net.Conn, incoming chan<- string, done chan struct{}) {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		incoming <- "RECV -> " + sc.Text() // Forward received message | ارسال پیام دریافتی
	}
	closeDone(done) // Connection closed | قطع اتصال
}

/*
closeDone safely closes the done channel only once.

این تابع کانال done را فقط یک‌بار و ایمن می‌بندد
*/
func closeDone(done chan struct{}) {
	select {
	case <-done:
		// already closed | قبلاً بسته شده
	default:
		close(done)
	}
}
