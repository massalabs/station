*** Settings ***
Documentation       Test suite for Thyra front ends.
...                 This test suite mainly focuses on making sure the front end endpoints are working as expected.

Library             RequestsLibrary
Library             SeleniumLibrary
Resource            ../variables.resource

Suite Teardown      Close All Browsers


*** Variables ***
${BROWSER}=                 headlessfirefox
${REMOTE_URL}=              ${NONE}
${DESIRED_CAPABILITIES}=    ${NONE}


*** Test Cases ***
GET /thyra/home/{resource}
    Open Thyra Page    ${API_URL}/thyra/home/
    Page Should Contain    Which Plugins
    Page Should Contain    Registry
    Page Should Contain    Web On Chain
    Page Should Contain    Manage plugin


*** Keywords ***
Open Thyra Page
    [Arguments]    ${PATH}
    Open Browser
    ...    ${PATH}
    ...    ${BROWSER}
    ...    remote_url=${REMOTE_URL}
    ...    desired_capabilities=${DESIRED_CAPABILITIES}
    ...    options=add_argument ( "--disable-dev-shm-usage" )
