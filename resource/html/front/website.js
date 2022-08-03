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

function successDeployWebsite(contract) {
	document.getElementsByClassName('alert-primary')[0].style.display = 'block';

	document.getElementsByClassName('alert-primary')[0].innerHTML =
		'Contract deployed to address ' + contract.address;

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
		.get('/uploadWeb')
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
			.post('/uploadWeb/' + dnsNameInputValue)
			.then((operation) => {
				document.getElementsByClassName('loading')[0].style.display = 'none';
				document.getElementsByClassName('loading')[1].style.display = 'none';
				successDeployWebsite(operation.data);
			})
			.catch((e) => {
				errorAlert(e.response.data.code);
			});
	}
}

function tableInsert(resp, count) {
	const tBody = document.getElementById('website-deployers-table').getElementsByTagName('tbody')[0];
	const row = tBody.insertRow(-1);

	const cell0 = row.insertCell();
	const cell1 = row.insertCell();
	const cell2 = row.insertCell();

	cell0.innerHTML = resp.name;
	cell1.innerHTML = resp.address;
	cell2.innerHTML =
		"<div><input id='fileid" +
		count +
		"' type='file' hidden/><button id='updload-website" +
		count +
		"'" +
		"class='primary-button' id='buttonid' type='button' value='Upload MB' >Upload</button></div>";

	document.getElementById(`updload-website${count}`).addEventListener('click', function () {
		document.getElementById(`fileid${count}`).value = null;
		document.getElementById(`fileid${count}`).click();
	});

	document.getElementById(`fileid${count}`).addEventListener('change', function (evt) {
		let files = evt.target.files; // get files
		console.log(evt.target.files);
		let f = files[0];
		const reader = new FileReader();

		reader.onload = (event) => uploadWebsite(event.target.result, count); // desired file content
		reader.onerror = (error) => reject(error);
		reader.readAsText(f);
	});
}

function uploadWebsite(file, count) {
	console.log(file);
	const formData = new FormData();
	formData.append('zipfile', file);
	axios
		.post(`/fillWeb/${deployers[count].address}`, formData, {
			headers: {
				'Content-Type': 'multipart/form-data',
			},
		})
		.then((operation) => {
			console.log('OK');
		})
		.catch((e) => {
			errorAlert(e.response.data.code);
		});
}
