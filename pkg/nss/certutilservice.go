package nss

// This file provides a CertUtilService struct that encapsulates operations on the Network Security Services (NSS)
// database using the certutil command. The NSS is a set of libraries designed to support cross-platform development
// of security-enabled client and server applications.
//
// CertUtilService uses a Runner to execute these commands and returns an error if the command execution fails.
//
// Future enhancements:
//  If more functionalities are required to be managed through the CertUtilService, additional methods can be added to
//  the CertUtilService struct.

var _ CertUtilServicer = &CertUtilService{}
