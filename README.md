# Budget Categorizer

I created this project for automating some tasks I was doing manually when taking care of my personal finances.  
Long story short, this is an ETL for tracking my expenses.  
The project consists of two lambda functions : [budgetcategorizer](https://github.com/jbleduigou/budgetcategorizer) and [budget2sheets](https://github.com/jbleduigou/budget2sheets)  
The transform part is performed by **budgetcategorizer** and has two main responsibilities : sanitize the expense description and assign a category to the expense.  
The load part is performed by **budget2sheets** which is going to upload the transcations (i.e. expenses) to Google Sheets.

## Overall Architecture
![Architecture Diagram](architecture_diagram.png)

## Getting Started

Clone the repo inside the following directory:
```
~/go/src/github.com/jbleduigou/

```
If you want to fork the repo, replace the latest path element with your GitHub handle.

### Prerequisites
You will need to have Go installed on your machine.  
Currently the project uses version 1.13

### Building
You will find a Makefile at the root of the project.
To run the full build and have zip file ready for AWS Lambda use:
```
make zip
```
If you only want to run the unit tests:
```
make test
```

## Improvements / Remaining Work

* extract logic to dedicated classes and write unit tests
* improve error handling
* fix weird behaviour for GitHub Actions upload artifact : https://github.com/actions/upload-artifact/issues/39 ?

## License
Licensed under the Apache License, Version 2.0.
Copyright 2019 Jean-Baptiste Le Duigou
