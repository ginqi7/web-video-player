function playVideo(e) {
    console.log(e.target)
    file_name = e.target.textContent.trim();
    relativePath = document.querySelector(".navigations").textContent.trim();
    full_file_name = relativePath + file_name;
    console.log(full_file_name)
    const url = `/play?path=${encodeURIComponent(full_file_name)}`;
    window.location.assign(url);
}
