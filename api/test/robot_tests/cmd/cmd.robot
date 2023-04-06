*** Settings ***
Documentation       This is a test suite for Thyra /cmd endpoints.

Library             RequestsLibrary
Library             Collections
Resource            keywords.resource
Resource            ../keywords.resource
Resource            ../variables.resource

Suite Setup         Suite Setup


*** Test Cases ***
POST a Smart Contract
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/../../testSC/build/main-testSC.wasm
    ${data}=    Create Dictionary    walletNickname=${WALLET_NICKNAME}
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST    ${API_URL}/cmd/deploySC    data=${data}    files=${file}    expected_status=${STATUS_OK}
    Should Contain    ${response.json()}    TestSC is deployed at :

    ${sc_address}=    Get SC address    ${response.json()}
    Set Global Variable    ${DEPLOYED_SC_ADDR}    ${sc_address}

POST /cmd/executeFunction
    ${randomID}=    Generate Random String    10
    ${argument}=    keywords.String To Arg    ${randomID}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=${DEPLOYED_SC_ADDR}
    ...    args=${argument}
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
    Log To Console    ${response.json()}
    Should Be Equal    ${response.json()}    I'm an event! My id is ${randomID}
