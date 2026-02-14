package zbar

type SymbolType int32

const (
	SymbolType_NONE        SymbolType = 0   /**< no symbol decoded */
	SymbolType_PARTIAL     SymbolType = 1   /**< intermediate status */
	SymbolType_EAN2        SymbolType = 2   /**< GS1 2-digit add-on */
	SymbolType_EAN5        SymbolType = 5   /**< GS1 5-digit add-on */
	SymbolType_EAN8        SymbolType = 8   /**< EAN-8 */
	SymbolType_UPCE        SymbolType = 9   /**< UPC-E */
	SymbolType_ISBN10      SymbolType = 10  /**< ISBN-10 (from EAN-13). @since 0.4 */
	SymbolType_UPCA        SymbolType = 12  /**< UPC-A */
	SymbolType_EAN13       SymbolType = 13  /**< EAN-13 */
	SymbolType_ISBN13      SymbolType = 14  /**< ISBN-13 (from EAN-13). @since 0.4 */
	SymbolType_COMPOSITE   SymbolType = 15  /**< EAN/UPC composite */
	SymbolType_I25         SymbolType = 25  /**< Interleaved 2 of 5. @since 0.4 */
	SymbolType_DATABAR     SymbolType = 34  /**< GS1 DataBar (RSS). @since 0.11 */
	SymbolType_DATABAR_EXP SymbolType = 35  /**< GS1 DataBar Expanded. @since 0.11 */
	SymbolType_CODABAR     SymbolType = 38  /**< Codabar. @since 0.11 */
	SymbolType_CODE39      SymbolType = 39  /**< Code 39. @since 0.4 */
	SymbolType_PDF417      SymbolType = 57  /**< PDF417. @since 0.6 */
	SymbolType_QRCODE      SymbolType = 64  /**< QR Code. @since 0.10 */
	SymbolType_SQCODE      SymbolType = 80  /**< SQ Code. @since 0.20.1 */
	SymbolType_CODE93      SymbolType = 93  /**< Code 93. @since 0.11 */
	SymbolType_CODE128     SymbolType = 128 /**< Code 128 */
)

type Config int32

const (
	Config_ENABLE      Config = 0x0
	Config_ADD_CHECK   Config = 0x1
	Config_EMIT_CHECK  Config = 0x2
	Config_ASCII       Config = 0x3
	Config_BINARY      Config = 0x4
	Config_NUM         Config = 0x5
	Config_MIN_LEN     Config = 0x20
	Config_MAX_LEN     Config = 0x21
	Config_UNCERTAINTY Config = 0x40
	Config_POSITION    Config = 0x80
	Config_X_DENSITY   Config = 0x100
	Config_Y_DENSITY   Config = 0x101
)
