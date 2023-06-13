*** Settings ***
Documentation       Keywords specific to event manager test suite

Library             RequestsLibrary
Resource            ../keywords.resource


*** Keywords ***
Suite Setup
    Basic Suite Setup
    Get Wallet Address
    Deploy testSC

Deploy testSC
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/../../testSC/build/main-testSC.wasm
    ${data}=    Create Dictionary    walletNickname=${WALLET_NICKNAME}
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST    ${API_URL}/cmd/deploySC    data=${data}    files=${file}    expected_status=${STATUS_OK}
    Should Contain    ${response.json()}    TestSC is deployed at :

    ${sc_address}=    Get SC address    ${response.json()}
    Set Global Variable    ${DEPLOYED_SC_ADDR}    ${sc_address}

Get Wallet Address
    ${response}=    GET
    ...    ${API_URL}/plugin/massalabs/wallet/api/accounts/${WALLET_NICKNAME}
    ...    expected_status=${STATUS_OK}
    ${wallet_addr}=    Get From Dictionary    ${response.json()}    address
    Set Global Variable    ${WALLET_ADDR}    ${wallet_addr}

Generate Event
    [Documentation]    Calls testSC event function with a string passed as an argument
    [Arguments]    ${str}
    ${argument}=    keywords.String To Arg    ${str}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=${DEPLOYED_SC_ADDR}
    ...    args=${argument}
    ...    stream=${False}

    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
