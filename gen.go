package resource

// API swagger
//nolint:lll
//go:generate swagger generate server --quiet --target api/swagger/server --name thyra-server --spec api/swagger/server/restapi/resource/swagger.yml --exclude-main

// Thyra react frontend
//go:generate npm install --prefix web/thyra-frontend
//go:generate npm run --prefix web/thyra-frontend build
