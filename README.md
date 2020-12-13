# Go-Bank

This a simple bank system written in Golang.

The following commands are avaliable:
* sign
* login
  * info
  * deposit
  * withdraw
  * transfer
  * logout
  * quit
* quit

The program uses two files overall to store the accounts, one is "accounts" and the other "accounts1". One of them will have all the accounts and the other will be empty and wait for every time the file need to be updated and when they do, line by line the program will read the line from the used file and write it to the unused file.

The program also uses https://github.com/mgutz/ansi.git ansi package to color the output.
