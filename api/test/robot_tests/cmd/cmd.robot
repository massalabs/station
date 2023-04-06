*** Settings ***
Documentation       This is a test suite for Thyra /cmd endpoints.

Library             RequestsLibrary
Resource            keywords.resource
Resource            ../variables.resource

Suite Setup         Suite Setup


*** Test Cases ***
POST a Smart Contract
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/main-websiteDeployer.wasm
    ${data}=    Create Dictionary    walletNickname=${WALLET_NICKNAME}
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST    ${API_URL}/cmd/deploySC    data=${data}    files=${file}    expected_status=${STATUS_OK}
    Should Contain    ${response.json()}    Website Deployer is deployed at :
