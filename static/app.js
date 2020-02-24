var id_token;

function setMessage(msg) {
  document.getElementById('message').innerHTML = msg;
}

function onSignIn(googleUser) {
  id_token = googleUser.getAuthResponse().id_token;
  var hd = document.head.querySelector("[name~=google-signin-hd][content]").content;

  if (googleUser.getHostedDomain() != hd) {
    setMessage("Please sign in with an account under the domain: " + hd);
    return
  }
  setMessage("Click on an emote to vote for it!")
  loadEmotes();
}

function whoami() {
  fetch('/whoami', {headers: {"Authorization": id_token}})
    .then(
      function (response) {
        if (response.status !== 200) {
          console.log('Looks like there was a problem. Status Code: ' +
            response.status);
          return;
        }

        // Examine the text in the response
        response.json().then(function (data) {
          console.log(data);
        });
      }
    )
    .catch(function (err) {
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
          setMessage("Click on an emote to vote for it!")
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
		newdiv.innerHTML = `<img src="/static/emotes/${ele.filename}"></img><div>:${ele.name}:</div>`;
    newdiv.classList.add("emote");
		listdiv.appendChild(newdiv);
	});
}
