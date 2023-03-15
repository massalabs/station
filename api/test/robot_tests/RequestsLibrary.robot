*** Settings ***
Library               RequestsLibrary
Resource              variables.resource
Documentation         This is a simple example using the RequestsLibrary https://github.com/MarketSquare/robotframework-requests

*** Variables ***
${URL}=                 http://localhost
${WALLET_NICKNAME}=     testnet


*** Test Cases ***
GET /plugin-manager
    ${response}=    GET  ${URL}/plugin-manager
    Status Should Be  ${OK}
    ${listLength}=  Get Length  ${response.json()}
    Should Be Equal As Integers  ${listLength}  1
    log to console         ${response.json()}

POST a Smart Contract
    ${sc}=          Get File For Streaming Upload      ${CURDIR}/websiteDeployer.wasm
    ${data}=        Create Dictionary   walletNickname=${WALLET_NICKNAME}
    ${file}=        Create Dictionary   smartContract=${sc}
    ${response}=    POST                ${URL}/cmd/deploySC    data=${data}    files=${file}
    log to console         ${response.json()}
