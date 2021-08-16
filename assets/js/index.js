let store = {}

let owner = {
    id: '',
    name: '',
}

let filter = {
    personalAlbums: true,
    sharedAlbums: true,
    date: {
        start: '',
        end: '',
    },
    owners: [],
}

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

let renderFilter = () => {
    $("#selectedOwnersFilter").empty();
    
    filter.owners.forEach(v => {
        $("#selectedOwnersFilter").append(renderOwnerPill(v));
    });
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

let renderOwnerPill = (owner) => {
    return `
        <div class="col-12 col-pill">
            <input type="hidden" value="` + owner.id + `"/>
            <span class="badge rounded-pill bg-secondary">` + owner.name + `
                <button type="button" class="btn-close btn-close-white"></button>
            </span>
        </div>
    `
}

let bindFilter = () => {
    // bind to filter event
    $("#personalAlbumCheck").on("change", () => {
        filter.personalAlbums =  $("#personalAlbumCheck").prop('checked');
    });
    
    $("#sharedAlbumCheck").on("change", () => {
        filter.sharedAlbums = $("#sharedAlbumCheck").prop('checked');
    });

    $("#selectOwner").on("change", () => {
        if ($("#selectOwner option:selected").val() !== "empty-value") {
            let newOwner = {
                id: $("#selectOwner option:selected").val(),
                name: $("#selectOwner option:selected").text(),
            }

            let exists = false;
            filter.owners.forEach(v => {
                if (v.id === newOwner.id) {
                    exists = true;
                }
            })

            if (!exists) {
                filter.owners.push(newOwner);
                renderFilter();
            }
        }
    });

    $("#selectedOwnersFilter").on('click', '.btn-close', (e) => {
        let parents = $(e.target).parents("div");
        let id = $(parents[0]).find('input').val();
        
        filter.owners.forEach((v,idx) => {
            if ( v.id === id ) {
                filter.owners.splice(idx, 1);
            }
        });

        if ( filter.owners.length === 0 ) {
            $("#selectOwner option:eq(0)").prop('selected', true);
        }

        renderFilter();
    });
}
