var id_token;

function onSignIn(googleUser) {
  // The ID token you need to pass to your backend:
  id_token = googleUser.getAuthResponse().id_token;

  loadEmotes();
}

function whoami() {
  fetch('/whoami', {headers: {"Authorization": id_token}})
    .then(
      function(response) {
        if (response.status !== 200) {
          console.log('Looks like there was a problem. Status Code: ' +
            response.status);
          return;
        }

        // Examine the text in the response
        response.json().then(function(data) {
          console.log(data);
        });
      }
    )
    .catch(function(err) {
      console.log('Fetch Error :-S', err);
    });
}

function loadEmotes() {
  fetch('/emotes', {headers: {"Authorization": id_token}})
    .then(
      function(response) {
        if (response.status !== 200) {
          console.log('Looks like there was a problem. Status Code: ' +
            response.status);
          return;
        }

        // Examine the text in the response
        response.json().then(function(data) {
          console.log(data);
	  renderEmotes(data);
        });
      }
    )
    .catch(function(err) {
      console.log('Fetch Error :-S', err);
    });
}

function renderEmotes(emotes) {
	var listdiv = document.getElementById("emotelist");
	listdiv.innerHTML = "";
	emotes.forEach((ele, index) => {
		var newdiv = document.createElement("div");
		newdiv.innerHTML = ele.name;
		listdiv.appendChild(newdiv);
	});
}
