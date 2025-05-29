package utils

//go:wasmexport GetConstWidth
func GetConstWidth() int32 {
	return W
}

//go:wasmexport GetConstHeight
func GetConstHeight() int32 {
	return H
}

//go:wasmexport GetAddrA_OUTPUT_1
func GetAddrA_OUTPUT_1() *[LEN]float32 {
	return &A_OUTPUT_1
}

//go:wasmexport GetAddrA_OUTPUT_2
func GetAddrA_OUTPUT_2() *[LEN]float32 {
	return &A_OUTPUT_2
}

//go:wasmexport GetAddrA_OUTPUT_3
func GetAddrA_OUTPUT_3() *[LEN]float32 {
	return &A_OUTPUT_3
}

//go:wasmexport GetAddrA_OUTPUT_4
func GetAddrA_OUTPUT_4() *[LEN]float32 {
	return &A_OUTPUT_4
}

//go:wasmexport GetAddrA_OUTPUT_5
func GetAddrA_OUTPUT_5() *[LEN]float32 {
	return &A_OUTPUT_5
}

//go:wasmexport GetAddrA_COLOR
func GetAddrA_COLOR() *[LEN]float32 {
	return &A_COLOR
}

//go:wasmexport GetAddrA_COLOG
func GetAddrA_COLOG() *[LEN]float32 {
	return &A_COLOG
}

//go:wasmexport GetAddrA_COLOB
func GetAddrA_COLOB() *[LEN]float32 {
	return &A_COLOB
}

//go:wasmexport GetAddrA_PRESS
func GetAddrA_PRESS() *[LEN]float32 {
	return &A_PRESS
}

//go:wasmexport GetAddrA_VEL_U
func GetAddrA_VEL_U() *[LEN]float32 {
	return &A_VEL_U
}

//go:wasmexport GetAddrA_VEL_V
func GetAddrA_VEL_V() *[LEN]float32 {
	return &A_VEL_V
}

//go:wasmexport GetAddrPIX_DATA
func GetAddrPIX_DATA() *[LEN_4]uint8 {
	return &PIX_DATA
}

//go:wasmexport GetAddrPIX_DATA_COPY
func GetAddrPIX_DATA_COPY() *[LEN_4]uint8 {
	return &PIX_DATA_COPY
}
