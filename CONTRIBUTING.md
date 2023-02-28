# Contributing to Thyra

Thank you for considering contributing to Thyra! We welcome contributions from the community and value the time and effort you put into making this project better.

## Getting Started

To contribute to Thyra, you should have a basic understanding of the Go programming language and Git version control system. If you're new to Go, we recommend checking out [A Tour of Go](https://tour.golang.org/welcome/1) to get started.

Before you can start contributing, you'll need to complete the following steps:

- Install Go: Go is required to build and run Thyra. You can install Go by following the instructions at [https://go.dev/doc/install](https://go.dev/doc/install).

- Install Swagger: Thyra uses Swagger to generate code from the API documentation. You can install Swagger by running the following command:

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

You can find more information about Swagger at [https://github.com/go-swagger/go-swagger](https://github.com/go-swagger/go-swagger).

- Build generated files: Thyra generates code using Stringer and Go Swagger. You can build all generated files by running the following command:

```bash 
go generate ./web/...
```

- Run from source: Once you've completed the above steps, you can run Thyra from source by running the following commands:

```bash
go build -o thyra-server cmd/thyra-server/main.go
sudo setcap CAP_NET_BIND_SERVICE=+eip thyra-server
./thyra-server
```


## Ways to Contribute

There are many ways to contribute to Thyra, including but not limited to:

- Reporting issues: If you encounter a bug or have a feature request, please open an issue on the [Thyra issue tracker](https://github.com/massalabs/thyra/issues) with as much detail as possible.

- Fixing issues: If you have the skills and time to fix an issue, please fork the repository, create a branch for your changes, and submit a pull request with your changes. Please make sure to follow the [contribution guidelines](#contribution-guidelines) when submitting your pull request.

- Adding new features: If you have an idea for a new feature or improvement, please open an issue on the [Thyra issue tracker](https://github.com/massalabs/thyra/issues) to discuss your idea with the community. Once the community agrees on the feature, you can follow the same process as fixing issues to submit your changes.

- Documentation: If you have a knack for writing documentation, please consider contributing to the Thyra documentation.

- Code reviews: Reviewing pull requests from other contributors is an important way to contribute to the project. Please review pull requests and provide constructive feedback to help improve the codebase.

## Contribution Guidelines

To ensure that your contribution is accepted, please follow these guidelines:

- Follow the Go coding style guide: Thyra follows the standard Go coding style guide. Please make sure to follow the same style when submitting your changes.

- Write tests: Thyra values code quality and reliability. Please make sure to write tests for your changes to ensure that they work as expected.

- Keep your pull request small: If possible, please keep your pull requests small and focused. This makes it easier for the community to review your changes and provide feedback.

- Be respectful: We value diversity and inclusivity in the Thyra community. Please be respectful and professional in your interactions with others. Any form of harassment or discrimination will not be tolerated.
