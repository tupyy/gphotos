$(function() {
    $.tags = {
        addTagModal: undefined,
        tags: {},
        tagParentElement: undefined,

        init: function() {
            this.bind();
            
            this.addTagModal = new bootstrap.Modal(document.getElementById('addTag'), {
                keyboard: false
            });

            const parts = location.href.split('/')
            this.albumID = parts[parts.length - 1];

            this.tagParentElement = $('.container-tags .row-tags');

            this.doReq();

        },

        // fetch them and render them
        doReq: function() {
            this.fetchTags(this.albumID).then(() => {
                this.render();
                this.bindCloseButtons();
            }).catch(error => {
                $.alert('Error fetching tags', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'error',
                    isOnly: false
                });
            });
        },

        bind: function() {
            $('.container .container-tags .tag-buttons .btn-associate').on('click', (e) => {
                if ( typeof this.addTagModal !== 'undefined' ) {
                    $('#addTagForm .modal-body').empty();

                    block = '';
                    $.each( this.tags, (_, tag) => {
                        if (tag.isOwner) {
                            block += `
                                <div class="form-check">
                                    <label class="form-check-label" for="` + tag.id +`">` + tag.name + `</label>
                                `
                            if (tag.associated) {
                                block += `
                                    <input type="checkbox" class="form-check-input" id="` + tag.id +`" value="" checked>
                                        <i class="fas fa-tag" style="color: ` + tag.color + `"></i>
                                    </input>
                                </div>
                                `
                            } else {
                                block += `
                                    <input type="checkbox" class="form-check-input" id="` + tag.id +`" value="">
                                        <i class="fas fa-tag" style="color: ` + tag.color + `"></i>
                                    </input>
                                </div>
                                `
                            }
                        }
                    });

                    $('#addTagForm .modal-body').html(block);
                    this.addTagModal.toggle();
                }
            });

            $('.container .container-tags .modal #submitButton').on('click', (e) => {
                e.preventDefault();

                promises = [];
                $('#addTagForm .modal-body').find('input').each( (_,input) => {
                    const tag = this.tags[input.id];
                    if (tag.associated != input.checked) {
                        promises.push(this.changeAssociation(this.albumID, tag.id, input.checked));
                    }
                });

                $.when(...promises).done(() => {
                    this.tags = {};
                    this.removeTags();

                    this.doReq(); 
                    this.addTagModal.toggle();
                })
                    
            });
        },

        bindCloseButtons: function() {
            $('.container-tags .row-tags .btn-close').on('click', (e) => {
                let parent = $(e.target).parent('span');
                const tagID = $(parent).find('.tag-id').val();

                this.dissociate(this.albumID, tagID).then(() => {
                    this.tags = {};
                    this.removeTags();

                    this.doReq(); 
                });
            });
        },

        fetchTags: function(albumID) {
            const deferred = $.Deferred();

            axios.all([axios.get('/api/albums/' + albumID + '/tags'),axios.get('/api/tags')]).
                then(axios.spread((albumTags, userTags) => {
                    $.each(albumTags.data.tags, (idx, tag) => {
                        this.tags[tag.id] = {
                            'name': tag.name,
                            'color': tag.color,
                            'id': tag.id,
                            'associated': true,
                            'isOwner': false
                        }
                    });

                    // true if the current user is the owner of the tag
                    $.each(userTags.data.tags, (_, userTag) => {
                        let tag = this.tags[userTag.id];

                        if (typeof tag !== 'undefined') {
                            tag.isOwner = true;
                        } else {
                            this.tags[userTag.id] = {
                                'name': userTag.name,
                                'color': userTag.color,
                                'id': userTag.id,
                                'associated': false,
                                'isOwner': true
                            }
                        }
                    });

                    deferred.resolve();
                })).catch(error => {
                    console.log(error);
                    deferred.reject(error);
                });
            return deferred;
        },

        render: function() {
            let tagBlock = '';
            $.each(this.tags, (_, tag) => {
                if (tag.associated) {
                    tagBlock = tagBlock + `
                        <div class="tag-container">
                            <span class="tag" style="background:` + tag.color + `">
                                <input type="hidden" class="tag-id" value="` + tag.id + `"/>
                                <i class="fas fa-tag"></i>
                                ` + tag.name
                    if (tag.isOwner) {
                        tagBlock += `
                                <button type="button" class="btn-close btn-close-white" aria-label="Close"></button>
                            </span>
                        </div>`
                    } else {
                        tagBlock += `
                            </span>
                        </div>
                        `
                    }
                }
            });

            $(this.tagParentElement).html(tagBlock);
        },

        removeTags: function() {
            $(this.tagParentElement).empty();
        },

        changeAssociation: function(albumID, tagID, associated) {
            if (associated) {
                return this.associate(albumID, tagID);
            } else {
                return this.dissociate(albumID, tagID)
            }
        },

        associate: function(albumID, tagID) {
            const deferred = $.Deferred();
            axios({
                method: 'post',
                url: '/api/albums/' + albumID + '/tag/' + tagID + '/associate'
            }).then(function(response) {
                $.alert('Tag added', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'success',
                    isOnly: false
                });
                deferred.resolve();
            }).catch(function(e) {
                $.alert('Error adding tag', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'error',
                    isOnly: false
                });
                deferred.reject();
            });

            return deferred;
        },

        dissociate: function(albumID, tagID) {
            const deferred = $.Deferred();
            return axios({
                method: 'delete',
                url: '/api/albums/' + albumID + '/tag/' + tagID + '/dissociate'
            }).then(function(response) {
                $.alert('Tag removed', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'success',
                    isOnly: false
                });
                deferred.resolve();
            }).catch(function(e) {
                $.alert('Error removing tag', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'error',
                    isOnly: false
                });
                deferred.reject();
            });

            return deferred;
        }
    }

    $.tags.init();
});
