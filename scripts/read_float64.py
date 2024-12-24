import struct

def read_float64_array(filename, num_chunks):
    floats = []
    with open(filename, 'rb') as f:
        for _ in range(num_chunks):
            chunk = f.read(8)
            if not chunk:
                break
            floats.append(struct.unpack('d', chunk)[0])
    return floats

def read_raw_ascii_strings(filename, num_chunks):
    raw_strings = []
    with open(filename, 'rb') as f:
        for _ in range(num_chunks):
            chunk = f.read(8)
            if not chunk:
                break
            raw_strings.append("".join(chr(b) if 32 <= b <= 126 else f"\\x{b:02x}" for b in chunk))
    return raw_strings

# Read float arrays
d1_floats = read_float64_array('device_11_memory', 3)
d2_floats = read_float64_array('device_12_memory', 3)
d3_floats = read_float64_array('device_13_memory', 3)

# Read raw ASCII strings
d1_raw_ascii = read_raw_ascii_strings('device_11_memory', 3)
d2_raw_ascii = read_raw_ascii_strings('device_12_memory', 3)
d3_raw_ascii = read_raw_ascii_strings('device_13_memory', 3)

# Print results
print(f'device 11 floats: {d1_floats}')
print(f'device 12 floats: {d2_floats}')
print(f'device 13 floats: {d3_floats}')

print(f'device 11 raw ascii: {d1_raw_ascii}')
print(f'device 12 raw ascii: {d2_raw_ascii}')
print(f'device 13 raw ascii: {d3_raw_ascii}')

