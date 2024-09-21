var last_subtitle = {}

function parseSRTTime(timeString) {
    // 将时间字符串拆分为小时、分钟、秒和毫秒
    const parts = timeString.split(/[:,]/); // 使用正则表达式分割

    const hours = parseInt(parts[0], 10);
    const minutes = parseInt(parts[1], 10);
    const seconds = parseInt(parts[2], 10);
    const milliseconds = parseInt(parts[3], 10);

    // 计算总秒数
    return (hours * 3600) + (minutes * 60) + seconds + milliseconds / 1000;
}

function fetchSTR(path, callback) {
    fetch(path)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.text();
        })
        .then(data => {
	    subtitles = parseSRT(data);
	    callback(subtitles)
	})
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
    
}

function findSTR(time, subtitles) {
    for (var subtile of subtitles) {
	if (subtile.start <= time &&
	    subtile.end >= time) {
	    return subtile;
	}
    }
}


function parseSRT(data) {
    const srtEntries = data.split('\n\n');
    const subtitles = [];

    srtEntries.forEach(entry => {
        const lines = entry.split('\n');
        if (lines.length >= 3) {
            const index = lines[0];
            const time = lines[1];
            const text = lines.slice(2).join('\n').trim();

            const times = time.split(' --> ');
            const start = parseSRTTime(times[0].trim());
            const end = parseSRTTime(times[1].trim());

            subtitles.push({
                index: parseInt(index, 10),
                start: start,
                end: end,
                text: text,
            });
        }
    });

    return subtitles;
}

function convertSRTDOM(subtitle) {
    if (subtitle && subtitle != last_subtitle) {
	last_subtitle = subtitle
	var subtitleDOM = document.querySelector(".subtitle")
	var words = subtitle.text.split(" ");
	subtitleDOM.innerHTML = ""
	for (var word of words) {
	    subtitleDOM.innerHTML += `<b onmouseover="overWord(event)" onmouseout="outWord(event)" onclick="overWord(event)" class="item cursor-pointer hover:bg-blue-500 m-2">${word}</b>`
	}
	window.scrollTo(0, document.body.scrollHeight);
    }
}

