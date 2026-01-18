#include <stdio.h>
#include <stdint.h>

void __coco_print_int(int64_t value)
{
    printf("%ld\n", value);
}

void __coco_print_float(double value)
{
    printf("%g\n", value);
}

void __coco_print_bool(int64_t value)
{
    if (value == 0)
    {
        printf("false\n");
    }
    else if (value == 1)
    {
        printf("true\n");
    }
    else
    {
        printf("%ld\n", value);
    }
}
