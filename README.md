# Budget Categorizer
![Go](https://github.com/jbleduigou/budgetcategorizer/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbleduigou/budgetcategorizer)](https://goreportcard.com/report/github.com/jbleduigou/budgetcategorizer)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=jbleduigou/budgetcategorizer)](https://dependabot.com)

I created this project for automating some tasks I was doing manually when taking care of my personal finances.  
Long story short, this is an ETL for tracking my expenses.  
The project consists of two lambda functions : [budgetcategorizer](https://github.com/jbleduigou/budgetcategorizer) and [budget2sheets](https://github.com/jbleduigou/budget2sheets)  
The transform part is performed by **budgetcategorizer** and has two main responsibilities : sanitize the expense description and assign a category to the expense.  
The load part is performed by **budget2sheets** which is going to upload the transactions (i.e. expenses) to Google Sheets.

## Overall Architecture

![Architecture Diagram](docs/architecture_diagram.png)

## Getting Started

Clone the repo inside the following directory:

```bash
~/go/src/github.com/jbleduigou/

```

If you want to fork the repo, replace the latest path element with your GitHub handle.

### Prerequisites

You will need to have Go installed on your machine.  
Currently the project uses version 1.22

### Building

You will find a Makefile at the root of the project.  
To run the full build and have zip file ready for AWS Lambda use:

```bash
make zip
```

If you only want to run the unit tests:

```bash
make test
```

## Deployment

For now deployment is made manually.  
It would be nice to have a cloudformation template at some point.

## Improvements / Remaining Work

* extract logic to dedicated classes and write unit tests
* improve error handling
* create cloud formation template
* fix weird behaviour for GitHub Actions upload artifact : https://github.com/actions/upload-artifact/issues/39 ?

## Configuration

Configuration in this application is based on the concepts exposed by [The Twelve Factor App](https://12factor.net/).  
The idea is to strictly separate config from code by using environment variables.  
Please read the page on [12 Factor Configuration](https://12factor.net/config) for more details.

### Environment Variables

The following environment variables should be declared within you lambda:

| Name                         | Description   | Sample Value  |
| ---------------------------- |:-------------:| :-----:|
| SQS_QUEUE_URL                | URL of the SQS queue where transactions are being pushed | https://sqs.eu-west-3.amazonaws.com/6698939/transactions |
| CONFIGURATION_FILE_BUCKET    | Name of the S3 bucket where the configuration file is stored      |   budgetcategorizer |
| CONFIGURATION_FILE_OBJECT_KEY| Object key of the configuration file  |    configuration.yml |

### Configuration File

It might seem like double duty to use both environment variables and a configuration file.  
However the configuration of categories and keywords can potentially be fairly complex.  
Because of that I decided to store this information in a dedicated configuration file.  

The configuration file is YAML formatted and should look like:  

```yaml
categories:
  - Courses Alimentation
  - Loyer

keywords:
  Express Proxi Saint Thonan: Courses Alimentation
  Agence Immo: Loyer
```

The first block declares all the categories.  
The second block declares a list of key/value pairs, associating a keyword with a category.  
Obviously you can have more than one keyword for a given category.

## Project Structure

The project structure was inspired by the project [Go DDD](https://github.com/marcusolsson/goddd).  
The entry point is located in folder cmd/budgetcategorizer.  
What it does is instantiating all the dependencies for the command.  
The business logic was separated by concerns and placed in dedicated folders.  
Interfaces were introduced to avoid tight coupling and therefore facilitate unit testing (amongst other benefits).  

### Data Models

The main data model is located at root of project in the file transaction.go  

```yaml
{
  "Date": "18/12/2019", // Date of transaction
  "Description": "Mmmh un donut!", // Description from bank statement
  "Comment": "", // Not used for now
  "Category": "Courses Alimentation", // The category assigned by budgetcategorizer
  "Value": 3.18 // Transaction amount, can be negative in case of refund for instance
}
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Jean-Baptiste Le Duigou** - *Initial work* - [jbleduigou](https://github.com/jbleduigou)
* **Cyrille Hemidy** - *Documentation improvement and some code refactoring* - [chemidy](https://github.com/chemidy)

See also the list of [contributors](https://github.com/jbleduigou/budgetcategorizer/contributors) who participated in this project.

## License

Licensed under the Apache License, Version 2.0.  
See [LICENSE.txt](LICENSE.txt) for more details.  
Copyright 2020 Jean-Baptiste Le Duigou
