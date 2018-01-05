package tea

type Options map[string]interface{}

const (

	//OptionKeepAlive用于设置TCP KeepAlive
	//Value是一个布尔值,默认是true
	OptionKeepAlive = "KeepAlive"

	//OptionNoDelay 没有延迟
	//Value是一个布尔值,默认是true
	OptionNoDelay = "No-Delay"

	//连接间隔
	OptionConnInterval = "Conn-Interval"

	//最大连接数
	OptionConnNum = "Conn-Number"

	//消息最大缓存数
	OptionMsgNum = "Msg-Number"

	//消息长度
	// OptioneMsgLen = "Msg-lenght"

	//最小消息长度
	OptionMinMsgLen = "Msg-Min-lenght"

	//最大消息长度
	OptionMaxMsgLen = "Msg-Max-lenght"

	//最大读写
	OptionMaxRW = "Max-ReadWrite"
	//小端
	OptionLittleEndian = "LittleEndian"
)

func (opts Options) SetOption(n string, v interface{}) error {
	switch n {
	case OptionNoDelay:
		fallthrough
	case OptionKeepAlive, OptionLittleEndian:
		switch v := v.(type) {
		case bool:
			opts[n] = v
			return nil
		default:
			return ErrBadValue
		}
	case OptionMinMsgLen, OptionMaxMsgLen:
		switch v := v.(type) {
		case uint32:
			opts[n] = v
			return nil
		default:
			return ErrBadValue
		}
	case OptionConnNum, OptionConnInterval, OptionMsgNum:
		switch v := v.(type) {
		case int:
			opts[n] = v
			return nil
		default:
			return ErrBadValue
		}
	}
	return ErrBadValue
}

func (opts Options) GetOption(n string) (interface{}, error) {
	v, ok := opts[n]
	if ok {
		return v, nil
	}
	return nil, ErrBadValue
}
