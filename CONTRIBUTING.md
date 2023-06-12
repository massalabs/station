# MassaStation Contributing Guide

Thank you for your interest in contributing to MassaStation! We welcome contributions from the community and value the time and effort you put into helping us make MassaStation better.

## Table of Contents
- [Reporting Issues](#reporting-issues)
- [Reviewing Pull Requests](#reviewing-pull-requests)
- [Contributing fixes / features](#contributing-fixes--features)
- [Development](#development)
  - [Getting Started](#getting-started)
  - [Setting Up Development Environment](#setting-up-development-environment)


## Reporting Issues

Before submitting an issue, please do a quick search to check if a similar issue has already been reported. This helps avoid duplicates and allows us to focus on resolving existing problems more efficiently.

If you couldn't find a similar issue, you can submit a new one by following these steps:
1. Click on the "Issues" tab in the project's GitHub repository and then click on the "New Issue" button. Or simply click [here](https://github.com/massalabs/thyra/issues/new/choose) to go directly to the issue creation page.
3. Select the type of issue you are reporting (bug, feature request, etc.).
4. A template will be automatically populated with the required information. Please fill out the template as completely as possible.

We really appreciate your contributions in reporting issues and helping us improve MassaStation!


## Reviewing Pull Requests

Reviewing pull requests is an excellent way to contribute to the project. It helps us ensure that the codebase is well-maintained and that new features and bug fixes are properly tested and documented. If you have the time and skills to review pull requests, please consider doing so.


## Contributing fixes / features

To ensure a smooth collaboration and avoid duplication of efforts, please follow these guidelines:

1. Before starting work on a new feature or bug fix, check the project's GitHub repository for existing issues.
2. If you find an issue you'd like to work on, comment on the issue to express your interest. This helps us track who is working on what and avoid multiple contributors tackling the same problem simultaneously. Additionally, the project maintainers can provide guidance and clarification on the issue, ensuring that everyone has a clear understanding of what needs to be done.
3. Once you have received approval from the project maintainers to work on an issue, you can start working on it. Make sure to fork the repository, create a new branch for your changes, and commit your work in logical and well-documented commits. This makes it easier for the project maintainers to review your code and understand the changes you've made.
4. When you are ready to submit your changes, open a pull request. Provide a clear and concise description of the changes you've made and reference the relevant issue number in your pull request. This helps us track the progress of the issue and ensures that your changes are properly reviewed.
5. The project maintainers and community members may provide feedback or ask questions on your pull request. Engage in the discussion and address any requested changes or concerns promptly. Collaboration and open communication are key to the success of the project.


## Development

### Getting Started

MassaStation is divided in two parts: the backend and the frontend. The backend is written in Go. It is responsible for handling the plugins, the communication with the blockchain, and serves an API and the frontend. It also provides a GUI as an icon in the system tray to interact with the application.
The frontend allows users to access and manage plugins, to upload and browse websites stored on the blockchain, and more thanks to the plugins. It is written in TypeScript and uses React.

To contribute to MassaStation backend, you should have a basic understanding of the Go programming language and Git version control system. If you're new to Go, we recommend checking out [A Tour of Go](https://tour.golang.org/welcome/1) to get started.

To contribute to MassaStation frontend, you should have a basic understanding of the TypeScript programming language and Git version control system. If you're new to TypeScript, we recommend checking out [TypeScript in 5 minutes](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html) to get started.


### Setting Up Development Environment

To contribute to MassaStation, you'll need to set up your development environment. Follow the steps below to get started:

1. **Install Go:** Go is required to build and run MassaStation. You can install Go by following the instructions at [https://golang.org/](https://golang.org/).

2. **Install Node.js and NPM:**
   - **Windows:** Download the Node.js installer from [https://nodejs.org/en/download/](https://nodejs.org/en/download/) and run the installer to install Node.js and NPM.
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

3. **Install Dependencies:**
   - **Ubuntu:** Install the following system dependencies using `apt`:
     ```bash
     sudo apt install -y build-essential libgl1-mesa-dev xorg-dev p7zip
     ```
   - **Windows:**
     - Install `mingw` by following the instructions at [https://www.mingw-w64.org/downloads](https://www.mingw-w64.org/downloads) to provide the necessary `gcc` compiler for building MassaStation.

4. **Install Go Swagger:** MassaStation uses Go Swagger to generate code from the API documentation. Install Go Swagger by running the following command:
   ```bash
   go install github.com/go-swagger/go-swagger/cmd/swagger@latest
   ```

5. **Install Go Stringer:** MassaStation utilizes Go Stringer to generate declarations for enum types. Install Go Stringer by running the following command:
   ```bash
   go install golang.org/x/tools/cmd/stringer@latest
   ```

Once you have completed the above steps, your development environment for MassaStation is set up and ready to go!


### Building MassaStation

To build MassaStation, follow these steps:

1. **Generate Code and Build Front End:** Run the following command to generate code using Go Swagger and build the front end to be served by the API:
     ```bash
     go generate ./...
     ```

   > **_NOTE:_** On Linux, you can add the capability to bind to a port lower than 1024 without the program being executed as root by running the following command:
   >
   > ```bash
   > sudo setcap CAP_NET_BIND_SERVICE=+eip massastation
   > ```

2. **Build the Project:** Once the code generation and front end build are complete, run the following command to build MassaStation:
     ```bash
     go build -o massastation cmd/massastation/main.go
     ```

3. **Running the Project:** Finally, to run MassaStation, execute the `massastation` binary.















___________

# Contributing to Thyra

Thank you for considering contributing to MassaStation! We welcome contributions from the community and value the time and effort you put into making this project better.

## Getting Started

To contribute to MassaStation, you should have a basic understanding of the Go programming language and Git version control system. If you're new to Go, we recommend checking out [A Tour of Go](https://tour.golang.org/welcome/1) to get started.

Before you can start contributing, you'll need to complete the following steps:

- Install the dependencies:
  - Ubuntu like:

```bash
  sudo apt update
  sudo apt install -y build-essential libgl1-mesa-dev xorg-dev p7zip
```

- Install Node.js and NPM:

  - Windows:
    - Download the Node.js installer from <https://nodejs.org/en/download/>.
    - Run the installer and follow the prompts to install Node.js and NPM.

  - macOS:

    - Install Homebrew by following the instructions at <https://brew.sh/>.

    - Run the following command to install Node.js and NPM:

        ```bash
        brew install node
        ```

  - Ubuntu:
    - Run the following command to update the package lists:

        ```bash
        sudo apt update

        ```

    - Run the following command to install Node.js and NPM:

        ```bash
        sudo apt install nodejs npm
        ```

- Install Go: Go is required to build and run MassaStation. You can install Go by following the instructions at [https://go.dev/doc/install](https://go.dev/doc/install).

- Install Swagger: MassaStation uses Swagger to generate code from the API documentation. You can install Swagger by running the following command:

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

- Install Stringer: it genreates declarations for enum types. You can install Stringer by running the following command:

```bash
go install golang.org/x/tools/cmd/stringer@latest
```

You can find more information about Swagger at [https://github.com/go-swagger/go-swagger](https://github.com/go-swagger/go-swagger).

- Build generated files: MassaStation generates code using Stringer and Go Swagger. You can build all generated files by running the following command:

```bash
go generate ./...
```

- Build from source: Once you've completed the above steps, you can build MassaStation from source by running the following commands:

```bash
go build -o massastation cmd/massastation/main.go
```

> **_NOTE:_** On Linux, you can add the possibility to bind to a port lower than 1024 without the program being executed as root by doing `sudo setcap CAP_NET_BIND_SERVICE=+eip massastation`

- Run MassaStation: you can finaly launch MassaStation by executing `massastation` binary.

## MassaStation frontend development

Navigate to <http://my.massa/massastation/index.html> to see MassaStation frontend.

You can run the ReactJS application with vite:

```bash
cd web/massastation
npm run dev
```

## Code Formatting

We take code formatting seriously in MassaStation to maintain a consistent code style. Please follow these guidelines to ensure that your code is properly formatted:

### golangci-lint

We use `golangci-lint` to run linters in parallel. We recommend installing it locally and running it on your source code before pushing any modifications, otherwise some potential lint errors will be caught by the CI pipeline.

To run `golangci-lint` locally:

```bash
golangci-lint run ./...
```

#### How to resolve golangci-lint recurring errors ?

- File is not `gofumpt`

gofumpt need to be installed locally `go install mvdan.cc/gofumpt@latest`

run gofumpt locally on your source code `gofumpt -l -w ./...`

- File is not `gci`

gci need to be installed locally `go install github.com/daixiang0/gci@latest`

run gci locally on your source code `gci --write ./...`

### Frontend code formatting

```bash
cd web/massastation
npm run fmt
```

## Code with auto-reload

You can run the application with this command: `air`.

It will generate and reload MassaStation each time a file is modified.

## Ways to Contribute

There are many ways to contribute to MassaStation, including but not limited to:

- Reporting issues: If you encounter a bug or have a feature request, please open an issue on the [MassaStation issue tracker](https://github.com/massalabs/thyra/issues) with as much detail as possible.

- Fixing issues: If you have the skills and time to fix an issue, please fork the repository, create a branch for your changes, and submit a pull request with your changes. Please make sure to follow the [contribution guidelines](#contribution-guidelines) when submitting your pull request.

- Adding new features: If you have an idea for a new feature or improvement, please open an issue on the [MassaStation issue tracker](https://github.com/massalabs/thyra/issues) to discuss your idea with the community. Once the community agrees on the feature, you can follow the same process as fixing issues to submit your changes.

- Documentation: If you have a knack for writing documentation, please consider contributing to the MassaStation documentation.

- Code reviews: Reviewing pull requests from other contributors is an important way to contribute to the project. Please review pull requests and provide constructive feedback to help improve the codebase.

## Contribution Guidelines

To ensure that your contribution is accepted, please follow these guidelines:

- Follow the Go coding style guide: MassaStation follows the standard [Go coding style guide](https://google.github.io/styleguide/go/guide) and [best practices](https://go.dev/doc/effective_go). Please make sure to follow the same style when submitting your changes.

- Write tests: MassaStation values code quality and reliability. Please make sure to write tests for your changes to ensure that they work as expected.

- Keep your pull request small: If possible, please keep your pull requests small and focused. This makes it easier for the community to review your changes and provide feedback.

- Be respectful: We value diversity and inclusivity in the MassaStation community. Please be respectful and professional in your interactions with others. Any form of harassment or discrimination will not be tolerated.
