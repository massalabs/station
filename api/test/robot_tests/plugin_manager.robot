*** Settings ***
Documentation       This is a test suite for Thyra Plugin Manager endpoints.

Library             RequestsLibrary
Resource            variables.resource


*** Test Cases ***
GET /plugin-manager
    ${response}=    GET    ${API_URL}/plugin-manager
    Status Should Be    ${STATUS_OK}
    ${listLength}=    Get Length    ${response.json()}
    Should Be Equal As Integers    ${listLength}    1

GET /plugin-manager/invalid
    ${response}=    GET    ${API_URL}/plugin-manager/invalid    expected_status=${STATUS_NOT_FOUND}
    Should Be Equal As Strings    ${response.json()['code']}    Plugin-0001

GET /thyra/plugin/invalid/invalid
    ${response}=    GET    ${API_URL}/thyra/plugin/invalid/invalid    expected_status=${STATUS_NOT_IMPLEMENTED}
    Should Be Equal As Strings    ${response.json()}    operation PluginRouter has not yet been implemented
