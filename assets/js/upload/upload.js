$(function () {
    $.widget('upload.uploadform', {
        options: {
            // input file field
            fileInput: undefined,
            // Define a list of dictionaries having the id of the fileUI widget, fileUI widget and file object
            filesUI: {},
            // fileUi container
            fileUIContainer: undefined,
            // filter the new add files against already added files,
            xhrOptions: undefined,

            jqXHR: undefined,

            filter: function (files) {
                let that = this;
                if (that.files === undefined) {
                    return files;
                } else {
                    let newFiles = [];
                    $.each(files, function (index, file) {
                        let isNewFile = true;
                        for (let i = 0; i < that.files.length; i++) {
                            if (that.files[i].name === file.name) {
                                isNewFile = false;
                            }
                        }
                        if (isNewFile) {
                            newFiles.push(file);
                        }
                    });
                    return newFiles;
                }
            }
        },
        _create: function (e) {
            let options = this.options;
            if (options.fileInput === undefined) {
                options.fileInput = this.element.find('input[type=file]');
                options.fileUIContainer = this.element.find('.file-list-container');
                this._on(this.options.fileInput, {
                    change: this._onChange
                });
                this._on(this.element.find('#submitButton'), {
                    click: function (e) {
                        e.preventDefault();
                        this._submit();
                    }
                });
                this._on(this.element.find('#abortButton'), {
                    'click': this._abort
                })
            }
        },

        /**
         * Bind to input change event
         */
        _onChange: function (event) {
            let that = this,
                options = this.options,
                data = {
                    fileInput: $(event.target)
                };
            let newFiles = $.makeArray(data.fileInput.prop('files'));
            if (newFiles.length > 0) {
                $.when(options.filter(newFiles)).then(function (newFiles) {
                        $.map(newFiles, function (file) {
                            let element = that._createUI(file);
                            that.options.filesUI[element.fileui('option', 'id')] = element;
                        });
                    }
                );
            }
        },

        /**
         * Return the newly created element for file
         */
        _createUI: function (file) {
            let that = this,
                options = this.options;
            let newElement = $("<div></div>").fileui({
                'filename': file.name,
                'file': file
            });
            that._on(newElement.fileui(), {
                'fileuidelete': that._deleteUI
            });
            newElement.appendTo(options.fileUIContainer);
            return newElement;
        },

        // Delete UI
        _deleteUI: function (event, data) {
            let options = this.options;
            $.each(options.filesUI, function (id, entry) {
                if (id === data.id) {
                    entry.fileui('destroy');
                    delete options.filesUI[id];
                    return false;
                }
            });
        },

        _initOptions: function (options) {
            if (!options) {
                options = {}
            }
            options.headers = {};
            options.type = {};
            options.data = {};
            options.url = "";
        },
        // create the ajax settings for signing the files
        _initDataforSigning: function (options) {
            this._initOptions(options);
            options.headers['Content-Type'] = 'application/json';
            options.type = 'POST';
            options.url = '/sign-s3';
            tempData = {};
            $.each(options.filesUI, (idx, value) => {
                tempData[value.fileui('option', 'id')] = {
                    'filename': value.fileui('option', 'filename'),
                    'filetype': value.fileui('option', 'file').type
                };
            });
            options.data = JSON.stringify(tempData)
        },
        _submit: function () {
            let self = this,
                o = this.options;
            self._initDataforSigning(o);
            this.jqXHR = $.ajax(o);
            this.jqXHR.done(function (result, textStatus, jqXHR) {
                self._initDataForAws(result);
            });
            this.jqXHR.then(function () {
                $.each(self.options.filesUI, (id, obj) => {
                    obj.fileui('send').then(function() {
                        alert('done');
                    })
                });
            })
        },
        _abort: function () {
            if (this.jqXHR) {
                return this.jqXHR.abort();
            }
        },
        _initDataForAws: function (signed_urls) {
            let o = this.options;
            for (let key in signed_urls) {
                if (key in o.filesUI) {
                    let item = o.filesUI[key];
                    item.fileui('setSignedUrl', signed_urls[key]);
                }
            }
        },

    });
});
