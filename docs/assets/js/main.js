import './wasm_exec.js';

const getPos = (e) => {
    const rect = e.target.getBoundingClientRect();
    const x = e.clientX - rect.left; //x position within the element.
    const y = e.clientY - rect.top; //y position within the element.
    return { x, y };
};

function fallbackCopyTextToClipboard(text) {
    var textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        var successful = document.execCommand('copy');
        var msg = successful ? 'successful' : 'unsuccessful';
        console.log('Fallback: Copying text command was ' + msg);
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
    }

    document.body.removeChild(textArea);
}

// https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript
function copyTextToClipboard(text) {
    if (!navigator.clipboard) {
        fallbackCopyTextToClipboard(text);
        return;
    }
    navigator.clipboard.writeText(text).then(function () {
        console.log('Async: Copying to clipboard was successful!');
    }, function (err) {
        console.error('Async: Could not copy text: ', err);
    });
}

const createImage = async (src) => {
    const img = document.createElement('img');
    img.setAttribute('width', 480);
    img.setAttribute('height', 480);
    img.setAttribute('src', src);
    img.setAttribute('alt', 'random image');
    // document.body.appendChild(img);
    // return img;
    return new Promise(resolve => img.addEventListener('load', () => resolve(img)));
};

const getImageData = (img) => {
    // console.log('img', img);
    const canvas = document.createElement('canvas');
    // document.body.appendChild(canvas);
    canvas.width = 480;
    canvas.height = 480;
    const ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);
    const img_data = ctx.getImageData(0, 0, 480, 480);
    return img_data;
};

async function main() {
    const img1 = await createImage('./assets/images/image1.png');
    const imgdata1 = getImageData(img1);
    // console.log('img1', img1, 'imgdata1', imgdata1);
    // const img2 = await createImage('./assets/images/image2.png');
    // const imgdata2 = getImageData(img2);
    // console.log('img2', img2, 'imgdata2', imgdata2);
    const W = 480;
    const H = W;
    const canvas_output = document.querySelector('#canvas-output');
    canvas_output.width = W;
    canvas_output.height = H;
    const canvas_output_ctx = canvas_output.getContext('2d');
    const img_data = canvas_output_ctx.getImageData(0, 0, W, H);
    const pix_data = img_data.data;
    // console.log('img_data', img_data);
    // console.log('pix_data', pix_data);


    // const button_play_pause = document.querySelector('#button-play-pause');
    let paused = false;
    let first = null;
    let offset = 2000.0;
    // button_play_pause.addEventListener('click', () => {
    //     paused = !paused;
    //     if (paused) {
    //         console.log('pausing');
    //         button_play_pause.textContent = '▶️ Play';
    //     } else {
    //         console.log('playing');
    //         requestAnimationFrame(draw);
    //         button_play_pause.textContent = '⏸ Pause';
    //     }
    // })

    // const pre = document.querySelector('#output-pre');
    // pre.textContent = 'Hi there!';

    // const copy_alert = document.querySelector('#copy-alert');
    // const button_copy = document.querySelector('#button-copy');
    // button_copy.addEventListener('click', () => {
    //     copyTextToClipboard(pre.textContent);
    //     copy_alert.classList.remove('hidden');
    //     setTimeout(() => copy_alert.classList.add('hidden'), 1000);
    // });

    const res = await fetch('assets/wasm/fluidsim.wasm');
    if (!res.ok) return console.error('failed to fetch the wasm module. status:', res.status);
    const moduleBytes = await res.arrayBuffer();
    const go = new Go();
    const module = await WebAssembly.instantiate(moduleBytes, go.importObject);
    // console.log('module', module);
    go.run(module.instance);

    // const decoder = new TextDecoder();
    const _W = module.instance.exports.GetConstWidth();
    const _H = module.instance.exports.GetConstHeight();
    // console.log('W', W, 'H', H, '_W', _W, '_H', _H);
    const LEN = W * H;
    const LEN_4 = LEN * 4;
    const get_arr = (s, arr = Float32Array, buf_len = LEN) => {
        const address = module.instance.exports[s]();
        // console.log('address', address);
        const mem = new arr(module.instance.exports.mem.buffer, address, buf_len);
        // const mem = new Uint8Array(module.instance.exports.mem.buffer, address, buf_len);
        return mem;
    };
    const A_COLOR = get_arr('GetAddrA_COLOR');
    // console.log('A_COLOR', A_COLOR);
    const A_COLOG = get_arr('GetAddrA_COLOG');
    // console.log('A_COLOG', A_COLOG);
    const A_COLOB = get_arr('GetAddrA_COLOB');
    // console.log('A_COLOB', A_COLOB);
    const A_PRESS = get_arr('GetAddrA_PRESS');
    // console.log('A_PRESS', A_PRESS);
    const A_VEL_U = get_arr('GetAddrA_VEL_U');
    // console.log('A_VEL_U', A_VEL_U);
    const A_VEL_V = get_arr('GetAddrA_VEL_V');
    // console.log('A_VEL_V', A_VEL_V);

    const PIX_DATA = get_arr('GetAddrPIX_DATA', Uint8ClampedArray, LEN_4);
    // console.log('PIX_DATA', PIX_DATA);
    const PIX_DATA_COPY = get_arr('GetAddrPIX_DATA_COPY', Uint8ClampedArray, LEN_4);
    // console.log('PIX_DATA_COPY', PIX_DATA_COPY);

    // const addr_A_COLOR = module.instance.exports.GetAddrA_COLOR();
    // const mem = new Uint8Array(module.instance.exports.mem.buffer, addr_A_COLOR, LEN);
    // console.log('mem', mem);

    let last_t = null;
    // const TIME_STEP = 1 * 0.0001;
    const draw = (_t) => {
        requestAnimationFrame(draw);
        const t = _t * 0.001;
        if (last_t < 0) last_t = t;
        const dt = t - last_t;
        if (dt < 0.001) return;
        last_t = t;

        // let t = 0.0001 * _t;
        // if (paused) { first = null; offset = last_t; return; }
        // if (!first) first = t;
        // requestAnimationFrame(draw);
        // t = t - first + offset;
        // if (!last_t) last_t = t;
        // const delta_t = t - last_t;
        // if (delta_t < TIME_STEP) return;
        // last_t = t;


        // module.instance.exports.Step(0.01 * t);
        module.instance.exports.Step(t, dt);
        // try {
        //     module.instance.exports.Step(t, delta_t);
        // } catch (e) {
        //     paused = true;
        //     console.error("Step failed", e);
        // }
        // console.log('step', t, delta_t);
        // pre.textContent = decoder.decode(mem);
        pix_data.set(PIX_DATA);
        canvas_output_ctx.putImageData(img_data, 0, 0);
    };
    for (let i = 0; i < LEN_4; i += 4) {
        const r = imgdata1.data[i + 0];
        const g = imgdata1.data[i + 1];
        const b = imgdata1.data[i + 2];
        const a = 255;
        PIX_DATA_COPY[i + 0] = r;
        PIX_DATA_COPY[i + 1] = g;
        PIX_DATA_COPY[i + 2] = b;
        PIX_DATA_COPY[i + 3] = a;
    }
    module.instance.exports.Setup();
    pix_data.set(PIX_DATA);
    canvas_output_ctx.putImageData(img_data, 0, 0);
    {
        const handle_ink_color_change = (color_ink) => {
            const value = color_ink.value.slice(1);
            const r = parseInt(value.slice(0, 2), 16) / 255;
            const g = parseInt(value.slice(2, 4), 16) / 255;
            const b = parseInt(value.slice(4, 6), 16) / 255;
            module.instance.exports.SetINK_COLOR_RGB(r, g, b);
        };

        canvas_output.addEventListener("mousemove", (event) => {
            const { x, y } = getPos(event);
            // LAST_MOUSE_X = MOUSE_X;
            // LAST_MOUSE_Y = MOUSE_Y;
            // MOUSE_X = x;
            // MOUSE_Y = y;
            module.instance.exports.SetMOUSE_XY(x, y);
        });
        canvas_output.addEventListener("mousedown", () => {
            // MOUSE_DOWN = true;
            module.instance.exports.SetMOUSE_DOWN(true);
        });
        document.body.addEventListener("mouseup", () => {
            // MOUSE_DOWN = false;
            module.instance.exports.SetMOUSE_DOWN(false);
        });
        const color_ink = document.querySelector("#color-ink");
        color_ink.addEventListener("change", () =>
            handle_ink_color_change(color_ink)
        );
        handle_ink_color_change(color_ink);
        const button_reset = document.querySelector("#button-reset");
        button_reset.addEventListener("click", () => module.instance.exports.Setup());
        // const button_fetch = document.querySelector("#button-fetch");
        // button_fetch.addEventListener("click", handle_button_fetch);
        // await handle_button_fetch();
    }

    requestAnimationFrame(draw);

    console.log('done');
}

main().catch(console.error);
