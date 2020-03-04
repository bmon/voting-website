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
  Promise.all([
    fetch('/emotes', {headers: {"Authorization": id_token}}),
    fetch('/votes', {headers: {"Authorization": id_token}}),
  ])
    .then((responses) => {
      var[emotes, votes] = responses;
      if (emotes.status !== 200 || votes.status !== 200) {
        console.log("failed to load emote data")
      }
        // Examine the text in the response
        Promise.all([emotes.json(), votes.json()]).then(data => {
          console.log(data);
          var[emoteData, voteData] = data;
	        renderEmotes(emoteData, voteData);
          setMessage("Click on an emote to vote for it!")
        });
      }
    )
    .catch(function(err) {
      console.log('Fetch Error :-S', err);
    });
}

function renderEmotes(emotes, votes) {
	var listdiv = document.getElementById("emotelist");
	listdiv.innerHTML = "";
  console.log(emotes)
	emotes.forEach((emote, index) => {
		var newdiv = document.createElement("div");
    newdiv.dataset.emotename = emote.name
		newdiv.innerHTML = `<img src="/static/emotes/${emote.filename}"></img><div>:${emote.name}:</div>`;
    newdiv.classList.add("emote");
    if (votes.includes(emote.name)) {
      newdiv.classList.add("voted");
    }

    newdiv.onclick = toggleVote;
		listdiv.appendChild(newdiv);
	});
}

function toggleVote() {
  if (this.classList.contains("voting")) return;

  this.classList.add("voting");
  var action = "add"
  if (this.classList.contains("voted")) {
    action = "retract";
  }

  fetch('/vote', {
    method: "POST",
    headers: {
      "Authorization": id_token,
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: `emote=${this.dataset.emotename}&action=${action}`
  }).then(resp => {
    this.classList.remove("voting");
    if (resp.status === 200) {
      if (action === "add") {
        this.classList.add("voted");
      } else {
        this.classList.remove("voted");
      }
    }
  }).catch(err => {
    console.log(err)
    this.classList.remove("voting");
  })
}
