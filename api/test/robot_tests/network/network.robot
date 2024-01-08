*** Settings ***
Library             RequestsLibrary
Resource            ../keywords.resource
Resource            ../variables.resource

*** Test Cases ***
GET /network
    ${response}=    GET
    ...    ${API_URL}/network
    ...    expected_status=${STATUS_OK}
    Should Be Equal    ${response.json()['currentNetwork']}    buildnet

GET /massa/node
    ${response}=    GET
    ...    ${API_URL}/massa/node
    ...    expected_status=${STATUS_OK}
    Should Be Equal    ${response.json()['network']}    buildnet

POST /network/mainnet
    ${response}=    POST
    ...    ${API_URL}/network/mainnet
    ...    expected_status=${STATUS_OK}
    Should Be Equal    ${response.json()['currentNetwork']}    mainnet

    # Switch back to buildnet for other tests
    Switch to buildnet
