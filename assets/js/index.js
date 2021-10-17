const albumsElementID = "#albums";
const baseURL = "/api/albums"

let store = {}

let filterSort = {
    personalAlbums: true,
    sharedAlbums: true,
    date: {
        start: '',
        end: '',
    },
    owners: [],
    sort: '',
    buildRequestURL: function(baseURL) {
        let reqUrl = baseURL + "?personal=" + this.personalAlbums + "&shared=" + this.sharedAlbums;
        
        if ( this.date.start !== "" ) {
            reqUrl = reqUrl + "&start_date=" + this.date.start;
        }

        if (this.date.end !== "") {
            reqUrl = reqUrl + "&end_date=" + this.date.end;
        }

        if (this.owners.length > 0) {
            this.owners.forEach(id => {
                reqUrl = reqUrl + "&owner=" + id
            });
        }

        if (this.sort !== '') {
            reqUrl = reqUrl + '&sort=' + this.sort;
        }

        return encodeURI(reqUrl);
    },
}


const init = () => {
    filterSort.personalAlbums =  $("#personalAlbumCheck").prop('checked');
    filterSort.sharedAlbums = $("#sharedAlbumCheck").prop('checked');
    filterSort.date.start = $('#startDate').val();
    filterSort.date.end = $('#endDate').val();
    filterSort.sort = $('#sortSelect').val();
}

const doReq = () => {
    showSpinner($("#albums"),true);
    $("#count_albums").parent().hide();

    clearAlbums();

    axios.get(filterSort.buildRequestURL(baseURL))
        .then(response => {

            store.data = response.data;
            store.albums = response.data.albums;
            store.username = response.data.username;

            render();

        })
        .catch(e => {
            console.log(e);
        })
        .then(() => {
            showSpinner($("#albums"), false);
            $("#count_albums").parent().show();
            bindToCardOwner();
        });
}


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
                <div class="card-header">
                    <div class="row row-owner">
                        <span id="owner">
                            ` + album.owner + `
                        </span>
                        <span class="location">
                            <i class="fas fa-map-marker-alt"></i>
                            ` + album.location + `
                        </span>
                    </div>
                    <div class="row row-date">
                        <span>
                            <i class="far fa-calendar-alt"></i>
                            ` + album.date + `
                        </span>
                    </div> 
                </div>
                <a href="/album/` + album.id + `">
                    <img src="` + album.thumbnail + `" class="card-img-top"/>
                </a>
                <div class="album">
                    <div class="card-body">
                        <h1 class="card-title title">
                            <a href="/album/` + album.id + `">` + album.name + `</a>
                        </h1>
                    </div>
                </div>
            </div>
        </div>
    `
}

const renderFilter = () => {
        $("#personalAlbumCheck").prop('checked', filterSort.personalAlbums);
        $("#sharedAlbumCheck").prop('checked', filterSort.sharedAlbums);

        
        $("#ownerFilter .form-check").each((_, e) => {
            id = $(e).find('input').val();
            find = false;
            filterSort.owners.forEach((i) => {
                if (i === id) {
                    find = true;
                }
            });
            $(e).find('input').prop('checked', find);
        });
}

const clearAlbums = () => {
    $(albumsElementID).empty();
}

const showSpinner = (parentElement, show) => {
    if (show) {
        parentElement.append(spinner())
    } else {
        $("#loadingSpinner").remove();
    }
}

const spinner = () => {
    return `
    <div class="d-flex justify-content-center" id="loadingSpinner">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
    `
}

const selectOwner = (ownerID) => {
    let exists = false;
    filterSort.owners.forEach(v => {
        if (v === ownerID) {
            exists = true;
        }
    })

    if (!exists) {
        filterSort.owners.push(ownerID);
    }
}

const removeOwner = (ownerID) => {
    filterSort.owners.forEach( (v,idx) => {
        if (v === ownerID) {
            filterSort.owners.splice(idx, 1);
            return false;
        }
    })
}

const bindToEvents = () => {
    // bind to filterSort event
    $("#personalAlbumCheck").on("change", () => {
        filterSort.personalAlbums =  $("#personalAlbumCheck").prop('checked');

        doReq();
    });
    
    $("#sharedAlbumCheck").on("change", () => {
        filterSort.sharedAlbums = $("#sharedAlbumCheck").prop('checked');
        
        doReq();
    });

    $("#startDate").on('change', (e) => {
        filterSort.date.start = $(e.target).val();

        doReq();
    });
    
    $("#endDate").on('change', (e) => {
        filterSort.date.end = $(e.target).val();

        doReq();
    });

    $("#ownerFilter").on("change","input", (e) => {
        checkElem = $(e.target)[0];

        if ( $(checkElem).prop('checked') ) {
            selectOwner($(checkElem)[0].id);
        } else {
            removeOwner($(checkElem)[0].id);
        }

        doReq();
    });

    $("#selectedOwnersFilter").on('click', '.btn-close', (e) => {
        let parents = $(e.target).parents("div");
        let id = $(parents[0]).find('input').val();
        
        filterSort.owners.forEach((v,idx) => {
            if ( v.id === id ) {
                filterSort.owners.splice(idx, 1);
            }
        });

        if ( filterSort.owners.length === 0 ) {
            $("#selectOwner option:eq(0)").prop('selected', true);
        }

        renderFilter();

        doReq();
    });

    $('#sortSelect').on('change', (e) => {
        filterSort.sort = $(e.target).val();

        doReq();
    });

    $('.main-container .col-filter .btn-close').on('click', () => {
       $('#filter-container').css('visibility','hidden'); 
        $('.main-container .col-filter .btn-close').css('visibility', 'hidden');
    });

    $('#btn-show-filter').on('click', () => {
        $('#filter-container').css('visibility', 'visible');
        $('.main-container .col-filter .btn-close').css('visibility', 'visible');
    });

}

const bindToCardOwner = () => {
    $('.container-album-card .card-header .row-owner').on('click',(e) => {
        ownerName = $(e.target).html().trim();

        if (ownerName == store.username) {
            filterSort.personalAlbums = true;
            filterSort.sharedAlbums = false;
            filterSort.owners = [];
        } else {
            $("#ownerFilter .form-check").each((_, e) => {
                if ($(e).find('label').text() === ownerName) {
                    id = $(e).find('input').val()
                    
                    filterSort.owners = [id];
                    filterSort.personalAlbums = false;
                }
            });
        }

        renderFilter();

        doReq();
    });
}

$(() => {
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
    
    $('#datepicker-mobil').datepicker({
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

    // init filterSort obj
    init();

    // get albums from server
    doReq();
    
    // bind to filterSorts controls
    bindToEvents();

});
