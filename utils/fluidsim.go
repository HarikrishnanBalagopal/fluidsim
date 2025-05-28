package utils

import "math"

func clamp(x float32, min_x float32, max_x float32) float32 {
	if x < min_x {
		return min_x
	}
	if x > max_x {

		return max_x
	}
	return x
}

func advect_vel_color(t float32, dt float32) {
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			i := y*W + x
			curr_u := A_VEL_U[i]
			curr_v := A_VEL_V[i]
			old_x := float32(x) - curr_u*dt
			old_y := float32(y) - curr_v*dt

			__low_x := int(math.Floor(float64(old_x)))
			__low_y := int(math.Floor(float64(old_y)))
			low_x := (__low_x + W) % W
			high_x := (low_x + 1) % W
			low_y := (__low_y + H) % H
			high_y := (low_y + 1) % H
			off_x := old_x - float32(__low_x)
			off_y := old_y - float32(__low_y)

			lylx_i := low_y*W + low_x
			lyhx_i := low_y*W + high_x
			hylx_i := high_y*W + low_x
			hyhx_i := high_y*W + high_x

			lylx_u := A_VEL_U[lylx_i]
			lylx_v := A_VEL_V[lylx_i]
			lylx_c := A_COLOR[lylx_i]
			lylx_g := A_COLOG[lylx_i]
			lylx_b := A_COLOB[lylx_i]

			lyhx_u := A_VEL_U[lyhx_i]
			lyhx_v := A_VEL_V[lyhx_i]
			lyhx_c := A_COLOR[lyhx_i]
			lyhx_g := A_COLOG[lyhx_i]
			lyhx_b := A_COLOB[lyhx_i]

			hylx_u := A_VEL_U[hylx_i]
			hylx_v := A_VEL_V[hylx_i]
			hylx_c := A_COLOR[hylx_i]
			hylx_g := A_COLOG[hylx_i]
			hylx_b := A_COLOB[hylx_i]

			hyhx_u := A_VEL_U[hyhx_i]
			hyhx_v := A_VEL_V[hyhx_i]
			hyhx_c := A_COLOR[hyhx_i]
			hyhx_g := A_COLOG[hyhx_i]
			hyhx_b := A_COLOB[hyhx_i]

			ly_u := (1-off_x)*lylx_u + off_x*lyhx_u
			ly_v := (1-off_x)*lylx_v + off_x*lyhx_v
			ly_c := (1-off_x)*lylx_c + off_x*lyhx_c
			ly_g := (1-off_x)*lylx_g + off_x*lyhx_g
			ly_b := (1-off_x)*lylx_b + off_x*lyhx_b

			hy_u := (1-off_x)*hylx_u + off_x*hyhx_u
			hy_v := (1-off_x)*hylx_v + off_x*hyhx_v
			hy_c := (1-off_x)*hylx_c + off_x*hyhx_c
			hy_g := (1-off_x)*hylx_g + off_x*hyhx_g
			hy_b := (1-off_x)*hylx_b + off_x*hyhx_b

			new_u := (1-off_y)*ly_u + off_y*hy_u
			new_v := (1-off_y)*ly_v + off_y*hy_v
			new_c := (1-off_y)*ly_c + off_y*hy_c
			new_g := (1-off_y)*ly_g + off_y*hy_g
			new_b := (1-off_y)*ly_b + off_y*hy_b

			A_OUTPUT_1[i] = new_u
			A_OUTPUT_2[i] = new_v
			A_OUTPUT_3[i] = clamp(new_c, 0, 1)
			A_OUTPUT_4[i] = clamp(new_g, 0, 1)
			A_OUTPUT_5[i] = clamp(new_b, 0, 1)
		}
	}
	// copy output
	for i := 0; i < LEN; i++ {
		A_VEL_U[i] = A_OUTPUT_1[i]
		A_VEL_V[i] = A_OUTPUT_2[i]
		A_COLOR[i] = A_OUTPUT_3[i]
		A_COLOG[i] = A_OUTPUT_4[i]
		A_COLOB[i] = A_OUTPUT_5[i]
	}
}

func divergence_vel(t float32, dt float32) {
	scale := (-2 * CELL_DIST * DENSITY) / dt
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			i := y*W + x
			cylx_i := y*W + ((x + W - 1) % W)
			cyhx_i := y*W + ((x + W + 1) % W)
			lycx_i := ((y+H-1)%H)*W + x
			hycx_i := ((y+H+1)%H)*W + x

			cylx_u := A_VEL_U[cylx_i]
			cyhx_u := A_VEL_U[cyhx_i]
			lycx_v := A_VEL_V[lycx_i]
			hycx_v := A_VEL_V[hycx_i]

			div := scale * (cyhx_u - cylx_u + hycx_v - lycx_v)

			A_OUTPUT_1[i] = div
		}
	}
}

func calc_pressure_jacobi(t float32, dt float32) {
	// 10 iterations of jacobi method
	for iter := 0; iter < 10; iter++ {
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				i := y*W + x
				// we can use the +/- 1 cells instead
				cylx_i := y*W + ((x + W - 2) % W)
				cyhx_i := y*W + ((x + W + 2) % W)
				lycx_i := ((y+H-2)%H)*W + x
				hycx_i := ((y+H+2)%H)*W + x

				div := A_OUTPUT_1[i]
				cylx_p := A_PRESS[cylx_i]
				cyhx_p := A_PRESS[cyhx_i]
				lycx_p := A_PRESS[lycx_i]
				hycx_p := A_PRESS[hycx_i]
				A_OUTPUT_2[i] = (div + cylx_p + cyhx_p + lycx_p + hycx_p) / 4
			}
		}
		// copy output
		for i := 0; i < LEN; i++ {
			A_PRESS[i] = A_OUTPUT_2[i]
		}
	}
}

func add_forces(t float32, dt float32) {
	if !MOUSE_DOWN {
		return
	}
	x := int(math.Floor(float64(MOUSE_X)))
	y := int(math.Floor(float64(MOUSE_Y)))

	u := float32(10 * (MOUSE_X - LAST_MOUSE_X))
	v := float32(10 * (MOUSE_Y - LAST_MOUSE_Y))
	// const amp = Math.sqrt(u * u + v * v);
	const sigma_sq float64 = 100

	for off_y := -MOUSE_SIZE; off_y <= MOUSE_SIZE; off_y++ {
		yy := y + off_y
		for off_x := -MOUSE_SIZE; off_x <= MOUSE_SIZE; off_x++ {
			xx := x + off_x
			i := yy*W + xx
			// const ex = Math.exp(-off_x * off_x);
			// const ey = Math.exp(-off_y * off_y);
			var ttt1 float64 = -float64(off_x*off_x+off_y*off_y) / sigma_sq
			exy := clamp(
				float32(math.Exp(ttt1)),
				0,
				1,
			)
			// A_VEL_U[i] = 10 * Math.exp(-off_x * off_x);
			// A_VEL_V[i] = 10 * Math.exp(-off_y * off_y);
			A_VEL_U[i] = u
			A_VEL_V[i] = v
			A_COLOR[i] = clamp((1-exy)*A_COLOR[i]+exy*INK_COLOR_R, 0, 1)
			A_COLOG[i] = clamp((1-exy)*A_COLOG[i]+exy*INK_COLOR_G, 0, 1)
			A_COLOB[i] = clamp((1-exy)*A_COLOB[i]+exy*INK_COLOR_B, 0, 1)
			// A_PRESS[i] = 1000 * dt;
			// A_PRESS[i] = 10*Math.exp(-(off_x * off_x + off_y * off_y)); // radius^2
		}
	}
}

func sub_gradient_pressure(t float32, dt float32) {
	scale := dt / (2 * CELL_DIST * DENSITY)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			i := y*W + x
			cylx_i := y*W + ((x + W - 1) % W)
			cyhx_i := y*W + ((x + W + 1) % W)
			lycx_i := ((y+H-1)%H)*W + x
			hycx_i := ((y+H+1)%H)*W + x
			A_OUTPUT_1[i] = A_VEL_U[i] - scale*(A_PRESS[cyhx_i]-A_PRESS[cylx_i])
			A_OUTPUT_2[i] = A_VEL_V[i] - scale*(A_PRESS[hycx_i]-A_PRESS[lycx_i])
		}
	}
	// copy output
	for i := 0; i < LEN; i++ {
		A_VEL_U[i] = A_OUTPUT_1[i]
		A_VEL_V[i] = A_OUTPUT_2[i]
	}
}

func draw() {
	for i := 0; i < LEN; i++ {
		j := i * 4
		// const r = A_VEL_U[i] * 2.55;
		// const g = A_VEL_V[i] * 2.55;
		// const g = A_PRESS[i];// * 255;
		// const r = g;
		r := A_COLOR[i] * 255
		g := A_COLOG[i] * 255
		b := A_COLOB[i] * 255
		PIX_DATA[j+0] = r
		PIX_DATA[j+1] = g
		PIX_DATA[j+2] = b
		PIX_DATA[j+3] = 255
	}
	// CANVAS_CTX.putImageData(IMG_DATA, 0, 0)
}

func setup() {
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			i := y*W + x
			xy := float64(x+y) / 100
			A_VEL_U[i] = float32(100 * math.Sin(xy))
			A_VEL_V[i] = float32(100 * math.Cos(xy))
			A_PRESS[i] = 0
			// const i1 = Math.floor((10 * x) / W) % 2;
			// const i2 = Math.floor((10 * y) / H) % 2;
			// const i3 = (i1 + i2) % 2;
			// A_COLOR[i] = clamp(i3, 0, 1);
			// A_COLOG[i] = clamp(i3, 0, 1);
			// A_COLOB[i] = clamp(i3, 0, 1);
			j := i * 4
			r := PIX_DATA_COPY[j+0] / 255
			g := PIX_DATA_COPY[j+1] / 255
			b := PIX_DATA_COPY[j+2] / 255
			A_COLOR[i] = clamp(r, 0, 1)
			A_COLOG[i] = clamp(g, 0, 1)
			A_COLOB[i] = clamp(b, 0, 1)
		}
	}
	// draw();
}
