var player = document.querySelector("#my-player")
var audio = document.querySelector("#my-audio")

document.addEventListener('keydown', function(event) {
    // 检查是否按下空格键（keyCode 32）
    if (event.code === 'Space') {
        event.preventDefault(); // 防止页面滚动
        if (player.paused) {
            player.play(); // 如果视频暂停，则播放
        } else {
            player.pause(); // 如果视频正在播放，则暂停
        }
    }
});

document.addEventListener("DOMContentLoaded", function() {
    const video = document.getElementById("my-player");
    const audio = document.getElementById("my-audio");

    video.addEventListener("play", function() {
        audio.play(); // Play the audio when the video starts
    });

    video.addEventListener("pause", function() {
        audio.pause(); // Pause the audio when the video pauses
    });

    video.addEventListener("timeupdate", function() {
        // Sync the audio time with the video time, if necessary
        // This can also help in seeking
        if (Math.abs(video.currentTime - audio.currentTime) > 0.1) {
            audio.currentTime = video.currentTime;
        }
    });
});


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

