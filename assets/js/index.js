const albumsElementID = "#albums";
const baseURL = "/api/albums"

const defaultPageSize = 20;

let store = {
    offset: 0,
    limit: defaultPageSize,
    countAlbums: 0,
}

let queryParams = {
    personalAlbums: true,
    sharedAlbums: true,
    date: {
        start: '',
        end: '',
    },
    searchExpression: '',
    owners: [],
    sort: '',
    buildRequestURL: function(baseURL) {
        let reqUrl = baseURL + "?personal=" + this.personalAlbums + "&shared=" + this.sharedAlbums;
        
        if (this.owners.length > 0) {
            this.owners.forEach(id => {
                reqUrl = reqUrl + "&owner=" + id
            });
        }

        if (this.searchExpression !== "") {
            console.log(encodeURI(this.searchExpression))
            reqUrl = reqUrl + "&filter=" + $.base64.btoa(encodeURI(this.searchExpression));
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
    queryParams.personalAlbums =  $("#personalAlbumCheck").prop('checked');
    queryParams.sharedAlbums = $("#sharedAlbumCheck").prop('checked');
    queryParams.sort = $('#sortSelect').val();

    // init controls
    $("#searchBar").val('');
}

const search = () => {
    queryParams.searchExpression = $("#searchBar").val();

    fetch();

    store.albums = [];

    clearAlbums();
    clearPagination();

    render();
}

const fetch = () => {
    showSpinner($("#albums"),true);
    $("#count_albums").parent().hide();

    clearAlbums();

    axios.get(queryParams.buildRequestURL(baseURL))
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
    let tags = ""
    
    if (typeof album.tags != 'undefined') {
        for(const [k,v] of Object.entries(album.tags)) {
            color = "black";
            if (v.color !== '') {
                color = v.color;
            }
            tags += `<span class="album-tag" style="background:` + color + `">
                <i class="fas fa-tag"></i>`+v.name+
                `</span>`;
        };
    }

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
                <div class="album-image">
                    <div class="row-tags">
                    ` + tags +`
                    </div>
                    <a href="/album/` + album.id + `">
                        <img src="` + album.thumbnail + `" class="card-img-top"/>
                    </a>
                </div>
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

    if (numberPages === 0) {
        return ""
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
        $("#personalAlbumCheck").prop('checked', queryParams.personalAlbums);
        $("#sharedAlbumCheck").prop('checked', queryParams.sharedAlbums);

        
        $("#ownerFilter .form-check").each((_, e) => {
            id = $(e).find('input').val();
            find = false;
            queryParams.owners.forEach((i) => {
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

// parse the text of searchBar into key=value. If the key is not provided, use "any".
const parseSearchInput = () => {
    let val = $(searchBar).val();
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
    queryParams.owners.forEach(v => {
        if (v === ownerID) {
            exists = true;
        }
    })

    if (!exists) {
        queryParams.owners.push(ownerID);
    }
}

const removeOwner = (ownerID) => {
    queryParams.owners.forEach( (v,idx) => {
        if (v === ownerID) {
            queryParams.owners.splice(idx, 1);
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

    fetch();

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

    fetch();

    clearAlbums();
    clearPagination();

    render();
}

const bindToEvents = () => {
    // bind to queryParams event
    $("#personalAlbumCheck").on("change", () => {
        queryParams.personalAlbums =  $("#personalAlbumCheck").prop('checked');

        fetch();
    });
    
    $("#sharedAlbumCheck").on("change", () => {
        queryParams.sharedAlbums = $("#sharedAlbumCheck").prop('checked');
        
        fetch();
    });

    $("#startDate").on('change', (e) => {
    });
    
    $("#endDate").on('change', (e) => {
        queryParams.searchExpression ='date > ' + $("#startDate").val() + ' & date < ' + $(e.target).val() + '';

        fetch();
    });

    $("#searchBar").on('keypress', (e) => {
        if (e.which == 13) {
            search();
        }

        $(".search-group button").removeClass('hidden');
    })

    $(".search-group button").on('click', (e) => {
        $("#searchBar").val('');
        $(".search-group button").addClass('hidden');

        search();
    })

    $("#ownerFilter").on("change","input", (e) => {
        checkElem = $(e.target)[0];

        if ( $(checkElem).prop('checked') ) {
            selectOwner($(checkElem)[0].id);
        } else {
            removeOwner($(checkElem)[0].id);
        }

        fetch();
    });

    $("#selectedOwnersFilter").on('click', '.btn-close', (e) => {
        let parents = $(e.target).parents("div");
        let id = $(parents[0]).find('input').val();
        
        queryParams.owners.forEach((v,idx) => {
            if ( v.id === id ) {
                queryParams.owners.splice(idx, 1);
            }
        });

        if ( queryParams.owners.length === 0 ) {
            $("#selectOwner option:eq(0)").prop('selected', true);
        }

        renderFilter();

        fetch();
    });

    $('#sortSelect').on('change', (e) => {
        queryParams.sort = $(e.target).val();

        fetch();
    });

    $('.main-container .col-filter .btn-close').on('click', () => {
       $('#filter-container').css('visibility','hidden'); 
        $('.main-container .col-filter .btn-close').css('visibility', 'hidden');
        $('.search-group').css('visibility', 'visible');
        $('#pagination').css('visibility', 'visible');
    });

    $('#btn-show-filter').on('click', () => {
        $('#filter-container').css('visibility', 'visible');
        $('.main-container .col-filter .btn-close').css('visibility', 'visible');
        $('.search-group').css('visibility', 'hidden');
        $('#pagination').css('visibility', 'hidden');
    });
}

const bindToCardOwner = () => {
    $('.container-album-card .card-header .row-owner').on('click',(e) => {
        ownerName = $(e.target).html().trim();

        if (ownerName == store.username) {
            queryParams.personalAlbums = true;
            queryParams.sharedAlbums = false;
            queryParams.owners = [];
        } else {
            $("#ownerFilter .form-check").each((_, e) => {
                if ($(e).find('label').text() === ownerName) {
                    id = $(e).find('input').val()
                    
                    queryParams.owners = [id];
                    queryParams.personalAlbums = false;
                }
            });
        }

        renderFilter();

        fetch();
    });
}

$(() => {
    $('#datepicker-start').datepicker({
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
    
    $('#datepicker-end').datepicker({
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

    // init queryParams obj
    init();

    // get albums from server
    fetch();
    
    // bind to queryParamss controls
    bindToEvents();

});
