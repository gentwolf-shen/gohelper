package iconv

func GB2312ToUTF8String(in string) (string, error) {
	return ConvertString(in, "GB2312", "UTF-8")
}

func GB2312ToUTF8(in []byte) ([]byte, error) {
	output := make([]byte, len(in)*2)
	_, outputLen, err := Convert(in, output, "GB2312", "UTF-8")

	if err != nil {
		return nil, err
	}

	return output[0:outputLen], nil
}

func GBKToUTF8(in []byte) ([]byte, error) {
	output := make([]byte, len(in)*2)
	_, outputLen, err := Convert(in, output, "GBK", "UTF-8")

	if err != nil {
		return nil, err
	}

	return output[0:outputLen], nil
}
