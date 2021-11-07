$(function () {
    $.photos = {
        data: {
            gallery: null,
            currentSelectedImages: [],
            baseURL: '/api/albums'
        },

        init: function() {
            $('.container .row-edit-tools input').prop('checked', false);
            this.reset();
            this.data.gallery = this.createGallery();
            this.bind();
        },

        reset: function() {
            this.data.currentSelectedImages = [];
            $('.container .row-edit-tools #btn-delete-photos').prop('disabled', true);
            $('.container .row-edit-tools #btn-set-thumbnail').prop('disabled', true);
        },

        bind: function() {
            const self = this;

            $('.container .row-edit-tools input').change(function () {
                if ($(this).is(':checked')) {
                    self.data.gallery.destroy();

                    // add class to each photo in the container
                    $('.container-photos #album-gallery img').each((_, elem) => {
                        $(elem).attr('class', 'selectable');
                        $(elem).on('click', null, self.data, self.onSelectImage);
                    })
                    
                } else {
                    self.data.gallery = self.createGallery(); 
                    
                    // add class to each photo in the container
                    $('.container-photos #album-gallery img').each((_, elem) => {
                        $(elem).removeAttr('class');
                        $(elem).parent().children('i').addClass('hidden');
                        $(elem).off('click');
                    });

                    self.reset();
                }
            });

            $('.container .row-edit-tools #btn-set-thumbnail').on('click', null, self.data, self.onSetThumbnail);
            $('.container .row-edit-tools #btn-delete-photos').on('click', null, self.data, self.onDeletePhotos);
        },

        // create gallery
        createGallery: function() {
            return window.lightGallery(
              document.getElementById("album-gallery"),
              {
                autoplayFirstVideo: false,
                pager: true,
                galleryId: "photos",
                plugins: [lgThumbnail],
                licenseKey: '0000-0000-000-0000',
                mobileSettings: {
                  controls: true,
                  showCloseIcon: true,
                  download: true,
                  rotate: true
                }
              }
            );
        },

        onSelectImage: function(e) {
            e.preventDefault();

            let data = e.data;

            if ($(e.currentTarget).hasClass('selected')) {
                $(e.currentTarget).parent().children('i').addClass('hidden');
                $(e.currentTarget).removeClass('selected');
        
                data.currentSelectedImages.forEach((elem, idx) => {
                    if ( elem === e.currentTarget ) {
                        data.currentSelectedImages.splice(idx, 1);
                        return;
                    }
                });
            } else {
                $(e.currentTarget).addClass('selected');
                $(e.currentTarget).parent().children('i').removeClass('hidden');
        
                data.currentSelectedImages.push(e.currentTarget);
            } 
        
            if (data.currentSelectedImages.length > 0) {
                $('.container .row-edit-tools #btn-delete-photos').removeAttr('disabled');
            } else {
                $('.container .row-edit-tools #btn-delete-photos').prop('disabled', true);
            }
        
            if (data.currentSelectedImages.length !== 1) {
                $('.container .row-edit-tools #btn-set-thumbnail').prop('disabled', true);
            } else {
                $('.container .row-edit-tools #btn-set-thumbnail').removeAttr('disabled');
            }
        },

        onSetThumbnail: function(e) {
            var url = location.href;
            const parts = url.split('/');

            let mySpinner = $.spinner('test');

            img = e.data.currentSelectedImages[0];
            a = $(img).parent('a');
            const imgParts = a[0].getAttribute('href').split('/');

            axios({
                method: 'post',
                url: e.data.baseURL + '/' + parts[parts.length-1] + '/album/thumbnail',
                data: {
                    image: imgParts[imgParts.length-2],
                }
            }).then(function(response) {
                mySpinner.remove();
            
                $.alert('Album cover set', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'success',
                    isOnly: false
                });
            });
        },

        onDeletePhotos: function(e) {
            console.log("delete photos");
        }
    }

    // init photo
    $.photos.init();
});


