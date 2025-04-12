## Very simple tool to extract BPK1 "archives" used by Swapdoodle.
## Performs no checksum validation.
## Thanks for the help to: https://www.3dbrew.org/wiki/Swapdoodle

import sys
import os

RECURSE_ARG = "--recurse"

def uint32(b: bytes):
    return int.from_bytes(b, "little")

def unpack(file: str):
    with open(file, "rb") as f:
        if f.read(4) != b"BPK1":
            print(f"Could not unpack {file}: Bad header")
            return []

        destination = os.path.splitext(file)[0]
        os.mkdir(destination)
            
        blocks = uint32(f.read(4))
        print(f"Read {blocks} blocks.")
        
        f.read(0x4 + 0x4 + 0x4 + 0x2c)

        block_names = {}
        out_file_names = []
        
        for _ in range(blocks):
            offset = uint32(f.read(4))
            size = uint32(f.read(4))
            checksum = uint32(f.read(4))
            name = f.read(8).decode().rstrip('\x00')
            print(f"Block {name}: offset {offset}, size {size}, checksum {checksum}")
            pos = f.tell()
            f.seek(offset)

            name_occurence = block_names.get(name, -1) + 1
            block_names[name] = name_occurence
            out_file_name = os.path.join(destination, f"{name}${name_occurence}.bin")
            out_file_names.append(out_file_name)

            with open(out_file_name, "wb") as out:
                out.write(f.read(size))
            f.seek(pos)
        
        return out_file_names

def recurse_unpack(file: str):
    files = unpack(file)
    for file in files:
        recurse_unpack(file)

if __name__ == "__main__":
    args = sys.argv[1:]
    unpack_fn = unpack
    if RECURSE_ARG in args:
        args.remove(RECURSE_ARG)
        unpack_fn = recurse_unpack
    for file in args:
        unpack_fn(file)