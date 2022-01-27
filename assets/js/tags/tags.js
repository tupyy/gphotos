$(function () {
    $.tags = {
        addTagModal: undefined,

        init: function() {
            this.bind();
        
            jscolor.presets.default = {    
                position: 'right',    
                closeButton: true,
                palette: [        
                    '#000000', '#7d7d7d', '#870014', '#ec1c23', '#ff7e26',        
                    '#fef100', '#22b14b', '#00a1e7', '#3f47cc', '#a349a4',        
                    '#ffffff', '#c3c3c3', '#b87957', '#feaec9', '#ffc80d',        
                    '#eee3af', '#b5e61d', '#99d9ea', '#7092be', '#c8bfe7',    
                ],    //paletteCols: 12,    //hideOnPaletteClick: true,};
            };
            
            this.addTagModal = new bootstrap.Modal(document.getElementById('addTag'), {
                keyboard: false
            });
        },

        bind: function() {
            $('#submitButton').on('click', (e) => {
                this.submit();
                e.preventDefault();
            });

            $('.container-buttons .btn-outline-primary').on('click', (e) => {
                $('#tagID').val('');
                $('#tagName').val('');
                
                this.addTagModal.toggle();
            }),

            $('.row-tag .btn-outline-primary').on('click', (e) => {
                const parent = $(e.target).parents('.row-tag')
                const id = $(parent[0]).find('input').val();
                const name = $(parent[0]).find('span').html();
                const color = $(parent[0]).find('i').css('color');

                $('#tagName').val(name);
                let colorElement = $('#tagColor');
                colorElement[0].jscolor.fromString(this.rgb2hex(color));
                $('#tagID').val(id)

                this.addTagModal.toggle();
            });

            $('.row-tag .btn-outline-danger').on('click', (e) => {
                const parent = $(e.target).parents('.row-tag')
                const id = $(parent).find('input').val();

                this.deleteTag(id);
            });
        },

        submit: function(isUpdate) {
            const name = $("#tagName").val();
            const color = $("#tagColor").val();
            const tagID = $("#tagID").val();

            if (name === '' || color === '') {
                $("#addTagForm").addClass("was-validated");
                return;
            }

            let hexColor = color;

            if (hexColor.indexOf('rgb') == 0) {
                hexColor = this.rgb2hex(color);
            }

            let promise;
            if (tagID == '') {
                promise = this.createTag(name, hexColor);
            } else {
                promise = this.updateTag(tagID, name, hexColor);
            }

            promise.then((response) => {
                this.addTagModal.hide();
                location.reload();
            }).catch ((error) => {
                if (error.response.status == 400) {
                    $('#tagNameFeedback').html(error.response.data.message);
                    $('#tagName').addClass('is-invalid');
                } else {
                    $.alert('Error creating tag', {
                        closeTime: 2000,
                        autoClose: true,
                        position: ['top-left'],
                        withTime: false,
                        type: 'error',
                        isOnly: false
                    });

                    this.addTagModal.hide();
                    location.reload();
                }
            });
        },

        createTag: function(name, color) {
            return axios.post('/api/tags', {
                name: name,
                color: color
            });
        },

        updateTag: function(id, name, color) {
            return axios.patch('/api/tags/' + id, {
                name: name,
                color: color
            });

        },

        deleteTag: function(id) {
            axios({
                method: 'delete',
                url:'/api/tags/' + id,
            }).then(function(response) {
                $.alert('Tag removed', {
                    closeTime: 2000,
                    autoClose: true,
                    position: ['top-left'],
                    withTime: false,
                    type: 'success',
                    isOnly: false
                });
                location.reload();
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
        },

        rgb2hex: function(rgb) {
            rgb = rgb.match(/^rgb\((\d+),\s*(\d+),\s*(\d+)\)$/);

            function hex(x) {
              return ("0" + parseInt(x).toString(16)).slice(-2);
            }
            
            // Check if rgb is null
            if (rgb == null ) {
              // You could repalce the return with a default color, i.e. the line below
              // return "#ffffff"
              return "#ffffff";
            }
            
            return "#" + hex(rgb[1]) + hex(rgb[2]) + hex(rgb[3]);
        }
    }

    $.tags.init();
});
