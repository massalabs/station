*** Settings ***
Documentation       This is a test suite for Thyra /cmd endpoints.

Library             RequestsLibrary
Resource            variables.resource
Resource            ../variables.resource


*** Test Cases ***
POST a Smart Contract
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/websiteDeployer.wasm
    ${data}=    Create Dictionary    walletNickname=${WALLET_NICKNAME}
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST    ${API_URL}/cmd/deploySC    data=${data}    files=${file}    expected_status=500
    log to console    ${response.json()}
