package errorkit

const (
	// Server Errors
	InternalServer     = "Terjadi kesalahan pada server, silakan coba beberapa saat lagi"
	NotImplemented     = "Fungsi ini belum diimplementasikan"
	ServiceUnavailable = "Layanan tidak tersedia, silakan coba lagi nanti"

	// Client Errors
	BadRequest          = "Permintaan tidak lengkap atau tidak valid"
	Unauthorized        = "Akses tidak diizinkan, silakan login terlebih dahulu"
	Forbidden           = "Anda tidak memiliki izin untuk mengakses sumber daya ini"
	NotFound            = "Sumber daya yang diminta tidak ditemukan"
	MethodNotAllowed    = "Metode HTTP yang digunakan tidak diizinkan untuk permintaan ini"
	Conflict            = "Terjadi konflik saat memproses permintaan, silakan coba lagi"
	UnprocessableEntity = "Entitas tidak dapat diproses, periksa data yang dikirim"

	// Validation Errors
	ValidationError    = "Data yang dikirim tidak valid"
	RequiredField      = "Kolom %s wajib diisi"
	InvalidEmail       = "Email tidak valid"
	InvalidPhoneNumber = "Nomor telepon tidak valid"
	PasswordTooWeak    = "Password terlalu lemah, silakan gunakan kombinasi yang lebih kuat"

	// Database Errors
	DatabaseError   = "Terjadi kesalahan saat mengakses basis data"
	RecordNotFound  = "Data yang dicari tidak ditemukan"
	DuplicateRecord = "Data sudah ada, tidak dapat membuat duplikat"

	// Payment Errors
	PaymentFailed     = "Pembayaran gagal, silakan coba lagi"
	InsufficientFunds = "Saldo tidak mencukupi untuk melakukan transaksi"

	// Other Errors
	Timeout      = "Permintaan telah kadaluarsa, silakan coba lagi"
	UnknownError = "Terjadi kesalahan yang tidak diketahui, silakan coba lagi"
)

type ErrorResponse struct {
	err          error
	errorMessage string
	statusCode   int
}

func Error(err error, errorMessage string, statusCode int) error {
	return &ErrorResponse{
		err:          err,
		errorMessage: errorMessage,
		statusCode:   statusCode,
	}
}

// Error is SystemError message
func (er *ErrorResponse) Error() string {
	return er.err.Error()
}

// ErrorMessage reason error
func (er *ErrorResponse) ErrorMessage() string {
	return er.errorMessage
}

// StatusCode http status code
func (er *ErrorResponse) StatusCode() int {
	return er.statusCode
}
