# Blog - A simple GoLang based blogging backend

Example code built to use Gin and SQLite3 with GoLang meant to teach how to use these modules.

This project uses the `crypt/pbkdf2` module, and by extension the [mktoken](https://github.com/greeneg/mktoken) tool to create a file that stores PBKDF2 psuedo-encrypted hashs as token strings. This file should NEVER be exposed to the network directly, as it contains the API tokens used by the eventual web UI for mutative actions in the database.
