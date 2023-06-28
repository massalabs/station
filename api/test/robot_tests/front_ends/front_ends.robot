*** Settings ***
Documentation       Test suite for Massa Station front ends.
...                 This test suite mainly focuses on making sure the front end endpoints are working as expected.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            ../variables.resource
Resource            keywords.resource
Resource            variables.resource

Suite Teardown      Close All Browsers


*** Test Cases ***
GET /
    Open Massa Station Page    ${API_URL}/
    Sleep    5    wait for redirection
    Page Should Contain    Decentralization made easy

GET /web/
    Open Massa Station Page    ${API_URL}/web/
    Sleep    5    wait for redirection
    Page Should Contain    Decentralization made easy

GET /web/index
    Open Massa Station Page    ${API_URL}/web/index
    Sleep    5    wait for redirection
    Page Should Contain    Decentralization made easy
