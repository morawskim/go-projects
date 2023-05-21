#include "mylib.h"

void print_string(char* str) {
    printf("string passed from Go: %s\n", str);
}

void print_person(person* person) {
    printf("person struct passed from Go\n");
    printf("Name: %s\n", person->firstName);
    printf("Surname: %s\n", person->lastName);
    printf("Age: %d\n", person->age);
    printf("Street: %s\n", person->address->street);
    printf("City: %s\n", person->address->city);
}
