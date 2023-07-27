*** Settings ***
Documentation       Keywords specific to the Web on Chain test suite

Library             ArchiveLibrary
Library             OperatingSystem
Resource            variables.resource
Resource            ../keywords.resource


*** Keywords ***
Suite Setup
    Basic Suite Setup
    Zip Test Website
    Set Suite Variables

Zip Test Website
    [Documentation]    Zip the test website in ./test_website and save it to the "${TEST_ZIP_FILE}" file

    Log To Console    Zipping test website
    Remove File    ${CURDIR}/${TEST_ZIP_FILE}
    Create Zip From Files In Directory    ${CURDIR}/test_website    ${CURDIR}/${TEST_ZIP_FILE}
    Archive Should Contain File    ${CURDIR}/${TEST_ZIP_FILE}    index.html

Set Suite Variables
    ${random_string}=    Generate Random String    8    [LOWER][NUMBERS]
    ${website_url}=    Set Variable    test${random_string}
    ${website_description}=    Set Variable    description${random_string}
    ${response}=    GET    
    ...    ${API_URL}/massa/node
    ...    expected_status=${STATUS_OK}
    ${current_dns}    Set Variable    ${response.json()['dns']}
    Set Global Variable    ${DNS}    ${current_dns}
    Set Suite Variable    ${website_url}    ${website_url}
    Set Suite Variable    ${website_description}    ${website_description}