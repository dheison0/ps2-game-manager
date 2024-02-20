#include "crc32.h"

u32 crc32(const char *string)
{
    int crc, table, count, byte;

    for (table = 0; table < 256; table++)
    {
        crc = table << 24;
        for (count = 8; count > 0; count--)
        {
            if (crc < 0)
                crc = crc << 1;
            else
                crc = (crc << 1) ^ 0x04C11DB7;
        }
        crctab[255 - table] = crc;
    }

    do
    {
        byte = string[count++];
        crc = crctab[byte ^ ((crc >> 24) & 0xFF)] ^ ((crc << 8) & 0xFFFFFF00);
    } while (string[count - 1] != 0);

    return crc;
}