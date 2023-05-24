# CGO

## GDB

To debug C code we need gdb.
At the moment the Delve debugger does not support C.

To build code call `make build`.
During build neither special flags are set.
In other project (more complex) maybe you should include at least these flags for GCC:
* `-g` - to include debug symbols
* `-O0` - to disable most optimizations

After build call `gdb ./cgo-example`. 
You should see prompt gdb.
Type `break print_person` to set breakpoint to stop execution when
the C function `print_person` is called.
Next type `run` to start program.

> Reading symbols from ./cgo-example...
Loading Go Runtime support.
(gdb) break print_person
Breakpoint 1 at 0x485050: file /home/marcin/projekty/go-projects/cgo/mylib.c, line 12.
(gdb) run
Starting program: /home/marcin/projekty/go-projects/cgo/cgo-example
.......
Thread 1 "cgo-example" hit Breakpoint 1, print_person (person=0xc000076020) at /home/marcin/projekty/go-projects/cgo/mylib.c:12
12          printf("person struct passed from Go\n");
(gdb)

To display source code around current position type `list`.
To display value of variable type `print VARIABLE_NAME`.
To display content of struct type `print *POINTER_TO_STRUCT`. 

For example:
> p person
$9 = (person *) 0xc00013c000
print *person
$10 = {firstName = 0x567c70 "Jan", lastName = 0x567c90 "Kowalski", age = 25, address = 0x567cb0}
