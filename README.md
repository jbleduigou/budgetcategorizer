# budgetcategorizer

This is a sample AWS Lambda function for transforming transactions for budgeting.
I was sick of having to do it manually :)


## improvements

* extract logic to dedicated classes and write unit tests
* improve error handling
* write to SQS instead of exporting a CSV to S3
* fix weird behaviour for GitHub Actions upload artifact : https://github.com/actions/upload-artifact/issues/39 ?

## License
Licensed under the Apache License, Version 2.0.
Copyright 2019 Jean-Baptiste Le Duigou
