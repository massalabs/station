*** Settings ***
Documentation       This is a test suite for Massa Station Plugin Store endpoints.

Library             RequestsLibrary
Resource            ../variables.resource
Resource            ../keywords.resource

Suite Setup         Basic Suite Setup


*** Test Cases ***
GET /plugin-store
    ${actual}=    GET    ${API_URL}/plugin-store
    ${actual}=    Set Variable    ${actual.json()}

    ${expected}=    GET     https://massa-station-assets.s3.eu-west-3.amazonaws.com/plugins/plugins.json
    ${expected}=    Set Variable    ${expected.json()}

    Should Be Equal As Strings    ${actual[0]['name']}    ${expected[0]['name']}
    Should Be Equal As Strings    ${actual[0]['description']}    ${expected[0]['description']}
