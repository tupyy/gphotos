$(function () {
    $.widget('upload.uploadform', {
        options: {
            // input file field
            fileInput: undefined,
            // Define a list of dictionaries having the id of the fileUI widget, fileUI widget and file object
            filesUI: new Map(),
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

                            element.fileui('setUrl',that._getUrl());
                            element.fileui('setCallback', that._onFinishedUpload());
                            that.options.filesUI.set(element.fileui('option', 'id'),element);
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
        _deleteUI: function (_, data) {
            let options = this.options;
            $.each(options.filesUI, function (id, entry) {
                if (id === data.id) {
                    entry.fileui('destroy');
                    delete options.filesUI[id];
                    return false;
                }
            });
        },
        _submit: function () {
            if (this.options.filesUI.size > 0) {
                const interator = this.options.filesUI.keys();

                let obj = this.options.filesUI.get(interator.next().value);

                obj.fileui('send', this._onFinishedUpload());
            }
        },
        _onFinishedUpload: function() {
            const options = this.options
            
            return (id) => {
                options.filesUI.delete(id);
                
                if (options.filesUI.size > 0) {
                    const interator = options.filesUI.keys();

                    let obj = options.filesUI.get(interator.next().value);

                    obj.fileui('send');
                }
            }
        },
        _abort: function () {
            if (this.jqXHR) {
                return this.jqXHR.abort();
            }
        },
        _getUrl: function() {
            let url = window.location.href;
            let parts = url.split('/');
            let id = parts[parts.length - 1];
            return '/api/albums/' + id + '/album/upload'
        }
    });
});
