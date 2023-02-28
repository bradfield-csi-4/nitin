# Physical Storage Exercise Solution

This package contains an implementation of a toy "Bradfield" file format. Bradfield files store row-oriented tuples and each file represents a variable number of rows that fit a predefined "schema" (ordered list of column names). Tuple values can only be strings and the format is designed to be written and read in a stream-oriented fashion, one tuple at a time.

The Bradfield file format begins with a variable sized integer that stores the length of the file "header". The header is stored as marshaled JSON and contains various metadata about the file, like its version number, the number of rows, and the "schema" (list of column names). While JSON may seem like an odd choice to embed in a custom file format (mostly due to poor performance), its a reasonable choice for storing the header information because it's very easy to add or remove fields, and it makes it easy to parse the header information in any programming language which makes inspecting the files in a variety of environments simple. In addition, the poor performance of JSON is a non-issue because the cost of parsing the JSON only needs to be paid once per file, so as long as we store a reasonable number of rows in each file, parsing the JSON will have no impact on our applications performance.

After the Bradfield file header, the file will contain a repeating sequence of tuples (rows) where each tuple consists of a repeating sequence of length-prefixed (variable sized integers) strings where the number of strings per tuple is equal to the number of columns specified in the header. This is not particularly efficient from a storage size or efficiency perspective, but it is good enough for our purposes.

(See `format.png` for a graphical description of the file format.)

## Instructions

This repository is built entirely using the Go standard library so you should be able to run the tests simply by cloning the repository and then running the tests in their editor or from the command line using the Go toolchain. If you're not familiar with Go and its toolchain, we recommend that you install visual studio code and the associated Go extension. This should allow you to run all of the tests in this repository by simply opening the test files and clicking "Run tests".
