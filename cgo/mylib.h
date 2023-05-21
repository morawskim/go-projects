#ifndef _MYLIB_H
#define _MYLIB_H

#include <stdio.h>
#include <stdlib.h>

typedef struct address {
    char *street;
    char *city;
} address;

typedef struct person {
    char *firstName;
    char *lastName;
    int age;
    address *address;
} person;

void print_string(char* str);

void print_person(person* person);

#endif
