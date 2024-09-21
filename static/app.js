var player = document.querySelector("#my-player")


window.onload = function() {
    fetchSTR(player.currentSrc.split('.').slice(0, -1).join('.')+".srt", (subtitles) => {
	setInterval(() => {
	    var subtitle = findSTR(player.currentTime, subtitles)
	    convertSRTDOM(subtitle)
	}, 50);
    })
    
    document.addEventListener('click', function(event) {
        console.log('全局点击事件触发！');
        console.log('点击位置:', event.clientX, event.clientY);
    });
};

function overWord(e) {
    e.preventDefault()
    player.pause()
    word = e.target.innerText
    translate(word, (data) => {
	console.log(data)
	data.word = word
	renderDictionary(data)
    })
}

function renderDictionary(data) {
    var dictionary = document.querySelector(".dictionary")
    dictionary.classList.remove("hidden")
    dictionary.innerHTML = ""
    if (data.word) {
	dictionary.innerHTML += `<h2>${data.word}</h2>`
    }
    if (data["usphone"]) {
	dictionary.innerHTML += `<div>美: ${data.usphone}</div>`
    }
    if (data["ukphone"]) {
	dictionary.innerHTML += `<div>英: ${data.ukphone}</div>`
    }
    for (var translate of data["trs"]) {
	dictionary.innerHTML += `<div>${translate}</div>`
    }
}

function outWord(e) {
    var dictionary = document.querySelector(".dictionary")
    dictionary.classList.add("hidden")
    
    player.play()
}

function toggleWord(e) {
    e.preventDefault()
    if (player.paused) {
	outWord(e)
    } else {
	overWord(e)
    }
}

