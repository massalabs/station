*** Settings ***
Documentation       This is a test suite for Thyra /cmd endpoints.

Library             AsyncLibrary
Library             Collections
Library             RequestsLibrary
Resource            keywords.resource
Resource            ../keywords.resource
Resource            ../variables.resource

Suite Setup         Suite Setup


*** Test Cases ***
GET /events/{str}/{caller} from contract
    ${response}=    GET
    ...    ${API_URL}/events/TestSC is deployed at/${WALLET_ADDR}
    ...    expected_status=${STATUS_OK}
    Should Contain    ${response.json()['address']}    ${WALLET_ADDR}
    Should Contain    ${response.json()['data']}    TestSC is deployed at

GET /events/{str}/{caller} from called function
    ${randomID}=    Generate Random String    10
    ${expected_event}=    Set Variable    I'm an event! My id is ${randomID}
    # Since the event manager starts checking for events only after the current Slot and Period,
    # we need to make this request before the event is generated.
    ${handle}=    Async Run    GET
    ...    ${API_URL}/events/${expected_event}/${WALLET_ADDR}
    ...    expected_status=${STATUS_OK}

    Generate Event    ${randomID}

    ${response}=    Async Get    ${handle}

    Should Contain    ${response.json()['address']}    ${WALLET_ADDR}
    Should Contain    ${response.json()['data']}    ${expected_event}
