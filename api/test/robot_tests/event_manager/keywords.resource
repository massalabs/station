*** Settings ***
Documentation       Keywords specific to event manager test suite

Library             RequestsLibrary
Library             ../libs/encode_base64.py
Resource            ../keywords.resource


*** Keywords ***
Suite Setup
    Basic Suite Setup
    Get Wallet Address
    Deploy testSC

Deploy testSC
    ${smartContractB64}=  Encode File To Base64    ${CURDIR}/../../testSC/build/testSC.wasm
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    smartContract=${smartContractB64}
    ...    maxCoins=3000000000000
    ...    coins=8500000
    ...    fee=10000000
    ${headers}=    Create Dictionary    Content-Type=application/json
    ${response}=    POST    ${API_URL}/cmd/deploySC    json=${data}    headers=${headers}    expected_status=any
    Log To Console    json response: ${response.json()}    # Print the response content to the test log for debugging

    Should Be Equal As Integers    ${response.status_code}    ${STATUS_OK}    # Assert the status code is 200 OK
    Should Contain    ${response.json()['firstEvent']['data']}    TestSC Constructor called

    ${sc_address}=    Set Variable  ${response.json()['firstEvent']['address']}
    Set Global Variable    ${DEPLOYED_SC_ADDR}    ${sc_address}

Get Wallet Address
    ${response}=    GET
    ...    ${API_URL}/plugin/massa-labs/massa-wallet/api/accounts/${WALLET_NICKNAME}
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
    ...    fee=10000000

    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
