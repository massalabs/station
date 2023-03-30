# API Tests

This directory contains tests for the API.

## Prerequisites

The tests are written using the [Robot Framework](https://robotframework.org/) and the [Requests Library](http://marketsquare.github.io/robotframework-requests/). To install the dependencies, you need to have [Python 3](https://www.python.org/downloads/) and [Pip](https://pip.pypa.io/en/stable/installation/) installed.

Once you have `Python 3` and `Pip` installed, you can install the others dependencies by running the following command:

```bash
pip install -r requirements.txt
```

## Running the tests

To run the tests, you need to have a running instance of the API. You can run the API locally. See the [CONTRIBUTING](../../CONTRIBUTING.md) for instructions on how to do that.

> ⚠️ Note that some tests might modify or delete from your computer some files such as Wallets and Plugins. Please make sure you have a backup of your files before running the tests. ⚠️

Once you have a running instance of the API, you can run the tests with:

```bash
robot robot_tests
```

To run a specific test suite, you can run the following command:

```bash
robot robot_tests/<scope>/<test_file>.robot
```

To run a specific test case, you can run the following command:

```bash
robot -t <test_case_name> robot_tests/<scope>/<test_file>.robot
```

> To know more about the Robot Framework, you can read the [User Guide](https://robotframework.org/robotframework/latest/RobotFrameworkUserGuide.html).

## Adding new tests

To add new tests, you can simply write your tests in the corresponding `.robot` file. If the corresponding file doesn't exist, create the corresponding directory and then create the corresponding `.robot` file by following the style of the other `.robot` tests files.

You can add new variables in a `variables.resource` file.
If the variable is specific to a given scope, you can add it in the `variables.resource` file of the corresponding scope directory.

> To learn more about the Request Library, you can read the [documentation](https://marketsquare.github.io/robotframework-requests/doc/RequestsLibrary.html).