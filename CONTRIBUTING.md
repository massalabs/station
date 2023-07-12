# MassaStation Contributing Guide

Thank you for your interest in contributing to MassaStation! We welcome contributions from the community and value the time and effort you put into helping us make MassaStation better.

## Table of Contents

- [MassaStation Contributing Guide](#massastation-contributing-guide)
  - [Table of Contents](#table-of-contents)
  - [Reporting Issues](#reporting-issues)
  - [Suggesting Features](#suggesting-features)
  - [Reviewing Pull Requests](#reviewing-pull-requests)
  - [Contributing to fixes \& features](#contributing-to-fixes--features)
  - [Development](#development)
    - [Getting Started](#getting-started)
    - [Setting Up Development Environment](#setting-up-development-environment)
    - [Building MassaStation](#building-massastation)
    - [Importing API Specifications into Postman](#importing-api-specifications-into-postman)
    - [Linting and Formatting](#linting-and-formatting)

## Reporting Issues

Before submitting an issue, please do a quick search to check if a similar issue has already been reported. This helps to avoid duplicates and allows us to focus on resolving existing problems more efficiently.

If you couldn't find a similar issue, you can submit a new one by following these steps:

1. Click [here](https://github.com/massalabs/station/issues/new/choose) to go directly to the issue creation page.
2. Select the "Bug Report" issue type.
3. A template will be automatically populated with the required information. Please fill out the template as completely as possible.

We really appreciate your contributions in reporting issues and helping us improve MassaStation!

## Suggesting Features

If you have an idea for a new feature or enhancement, we'd love to hear about it! To suggest a new feature, follow these steps:

1. First, check if a similar feature request has already been submitted by searching the [issues](https://github.com/massalabs/station/issues) page. If you find a similar feature request, you can upvote it using the 👍 reaction.

If you couldn't find a similar feature request, you can submit a new one by following these steps:

2. Click [here](https://github.com/massalabs/station/issues/new/choose) to go directly to the issue creation page.
3. Select the "Task" issue type.
4. A template will be automatically populated with the required information. Please fill out the template as completely as possible.

We really appreciate your contributions in suggesting new features and making MassaStation better!

## Reviewing Pull Requests

Reviewing pull requests is an excellent way to contribute to the project. It helps us ensure that the codebase is well-maintained and that new features and bug fixes are properly tested and documented. If you have the time and skills to review pull requests, please consider doing so.

You can find a list of open pull requests [here](https://github.com/massalabs/station/pulls).

Thank you for your contributions in reviewing pull requests and helping us maintain a high-quality codebase!

## Contributing to fixes & features

To ensure a smooth collaboration and avoid duplication of efforts, please follow these guidelines:

1. Before starting work on a new feature or bug fix, check the project's GitHub repository for existing issues. You can start with some [Good First Issues](https://github.com/massalabs/station/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22), which are issues that are relatively easy to fix and are a good starting point for new contributors.
2. If you find an issue you'd like to work on, comment on the issue to express your interest. This helps us track who is working on what and avoid multiple contributors tackling the same problem simultaneously. Additionally, the project maintainers can provide guidance and clarification on the issue, ensuring that everyone has a clear understanding of what needs to be done.
3. Once you have received approval from the project maintainers to work on an issue, you can start working on it. Make sure to fork the repository, create a new branch for your changes, and commit your work in logical and well-documented commits. This makes it easier for the project maintainers to review your code and understand the changes you've made.
4. When you are ready to submit your changes, open a pull request. Provide a clear and concise description of the changes you've made and reference the relevant issue number in your pull request. This helps us track the progress of the issue and ensures that your changes are properly reviewed.
5. The project maintainers and community members may provide feedback or ask questions on your pull request. Engage in the discussion and address any requested changes or concerns promptly. Collaboration and open communication are key to the success of the project.

We really appreciate your contributions in fixing bugs and adding new features to MassaStation !

## Development

### Getting Started

MassaStation is divided in two parts: the backend and the frontend. The backend is written in Go. It is responsible for handling modules (aka "plugins"), the communication with the blockchain, and serves an API and the frontend. It also provides a GUI as an icon in the system tray to interact with the application.
The frontend allows users to access and manage modules, to upload and browse websites stored on the blockchain, and more thanks to the modules. It is written in TypeScript and uses React.

To contribute to MassaStation backend, you should have a basic understanding of the Go programming language and Git version control system. If you're new to Go, we recommend checking out [A Tour of Go](https://tour.golang.org/welcome/1) to get started.

To contribute to MassaStation frontend, you should have a basic understanding of the TypeScript programming language and Git version control system. If you're new to TypeScript, we recommend checking out [TypeScript in 5 minutes](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html) to get started.

### Setting Up Development Environment

To contribute to MassaStation, you'll need to set up your development environment. Follow the steps below to get started:

1. **Install Go:** Go is required to build and run MassaStation. You can install Go by following the instructions from [Go installation instructions](https://go.dev/doc/install).

2. **Install Node.js and NPM:**
   - **Windows:** Download the Node.js installer from [nodejs.org/](https://nodejs.org/en/download/) and run the installer to install Node.js and NPM.
   - **macOS:**
     - Install Homebrew by following the instructions at [https://brew.sh/](https://brew.sh/).
     - Run the following command to install Node.js and NPM:

       ```bash
       brew install node
       ```

   - **Ubuntu:**
     - Run the following command to update the package lists:

       ```bash
       sudo apt update
       ```

     - Run the following command to install Node.js and NPM:

       ```bash
       sudo apt install nodejs npm
       ```

3. Install Task
Follow the installation instructions here: [task-install](https://taskfile.dev/installation/).

On Windows, we recommend to run `go install github.com/go-task/task/v3/cmd/task@latest` and yo use task commands in a git bash terminal.

4. **Install Dependencies:**
   - **Ubuntu:** Install the following system dependencies using `apt`:

     ```bash
     sudo apt install -y build-essential libgl1-mesa-dev xorg-dev p7zip
     ```

   - **Windows:**
     - Install `mingw` by following the instructions from [mingw-w64.org](https://www.mingw-w64.org/downloads) to provide the necessary `gcc` compiler for building MassaStation.

5. **Install Go dependencies :** MassaStation uses Go Swagger to generate code from the API documentation and Go Stringer to generate declarations for enum types.

   ```bash
   task install
   ```

Once you have completed the above steps, your development environment for MassaStation is set up and ready to go!

### Building MassaStation

To build MassaStation, follow these steps:

1. **Generate Code and Build Front End:** Run the following command to generate code using Go Swagger and build the front end to be served by the API:

     ```bash
     task generate
     ```

2. **Build the Project:** Once the code generation and front end build are complete, run the following command to build MassaStation:

     ```bash
     task build
     ```

3. **Running the Project:** Finally, to run MassaStation, execute the `massastation` binary. Navigate to <http://station.massa> to see MassaStation frontend.

### Importing API Specifications into Postman

To use the MassaStation endpoints in Postman, import the API specifications Swagger file by following the steps outlined in the [Postman Import document](docs/postman_swagger.md).



### Linting and Formatting

To ensure consistent code style and maintain code quality, MassaStation follows specific linting and formatting guidelines. Please follow these guidelines when contributing to the project:

For the Go back-end code, we use `golangci-lint` to perform linting and ensure code quality. Follow these steps:

1. **Install `golangci-lint`:**
   - Follow the installation instructions provided from [here](https://golangci-lint.run/usage/install/#local-installation) to install `golangci-lint` locally on your system.

2. **format and lint your code**

    ```bash
    task fmt
    ```

Following these linting and formatting guidelines will ensure a consistent code style and maintain the overall code quality of MassaStation.
