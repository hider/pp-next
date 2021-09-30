function initialize() {
    showToast();
    setInterval(function() {
        syncTimer();
        syncEvents();
        syncVotes();
        syncResults();
    }, 1000);    
}

function syncTimer() {

    var start = new Date(room.ts);
    var now = new Date();        
    var diff = now - start;
    var mins = Math.floor(diff / 60000);
    var secs = Math.floor((diff/1000)%60);
    var timer = padZero(""+mins)+":"+padZero(""+secs);
    var elem = document.getElementById("timer");
    elem.innerHTML = timer;

    function padZero(s) {
        if (s.length < 2) {
            return "0" + s;
        }
        return s;
    }
}

function syncEvents() {
    fetch("/rooms/" + room.name + "/events")
        .then(r => r.json())
        .then(d => {
            if (d.revealed && !room.revealed || d.reset && !room.resetBy) {
                reload();
            }
        });
}

function syncVotes() {
    fetch("/rooms/" + room.name + "/userlist")
        .then(r => r.text())
        .then(s => {
            var el = document.getElementById("userlist");
            el.innerHTML = s;
        });
}

function syncResults() {
    if (room.revealed) {
        fetch("/rooms/" + room.name + "/results")
            .then(r => r.text())
            .then(s => {
                var el = document.getElementById("results");
                el.innerHTML = s;
            });
    }
}

function vote(v) {
    fetch("/rooms/" + room.name + "/vote", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: v
    });
}

function showToast() {
    
    if (room.revealed) {
        makeToast(room.revealedBy, "revealed the votes");            
    } else if (room.resetBy) {
        makeToast(room.resetBy, "started a new story")
    }    


    function makeToast(name, action) {
        M.toast({ 
            html: "<span class='lime-text'>"
                + name + "</span>&nbsp;" + action,
            displayLength: 15000,
         });
    }

}

function reveal() {
    fetch("/rooms/" + room.name + "/reveal", {
        method: "POST"
    }).then(reload());
}

function reset() {
    fetch("/rooms/" + room.name + "/reset", {
        method: "POST"
    }).then(reload());
}

function reload() {
    window.location.reload();
}


