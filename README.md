## A simple web application with Go backend that takes in data from a form page and shows all the entries in a datatable in another page.

____

# Run
- Create `.env` in the root
- Put `PORT=8092`
- `go run .`

----

- `static` directory contains .html files to serve the views
  - http://localhost:8092/form.html
  - http://localhost:8092/table.html
- `data` contains a json file to store user inputs

