# factorial

This is a GO version of the Parallel Prime Swing method of calculating factorials.
Makes use of all the cores your CPU provides. The more cores utilized the faster the calculation.
After the factorial is calculated the code will dump the actual numeric result string into a text file.

I've had this code calculate the factorial of 100,000,000 and it takes about 37 minutes total. Surprisingly, the bulk of the time is used during string length calculation and write to file operations.
