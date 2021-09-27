function deleteAlbum(url) {
    res = await axios.delete(url);

    console.log(res);
}

