*** Settings ***
Library     SeleniumLibrary
Resource    variables.resource

*** Keywords ***
Open Thyra Page
    [Arguments]    ${PATH}
    Open Browser
    ...    ${PATH}
    ...    ${BROWSER}
    ...    remote_url=${REMOTE_URL}
    ...    desired_capabilities=${DESIRED_CAPABILITIES}
    ...    options=add_argument ( "--disable-dev-shm-usage" )
