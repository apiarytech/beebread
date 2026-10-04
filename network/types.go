/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package network

import "github.com/apiarytech/royaljelly/iec"

// NW_BUF_LONG is the data of a NETWORK_BUFFER.
type NW_BUF_LONG [NETWORK_BUFFER_LONG_SIZE + 1]iec.BYTE

// NW_BUF_SHORT is the data of a NETWORK_BUFFER_SHORT.
type NW_BUF_SHORT [NETWORK_BUFFER_SHORT_SIZE + 1]iec.BYTE

// NETWORK_BUFFER is a buffer of data to send or received: SIZE bytes of
// BUFFER.
type NETWORK_BUFFER struct {
	SIZE   iec.UINT
	BUFFER NW_BUF_LONG
}

// NETWORK_BUFFER_SHORT is a short NETWORK_BUFFER.
type NETWORK_BUFFER_SHORT struct {
	SIZE   iec.UINT
	BUFFER NW_BUF_SHORT
}

// IP_FIFO_DATA is the queue IP_FIFO keeps of the blocks that share an
// IP_CONTROL.
type IP_FIFO_DATA struct {
	X      [128]iec.BYTE // ARRAY[1..128]: the queue of IDs
	Y      [128]iec.BYTE // ARRAY[1..128]: the entries of each ID
	ID     iec.BYTE      // the highest ID given
	MAX_ID iec.BYTE      // the entries an ID may have
	INIT   iec.BOOL      // set up
	EMPTY  iec.BOOL
	FULL   iec.BOOL
	TOP    iec.INT // the size of the queue
	NW     iec.INT // where to write
	NR     iec.INT // where to read
}

// IP_C is the interface between IP_CONTROL and the blocks that use its
// connection: the connection they ask for (W), its state (R), and IP_FIFO's
// queue (I).
type IP_C struct {
	C_MODE     iec.BYTE // W: 0 TCP client, 1 UDP client, 2 TCP server, 3 UDP server, 4 TCP server for any remote, 5 UDP server for any remote
	C_PORT     iec.WORD // W: the port
	C_IP       iec.DWORD
	C_STATE    iec.BYTE // R: 0 off, 1 just off, 254 just on, 255 on
	C_ENABLE   iec.BOOL // W: connect
	R_OBSERVE  iec.BOOL // W: watch the time since data were received
	TIME_RESET iec.BOOL // W: reset the timers
	ERROR      iec.DWORD
	FIFO       IP_FIFO_DATA
	MAILBOX    [16]iec.BYTE // ARRAY[1..16]: data the blocks exchange
}

// PRINTF_DATA are the arguments ~1..~11 of a message of LOG_MSG.
type PRINTF_DATA [11]iec.STRING // ARRAY[1..11] OF STRING(LOG_SIZE)

// LOG_CONTROL is a log of messages, a ring of LOG_MAX + 1 messages.
type LOG_CONTROL struct {
	NEW_MSG        iec.STRING // STRING(LOG_SIZE): the message to add
	NEW_MSG_OPTION iec.DWORD
	LEVEL          iec.BYTE
	SIZE           iec.INT // default LOG_MAX: the number of messages
	RESET          iec.BOOL
	PRINTF         PRINTF_DATA
	MSG            [LOG_MAX + 1]iec.STRING // ARRAY[0..LOG_MAX] OF STRING(LOG_SIZE)
	MSG_OPTION     [LOG_MAX + 1]iec.DWORD  // ARRAY[0..LOG_MAX]
	UPDATE_COUNT   iec.UINT
	IDX            iec.INT  // the current message
	RING_MODE      iec.BOOL // the log has wrapped
}

// NewLOG_CONTROL returns a LOG_CONTROL with its initial values.
func NewLOG_CONTROL() LOG_CONTROL { return LOG_CONTROL{SIZE: LOG_MAX} }

// URL is a URL split into its parts.
type URL struct {
	PROTOCOL iec.STRING // STRING(10)
	USER     iec.STRING // STRING(32)
	PASSWORD iec.STRING // STRING(32)
	DOMAIN   iec.STRING // STRING(80)
	PORT     iec.WORD
	PATH     iec.STRING // STRING(80)
	QUERY    iec.STRING // STRING(STRING_LENGTH)
	ANCHOR   iec.STRING // STRING(40)
	HEADER   iec.STRING // STRING(160)
}

// FILE_PATH_DATA is a file path split into its parts.
type FILE_PATH_DATA struct {
	DRIVE     iec.STRING // STRING(3)
	DIRECTORY iec.STRING // STRING(STRING_LENGTH)
	FILENAME  iec.STRING
}

// FILE_SERVER_RUNTIME_DATA are the current, highest and lowest times, in
// ms, of FILE_SERVER's operations.
type FILE_SERVER_RUNTIME_DATA struct {
	TIME_FILE_OPEN_CUR, TIME_FILE_OPEN_MAX, TIME_FILE_OPEN_MIN       iec.UDINT
	TIME_FILE_CLOSE_CUR, TIME_FILE_CLOSE_MAX, TIME_FILE_CLOSE_MIN    iec.UDINT
	TIME_FILE_READ_CUR, TIME_FILE_READ_MAX, TIME_FILE_READ_MIN       iec.UDINT
	TIME_FILE_WRITE_CUR, TIME_FILE_WRITE_MAX, TIME_FILE_WRITE_MIN    iec.UDINT
	TIME_FILE_SEEK_CUR, TIME_FILE_SEEK_MAX, TIME_FILE_SEEK_MIN       iec.UDINT
	TIME_FILE_TELL_CUR, TIME_FILE_TELL_MAX, TIME_FILE_TELL_MIN       iec.UDINT
	TIME_FILE_DELETE_CUR, TIME_FILE_DELETE_MAX, TIME_FILE_DELETE_MIN iec.UDINT
}

// FILE_SERVER_DATA is the interface to FILE_SERVER: the file, what to do
// with it, and its state.
type FILE_SERVER_DATA struct {
	FILE_OPEN iec.BOOL
	FILENAME  iec.STRING
	MODE      iec.BYTE
	OFFSET    iec.UDINT
	FILE_SIZE iec.UDINT
	ERROR     iec.BYTE
	RUNTIME   FILE_SERVER_RUNTIME_DATA
}

// UNI_CIRCULAR_BUF is the data of a UNI_CIRCULAR_BUFFER_DATA.
type UNI_CIRCULAR_BUF [8192]iec.BYTE

// UNI_CIRCULAR_BUFFER_DATA is a ring buffer of values of several types, the
// interface to UNI_CIRCULAR_BUFFER.
type UNI_CIRCULAR_BUFFER_DATA struct {
	D_MODE        iec.INT
	D_SIZE        iec.INT
	D_HEAD        iec.WORD
	D_STRING      iec.STRING // STRING(STRING_LENGTH)
	D_REAL        iec.REAL
	D_DWORD       iec.DWORD
	BUF_SIZE      iec.UINT
	BUF_COUNT     iec.UINT
	BUF_USED      iec.USINT
	BUF_USED_MAX  iec.USINT
	BUF_DATA_CNT  iec.UDINT
	BUF_DATA_LOST iec.UDINT
	BUF           UNI_CIRCULAR_BUF
	GETSTART_     iec.UINT // _GetStart
	GETEND_       iec.UINT // _GetEnd
	LAST_         iec.UINT // _Last
	FIRST_        iec.UINT // _First
}

// DLOG_DATA is the interface between the DLOG blocks of a data log.
type DLOG_DATA struct {
	STORE_TYPE     iec.BYTE
	ADD_COM        iec.INT
	ADD_DATA_REQ   iec.BOOL
	CLOCK_TRIG     iec.BOOL
	ID_MAX         iec.USINT
	DTI            iec.DT
	UCB            UNI_CIRCULAR_BUFFER_DATA
	LOAD_TIME_MAX  iec.TIME
	NEW_FILE       iec.STRING
	NEW_FILE_SIZE  iec.UDINT
	NEW_FILE_RTRIG iec.BOOL
}

// DLOG_CRON_ELEMENT is a field of a cron table: the values it selects.
type DLOG_CRON_ELEMENT struct {
	ELEMENTS     [60]iec.BOOL // ARRAY[0..59]
	VALUE        iec.INT
	VALUE_MIN    iec.INT
	VALUE_MAX    iec.INT
	ALL_SELECTED iec.BOOL
}

// DLOG_CRON_DATA are the six fields of a cron table.
type DLOG_CRON_DATA struct {
	CE [6]DLOG_CRON_ELEMENT // ARRAY[0..5]
}

// DLOG_REAL_ARRAY_ELEMENT is a value DLOG_REAL_ARRAY logs.
type DLOG_REAL_ARRAY_ELEMENT struct {
	VALUE       iec.REAL
	DELTA       iec.REAL
	COLUMN      iec.STRING // STRING(10)
	DELTA_LAST_ iec.REAL   // _DELTA_LAST
	VALUE_LAST_ iec.REAL   // _VALUE_LAST
}

// DLOG_REAL_ARRAY_DATA are the values DLOG_REAL_ARRAY logs.
type DLOG_REAL_ARRAY_DATA [64]DLOG_REAL_ARRAY_ELEMENT // ARRAY[1..64]

// DLOG_RETAIN is the state a data log keeps across a restart.
type DLOG_RETAIN struct {
	FN_REM         iec.STRING
	COLOR          iec.BOOL
	HEAD           iec.BOOL
	TRIG_CNT       iec.UDINT
	TRIG_CNT_TOTAL iec.UDINT
}

// DLOG_SAVE is the state a data log keeps across a restart.
type DLOG_SAVE struct {
	FN_REM         iec.STRING
	COLOR          iec.BOOL
	HEAD           iec.BOOL
	TRIG_CNT       iec.UDINT
	TRIG_CNT_TOTAL iec.UDINT
}

// IP2GEO_DATA is the location of an IP address.
type IP2GEO_DATA struct {
	STATE          iec.BOOL
	IP4            iec.DWORD
	COUNTRY_CODE   iec.STRING // STRING(2)
	COUNTRY_NAME   iec.STRING // STRING(20)
	REGION_CODE    iec.STRING // STRING(2)
	REGION_NAME    iec.STRING // STRING(20)
	CITY           iec.STRING // STRING(20)
	GEO_LATITUDE   iec.REAL
	GEO_LONGITUDE  iec.REAL
	TIME_ZONE_NAME iec.STRING // STRING(20)
	GMT_OFFSET     iec.INT
	IS_DST         iec.BOOL
}

// MYSQL_COM is the interface between MYSQL_CONTROL and its users.
type MYSQL_COM struct {
	S_BUF         NETWORK_BUFFER
	R_BUF         NETWORK_BUFFER
	SQL_CON       iec.BOOL
	SQL_URL       iec.STRING // STRING(STRING_LENGTH)
	SQL_PACKET_NO iec.BYTE
	TIMEOUT       iec.TIME
	DNS_IP4       iec.DWORD
	SQL_RCV_STATE iec.BYTE
	ERROR_C       iec.DWORD
	ERROR_T       iec.BYTE
}

// MYSQL_INFO is the state of a MySQL connection.
type MYSQL_INFO struct {
	SERVER_PROTOCOL_VERSION iec.BYTE
	SERVER_CAPABILITIES     iec.WORD
	SERVER_STATUS           iec.WORD
	SERVER_LANGUAGE         iec.BYTE
	SQL_CONNECTED           iec.BOOL
	SQL_ERROR               iec.STRING // STRING(STRING_LENGTH)
	DATA_INSERT_AKTIV       iec.BOOL
	DATA_INSERT_OK_CNT      iec.UDINT
	DATA_INSERT_NOK_CNT     iec.UDINT
}

// NET_VAR_DATA is the interface between NET_VAR_CONTROL and the NET_VAR
// blocks.
type NET_VAR_DATA struct {
	CYCLE    iec.UDINT
	STATE    iec.BYTE
	INDEX    iec.INT
	ID_MAX   iec.USINT
	ERROR_ID iec.BYTE
	BUF_SIZE iec.UINT
	S_BUF    NETWORK_BUFFER
	R_BUF    NETWORK_BUFFER
}

// VMAP_DATA maps Modbus addresses to the process image.
type VMAP_DATA struct {
	FC       iec.DWORD
	V_ADR    iec.INT
	V_SIZE   iec.INT
	P_ADR    iec.INT
	TIME_OUT iec.TIME
}

// XML_CONTROL is the interface to XML_READER: the command and the element
// found.
type XML_CONTROL struct {
	COMMAND      iec.WORD
	WATCHDOG     iec.TIME
	START_POS    iec.UINT
	STOP_POS     iec.UINT
	COUNT        iec.UINT
	TYP          iec.INT
	LEVEL        iec.UINT
	PATH         iec.STRING // STRING(STRING_LENGTH)
	ELEMENT      iec.STRING // STRING(STRING_LENGTH)
	ATTRIBUTE    iec.STRING // STRING(STRING_LENGTH)
	VALUE        iec.STRING // STRING(STRING_LENGTH)
	BLOCK1_START iec.UINT
	BLOCK1_STOP  iec.UINT
	BLOCK2_START iec.UINT
	BLOCK2_STOP  iec.UINT
}

// WORLD_WEATHER_CUR is the current weather of WORLD_WEATHER.
type WORLD_WEATHER_CUR struct {
	OBSERVATION_TIME iec.STRING // STRING(8)
	TEMP_C           iec.INT
	WEATHER_CODE     iec.INT
	WEATHER_DESC     iec.STRING // STRING(60)
	WEATHER_ICON     iec.INT
	WIND_SPEED_MILES iec.INT
	WIND_SPEED_KMPH  iec.INT
	WIND_DIR_DEGREE  iec.INT
	WIND_DIR16POINT  iec.STRING // STRING(3)
	PRECIPMM         iec.REAL
	HUMIDITY         iec.INT
	VISIBILITY       iec.INT
	PRESSURE         iec.INT
	CLOUDOVER        iec.INT
}

// WORLD_WEATHER_DAY is the weather forecast of a day of WORLD_WEATHER.
type WORLD_WEATHER_DAY struct {
	DATE_OF_DAY      iec.STRING // STRING(10)
	TEMP_MAX_C       iec.INT
	TEMP_MAX_F       iec.INT
	TEMP_MIN_C       iec.INT
	TEMP_MIN_F       iec.INT
	WIND_SPEED_MILES iec.INT
	WIND_SPEED_KMPH  iec.INT
	WIND_DIR_DEGREE  iec.INT
	WIND_DIR16POINT  iec.STRING // STRING(3)
	WEATHER_CODE     iec.INT
	WEATHER_DESC     iec.STRING // STRING(60)
	WEATHER_ICON     iec.INT
	PRECIPMM         iec.REAL
}

// WORLD_WEATHER_DATA is the weather WORLD_WEATHER reads.
type WORLD_WEATHER_DATA struct {
	CUR WORLD_WEATHER_CUR
	DAY [5]WORLD_WEATHER_DAY // ARRAY[0..4]
}

// YAHOO_WEATHER_FORECAST_DAY is the weather forecast of a day of
// YAHOO_WEATHER.
type YAHOO_WEATHER_FORECAST_DAY struct {
	LOW_TEMP  iec.INT
	HIGH_TEMP iec.INT
	TEXT      iec.STRING // STRING(40)
	DATE_LONG iec.STRING // STRING(20)
	CODE      iec.INT
	ICON      iec.INT
}

// YAHOO_WEATHER_DATA is the weather YAHOO_WEATHER reads.
type YAHOO_WEATHER_DATA struct {
	LOCATION_CITY            iec.STRING // STRING(40)
	LOCATION_REGION          iec.STRING // STRING(20)
	LOCATION_COUNTRY         iec.STRING // STRING(20)
	UNIT_TEMPERATURE         iec.STRING // STRING(1): F or C
	UNIT_DISTANCE            iec.STRING // STRING(2): mi or km
	UNIT_PRESSURE            iec.STRING // STRING(2): in or mb
	UNIT_SPEED               iec.STRING // STRING(4): mph or kph
	WIND_CHILL               iec.INT
	WIND_DIRECTION           iec.INT
	WIND_SPEED               iec.REAL
	ATMOSPHERE_HUMIDITY      iec.INT
	ATMOSPHERE_PRESSURE      iec.REAL
	ATMOSPHERE_VISIBILITY    iec.REAL
	ATMOSPHERE_RISING        iec.INT    // 0 steady, 1 rising, 2 falling
	ASTRONOMY_SUNRISE        iec.STRING // STRING(10)
	ASTRONOMY_SUNSET         iec.STRING // STRING(10)
	GEO_LATITUDE             iec.REAL
	GEO_LONGITUDE            iec.REAL
	CUR_CONDITIONS_TEMP      iec.INT
	CUR_CONDITIONS_TEXT      iec.STRING // STRING(40)
	CUR_CONDITIONS_CODE      iec.INT
	CUR_CONDITIONS_ICON      iec.INT
	FORECAST_TODAY_LOW_TEMP  iec.INT
	FORECAST_TODAY_HIGH_TEMP iec.INT
	FORECAST_TODAY_TEXT      iec.STRING // STRING(40)
	FORECAST_TODAY_CODE      iec.INT
	FORECAST_TODAY_ICON      iec.INT
	FORECAST_TODAY_DATE_LONG iec.STRING                    // STRING(20)
	FORECAST_DAY             [9]YAHOO_WEATHER_FORECAST_DAY // ARRAY[1..9]
}

// US_LOG_VIEWPORT is the state of LOG_VIEWPORT. OSCAT names it
// us_LOG_VIEWPORT.
type US_LOG_VIEWPORT struct {
	LINE_ARRAY   [40]iec.INT // ARRAY[1..40]
	COUNT        iec.INT     // the visible messages
	UPDATE_COUNT iec.UINT
	MOVE_TO_X    iec.INT
	UPDATE       iec.BOOL // redraw
}

// US_TN_SCREEN is the screen of a telnet session, 24 rows of 80 columns.
// OSCAT names it us_TN_SCREEN.
type US_TN_SCREEN struct {
	BYA_CHAR             [1920]iec.BYTE // ARRAY[0..1919]
	BYA_COLOR            [1920]iec.BYTE // ARRAY[0..1919]
	BYA_BACKUP           [1920]iec.BYTE // ARRAY[0..1919]
	BYA_LINE_UPDATE      [24]iec.BOOL   // ARRAY[0..23]
	BY_INPUT_EXTEN_CODE  iec.BYTE
	BY_INPUT_ASCII_CODE  iec.BYTE
	BO_INPUT_ASCII_ISNUM iec.BOOL
	IN_PAGE_NUMBER       iec.INT
	IN_CURSOR_X          iec.INT
	IN_CURSOR_Y          iec.INT
	IN_EOS_OFFSET        iec.INT
	BY_CLEAR_SCREEN_ATTR iec.BYTE
	BO_CLEAR_SCREEN      iec.BOOL
	BO_MODAL_DIALOG      iec.BOOL
	BO_MENUE_BAR_DIALOG  iec.BOOL
}

// US_TN_INPUT_CONTROL_DATA is an input element of a telnet screen. OSCAT
// names it us_TN_INPUT_CONTROL_DATA.
type US_TN_INPUT_CONTROL_DATA struct {
	BY_INPUT_EXTEN_CODE  iec.BYTE
	BY_INPUT_ASCII_CODE  iec.BYTE
	BO_INPUT_ASCII_ISNUM iec.BOOL
	IN_TITLE_X_OFFSET    iec.INT
	IN_TITLE_Y_OFFSET    iec.INT
	BY_TITLE_ATTR        iec.BYTE
	ST_TITLE_STRING      iec.STRING
	IN_CURSOR_X          iec.INT
	IN_CURSOR_Y          iec.INT
	IN_TYPE              iec.INT
	IN_X                 iec.INT
	IN_Y                 iec.INT
	IN_CURSOR_POS        iec.INT
	BY_ATTR_MF           iec.BYTE // the attribute with the focus
	BY_ATTR_OF           iec.BYTE // the attribute without the focus
	IN_SELECTED          iec.INT
	ST_INPUT_MASK        iec.STRING
	ST_INPUT_DATA        iec.STRING // STRING(STRING_LENGTH)
	ST_INPUT_STRING      iec.STRING
	ST_INPUT_TOOLTIP     iec.STRING
	IN_INPUT_OPTION      iec.INT
	BO_INPUT_ENTERED     iec.BOOL
	BO_INPUT_HIDDEN      iec.BOOL
	BO_INPUT_ONLY_NUM    iec.BOOL
	BO_FOCUS             iec.BOOL
	BO_UPDATE_INPUT      iec.BOOL
	BO_UPDATE_ALL        iec.BOOL
}

// US_TN_INPUT_CONTROL are the input elements of a telnet screen. OSCAT
// names it us_TN_INPUT_CONTROL.
type US_TN_INPUT_CONTROL struct {
	BO_ENABLE                 iec.BOOL
	BO_UPDATE_ALL             iec.BOOL
	BO_RESET_FOKUS            iec.BOOL
	IN_FOCUS_AT               iec.INT
	IN_COUNT                  iec.INT
	IN_TOOLTIP_X              iec.INT
	IN_TOOLTIP_Y              iec.INT
	BY_TOOLTIP_ATTR           iec.BYTE
	IN_TOOLTIP_SIZE           iec.INT
	USA_TN_INPUT_CONTROL_DATA [20]US_TN_INPUT_CONTROL_DATA // ARRAY[1..20]
}

// US_TN_MENU is the menu bar of a telnet screen. OSCAT names it us_TN_MENU.
type US_TN_MENU struct {
	ST_MENU_TEXT     iec.STRING // STRING(STRING_LENGTH)
	IN_MENU_E_COUNT  iec.INT
	IN_Y             iec.INT
	IN_X             iec.INT
	BY_ATTR_MF       iec.BYTE
	BY_ATTR_OF       iec.BYTE
	IN_X_SM_NEW      iec.INT
	IN_Y_SM_NEW      iec.INT
	IN_X_SM_OLD      iec.INT
	IN_Y_SM_OLD      iec.INT
	IN_CUR_MENU_ITEM iec.INT
	IN_CUR_SUB_ITEM  iec.INT
	IN_STATE         iec.INT
	IN_MENU_SELECTED iec.INT
	BO_CREATE        iec.BOOL
	BO_DESTROY       iec.BOOL
	BO_UPDATE        iec.BOOL
}

// US_TN_MENU_POPUP is a popup menu of a telnet screen. OSCAT names it
// us_TN_MENU_POPUP.
type US_TN_MENU_POPUP struct {
	ST_MENU_TEXT        iec.STRING // STRING(STRING_LENGTH)
	IN_MENU_E_COUNT     iec.INT
	IN_X                iec.INT
	IN_Y                iec.INT
	IN_COLS             iec.INT
	IN_ROWS             iec.INT
	IN_CUR_ITEM         iec.INT
	BY_ATTR_MF          iec.BYTE
	BY_ATTR_OF          iec.BYTE
	BY_INPUT_EXTEN_CODE iec.BYTE
	BO_CREATE           iec.BOOL
	BO_DESTROY          iec.BOOL
	BO_UPDATE           iec.BOOL
	BO_ACTIV            iec.BOOL
}
