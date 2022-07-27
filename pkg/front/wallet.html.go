// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: ..\..\pkg\front\wallet.html
// generated:  Wed Jul 27 13:15:43 CEST 2022

package front

const WalletHtml = `<!DOCTYPE html>
<html>
	<head>
		<link
			href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css"
			rel="stylesheet"
			integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor"
			crossorigin="anonymous"
		/>
		<link rel="stylesheet" href="./index.css" />
	</head>
	<style></style>
	<body>
		<nav class="navbar navbar-expand-lg navbar-dark">
			<div class="container">
			  <a class="navbar-brand" href="#"><img src="./logo_massa.webp" class="massa-logo" alt="Massa logo" /></a>
			  <h2>Thyra</h1>
			  <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			  </button>

			  <div class="collapse navbar-collapse justify-content-end " id="navbarNav">
				<ul class="navbar-nav ">
				  <li class="nav-item">
					<a class="nav-link active" aria-current="page" href="#">Wallet</a>
				  </li>
				  <li class="nav-item">
					<a class="nav-link" href="/webuploader.mythyra.massa/index.html">Websites</a>
				  </li>
				  <li class="nav-item">
					<a class="nav-link" href="#">Pricing</a>
				  </li>
				</ul>
			  </div>
			</div>
		  </nav>
			<div class="alert alert-danger" role="alert">
				This is a danger alert check it out!
			</div>

		<div class="container">

			<div class='row mt-5'>
				<div class='col-6'>
					
					<div class="row mb-5">
						<div class="col">
							<input id='fileid' type='file' hidden/>
							<button id="import-wallet"class="primary-button" id='buttonid' type='button' value='Upload MB' >Load Wallet</button>
						</div>
					</div>

					<div class="row">
						<div id="create-form">
							<h2 class="mb-4">Create Wallet</h1>
		
							<div class="form-floating mb-3 col-md-10">
								<input
									class="form-control"
									placeholder="Nickname"
									id="nicknameCreate"
									name="nicknameCreate"
									type="text"
								/>
								<label for="nicknameCreate">Nickname</label>
							</div>
		
							<div class="form-floating mb-4  col-md-10">
								<input
									class="form-control"
									placeholder="Password"
									id="password"
									name="password"
									type="password"
								/>
								<label for="password">Password</label>
							</div>
						</div>
		
						<div class="row mb-5">
							<div class="col">
								<button type="button" class="primary-button" onClick="createWallet()">Add</button>
							</div>
						</div>
					</div>

				</div>

				<div class='col-6'>
					<table id="user-wallet-table"class="table table-striped">
						<thead>
							<tr>
								<th scope="col">Address</th>
								<th scope="col">Nickname</th>
								<th scope="col">Balance</th>
								<th scope="col"></th>
							</tr>
						</thead>
						<tbody></tbody>
					</table>
				</div>

			</div>

		</div>
	</body>
	<script
		src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js"
		integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2"
		crossorigin="anonymous"
	></script>
	<script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="index.js"></script>
</html>
`
