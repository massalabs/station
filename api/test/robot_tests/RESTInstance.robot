*** Settings ***
Library             REST   http://localhost
Library             OperatingSystem
Resource            variables.resource
Documentation       This is a simple example using the RESTInstance library https://github.com/asyrjasalo/RESTinstance

*** Variables ***
${WALLET_NICKNAME}=     testnet

*** Test Cases ***
GET /plugin-manager
    GET             /plugin-manager
    Integer         response status            ${OK}
    Array           response body              minItems=1  maxItems=1  uniqueItems=true
    [Teardown]      Output   response body

POST a Smart Contract
    # Fix "no multipart boundary param in Content-Type"
    ${header}       Create Dictionary   Content-Type=multipart/form-data    boundary=""
    ${sc}=          Get Binary File     ${CURDIR}/websiteDeployer.wasm
    ${body}=        Create Dictionary   walletNickname=${WALLET_NICKNAME}
    ${data}=        Create Dictionary   smartContract=${sc}
    POST            /cmd/deploySC       headers=${header}    body=${body}    data=${data}
    [Teardown]      Output              response body
