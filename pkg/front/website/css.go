// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: html\front\website.css
// generated:  Fri Sep 30 17:55:38 CEST 2022

package website

const CSS = `body {
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu',
		'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
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
	width: fit-content;
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

.small-button {
	font-size: 16px;
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
	white-space: nowrap;
	overflow: hidden;
	display: block;
	text-overflow: ellipsis;
	padding-left: 14px;
	padding-right: 14px;
	background-color: #e74e4e;
}

.popover__disabled {
	pointer-events: none;
}

.popover__content {
	opacity: 0;
	display: none;
	background-color: #e74e4e;
	width: 130px;
}
.popover__content:before {
	z-index: -1;
	content: '';
}
.popover__wrapper:hover .popover__content {
	z-index: 10;
	opacity: 1;
	display: inline-block;
}

.wallet-item {
	padding: 2px 3px;
	white-space: nowrap;
	overflow: hidden;
	display: block;
	text-overflow: ellipsis;
}

.wallet-link {
	text-decoration: none;
	color: white;
}

#wallet-list {
	list-style-type: none;
	padding-left: 8px;
	padding-right: 8px;
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

/* -------------------------------------------------------------------------
  VERTICAL STEPPER
-------------------------------------------------------------------------- */

.stepper {
	display: none;
	background-color: #33506a;
	padding: 50px;
	border: 3px solid #808080;
	border-radius: 30px;
	width: 60%;
	max-width: 720px;
	margin: auto;
}

/* Steps */
.step {
	position: relative;
	min-height: 2em;
}
.step + .step {
	margin-top: 2em;
}
.step > div:first-child {
	position: static;
	height: 0;
}
.step > div:not(:first-child) {
	margin-left: 3em;
	padding-left: 1em;
}
.step.step-active {
	color: #808080;
}
.step.step-active .circle {
	background-color: #808080;
}

/* Circle */
.circle {
	background: #808080;
	position: relative;
	top: 7px;
	width: 2.5em;
	height: 2.5em;
	line-height: 2.5em;
	border-radius: 100%;
	color: white;
	text-align: center;
}

/* Vertical Line */
.circle:after {
	content: '';
	position: absolute;
	display: block;
	top: 1px;
	right: 50%;
	bottom: 1px;
	left: 50%;
	height: 90%;
	width: 2px;
	transform: scale(1, 2);
	transform-origin: 50% -100%;
	background-color: #808080;
	z-index: 0;
	font-size: 28px;
}

.bi-check {
	font-size: 28px;
}

.step:last-child .circle:after {
	display: none;
}

.step:first-child .circle:after {
	height: 95%;
	z-index: 0;
}

/* Stepper Titles */
.title {
	position: relative;
	color: #e6e6e6;
	line-height: 2.5em;
	font-weight: bold;
	font-size: 22px;
}

/* loading dots */

.loading-dots:after {
	content: ' .';
	animation: dots 1s steps(5, end) infinite;
}

@keyframes dots {
	0%,
	20% {
		color: rgba(0, 0, 0, 0);
		text-shadow: 0.25em 0 0 rgba(0, 0, 0, 0), 0.5em 0 0 rgba(0, 0, 0, 0);
	}
	40% {
		color: white;
		text-shadow: 0.25em 0 0 rgba(0, 0, 0, 0), 0.5em 0 0 rgba(0, 0, 0, 0);
	}
	60% {
		text-shadow: 0.25em 0 0 white, 0.5em 0 0 rgba(0, 0, 0, 0);
	}
	80%,
	100% {
		text-shadow: 0.25em 0 0 white, 0.5em 0 0 white;
	}
}
/* Website creation  */

.website-card {
	display: flex;
	align-items: center;
	text-align: center;
	flex-direction: column;
	background: #33506a;

	border: 3px solid #808080;
	border-radius: 30px;
	padding: 24px;
	width: 700px;
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

.website-card-label {
	font-size: 1.2em;
	text-align: left;
	margin-right: 16px;
	margin-top: 8px;
	display: flex;
	flex-direction: row;
}

.website-format-info {
	font-size: 0.95 em;
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

#fileError {
	display: none;
	font-size: 14px;
	color: red;
}

.dns-error {
	display: none;
	font-size: 14px;
	color: red;
	margin-top: -20px;
}

#website-upload {
	margin-top: 32px;
}

#website-upload-refuse {
	background-color: #501313;
	border: none;
	display: none;
	margin-top: 32px;
}

#website-upload-refuse:hover {
	cursor: default;
	transform: scale(1);
}

.website-info-display {
	display: flex;
	font-size: 11px;
	margin-top: -5px;
	color: grey;
}

::placeholder {
	color: blue;
	font-size: 13px;
}

.website-info:hover {
	display: flex;
}

.align-right {
	width: 100%;
	text-align: right;
	justify-content: right;
	align-items: right;
	display: flex;
	padding-right: 60px;
}

#website-wallet {
	font-size: 16px;
	background-color: #6c757d;
	white-space: nowrap;
	overflow: hidden;
	display: block;
	text-overflow: ellipsis;
	padding-left: 14px;
	padding-right: 14px;
}

.file-select-button-text {
	white-space: nowrap;
	overflow: hidden;
	display: block;
	text-overflow: ellipsis;
	width: 150px;
}

.website-dns {
	margin-top: 10px;
}

#file-select-button {
	background-color: #808080;
	border: none;
	overflow: hidden;
	text-overflow: ellipsis;
	max-width: 400px;
	height: 36px;
}
`
