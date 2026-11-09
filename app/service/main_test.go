package service

// =============================================================================
// TestMain — service layer memakai logger singleton `log.Logging`
// (di-init di main.go via log.InitLog()). Di test kita init sekali di sini
// supaya code path yang memanggil log.Logging.* tidak panic nil pointer.
// Catatan: InitLogger menulis ke logs/access.log & logs/error.log
// (pastikan folder logs/ ada).
//
// Pola stub: embed interface (nil) lalu implement HANYA method yang dipakai
// skenario — lihat auth_service_test.go.
// =============================================================================

import (
	"os"
	"testing"

	servicelog "dc-transaction-service/log"
)

func TestMain(m *testing.M) {
	servicelog.InitLog()
	os.Exit(m.Run())
}
