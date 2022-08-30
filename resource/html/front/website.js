getWallets();
getWebsiteDeployerSC();

let wallets = [];
let deployers = [];

// open file upload
function openDialog(count) {
	console.log(count);
	document.getElementById('fileid0').value = null;
	document.getElementById('fileid0').click();
}

// Handle event on file selecting
function handleFileSelect(evt) {
	let files = evt.target.files; // get files
	let f = files[0];
	const reader = new FileReader();
	reader.onload = (event) => importWallet(JSON.parse(event.target.result)); // desired file content
	reader.onerror = (error) => reject(error);
	reader.readAsText(f);
}

function errorAlert(error) {
	document.getElementsByClassName('alert-danger')[0].style.display = 'block';

	document.getElementsByClassName('alert-danger')[0].innerHTML = error;

	setTimeout(function () {
		document.getElementsByClassName('alert-danger')[0].style.display = 'none';
	}, 5000);
}

function successMessage(message) {
	document.getElementsByClassName('alert-primary')[0].style.display = 'block';

	document.getElementsByClassName('alert-primary')[0].innerHTML = message;

	setTimeout(function () {
		document.getElementsByClassName('alert-primary')[0].style.display = 'none';
	}, 5000);
}
// Import a wallet through PUT query
async function feedSelectOption(w) {
	const formSelect = document.getElementsByClassName('form-select');
	for (let i = 0; i < w.length; i++) {
		const opt = document.createElement('option');
		opt.value = i;
		opt.text = w[i].nickname;
		formSelect[0].append(opt);
	}
}

// Create a wallet through POST query
async function getWallets() {
	axios
		.get('/mgmt/wallet')
		.then((resp) => {
			if (resp) {
				const data = resp.data;
				wallets = data;
				feedSelectOption(wallets);
			}
		})
		.catch((e) => {
			errorAlert(e);
		});
}

async function getWebsiteDeployerSC() {
	axios
		.get('/my/domains')
		.then((websites) => {
			let count = 0;
			for (const website of websites.data) {
				tableInsert(website, count);
				count++;
			}
			deployers = websites.data;
		})
		.catch((e) => {
			errorAlert(e.response.data.code);
		});
}

async function deployWebsiteDeployerSC() {
	const dnsNameInputValue = document.getElementById('websiteName').value;

	if (dnsNameInputValue == '') {
		console.log(dnsNameInputValue == '');
		errorAlert('Input a DNS name');
	} else {
		document.getElementsByClassName('loading')[0].style.display = 'inline-block';
		document.getElementsByClassName('loading')[1].style.display = 'inline-block';
		axios
			.put('/websiteCreator/prepare/', { url: dnsNameInputValue })
			.then((operation) => {
				document.getElementsByClassName('loading')[0].style.display = 'none';
				document.getElementsByClassName('loading')[1].style.display = 'none';
				successMessage('Contract deployed to address ' + operation.data.address);
			})
			.catch((e) => {
				errorAlert(e.response.data.code);
			});
	}
}

function tableInsert(resp, count) {
	const tBody = document.getElementById('website-deployers-table').getElementsByTagName('tbody')[0];
	const row = tBody.insertRow(-1);
	const url = 'http://' + resp.name + '.massa/';

	const cell0 = row.insertCell();
	const cell1 = row.insertCell();
	const cell2 = row.insertCell();
	const cell3 = row.insertCell();

	cell0.innerHTML = resp.name;
	cell1.innerHTML = resp.address;
	cell2.innerHTML = "<a href='" + url + "'>" + url + '</a>';
	cell3.innerHTML =
		"<div><input id='fileid" +
		count +
		"' type='file' hidden/><button id='upload-website" +
		count +
		"'" +
		"class='primary-button' id='buttonid' type='button' value='Upload MB' >Upload</button><span style='display: none' class='spinner-border loading" +
		count +
		"' role='status'><img src='./logo.png' class='massa-logo-spinner' alt='Massa logo' /></span></div> ";

	document.getElementById(`upload-website${count}`).addEventListener('click', function () {
		document.getElementById(`fileid${count}`).value = null;
		document.getElementById(`fileid${count}`).click();
	});

	document.getElementById(`fileid${count}`).addEventListener('change', function (evt) {
		let files = evt.target.files; // get files
		let f = files[0];
		uploadWebsite(f, count);
	});
}

function uploadWebsite(file, count) {
	const bodyFormData = new FormData();
	bodyFormData.append('zipfile', file);
	bodyFormData.append('address', deployers[count].address);
	document.getElementsByClassName('loading' + count)[0].style.display = 'inline-block';
	axios({
		url: `/websiteCreator/upload`,
		method: 'POST',
		data: bodyFormData,
		headers: {
			'Content-Type': 'multipart/form-data',
		},
	})
		.then((operation) => {
			document.getElementsByClassName('loading' + count)[0].style.display = 'none';
			successMessage('Website uploaded with operation ID : ' + operation);
		})
		.catch((e) => {
			errorAlert(e.response.data.code);
		});
}
