import numpy as np
import struct

def write_raw_float64_to_file(filename, value):
    if not isinstance(value, float):
        raise TypeError("Value must be a float.")
    with open(filename, 'wb') as f:
        f.write(bytearray(struct.pack('d', value)))

write_raw_float64_to_file('device_11_memory', 1.0)
write_raw_float64_to_file('device_12_memory', 2.0)
