package zbar

type SymbolType int32

const (
	SymbolType_NONE SymbolType = iota
	SymbolType_PARTIAL
	SymbolType_EAN8
	SymbolType_UPCE
	SymbolType_ISBN10
	SymbolType_UPCA
	SymbolType_EAN13
	SymbolType_ISBN13
	SymbolType_I25
	SymbolType_CODE39
	SymbolType_PDF417
	SymbolType_QRCODE
	SymbolType_CODE128
	SymbolType_SYMBOL
	SymbolType_ADDON2
	SymbolType_ADDON5
	SymbolType_ADDON
)

type Config int32

const (
	Config_ENABLE Config = iota
	Config_ADD_CHECK
	Config_EMIT_CHECK
	Config_ASCII
	Config_NUM
	Config_MIN_LEN
	Config_MAX_LEN
	Config_POSITION
	Config_X_DENSITY
	Config_Y_DENSITY
)
