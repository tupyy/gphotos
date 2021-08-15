let store = {}

const albumsElementID = "#albums";

const requestURL = "/api/albums"

$(() => {
    // get albums from server
    $('#datepicker').datepicker({
        format: "dd/mm/yyyy",
        weekStart: 1,
        autoclose: true,
        clearBtn: true,
        todayHighlight: true,
        beforeShowMonth: function(date){
              if (date.getMonth() == 8) {
                return false;
              }
            },
        beforeShowYear: function(date){
              if (date.getFullYear() == 2007) {
                return false;
              }
            }
    });

    bindFilter();

    showSpinner($("#albums"),true);
    clearAlbums();

    axios.get(requestURL)
        .then(response => {

            store.data = response.data;
            store.albums = _.concat(response.data.personal_albums, response.data.shared_albums);

            render();

        })
        .catch(e => {
            console.log(e);
        })
        .then(() => {
            showSpinner($("#albums"), false);
        });

});

let render = () => {
    store.albums.forEach(album => {
        $(albumsElementID).append(renderAlbum(album));
    });
    
    $("#count_albums").html(store.albums.length);
}

let renderAlbum = (album) => {
    return `
        <div class="album-col col-2" id="` + album.id + `">
            <div class="container-album-card card">
                <a href="/album/` + album.id + `">
                    <img src="/static/img/eiffeltoren.jpg" class="card-img-top"/>
                </a>
                <div class="album">
                    <div class="card-body">
                        <span class="author">` + album.owner + `</span>
                        <h1 class="card-title title">
                            <a href="/album/` + album.id + `">` + album.name + `</a>
                        </h1>
                        <p class="card-text">` + album.description + `</p>
                    </div>
                </div>
            </div>
        </div>
    `
}

let clearAlbums = () => {
    $(albumsElementID).empty();
}

let showSpinner = (parentElement, show) => {
    if (show) {
        parentElement.append(spinner())
    } else {
        $("#loadingSpinner").remove();
    }
}
let spinner = () => {
    return `
    <div class="d-flex justify-content-center" id="loadingSpinner">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
    `
}

let bindFilter = () => {
    // bind to filter event
    $("#personalAlbumCheck").on("change", () => {
        clearAlbums();
        store.albums = [];

        if ($("#personalAlbumCheck").prop('checked')) {
            store.albums = _.concat(store.albums, store.data.personal_albums);
        }

        if ($("#sharedAlbumCheck").prop('checked')) {
            store.albums = _.concat(store.albums, store.data.shared_albums);
        }
        
        render();
    });
    
    $("#sharedAlbumCheck").on("change", () => {
        clearAlbums();
        store.albums = [];

        if ($("#personalAlbumCheck").prop('checked')) {
            store.albums = _.concat(store.albums, store.data.personal_albums);
        }

        if ($("#sharedAlbumCheck").prop('checked')) {
            store.albums = _.concat(store.albums, store.data.shared_albums);
        }

        render();
    });
}
