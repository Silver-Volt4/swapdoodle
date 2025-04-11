## Very simple tool to extract BPK1 "archives" used by Swapdoodle.
## Performs no checksum validation.
## Thanks for the help to: https://www.3dbrew.org/wiki/Swapdoodle

import sys

file = sys.argv[1]

def uint32(b: bytes):
    return int.from_bytes(b, "little")

with open(file, "rb") as f:
    if f.read(4).decode() != "BPK1":
        print("Bad header")
        exit()
        
    blocks = uint32(f.read(4))
    print(f"Read {blocks} blocks.")
    
    f.read(0x4 + 0x4 + 0x4 + 0x2c)
    
    for _ in range(blocks):
        offset = uint32(f.read(4))
        size = uint32(f.read(4))
        checksum = uint32(f.read(4))
        name = f.read(8).decode().rstrip('\x00')
        print(f"Block {name}: offset {offset}, size {size}, checksum {checksum}")
        pos = f.tell()
        f.seek(offset)
        with open(f"{name}.bin", "wb") as out:
            out.write(f.read(size))
        f.seek(pos)
        