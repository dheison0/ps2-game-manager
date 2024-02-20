#ifndef CRC32_H
#define CRC32_H

typedef unsigned int u32;
u32 crctab[0x400];
u32 crc32(const char *string);

#endif