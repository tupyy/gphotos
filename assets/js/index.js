const albumsElementID = "#albums";
const baseURL = "/api/albums"

const defaultPageSize = 20;

let store = {
    offset: 0,
    limit: defaultPageSize,
    countAlbums: 0,
}

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

        // add offset
        reqUrl = reqUrl + '&offset=' + store.offset;

        // add limit
        reqUrl = reqUrl + '&limit=' + store.limit;

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
            store.totalAlbums = response.data.count;

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
    $("#pagination").html(renderPagination(store.limit,store.offset, store.totalAlbums));
}

let renderAlbum = (album) => {
    return `
        <div class="album-col col-6" id="` + album.id + `">
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

const renderPagination = (pageSize, offset, total) => {
    numberPages = Math.floor(total / pageSize);
    if (total % pageSize > 0) {
        numberPages++;
    }

    currentPage = Math.floor(offset / pageSize) + 1;
    let prevDisabled = '';
    if (currentPage == 1) {
        prevDisabled = 'disabled';
    }

    let nextDisabled = '';
    if (currentPage == numberPages) {
        nextDisabled = 'disabled';
    }

    head =  `
        <nav aria-label="Page navigation">
          <ul class="pagination">
            <li class="page-item ` + prevDisabled + `">
              <a class="page-link `+prevDisabled+`" href="javascript:doPagination(false)" aria-label="Previous">
                <span aria-hidden="true">&laquo;</span>
                <span class="sr-only">Previous</span>
              </a>
            </li>`
    footer = `
            <li class="page-item ` + nextDisabled + `">
              <a class="page-link ` + nextDisabled + `" href="javascript:doPagination(true)" aria-label="Next">
                <span aria-hidden="true">&raquo;</span>
                <span class="sr-only">Next</span>
              </a>
            </li>
          </ul>
        </nav>
    `
    let body = '';
    for (let i = 1; i <= numberPages; i++) {
        if (currentPage == i) {
            body += '<li class="page-item active"><a class="page-link active" href="javascript:gotoPage('+i+')">' + i + '</a></li>'
        } else {
            body += '<li class="page-item"><a class="page-link" href="javascript:gotoPage('+i+')">' + i + '</a></li>'
        }
    }

    return head + body + footer
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

const clearPagination = () => {
    $('#pagination').empty();
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

const doPagination = (increase) => {
    if (increase) {
        store.offset += store.limit;
    } else {
        store.offset -= store.limit;
    }

    store.albums = [];

    doReq();

    clearAlbums();
    clearPagination();
    render();
}

const gotoPage = (page) => {
    if (page <= 0) {
        return;
    }

    // backend is base 0
    store.offset = defaultPageSize * (page-1);
    store.albums = [];

    doReq();

    clearAlbums();
    clearPagination();
    render();
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
