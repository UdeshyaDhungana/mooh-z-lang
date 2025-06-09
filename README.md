# The Muji Programming Language

Why Muji?

- Why not?

Muji is a relatively small programming language (implemented in ~3500 loc) that began as my side project.

Since the language is small, here is the syntax:

### Assign a variable
```muji
thoos_muji x = 43;
```

We currently support the following data types:
- Integer
- Floats
- Strings
- Arrays
- Boolean
- Hashmaps

Integers, floats, and strings are defined the usual way. Please note that we do not yet support scientific notation for floats.

### Functions

Defining a function is as easy as
```muji
thoos_muji add = kaam_gar_muji(x, y) {
    patha_muji x + y;
};
add(6, 9);
```

Yes, functions should be assigned to a variable by writing out a function expression. Muji lang supports only one return value per function which is returned using the `patha_muji` keyword.
Since functions are first class citizens, you can easily pass functions into another function.

### Arrays
Arrays are defined the usual way.
```muji
thoos_muji myArr = [69, 420, 666];
```
Muji lang supports all types of data structures in an array like Javascript does (our inspiration)
```muji
thoos_muji myArr = [1, "newstring", 69.69];
```

Arrays can be indexed and modified as you would do in the English version of this language(JS).

### Boolean
The two truth values in Muji lang are:
```muji
thoos_muji yes = sacho_muji;
thoos_muji no = jhut_muji;
```

### Conditionals
Only the vanilla if/if else statement (with blocks) is supported.
```muji
$ No alternative $
yedi_muji(condition) {
    ...
}

$ With alternative $
yedi_muji(condition) {
    ...
} nabhae_muji(condition 2) {
    ...
} nabhae_muji (condition 3) {
    ...
} nabhae_chikne {
    ...
}
```

### Comments
Comments should begin and end with `$` as shown above.

### Loops
Muji lang supports two kinds of loops. The `jaba_samma_muji` loop (your traditional `while` loop), and `ghuma_muji` loop (also called the `for` loop traditionally).

```muji
thoos_muji i = 0;
thoos_muji sum = 0;
jaba_samma_muji(i < 420) {
    sum = sum + i;
    i = i + 1;
}
```

```muji
thoos_muji i = 0;
thoos_muji sum = 0;
ghuma_muji(i = 0; i < 42; i = i + 1) {
    sum = sum + i;
}
```

### Hashmaps
Muji also supports hashmaps. The keys of the hashmap should be strings. Value can be anything.

```muji
thoos_muji country_codes = { "NP": "+977", "IN": "+91" };
```

As with arrays, hashmaps can be indexed.

### Builtins
We support a few builtin functions as of now:

#### `lambai_muji`
Applicable to:
- Arrays
- Hashmaps
- Strings

Returns the length (integer) in each case.

Use case
```muji
thoos_muji x = [1,2,3];
lambai_muji(x);
```

#### `khaad_muji`
Applicable to:
- Arrays

Python's `append` equivalent
```muji
thoos_muji myarr = [1,2,3];
khaad_muji(myarr)
```

#### `udaa_muji`
The first argument should be an array, the second argument is index, and is optional.
If supplied, removes the element at the given index, and returns the removed object. If not supplied, does the same thing to the last element of the array

#### `bhan_muji`
Your `print()` equivalent.
```muji
bhan_muji("Hello, world!");
```

#### `abs`
Applicable to floats and intgers. Calcualtes the absolute value
```muji
thoos_muji x = -69;
thoos_muji y = abs(x);
```

Please check out the `example-programs` to know more.

Please note that the language is in the pre-alpha stage. You may encounter bugs. We encourage you to report any issues you find. 

Please help us grow by contributing! If you wish to add more features, please write corresponding tests as well!

Next step:
- [ ] Support for block level scoping

Build using

```go
make
```

Run the REPL

```bash
./build/muji
```
Or, run the interpreter on a file

```bash
./build/muji example.muji
```
