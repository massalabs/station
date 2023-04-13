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

Zip Test Website
    [Documentation]    Zip the test website in ./test_website and save it to the "${TEST_ZIP_FILE}" file

    Log To Console    Zipping test website
    Remove File    ${CURDIR}/${TEST_ZIP_FILE}
    Create Zip From Files In Directory    ${CURDIR}/test_website    ${CURDIR}/${TEST_ZIP_FILE}
    Archive Should Contain File    ${CURDIR}/${TEST_ZIP_FILE}    index.html
