function translate(word, callback) {    
    fetch(`/translate?word=${word}`)
	.then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
	    callback(data)
	})
        .catch(error => {
            console.error('There was a problem with the fetch operation:', error);
        });
}
