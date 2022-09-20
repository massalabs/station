let gWallets = [];
let deployers = [];
let actualTxType = '';
let nextFileToUpload;
let uploadable = false;

// INIT
getWallets();
getWebsiteDeployerSC();
initializeDefaultWallet();
setupModal();

const eventManager = new EventManager();

// Write the default wallet text in wallet popover component
async function getWebsiteDeployerSC() {
	let defaultWallet = getDefaultWallet();

	$('#website-deployers-table tbody tr').remove();

	axios
		.get('/my/domains/' + defaultWallet)
		.then((websites) => {
			let count = 0;
			for (const website of websites.data) {
				tableInsert(website, count);
				count++;
			}
			deployers = websites.data;
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});
}

// Write the default wallet text in wallet popover component
function initializeDefaultWallet() {
	let defaultWallet = getDefaultWallet();
	if (defaultWallet === '') {
		defaultWallet = 'Connect';
	}
	$('.popover__title').html(defaultWallet);
}

// Retrieve the default wallet nickname in cookies
function getDefaultWallet() {
	let defaultWallet = '';
	const cookies = document.cookie.split(';');
	cookies.forEach((cookie) => {
		const keyValue = cookie.split('=');
		if (keyValue[0] === 'defaultWallet') {
			defaultWallet = keyValue[1];
		}
	});
	return defaultWallet;
}

function getWallet(nickname) {
	return gWallets.find((w) => w.nickname === nickname);
}

function getDeployerByAddress(contractAddress) {
	return deployers.find((c) => c.address === contractAddress);
}

function getDeployerByDns(dns) {
	return deployers.find((c) => c.name === dns);
}

function setupModal() {
	$('#passwordModal').on('shown.bs.modal', function () {
		$('#passwordModal').trigger('focus');
	});
}

function setTxType(txType) {
	actualTxType = txType;
}

async function callTx() {
	const passwordValue = $('#walletPassword').val();

	if (actualTxType === 'deployWebsiteAndUpload') {
		deployWebsiteAndUpload(passwordValue);
	}
	if (actualTxType.includes('uploadWebsiteCreator')) {
		const websiteIndex = actualTxType.split('uploadWebsiteCreator')[1];
		uploadWebsite(nextFileToUpload, websiteIndex, passwordValue);
	}
}

// open file upload
function openDialog() {
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

// Append wallet accounts in popover component list
async function feedWallet(w) {
	let counter = 0;
	for (const wallet of w) {
		$('#wallet-list').append(
			"<li class='wallet-item'><a class='wallet-link' id='wallet-link-" +
				counter +
				"' onclick='changeDefaultWallet(event)' href='#'>" +
				wallet.nickname +
				'</a></li>'
		);
		counter++;
	}
}

// Handle popover click & update default wallet in cookies
function changeDefaultWallet(event) {
	const idElementClicked = event.target.id;
	const newDefaultWalletId = idElementClicked.split('-')[2];
	const walletName = gWallets[newDefaultWalletId].nickname;

	document.cookie = 'defaultWallet=' + walletName;
	$('.popover__title').html(walletName);

	getWebsiteDeployerSC();
}

async function getWallets() {
	axios
		.get('/mgmt/wallet')
		.then((resp) => {
			if (resp) {
				gWallets = resp.data;
				feedWallet(gWallets);
			}
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});
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
	cell2.innerHTML =
		"<a href='" + url + "'target = '_blank' rel='noopener noreferrer'>" + url + '</a>';
	cell3.innerHTML =
		"<div><input id='fileid" +
		count +
		"' type='file' hidden/><button id='upload-website" +
		count +
		"' class='primary-button' id='buttonid' type='button' value='Upload MB' >Edit</button><img src='./logo.png' style='display: none' class='massa-logo-spinner loading" +
		count +
		" alt='Massa logo' /></span></div>";

	document.getElementById(`upload-website${count}`).addEventListener('click', function () {
		document.getElementById(`fileid${count}`).value = null;
		document.getElementById(`fileid${count}`).click();
	});

	document.getElementById(`fileid${count}`).addEventListener('change', function (evt) {
		let files = evt.target.files; // get files
		nextFileToUpload = files[0];

		setTxType('uploadWebsiteCreator' + count);
		$('#passwordModal').modal('show');
	});
}

function uploadWebsite(file, count, password) {
	let defaultWallet = getDefaultWallet();
	const bodyFormData = new FormData();

	const address = deployers[count].address;
	bodyFormData.append('zipfile', file);
	bodyFormData.append('address', address);
	bodyFormData.append('nickname', defaultWallet);

	const reader = new FileReader();
	reader.readAsDataURL(file);
	reader.onloadend = (_) => {
		const result = reader.result.length;

		const chunkSize = Math.floor(result / 260_000) + 1;

		stepper(getDeployerByAddress(address).name, chunkSize, false);

		axios({
			url: `/websiteCreator/upload`,
			method: 'POST',
			data: bodyFormData,
			headers: {
				'Content-Type': 'multipart/form-data',
				Authorization: password,
			},
		})
			.then((operation) => {
				document.getElementsByClassName('loading' + count)[0].style.display = 'none';
				successMessage('Website uploaded to address : ' + operation.data.address);
			})

			.catch((e) => {
				document.getElementsByClassName('loading' + count)[0].style.display = 'none';
				errorAlert(getErrorMessage(e.response.data.code));
			});
	};
}

$('#file-select-button').click(function () {
	$('.upload input').click();
});

// change button text with file name
$('.upload input').on('change', function () {
	let str = $('.upload input').val();

	let n = str.lastIndexOf('\\');

	let result = str.substring(n + 1);

	$('#file-select-button').html(result);
});

//check if file is .zip
$('.upload input').on('change', function () {
	let str = $('.upload input').val();

	let n = str.lastIndexOf('.');

	let result = str.substring(n + 1);

	if (result != 'zip' && $('.upload input').val() != '') {
		uploadable = false;

		document.getElementsByClassName('fileError')[0].style.display = 'flex';
		document.getElementById('website-upload').style.display = 'none';
		document.getElementById('website-upload-refuse').style.display = 'flex';
	} else {
		uploadable = true;
		document.getElementsByClassName('fileError')[0].style.display = 'none';
		document.getElementById('website-upload').style.display = 'flex';
		document.getElementById('website-upload-refuse').style.display = 'none';
	}
});

//remove label of input website name on focus
$('.website-dns input').on('focus', function () {
	document.getElementById('website-info-display').style.visibility = 'hidden';
});

//check if input string is valid
$('.website-dns input').on('change', function () {
	let str = $('.website-dns input').val();
	let pattern = new RegExp('^[a-z0-9]+$');
	let testPattern = pattern.test(str);

	if (testPattern == false) {
		uploadable = false;
		document.getElementsByClassName('dns-error')[0].style.display = 'flex';
		document.getElementById('website-upload').style.display = 'none';
		document.getElementById('website-upload-refuse').style.display = 'flex';
	} else {
		uploadable = true;
		document.getElementsByClassName('dns-error')[0].style.display = 'none';
		document.getElementById('website-upload').style.display = 'flex';
		document.getElementById('website-upload-refuse').style.display = 'none';
	}
});

function deployWebsiteAndUpload(password) {
	if (uploadable == true) {
		let defaultWallet = getDefaultWallet();
		const dnsNameInputValue = document.getElementById('websiteName').value;

		const file = document.querySelector('.upload input').files[0];
		const bodyFormData = new FormData();
		bodyFormData.append('url', dnsNameInputValue);
		bodyFormData.append('nickname', defaultWallet);
		bodyFormData.append('zipfile', file);

		const reader = new FileReader();
		reader.readAsDataURL(file);
		reader.onloadend = (_) => {
			const result = reader.result.length;

			const chunkSize = Math.floor(result / 260_000) + 1;

			stepper(dnsNameInputValue, chunkSize, true);

			axios({
				url: `/websiteCreator/prepare`,
				method: 'put',
				data: bodyFormData,
				headers: {
					'Content-Type': 'multipart/form-data',
					Authorization: password,
				},
			})
				.then((operation) => {
					successMessage('Website uploaded to address : ' + operation.data.address);
					getWebsiteDeployerSC();
				})
				.catch((e) => {
					errorAlert(getErrorMessage(e.response.data.code));
				});
		};
	}
}

async function stepper(dnsName, totalChunk, isFullProcess) {
	initStepper(dnsName, totalChunk);
	if (isFullProcess) {
		step1(dnsName, totalChunk);
	} else {
		$('.circle').eq(0).empty();
		$('.circle').eq(0).append('<i class="bi bi-check">');
		$('.title').eq(0).removeClass('loading-dots');
		step3(getDeployerByDns(dnsName).address, totalChunk);
	}
}

function initStepper(dnsName, totalChunk) {
	$('.website-card').hide();
	$('.stepper').show();

	$('.stepper-title').html('Deployment of ' + dnsName);
	$('.title').eq(0).addClass('loading-dots');

	$('.title')
		.eq(2)
		.text('Chunk upload ' + 1 + ' on ' + totalChunk);
}

function step1(dnsName, totalChunk) {
	eventManager.subscribe(
		`Website Deployer is deployed at :`,
		getWallet(getDefaultWallet()).address,
		(resp) => step2(dnsName, resp.data.data.split(':')[1], totalChunk)
	);
}

function step2(dnsName, contractAddress, totalChunk) {
	eventManager.subscribe(
		`Resolver set to record key :record${dnsName}at address `,
		getWallet(getDefaultWallet()).address,
		(_) => {
			step3(contractAddress, totalChunk);
		}
	);

	$('.circle').eq(0).empty();
	$('.circle').eq(0).append('<i class="bi bi-check">');

	$('.title').eq(0).removeClass('loading-dots');
	$('.title').eq(1).addClass('loading-dots');
}

function step3(contractAddress, totalChunk) {
	let actualChunk = 1;

	for (let i = 0; i < totalChunk; i++) {
		if (i === 0) {
			eventManager.subscribe(
				`First chunk deployed to ${contractAddress}`,
				getWallet(getDefaultWallet()).address,
				(_) => {
					actualChunk++;
					$('.title')
						.eq(2)
						.text('Chunk upload ' + actualChunk + ' on ' + totalChunk);
					$('.title').eq(2).addClass('loading-dots');

					if (totalChunk === 1) {
						resetStepper();
					}
				}
			);
		} else if (i == totalChunk - 1) {
			eventManager.subscribe(
				`Append chunk deployed to ${contractAddress} : ${totalChunk - 1}`,
				getWallet(getDefaultWallet()).address,
				(_) => {
					resetStepper();
				}
			);
		} else {
			eventManager.subscribe(
				`Append chunk deployed to ${contractAddress} : ${i}`,
				getWallet(getDefaultWallet()).address,
				(_) => {
					actualChunk++;
					$('.title')
						.eq(2)
						.text('Chunk upload ' + actualChunk + ' on ' + totalChunk);
					$('.title').eq(2).addClass('loading-dots');
				}
			);
		}
	}

	$('.circle').eq(1).empty();
	$('.circle').eq(1).append('<i class="bi bi-check">');

	$('.title').eq(1).removeClass('loading-dots');
	$('.title').eq(2).addClass('loading-dots');
}

function resetStepper() {
	$('.website-card').show();
	$('.stepper').hide();

	$('.circle').empty();
	$('.circle').eq(0).html('1');
	$('.circle').eq(1).html('2');
	$('.circle').eq(2).html('3');

	$('.title').eq(2).html('Chunk upload');

	$('.title').eq(2).removeClass('loading-dots');
}
