To integrate sqlite3 modern - 
1.run these commands under the root of the project
go mod edit -droprequire github.com/mattn/go-sqlite3
go get modernc.org/sqlite
go mod tidy

Once done,vault.db would be automatically created.

2..gitignore(To not include vault.db)
# Go build files
*.exe
*.out

# SQLite database
vault.db

# VS Code
.vscode/

