// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: html\front\website.css
// generated:  Mon Sep 12 11:55:24 CEST 2022

package website

const CSS = `body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
    "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #fff;

  background: linear-gradient(90deg, #172230 50.03%, #0c1219 85.82%);
}

a {
  color: white;
}

a:hover {
  color: rgb(164, 164, 164);
}

.massa-logo-banner {
  width: 170px;
}

.massa-logo-spinner {
  width: 25px;
}
.spinner-border {
  color: transparent;
}

.loading {
  display: none;
}

.title-wallet {
  margin-bottom: 45px;
}
option {
  background-color: transparent;
}
select.form-select:focus {
  background-color: black;
  color: white;
}

.form-control,
.form-control:focus {
  color: white;
  border-radius: 0;
  border: none;
  border-bottom: 2px solid hsla(0, 0%, 100%, 0.2);
  background-color: transparent;
  box-shadow: none;
}
.form-control:hover {
  cursor: pointer;
}

.form-select {
  background-color: transparent;
  color: white;
}

h1 {
  font-size: 3.5rem;
  color: hsla(0, 0%, 100%, 0.2);
}
h2 {
  color: hsla(0, 0%, 100%, 0.2);
}

button {
  border: none;
  padding: 3px 30px;
  border-radius: 100px;
  color: #fff;
}

.primary-button {
  font-size: 1.3rem;
  background-color: #e74e4e;
  border: 3px solid #e74e4e;
}

.primary-button:hover {
  background-color: #501313;
  border: 3px solid #501313;
  transform: scale(1.05);
  transition: 0.5s;
}

.alert-danger {
  display: none;
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translate(-50%, 0);
  margin: auto;
  width: 300px;
  padding: 10px 0;
  background-color: #e74e4e;
  border: none;
  text-align: center;
  color: white;
}

.alert-primary {
  display: none;
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translate(-50%, 0);
  margin: auto;
  width: 600px;
  padding: 10px 0;
  background-color: hsla(0, 0%, 100%, 0.05);
  border: none;
  text-align: center;
  color: white;
}
.clipboard {
  cursor: pointer;
  margin-left: 10px;
  margin-bottom: 4px;
}

.quit-button {
  cursor: pointer;
}

#website-deployers-table {
  color: white;
}

.table {
  color: white;
}

tbody,
td,
tfoot,
th,
thead,
tr {
  border-color: rgb(86, 86, 86);
}

.table > :not(caption) > * > * {
  padding: 18px 1.5rem;
}

.table-striped > tbody > tr:nth-of-type(odd) > * {
  background-color: hsla(0, 0%, 100%, 0.05);
  color: white;
}

/* Popover css */

.popover__wrapper {
  /* right: 20px;
  top: 50px;
  position: fixed; */
  width: 100px;
  text-align: center;
  text-decoration: none;
  margin-left: 16px;
}

.wallet_button {
  text-decoration: none;
}
.popover__title {
  text-align: center;
  font-size: 14px;
  color: white;
  border-radius: 100px;
  padding: 7px 0px;
  background-color: #e74e4e;
}

.popover__content {
  opacity: 0;
  display: none;
  background-color: #e74e4e;
  width: 130px;
}
.popover__content:before {
  z-index: -1;
  content: "";
}
.popover__wrapper:hover .popover__content {
  z-index: 10;
  opacity: 1;
  display: inline-block;
}

.wallet-item {
  padding: 2px 3px;
}

.wallet-link {
  text-decoration: none;
  color: white;
}

#wallet-list {
  list-style-type: none;
  padding: 0px;
}

td,
th {
  text-align: center;
  vertical-align: middle;
}
.massa-logo-spinner {
  width: 25px;
  animation: spin 2s infinite linear;
  -webkit-animation: spin 2s infinite linear;
}
@-webkit-keyframes spin {
  0% {
    -webkit-transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
  }
}

/* PASSWORD MODAL */
.modal-content {
  background: linear-gradient(90deg, #172230 50.03%, #0c1219 85.82%);
}

.close {
  background: transparent;
}
.modal-footer {
  border-top: none;
}

/* Website creation  */

.website-card {
  display: flex;
  align-items: center;
  text-align: center;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.274);

  border: 1px solid lightgray;
  border-radius: 32px;
  padding: 16px;
  width: 500px;
}

.website-line {
  display: grid;
  grid-template-columns: repeat(2, auto);
  margin-top: 16px;
  width: 80%;
}

.popover__wrapper-website {
  width: 150px;
  text-align: center;
  text-decoration: none;
}

.popover-title-website {
  text-align: center;
  font-size: 24px;
  color: white;
  border-radius: 100px;
  padding: 7px 0px;
  background-color: #e74e4e;
}

.website-card-label {
  font-size: 1.2em;
  color: black;
  text-align: left;
  margin-right: 16px;
  margin-top: 8px;
}

.website-file-input {
  position: absolute;
  z-index: -1;
  height: 1px;
  width: 1px;
}

.spacer {
  margin-top: 16px;
}

.website-centering {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-items: center;
  text-align: center;
  margin: auto auto;
}

.fileError {
  display: none;
  font-size: 14px;
  color: red;
}

.dns-error {
  display: none;
  font-size: 14px;
  color: red;
}

#website-upload-refuse {
  background-color: #501313;
  border: none;
  display: none;
}

#website-upload-refuse:hover {
  cursor: default;
  transform: scale(1);
}
`
