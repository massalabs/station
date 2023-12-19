*** Settings ***
Documentation       This is a test suite for Massa Station /cmd endpoints.

Library             RequestsLibrary
Library             Collections
Resource            keywords.resource
Resource            ../keywords.resource
Resource            ../variables.resource

Suite Setup         Suite Setup


*** Test Cases ***
POST a Smart Contract
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/../../testSC/build/main-testSC.wasm
    ${data}=    Create Dictionary
    ...    walletNickname=${WALLET_NICKNAME}
    ...    coins=3000000000
    ...    fee=1000
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST    ${API_URL}/cmd/deploySC    data=${data}    files=${file}    expected_status=any
    Log To Console    json response: ${response.json()}    # Print the response content to the test log for debugging

    Should Be Equal As Integers    ${response.status_code}    ${STATUS_OK}    # Assert the status code is 200 OK
    Should Contain    ${response.json()['firstEvent']['data']}    TestSC is deployed at

    ${sc_address}=    Get SC address    ${response.json()['firstEvent']['data']}
    Set Global Variable    ${DEPLOYED_SC_ADDR}    ${sc_address}

POST /cmd/executeFunction sync
    ${randomID}=    Generate Random String    10
    ${argument}=    keywords.String To Arg    ${randomID}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=${DEPLOYED_SC_ADDR}
    ...    args=${argument}
    ...    fee=1000
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=any
    Log To Console    ${response.json()}
    Should Be Equal As Integers    ${response.status_code}    ${STATUS_OK}    # Assert the status code is 200 OK
    Should Be Equal    ${response.json()['firstEvent']['data']}    I'm an event! My id is ${randomID}

POST /cmd/executeFunction sync without arguments
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=test
    ...    at=${DEPLOYED_SC_ADDR}
    ...    args=
    ...    fee=1000
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
    Log To Console    ${response.json()}
    Should Be Equal    ${response.json()['firstEvent']['data']}    TestSC test() called

POST /cmd/executeFunction async
    ${randomID}=    Generate Random String    10
    ${argument}=    keywords.String To Arg    ${randomID}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=${DEPLOYED_SC_ADDR}
    ...    args=${argument}
    ...    async=${True}
     ...    fee=1000
  ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_OK}
    Log To Console    ${response.json()}
    Should Be Equal
    ...    ${response.json()['firstEvent']['data']}
    ...    Function called successfully but did not wait for event


# Error cases

POST /cmd/deploySC with invalid datastore
    ${sc}=    Get File For Streaming Upload    ${CURDIR}/../../testSC/build/main-testSC.wasm
    ${data}=    Create Dictionary    walletNickname=${WALLET_NICKNAME}    datastore=invalid    fee=1000
    ${file}=    Create Dictionary    smartContract=${sc}
    ${response}=    POST
    ...    ${API_URL}/cmd/deploySC
    ...    data=${data}
    ...    files=${file}
    ...    expected_status=${STATUS_BAD_REQUEST}
    Should Be Equal    ${response.json()["message"]}    illegal base64 data at input byte 4

POST /cmd/executeFunction with invalid address
    ${randomID}=    Generate Random String    10
    ${argument}=    keywords.String To Arg    ${randomID}
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=invalid
    ...    args=${argument}
    ...    fee=1000
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_INTERNAL_SERVER_ERROR}
    Should Be Equal    ${response.json()["code"]}    Execute-0001
    Should Contain
    ...    ${response.json()["message"]}
    ...    Error: callSC failed: calling EstimateGasCost for function 
    Should Contain
    ...    ${response.json()["message"]}
    ...    receiving execute_read_only_call with
    Should Contain
    ...    ${response.json()["message"]}
    ...    Invalid params

POST /cmd/executeFunction with invalid arguments
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=event
    ...    at=AS12YBWcNcmN8wugh8xTZiyt48JjHqrNtem96jiCoGEZFGZPUyei6
    ...    args=invalid
    ...    fee=1000
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_UNPROCESSABLE_ENTITY}
    Should Be Equal    ${response.json()["message"]}    illegal base64 data at input byte 4

POST /cmd/executeFunction with invalid function name
    [Documentation]    Checks that we receive error messages from the node
    ${data}=    Create Dictionary
    ...    nickname=${WALLET_NICKNAME}
    ...    name=invalid
    ...    at=${DEPLOYED_SC_ADDR}
    ...    fee=1000
    ${response}=    POST
    ...    ${API_URL}/cmd/executeFunction
    ...    json=${data}
    ...    expected_status=${STATUS_INTERNAL_SERVER_ERROR}
    # Should Contain must be divided into multiple lines because the error message contains unknown values (e.g. operation id)
    Should Contain
    ...    ${response.json()["message"]}
    ...    Error: callSC failed
    Should Contain
    ...    ${response.json()["message"]}
    ...    ReadOnlyCall error: readonly call failed
    Should Contain
    ...    ${response.json()["message"]}
    ...    VM execution error: Missing export invalid
