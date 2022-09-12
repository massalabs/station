// GENERATED BY textFileToGoConst
// GitHub:     github.com/logrusorgru/textFileToGoConst
// input file: html\front\website.html
// generated:  Mon Sep 12 14:36:34 CEST 2022

package website

const HTML = `<!DOCTYPE html>
<html>

<head>

	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet"
		integrity="sha384-0evHe/X+R7YkIZDRvuzKMRqM+OrBnVFBL6DOitfPri4tjfHxaWutUpFmBp4vmVor" crossorigin="anonymous" />
	<link rel="stylesheet" href="./website.css" />
</head>
<style></style>

<body>
	<!-- NAVIGATION / HEADER -->
	<nav class="navbar navbar-expand-lg navbar-dark">
		<div class="container">
			<a class="navbar-brand" href="#"><img src="./logo_banner.webp" class="massa-logo-banner"
					alt="Massa logo" /></a>
			<h2>Thyra</h1>
				<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
					aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
					<span class="navbar-toggler-icon"></span>
				</button>

				<div class="collapse navbar-collapse justify-content-end " id="navbarNav">
					<ul class="navbar-nav ">
						<li class="nav-item">
							<a class="nav-link " href="/thyra/wallet/index.html">Wallet</a>
						</li>
						<li class="nav-item">
							<a class="nav-link active" aria-current="page" href="#">Websites</a>
						</li>
						<li class="nav-item">
							<a class="nav-link" href="#">Pricing</a>
						</li>
						<li>
							<div class="popover__wrapper">
								<a class="wallet_button" href="#">
									<h2 class="popover__title"></h2>
								</a>
								<div class="popover__content">
									<ul id="wallet-list">
									</ul>
								</div>
							</div>
						</li>
					</ul>
				</div>
		</div>
	</nav>


	<!-- Modal -->
	<div class="modal fade" id="passwordModal" tabindex="-1" role="dialog" aria-labelledby="passwordModalLabel"
		aria-hidden="true">
		<div class="modal-dialog modal-dialog-centered" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="passwordModalLabel">Input your password</h5>
					<button type="button" class="close" data-bs-dismiss="modal" aria-label="Close">
						<span aria-hidden="true">&times;</span>
					</button>
				</div>
				<div class="modal-body">
					<div class='col-6'>
						<div id="password-form">
							<div class="form-floating">
								<input class="form-control" placeholder="Wallet Password" id="walletPassword"
									name="walletPassword" type="password" />
								<label for="wallet-password">Wallet Password</label>
							</div>
						</div>
					</div>
				</div>
				<div class="modal-footer">
					<button type="button" data-bs-dismiss="modal" onClick="callTx('deployWebsiteCreator')"
						class="primary-button process-post-modal">Call</button>
				</div>
			</div>
		</div>
	</div>
	<!-- ALERTS -->

	<div class="alert alert-danger" role="alert"></div>
	<div class="alert alert-primary" role="alert"></div>

	<!-- DEPLOY WEBSITE CREATOR FORM-->
	<div class="container">
		<div class="website-centering">
			<h2 class="mb-4 mt-5">Decentralized website storing</h2>



			<div class="website-card">

				<div class="website-line">
					<h4 class="website-card-label">On wallet</h4>

					<h2 class="popover__title" style="font-size: 24px; background-color: darkblue;"></h2>


				</div>

				<div class="website-line">
					<h4 class="website-card-label">Website Name</h4>

					<div class="website-dns">
						<input class="form-control" id="websiteName" name="websiteName" type="text" />
					</div>
				</div>
				<h4 class="dns-error">Website Name must be only lowercase letters and numbers</h4>
				<div class="website-line">
					<h4 class="website-card-label">Website File</h4>
					<button class="primary-button" style="background-color : darkblue ; border : none"
						id="file-select-button">
						Import From

					</button>

					<div class="upload">
						<input class="website-file-input" type="file" accept=".zip" />
					</div>

					<h4 class="fileError">File type needs to be .zip only</h4>
				</div>
				<div class="spacer">
					<button class="primary-button" id="website-upload" onClick="setTxType('deployWebsiteAndUpload')"
						class="primary-button me-5" type='button' value='Upload MB' data-bs-toggle="modal"
						data-bs-target="#passwordModal">Upload</button>

					<button class="primary-button" id="website-upload-refuse">Upload</button>

					<img src="./logo.png" class="massa-logo-spinner loading" alt="Massa logo" />
				</div>

			</div>

		</div>


		<div class='row mt-5'>
			<div class='col'>
				<table id="website-deployers-table" class="table table-striped">
					<thead>
						<tr>
							<th scope="col">Name</th>
							<th scope="col">Address</th>
							<th scope="col">Url</th>
							<th scope="col"></th>
						</tr>
					</thead>
					<tbody></tbody>
				</table>
			</div>
		</div>
	</div>
</body>
<script src="errors.js"></script>
<script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.11.5/dist/umd/popper.min.js"
	crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0-beta1/dist/js/bootstrap.bundle.min.js"
	integrity="sha384-pprn3073KE6tl6bjs2QrFaJGz5/SUsLqktiwsUTF55Jfv3qYSDhgCecCxMW52nD2"
	crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>
<script src="https://code.jquery.com/jquery-3.6.0.min.js"
	integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
<script src="website.js"></script>

</html>`
