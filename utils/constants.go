package utils

const W = 480
const H = W
const LEN = W * H
const LEN_4 = LEN * 4
const CELL_DIST = 1
const DENSITY = 1

var INK_COLOR_R float32 = 0
var INK_COLOR_G float32 = 0
var INK_COLOR_B float32 = 0

var A_OUTPUT_1 = [LEN]float32{}
var A_OUTPUT_2 = [LEN]float32{}
var A_OUTPUT_3 = [LEN]float32{}
var A_OUTPUT_4 = [LEN]float32{}
var A_OUTPUT_5 = [LEN]float32{}

var A_COLOR = [LEN]float32{}
var A_COLOG = [LEN]float32{}
var A_COLOB = [LEN]float32{}
var A_PRESS = [LEN]float32{}
var A_VEL_U = [LEN]float32{}
var A_VEL_V = [LEN]float32{}

var PIX_DATA = [LEN_4]uint8{}
var PIX_DATA_COPY = [LEN_4]uint8{}

const MOUSE_SIZE = 20

var MOUSE_DOWN = false
var LAST_MOUSE_X = 0
var LAST_MOUSE_Y = 0
var MOUSE_X = 0
var MOUSE_Y = 0
var REQUESTED_ANIMATION = false
