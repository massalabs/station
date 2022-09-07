document.getElementById('import-wallet').addEventListener('click', openDialog);
document.getElementById('fileid').addEventListener('change', handleFileSelect, true);

getWallets();

let wallets = [];

// open file upload
function openDialog() {
	document.getElementById('fileid').value = null;
	document.getElementById('fileid').click();
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

// Import a wallet through PUT query
async function importWallet(wallet) {
	axios
		.put('/mgmt/wallet', wallet)
		.then((_) => {
			tableInsert(wallet);
			wallets.push(wallet);
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});
}

// Create a wallet through POST query
async function getWallets() {
	axios
		.get('/mgmt/wallet')
		.then((resp) => {
			if (resp) {
				const data = resp.data;
				for (const wallet of data) {
					tableInsert(wallet);
				}
				wallets = data;
			}
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});
}

// Create a wallet through POST query
function createWallet() {
	const nicknameCreate = document.getElementById('nicknameCreate').value;
	const password = document.getElementById('password').value;

	axios
		.post('/mgmt/wallet', {
			nickname: nicknameCreate,
			password: password,
		})
		.then((resp) => {
			tableInsert(resp.data);
			wallets.push(resp.data);
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});
}

function errorAlert(error) {
	document.getElementsByClassName('alert-danger')[0].style.display = 'block';

	document.getElementsByClassName('alert-danger')[0].innerHTML = error;

	setTimeout(function () {
		document.getElementsByClassName('alert-danger')[0].style.display = 'none';
	}, 5000);
}

function ellipsis(str) {
	return str.substr(0, 5) + '...' + str.substr(str.length - 5, str.length);
}

function tableInsert(resp) {
	const tBody = document.getElementById('user-wallet-table').getElementsByTagName('tbody')[0];
	const row = tBody.insertRow(-1);

	const cell0 = row.insertCell();
	const cell1 = row.insertCell();
	const cell2 = row.insertCell();
	const cell3 = row.insertCell();

	cell0.innerHTML =
		ellipsis(resp.address) +
		'<svg class="clipboard" onclick="copyToClipboard(this)" xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-external-link"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>';
	cell1.innerHTML = resp.nickname;
	cell2.innerHTML = 0;
	cell3.innerHTML =
		'<svg class="quit-button" onclick="deleteRow(this)" xmlns="http://www.w3.org/2000/svg" width="24" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line> <line x1="6" y1="6" x2="18" y2="18"></line></svg>';
}

function copyToClipboard(element) {
	const rowIndex = element.parentNode.parentNode.rowIndex;

	const tBody = document.getElementById('user-wallet-table').getElementsByTagName('tbody')[0];
	const address = tBody.rows[rowIndex - 1].cells[0].innerHTML;

	navigator.clipboard.writeText(address);
}

function deleteRow(element) {
	const rowIndex = element.parentNode.parentNode.rowIndex;

	const tBody = document.getElementById('user-wallet-table').getElementsByTagName('tbody')[0];
	const nickname = tBody.rows[rowIndex - 1].cells[1].innerHTML;

	axios
		.delete('/mgmt/wallet/' + nickname)
		.then((_) => {
			wallets = wallets.filter((wallet) => wallet.nickname != nickname);
		})
		.catch((e) => {
			errorAlert(getErrorMessage(e.response.data.code));
		});

	document.getElementById('user-wallet-table').deleteRow(rowIndex);
}
