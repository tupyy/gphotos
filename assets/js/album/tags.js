$(function() {
    $.tags = {
        init: function() {
            this.bind();
        },

        bind: function() {
            const self = this;

            $('.container .container-tags .tag button').on('click', (e) => {
                let parent = $(e.target).parent('span');
                const albumID = $(parent).find('#albumID').val();
                const tagID = $(parent).find('#tagID').val();

                self.dissociate(albumID, tagID);
            });
        },

        dissociate: function(albumID, tagID) {
            axios({
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
            }).catch(function(e) {
                $.alert('Error removing tag', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'error',
                    isOnly: false
                });
            });
        }
    }

    // init tags
    $.tags.init();
});
