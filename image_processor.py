import io
import os
import sys
import json
import pyvips
import numpy as np


def resize_image(image, width, height):
    # Use the smaller scaling factor to maintain aspect ratio
    scale_x = width / image.width
    scale_y = height / image.height

    # Resize the image using libvips
    return image.resize(min(scale_x, scale_y))


def resize_fixed(image, width, height):
    # Use the smaller scaling factor to maintain aspect ratio
    scale_x = width / image.width
    scale_y = height / image.height

    # Resize the image using libvips
    return image.resize(scale_x, vscale=scale_y, gap=8)


def mesh_gradient(image):
    glyphs = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
    text = ""

    array = image.write_to_memory()
    np_array = np.frombuffer(array, dtype=np.uint8)

    # sys.stderr.write("np size: " + str(np_array.size) + "\n")
    # sys.stderr.flush()

    np_array = np_array.reshape(image.height, image.width, image.bands)

    # sys.stderr.write("np size: " + str(np_array.size) + "\n")
    # sys.stderr.flush()

    for y in range(4):
        for x in range(4):
            # sys.stderr.write(str(x) + ":" + str(y) + " - ")
            # sys.stderr.flush()
            r, g, b = np_array[y, x, :3]
            # sys.stderr.write(str(r) + "/" + str(g) + "/" + str(b) + "\n")
            # sys.stderr.flush()
            color_val = ((r & 0xF0) << 4) | (g & 0xF0) | ((b & 0xF0) >> 4)
            text += glyphs[color_val >> 6]
            text += glyphs[color_val & 63]

    # sys.stderr.write(text + "\n")
    # sys.stderr.flush()
    return text


def process_image(image_info):
    # sys.stderr.write("py - load file\n")
    if image_info['input_file'].lower().endswith((".heif", ".heic")):
        image = pyvips.Image.new_from_file(image_info['input_file'], memory=True, unlimited=True)
    else:
        image = pyvips.Image.new_from_file(image_info['input_file'], memory=True)

    # sys.stderr.write("py - autorot\n")
    if image_info['orientation'] > 1:
        image = image.autorot()

    true_width = image.width
    true_height = image.height

    # sys.stderr.write("py - save o\n")
    # sys.stderr.flush()
    image.jxlsave(image_info['output_file'] + '.o.jxl', Q=75, strip=True, effort=4)
    # sys.stderr.write("py - u\n")
    # sys.stderr.flush()
    os.utime(image_info['output_file'] + ".o.jxl", (image_info['mod_time'], image_info['mod_time']))
    # sys.stderr.write("py - resize\n")
    # sys.stderr.flush()
    image = resize_image(image, 2048, 2048)
    # sys.stderr.write("py - save h\n")
    # sys.stderr.flush()
    image.jxlsave(image_info['output_file'] + '.h.jxl', Q=60, strip=True, effort=5)
    # sys.stderr.write("py - u\n")
    # sys.stderr.flush()
    os.utime(image_info['output_file'] + ".h.jxl", (image_info['mod_time'], image_info['mod_time']))
    # sys.stderr.write("py - resize\n")
    # sys.stderr.flush()
    image = resize_image(image, 400, 200)
    # sys.stderr.write("py - save s\n")
    # sys.stderr.flush()
    image.jxlsave(image_info['output_file'] + '.s.jxl', Q=20, strip=True, effort=5)
    # sys.stderr.write("py - u\n")
    # sys.stderr.flush()
    os.utime(image_info['output_file'] + ".s.jxl", (image_info['mod_time'], image_info['mod_time']))
    # sys.stderr.write("py - resize\n")
    # sys.stderr.flush()
    image = resize_fixed(image, 4, 4)
    # sys.stderr.write("py - srgb\n")
    # sys.stderr.flush()
    image = image.colourspace("srgb")

    # sys.stderr.write("py - json\n")
    # sys.stderr.flush()
    return {
        "width": true_width,
        "height": true_height,
        "color": mesh_gradient(image)
    }


if __name__ == "__main__":
    pyvips.voperation.cache_set_max_mem(2048)
    stdin_utf8 = io.TextIOWrapper(sys.stdin.buffer, encoding='utf-8')
    stdin_data = stdin_utf8.read()
    data = json.loads(stdin_data)
    sys.stdout.write(json.dumps(process_image(data)))