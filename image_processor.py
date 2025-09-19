import base64
import io
import os
import sys
import json
import pyvips
import numpy as np
import subprocess
import tempfile
import xml.etree.ElementTree as et

def resize_image(image, width, height):
    # Use the smaller scaling factor to maintain aspect ratio
    scale_x = width / image.width
    scale_y = height / image.height
    scale = min(scale_x, scale_y)
    # Resize the image using libvips if needed
    if scale < 1:
        return image.resize(scale)
    else:
        return image

def resize_fixed(image, width, height):
    # Use the smaller scaling factor to maintain aspect ratio
    scale_x = width / image.width
    scale_y = height / image.height
    # Resize the image using libvips
    return image.resize(scale_x, vscale=scale_y, gap=8)

def mesh_gradient(image):
    red_levels = np.array([0, 51, 102, 153, 204, 255])
    green_levels = np.array([0, 36, 73, 109, 146, 182, 219, 255])
    blue_levels = np.array([0, 64, 128, 192, 255])

    palette = np.array([(r, g, b) for r in red_levels for g in green_levels for b in blue_levels], dtype=np.uint8)

    array = image.write_to_memory()
    np_array = np.frombuffer(array, dtype=np.uint8)
    np_array = np_array.reshape(image.height, image.width, image.bands)

    indices = []
    for y in range(8):
        for x in range(8):
            pixel = np_array[y, x, :3]
            # Calculate squared Euclidean distance to each palette color
            diffs = palette.astype(np.int32) - pixel.astype(np.int32)
            dists = np.sum(diffs ** 2, axis=1)
            idx = np.argmin(dists)
            indices.append(idx)

    byte_seq = bytes(indices)
    return base64.b64encode(byte_seq).decode('ascii')

def get_duration(input_file):
    cmd = [
        'ffprobe', '-v', 'error', '-show_entries', 'format=duration',
        '-of', 'default=noprint_wrappers=1:nokey=1', input_file
    ]
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        return 0
    try:
        return int(float(result.stdout.strip()))
    except ValueError:
        return 0

def generate_video_frame(input_file, duration, temp_dir):
    seek_time = duration * 0.1
    output_jxl = os.path.join(temp_dir, 'frame.jxl')
    cmd = [
        'ffmpeg', '-i', input_file, '-ss', str(seek_time), '-vframes', '1',
        '-q:v', '100', output_jxl
    ]
    subprocess.run(cmd, check=True)
    return output_jxl

def generate_audio_waveform(input_file, temp_dir):
    output_png = os.path.join(temp_dir, 'waveform.png')
    cmd = [
        'ffmpeg', '-i', input_file, '-f', 'lavfi', '-i', 'color=c=#000000:s=640x120',
        '-filter_complex',
        '[0:a] aformat=channel_layouts=mono,showwavespic=s=640x120:colors=#808080:filter=peak:scale=sqrt [pk]; '
        '[0:a] aformat=channel_layouts=mono,showwavespic=s=640x120:colors=#ffffff:scale=sqrt [rms], '
        '[pk] [rms] overlay=format=auto [nobg], [1:v] [nobg] overlay=format=auto',
        '-frames:v', '1', '-update', 'true', output_png
    ]
    subprocess.run(cmd, check=True)
    return output_png

def generate_video_preview(input_file, output_file):
    # Step 1: Get input framerate using ffprobe
    ffprobe_cmd = [
        'ffprobe', '-v', 'error', '-select_streams', 'v:0',
        '-show_entries', 'stream=avg_frame_rate', '-of',
        'default=noprint_wrappers=1:nokey=1', input_file
    ]
    result = subprocess.run(ffprobe_cmd, capture_output=True, text=True, check=True)
    framerate_str = result.stdout.strip()  # e.g., "30/1"
    if not framerate_str:
        framerate_str = "1/1"

    # Parse framerate to float (e.g., "30/1" -> 30.0)
    num, den = map(int, framerate_str.split('/'))
    input_fps = num / den

    # Step 2: Build FFmpeg command
    preview_file = output_file + '.mkv'
    cmd = [
        'ffmpeg', '-i', input_file, '-c:v', 'libsvtav1', '-preset', '5',
        '-crf', '56', '-profile:v', 'main', '-level:v', '5.1',
        '-c:a', 'libopus', '-b:a', '16k', '-ac', '1', '-vbr', 'on'
    ]
    vf_filters = ['scale=480:-1']
    if input_fps > 30:
        vf_filters.insert(0, 'fps=fps=source_fps/2')
    cmd.extend(['-vf', ','.join(vf_filters)])
    cmd.append(preview_file)

    # Run the command
    subprocess.run(cmd, check=True)

def generate_audio_preview(input_file, output_file):
    preview_file = output_file + '.mkv'
    cmd = [
        'ffmpeg', '-i', input_file, '-c:a', 'libopus', '-b:a', '16k',
        '-ac', '1', '-vbr', 'on', '-vn', preview_file
    ]
    subprocess.run(cmd, check=True)

def is_spherical(image):
    if image.width / image.height != 2:
        return False
    xmp_bytes = None
    try:
        if 'xmp-data' in image.get_fields():
            xmp_bytes = image.get('xmp-data')
        elif 'xmp' in image.get_fields():
            xmp_bytes = image.get('xmp')
    except Exception:
        return False
    if not xmp_bytes:
        return False
    namespaces = {'x': 'adobe:ns:meta/', 'rdf': 'http://www.w3.org/1999/02/22-rdf-syntax-ns#', 'GPano': 'http://ns.google.com/photos/1.0/panorama/'}
    try:
        xmp_str = xmp_bytes.decode('utf-8')
        root = et.fromstring(xmp_str)
        rdf_descr = root.find('.//rdf:Description', namespaces)
        if rdf_descr is not None:
            proj = rdf_descr.attrib.get('{http://ns.google.com/photos/1.0/panorama/}ProjectionType')
            if proj and proj == 'equirectangular':
                return True
    except Exception:
        pass
    return False

def generate_cubemap(input_file, output_file, size):
    cubemap_file = output_file + '.c.jxl'
    cmd = ['kubi', '-l', 'row', '-s', str(size), '--order', '1', '4', '0', '5', '3', '2', input_file, cubemap_file]  # Adjust if kubi args differ
    subprocess.run(cmd, check=True)

def process_image(image_info):
    input_file = image_info['input_file']
    output_file = image_info['output_file']
    orientation = image_info.get('orientation', 0)

    _, ext = os.path.splitext(input_file.lower())

    video_exts = [".3gp", ".flv", ".mov", ".qt", ".m2ts", ".mts", ".divx", ".vob", ".webm", ".mkv", ".mka", ".wmv", ".avi", ".mp4", ".mpg", ".mpeg", ".ps", ".ts", ".rm", ".ogv", ".dv"]
    audio_exts = [".mp3", ".wav", ".opus", ".aac", ".ogg", ".wma", ".m4a", ".flac", ".alac", ".mka"]
    is_video = ext in video_exts
    is_audio = ext in audio_exts and not is_video

    duration = None
    if is_video or is_audio:
        duration = get_duration(input_file)

        with tempfile.TemporaryDirectory() as temp_dir:
            if is_video:
                temp_file = generate_video_frame(input_file, duration, temp_dir)
                generate_video_preview(input_file, output_file)
            else:
                temp_file = generate_audio_waveform(input_file, temp_dir)
                generate_audio_preview(input_file, output_file)

            # Load the generated thumbnail for processing
            image = pyvips.Image.new_from_file(temp_file, memory=True)
    else:
        if ext in (".heif", ".heic"):
            image = pyvips.Image.new_from_file(input_file, memory=True, unlimited=True)
        else:
            image = pyvips.Image.new_from_file(input_file, memory=True)

    if orientation > 1:
        image = image.autorot()

    true_width = image.width
    true_height = image.height

    is_sph = False if is_video or is_audio else is_spherical(image)
    if is_sph:
        generate_cubemap(input_file, output_file, min(true_width / 4, 1024))

    image.jxlsave(output_file + '.o.jxl', Q=75, strip=True, effort=4)

    image = resize_image(image, 2048, 2048)
    image.jxlsave(output_file + '.h.jxl', Q=60, strip=True, effort=5)

    image = resize_image(image, 400, 200)
    image.jxlsave(output_file + '.s.jxl', Q=20, strip=True, effort=5)

    image = resize_fixed(image, 8, 8)
    image = image.colourspace("srgb")

    result = {
        "width": true_width,
        "height": true_height,
        "color": mesh_gradient(image)
    }
    if is_sph:
        result["pano"] = True
    if duration is not None:
        result["duration"] = duration

    return result

if __name__ == "__main__":
    pyvips.voperation.cache_set_max_mem(2048)
    stdin_utf8 = io.TextIOWrapper(sys.stdin.buffer, encoding='utf-8')
    stdin_data = stdin_utf8.read()
    data = json.loads(stdin_data)
    sys.stdout.write(json.dumps(process_image(data)))
