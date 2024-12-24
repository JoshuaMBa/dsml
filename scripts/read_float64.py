import struct

def read_raw_float64_from_file(filename):
    with open(filename, 'rb') as f:
        data = f.read(8)
        return struct.unpack('d', data)[0]

d1 = read_raw_float64_from_file('device_11_memory')
d2 = read_raw_float64_from_file('device_12_memory')
d3 = read_raw_float64_from_file('device_13_memory')

print(f'device 11: {d1}')
print(f'device 12: {d2}')
print(f'device 13: {d3}')
