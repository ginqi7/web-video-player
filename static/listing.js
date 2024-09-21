function playVideo(e) {
    file_name = e.target.textContent.trim();
    relativePath = removeAllWhitespace(document.querySelector(".navigations").textContent.trim());
    full_file_name = relativePath + file_name;
    const url = `/play/${full_file_name}`;
    window.location.assign(url);
}

function removeAllWhitespace(str) {
    return str.replace(/\s+/g, '');
}
